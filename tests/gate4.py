#!/usr/bin/env python3
"""Run the Gate 4 acceptance fixtures in OpenJDK 21 and Catty."""

import argparse
from dataclasses import dataclass
import json
import os
from pathlib import Path
import re
import shutil
import sys
import tempfile

from acceptance import DATA, MAIN_DESC, ROOT, Suite, method_id, require
from classfile import ClassFile
from gate2 import java_environment
from gate3 import GATE3, check_toolchains, compile_sources, make_jar


FIXTURES = DATA / "gate4"
PROFILE = "catty-jvm21-gate4"
MODES = ("interp", "mixed", "aot")


@dataclass(frozen=True)
class Case:
    main: str
    mixed_method: tuple


CASES = {
    "lifecycle": Case("Lifecycle", ("Lifecycle$Worker", "run", "()V")),
    "monitors": Case("Monitors", ("Monitors", "nested", "()V")),
    "wait-set": Case("WaitSet", ("WaitSet$Trial$Waiter", "run", "()V")),
    "initialization": Case("Initialization", ("Initialization$Ready", "<clinit>", "()V")),
    "memory": Case("Memory", ("Memory$Publisher", "run", "()V")),
    "launcher": Case("Launcher", ("Launcher$Worker", "run", "()V")),
    "plugin": Case("app.Main", ("app/Main$Registrar", "run", "()V")),
}


def prepare(suite, label):
    if label == "plugin":
        stub = compile_sources(suite, "javac-stub", suite.classes("stub"),
                               sorted((GATE3 / "stub").rglob("*.java")))
        adapter = compile_sources(suite, "javac-adapter", suite.classes("adapter"),
                                  sorted((GATE3 / "adapter").rglob("*.java")))
        return compile_sources(suite, "javac-plugin-app", suite.classes(label),
                               sorted((FIXTURES / label / "app").rglob("*.java")), [stub]), adapter
    return compile_sources(suite, f"javac-{label}", suite.classes(label),
                           [FIXTURES / "Check.java", *sorted((FIXTURES / label).glob("*.java"))]), None


def fixture_runs(suite, label):
    if label == "plugin":
        worker = compile_sources(suite, "javac-plugin-worker", suite.classes("plugin-worker"),
                                 sorted((FIXTURES / label / "worker").rglob("*.java")))
        jar = make_jar(suite, "jar-plugin-worker", worker, suite.work / "plugin-worker.jar")
        return [("normal", [jar], b"plugin\n", 0, ())]
    if label == "launcher":
        return [
            ("worker-uncaught", [], b"main\nworker\ndescendant\n", 0,
             (rb'Exception in thread "Thread-[0-9]+" java\.lang\.RuntimeException',
              rb"Launcher\$Worker\.run")),
            ("main-uncaught", ["fail-main"], b"main\nworker\ndescendant\n", 1,
             (rb'Exception in thread "main" java\.lang\.RuntimeException',
              rb"Launcher\.main")),
        ]
    return [("normal", [], (label + "\n").encode(), 0, ())]


def check_result(result, stdout, status, diagnostics):
    require(result.status == status and result.stdout == stdout,
            f"expected status {status} and stdout {stdout!r}. See {result.log}")
    require(b"DATA RACE" not in result.stderr and b"panic:" not in result.stderr,
            f"Go runtime failure. See {result.log}")
    if diagnostics:
        require(result.stderr.count(b"Exception in thread ") == 1
                and all(re.search(pattern, result.stderr) for pattern in diagnostics),
                f"incorrect uncaught exception diagnostic. See {result.log}")
    else:
        require(not result.stderr, f"unexpected stderr. See {result.log}")


def build_image(suite, label, case, app, mode, race, expected_errors=()):
    tag = f"{label}-{mode}" + ("-race" if race else "")
    folder = suite.work / "images" / tag
    folder.mkdir(parents=True)
    binary, report = folder / "app", folder / "report.json"
    generated = folder / "source"
    args = [suite.bin, "build", "--enable-threads", "--classpath", app, "--main", case.main,
            "--mode", "auto" if mode == "mixed" else mode, "--report", report,
            "--generated-source", generated, "--output", binary]
    if mode == "mixed":
        owner, name, descriptor = case.mixed_method
        policy = folder / "policy.json"
        policy.write_text(json.dumps({"methods": [{"owner": owner, "name": name,
                           "descriptor": descriptor, "engine": "interp"}]}) + "\n", encoding="utf-8")
        args += ["--policy", policy]
    env = {"GOFLAGS": "-race", "CGO_ENABLED": "1"} if race else {"GOFLAGS": ""}
    build = suite.run(f"build-{tag}", args, env=env)
    if expected_errors:
        build.failed(*expected_errors)
        require(not binary.exists() and not report.exists(), f"rejected build published artifacts in {folder}")
        return
    build.succeeded()
    if race:
        metadata = suite.run(f"race-metadata-{tag}", ["go", "version", "-m", binary]).succeeded()
        require(b"-race=true" in metadata.stdout, f"{binary} lacks race instrumentation")

    data = json.loads(report.read_text(encoding="utf-8"))
    require((data["profile"], data["imageABI"], data["plugins"], data["threads"])
            == (PROFILE, 4, True, True), f"wrong Gate 4 profile in {report}")
    require({"monitorenter", "monitorexit"} <= set(data["opcodes"]),
            f"missing Gate 4 monitor opcodes in {report}")
    require({"threads", "monitors", "synchronized-memory"} <= set(data["capabilities"]),
            f"missing Gate 4 capabilities in {report}")
    require(all(record["owner"] != "ext/Worker" for record in data["methods"]),
            f"late plugin leaked into {report}")
    methods = {method_id(record): record["engine"] for record in data["methods"]}
    main = (case.main.replace(".", "/"), "main", MAIN_DESC)
    require(methods[main] == ("interp" if mode == "interp" else "aot")
            and methods[case.mixed_method] == ("aot" if mode == "aot" else "interp"),
            f"unexpected engine selection in {report}")
    return tag, binary


def check_trace(path, result, label, case, mode, run):
    events = [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines()]
    stacks, starts, ends = {}, {}, {}
    run_entries = set()
    for event in events:
        thread, kind = event["thread"], event["event"]
        require(event["pid"] == result.pid and type(thread) is int and thread > 0,
                f"invalid process or logical thread on {event} in {path}")
        if kind == "thread_start":
            parent = event.get("parentThread")
            require(thread not in starts and (parent is None or (parent in starts and parent not in ends))
                    and type(event["daemon"]) is bool, f"invalid start on {event} in {path}")
            starts[thread] = event
            stacks[thread] = []
            continue
        require(thread in starts and thread not in ends, f"event outside thread lifetime in {path}: {event}")
        stack = stacks[thread]
        if kind == "thread_end":
            require(not stack and event["result"] in ("terminated", "uncaught"),
                    f"invalid termination on {event} in {path}")
            ends[thread] = event["result"]
        elif kind == "enter":
            key = method_id(event) + (event["loaderId"],)
            require(event["depth"] == len(stack), f"wrong enter depth in {path}: {event}")
            stack.append(key)
            if method_id(event) == case.mixed_method:
                require(event["engine"] == ("aot" if mode == "aot" else "interp"),
                        f"wrong executing engine in {path}: {event}")
                run_entries.add(thread)
        elif kind in ("return", "unwind"):
            require(stack and stack.pop() == method_id(event) + (event["loaderId"],)
                    and event["depth"] == len(stack), f"unbalanced {kind} in {path}: {event}")
        elif kind == "catch":
            require(stack and stack[-1] == method_id(event) + (event["loaderId"],)
                    and all(type(event[k]) is int and event[k] >= 0 for k in ("throwPC", "handlerPC")),
                    f"invalid catch in {path}: {event}")

    roots = [thread for thread, event in starts.items() if event.get("parentThread") is None]
    require(len(roots) == 1 and len(starts) > 1 and run_entries, f"missing thread or method evidence in {path}")
    main = roots[0]
    require(not starts[main]["daemon"] and starts[main]["name"] == "main"
            and ends.get(main) == ("uncaught" if run == "main-uncaught" else "terminated"),
            f"incorrect main lifecycle in {path}")
    require(all(thread in ends or event["daemon"] for thread, event in starts.items()),
            f"non-daemon threads still active in {path}")
    uncaught_children = sum(result == "uncaught" and thread != main for thread, result in ends.items())
    require(uncaught_children == int(run == "worker-uncaught"), f"incorrect child failure count in {path}")

    if label == "monitors":
        nested_throws = [e["throwable"] for e in events
                         if e["event"] == "unwind" and method_id(e) == case.mixed_method]
        caught = [e["throwable"] for e in events
                  if e["event"] == "catch" and method_id(e) == ("Monitors", "main", MAIN_DESC)]
        require(len(nested_throws) == 1 and nested_throws[0] in caught,
                f"monitor unwinding lost throwable identity in {path}")
        monitor_ops = [e for e in events if e["event"] == "monitor" and "pc" in e]
        require({e["phase"] for e in monitor_ops} >= {"enter", "exit"}
                and any(e["phase"] == "enter" and e["depth"] >= 3 for e in monitor_ops)
                and all(type(e["pc"]) is int and e["pc"] >= 0 for e in monitor_ops),
                f"missing explicit monitor PC or reentry evidence in {path}")
    if label == "wait-set":
        resumes = [e for e in events if e["event"] == "monitor" and e["phase"] == "resume"]
        require({e["reason"] for e in resumes} >= {"notified", "timed_out", "interrupted"}
                and any(e["depth"] == 2 and e["reason"] == "interrupted" for e in resumes),
                f"missing wait outcome or restored recursion depth in {path}")
    if label == "initialization":
        for klass, failed in (("Initialization$Ready", False), ("Initialization$Broken", True)):
            records = [e for e in events if e["event"] in ("class_init", "class_init_wait")
                       and e["class"] == klass]
            require(len(records) >= 2 and records[0].get("phase") == "start"
                    and records[-1].get("phase") == "complete",
                    f"missing ordered initialization transitions for {klass} in {path}")
            begun, finished = records[0], records[-1]
            owner = begun["thread"]
            require(owner == finished["thread"] != main and finished.get("failed", False) == failed
                    and all(e["event"] == "class_init_wait" and e["ownerThread"] == owner
                            and e["thread"] != owner and e["loaderId"] == begun["loaderId"]
                            for e in records[1:-1]),
                    f"incorrect initialization ownership or outcome for {klass} in {path}")
    if label == "plugin":
        admitted = [e for e in events if e["event"] == "plugin_admission" and e["result"] == "accepted"]
        worker_runs = [e for e in events if e["event"] == "enter"
                       and method_id(e) == ("ext/Worker", "run", "()V")]
        loaders = {e["loaderId"] for e in admitted}
        plugin_locks = [e for e in events if e["event"] == "monitor" and e["phase"] == "enter"
                        and e.get("method") and method_id(e["method"]) == ("ext/Worker", "record", "()V")]
        require(len(admitted) == len(loaders) == len(worker_runs) == 2
                and len({e["thread"] for e in admitted}) == 2
                and {e["loaderId"] for e in worker_runs} == loaders
                and all(e["engine"] == "interp" and e["depth"] == 0
                        and starts[e["thread"]]["parentThread"] in run_entries for e in worker_runs),
                f"late plugin workers lost engine, loader or thread identity in {path}")
        require(len(plugin_locks) == len({e["object"] for e in plugin_locks}) == 2
                and {(e["method"]["loaderId"], e["thread"]) for e in plugin_locks}
                == {(e["loaderId"], e["thread"]) for e in worker_runs},
                f"plugin static monitors were shared across defining loaders in {path}")
    if label == "launcher":
        daemons = {thread for thread, event in starts.items() if event["daemon"]}
        require(len(daemons) == 1 and daemons.isdisjoint(ends),
                f"launcher must exit while the daemon is still waiting in {path}")
        descendants = [event for event in starts.values() if event.get("parentThread") in run_entries]
        require(len(descendants) == 1 and not descendants[0]["daemon"]
                and ends.get(descendants[0]["thread"]) == "terminated",
                f"launcher failed to account for a non-daemon started after main in {path}")


def exercise(suite, label, prepare_only, race):
    case = CASES[label]
    app, adapter = prepare(suite, label)
    images = [] if prepare_only else [build_image(suite, label, case, app, mode, False) for mode in MODES]
    race_images = [build_image(suite, label, case, app, mode, True) for mode in MODES] if race else []
    classpath = os.pathsep.join(str(path) for path in (adapter, app) if path)
    for name, args, stdout, status, diagnostics in fixture_runs(suite, label):
        oracle = suite.run(f"oracle-{label}-{name}",
                           ["java", "-Xverify:all", "-cp", classpath, case.main, *args], timeout=30)
        check_result(oracle, stdout, status, diagnostics)
        for mode, (tag, binary) in zip(MODES, images):
            trace = suite.work / f"{tag}-{name}.jsonl"
            result = suite.run(f"launch-{tag}-{name}", [binary, *args], cwd=binary.parent,
                               env={"CATTY_TRACE": str(trace), "PATH": "", "JAVA_HOME": "", "JDK_HOME": ""},
                               timeout=30)
            check_result(result, oracle.stdout, oracle.status, diagnostics)
            check_trace(trace, result, label, case, mode, name)
        for tag, binary in race_images:
            result = suite.run(f"launch-{tag}-{name}", [binary, *args], cwd=binary.parent,
                               env={"GORACE": "atexit_sleep_ms=0 halt_on_error=1"}, timeout=60)
            check_result(result, oracle.stdout, oracle.status, diagnostics)
    print(f"{label} passed", flush=True)


def profile_checks(suite, prepare_only):
    src = FIXTURES / "profile" / "BlockOnly.java"
    block = compile_sources(suite, "javac-profile-monitor", suite.classes("profile-monitor"), [src])
    class_file = ClassFile(block / "BlockOnly.class")
    code = class_file.methods[("main", MAIN_DESC)].code
    start = code.start
    # Keep the stored monitor reference and put an int on the operand stack at monitorenter.
    require(class_file.data[start:start + 4] == bytes.fromhex("2a594cc2"),
            "unexpected BlockOnly monitorenter sequence")
    bad = suite.classes("bad-monitor-operand")
    bad.mkdir(parents=True)
    class_file.data[start:start + 3] = bytes.fromhex("2a4c03")
    (bad / "BlockOnly.class").write_bytes(class_file.data)
    suite.run("oracle-bad-monitor-operand",
              ["java", "-Xverify:all", "-cp", bad, "BlockOnly"]).failed(r"VerifyError", r"monitorenter")
    if prepare_only:
        return

    src = DATA / "packaging" / "unsupported" / "Locked.java"
    locked = compile_sources(suite, "javac-legacy-synchronized", suite.classes("legacy-synchronized"), [src])
    for profile, enable_plugins in (("gate2", False), ("gate3", True)):
        suite.build(f"{profile}-reject-synchronized", [locked], "Locked",
                    expected_errors=(r"synchronized",), enable_plugins=enable_plugins)
        suite.build(f"{profile}-reject-monitor", [block], "BlockOnly",
                    expected_errors=(r"monitor",), enable_plugins=enable_plugins)
    build_image(suite, "bad-monitor-operand", Case("BlockOnly", ()), bad, "interp", False,
                expected_errors=(r"BlockOnly\.main", r"pc 3\b", r"reference"))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--java-home", type=Path, default=os.environ.get("JAVA_HOME"))
    parser.add_argument("--prepare-only", action="store_true", help="run only the OpenJDK 21 fixtures")
    parser.add_argument("--case", choices=sorted(CASES))
    parser.add_argument("--race", action="store_true", help="also run race-instrumented Catty images")
    parser.add_argument("--keep", action="store_true", help="keep logs and build artifacts on success")
    args = parser.parse_args()
    if args.prepare_only and args.race:
        parser.error("--race requires Catty execution")

    work = Path(tempfile.mkdtemp(prefix="catty-gate4-"))
    print(f"Gate 4 artifacts in {work}", flush=True)
    try:
        suite = Suite(work, env=java_environment(args.java_home))
        check_toolchains(suite)
        if not args.prepare_only:
            suite.run("build-catty", ["go", "build", "-o", suite.bin, "./cmd/catty"]).succeeded()
        if not args.case:
            profile_checks(suite, args.prepare_only)
        for label in ((args.case,) if args.case else CASES):
            exercise(suite, label, args.prepare_only, args.race)
    except (AssertionError, KeyError, OSError, ValueError) as error:
        print(f"Gate 4 failed. {error}\nArtifacts retained in {work}", file=sys.stderr)
        return 1
    if not args.keep:
        shutil.rmtree(work)
    print("Gate 4 OpenJDK 21 reference passed." if args.prepare_only else "Gate 4 differential tests passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
