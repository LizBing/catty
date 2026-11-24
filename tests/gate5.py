#!/usr/bin/env python3
"""Run the Gate 5 application-lifecycle fixtures in OpenJDK 21 and Catty."""

import argparse
from collections import Counter
from dataclasses import dataclass, field
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import tempfile

from acceptance import DATA, MAIN_DESC, ROOT, Suite, method_id, require
from gate2 import java_environment
from gate3 import GATE3, check_toolchains, compile_sources, make_jar


FIXTURES = DATA / "gate5"
PROFILE = "catty-jvm21-gate5"
ABI = 5
MODES = ("interp", "mixed", "aot")


@dataclass
class Case:
    main: str
    mixed: tuple = None
    paths: list = field(default_factory=list)
    extra_args: tuple = ()


CASES = {
    "startup": Case("Startup", ("Startup", "step", "(LStartup$Box;)V")),
    "natural": Case("Natural"),
    "registration": Case("Registration"),
    "hooks": Case("Hooks"),
    "exit": Case("ExitCase", ("ExitCase", "reflectiveExit", "(Ljava/lang/Runtime;I)V")),
    "exit-contention": Case("ExitContention"),
    "halt": Case("HaltCase"),
    "late-hook": Case("app.Main", None),
}


# ---------------------------------------------------------------- utilities

def posix_status(signed):
    return signed & 0xFF if os.name == "posix" else signed


def run_suffix(args):
    if not args:
        return "plain"
    return "-".join(re.sub(r"[^A-Za-z0-9_.-]", "_", str(arg))[:48] for arg in args)


def stderr_has(result, patterns):
    return all(re.search(pattern, result.stderr) for pattern in patterns)


def launch(suite, label, binary, args=(), *, timeout=30, env=None, isolated=True, trace=True):
    path = suite.work / "traces" / f"{label}.jsonl"
    path.parent.mkdir(parents=True, exist_ok=True)
    variables = {"GORACE": "atexit_sleep_ms=0 halt_on_error=1"}
    if trace:
        variables["CATTY_TRACE"] = str(path)
    if isolated:
        variables.update(PATH="", JAVA_HOME="", JDK_HOME="")
    variables.update(env or {})
    result = suite.run(label, [binary, *args], cwd=binary.parent, env=variables, timeout=timeout)
    return result, path


def run_allow_timeout(suite, label, command, *, timeout, env=None, cwd=None):
    """Run a child and capture stdout/stderr even when the watchdog fires."""
    options = os.environ.copy()
    for key in ("CATTY_TRACE", "CATTY_METRICS", "JAVA_TOOL_OPTIONS", "JDK_JAVA_OPTIONS",
                "_JAVA_OPTIONS", "CLASSPATH"):
        options.pop(key, None)
    options["LC_ALL"] = "C"
    options.setdefault("GOCACHE", str(suite.work / "go-cache"))
    options.update(suite.env)
    options.update(env or {})
    log = suite.work / "logs" / f"{label}.json"
    log.parent.mkdir(parents=True, exist_ok=True)
    details = {"command": [str(c) for c in command]}
    timed_out = False
    proc = subprocess.Popen([str(c) for c in command], cwd=cwd or ROOT, env=options,
                            stdin=subprocess.DEVNULL, stdout=subprocess.PIPE,
                            stderr=subprocess.PIPE)
    try:
        stdout, stderr = proc.communicate(timeout=timeout)
    except subprocess.TimeoutExpired:
        proc.kill()
        stdout, stderr = proc.communicate()
        timed_out = True
    details.update(pid=proc.pid, status=proc.returncode, timed_out=timed_out,
                   stdout=stdout.decode("utf-8", errors="replace"),
                   stderr=stderr.decode("utf-8", errors="replace"))
    log.write_text(json.dumps(details, indent=2) + "\n", encoding="utf-8")

    class Result:
        pass

    result = Result()
    result.pid = proc.pid
    result.status = proc.returncode
    result.stdout = stdout
    result.stderr = stderr
    result.log = log
    result.timed_out = timed_out
    return result


# ------------------------------------------------------------------- builds

def prepare(suite, label, case):
    if label == "late-hook":
        stub = compile_sources(suite, "javac-stub", suite.work / "stub",
                               sorted((GATE3 / "stub").rglob("*.java")))
        adapter = compile_sources(suite, "javac-adapter", suite.work / "adapter",
                                  sorted((GATE3 / "adapter").rglob("*.java")))
        app = compile_sources(suite, "javac-app-late-hook", suite.classes("late-hook-app"),
                              sorted((FIXTURES / "late-hook" / "app").rglob("*.java")), [stub])
        plugin = compile_sources(suite, "javac-plugin-late-hook", suite.classes("late-hook-plugin"),
                                 sorted((FIXTURES / "late-hook" / "plugin").rglob("*.java")))
        # The application image is built before this JAR is created.
        jar = make_jar(suite, "jar-late-hook", plugin, suite.work / "late-hook.jar")
        return [app], adapter, jar
    classes = compile_sources(suite, f"javac-{label}", suite.classes(label),
                              [FIXTURES / "Check.java", *sorted((FIXTURES / label).glob("*.java"))])
    return [classes], None, None


def policy_for(folder, case):
    owner, name, descriptor = case.mixed
    path = folder / "policy.json"
    path.write_text(json.dumps({"methods": [{"owner": owner, "name": name,
                       "descriptor": descriptor, "engine": "interp"}]}) + "\n", encoding="utf-8")
    return path


def build_image(suite, label, case, paths, mode, *, race=False, generated=False):
    suffix = "race" if race else "plain"
    tag = f"{label}-{mode}-{suffix}"
    folder = suite.work / "images" / tag
    folder.mkdir(parents=True)
    binary, report = folder / "app", folder / "report.json"
    source = folder / "source"
    args = [suite.bin, "build", "--enable-lifecycle",
            "--classpath", os.pathsep.join(str(p) for p in paths),
            "--main", case.main, "--mode", "auto" if mode == "mixed" else mode,
            "--report", report, "--output", binary]
    if generated:
        args += ["--generated-source", source]
    if mode == "mixed" and case.mixed:
        args += ["--policy", policy_for(folder, case)]
    env = {"GOFLAGS": "-race", "CGO_ENABLED": "1"} if race else {"GOFLAGS": ""}
    build = suite.run(f"build-{tag}", args, env=env)
    build.succeeded()
    if race:
        metadata = suite.run(f"race-metadata-{tag}", ["go", "version", "-m", binary]).succeeded()
        require(b"-race=true" in metadata.stdout, f"{binary} lacks race instrumentation")
    data = json.loads(report.read_text(encoding="utf-8"))
    require((data["profile"], data["imageABI"], data["lifecycle"], data["plugins"], data["threads"])
            == (PROFILE, ABI, True, True, True), f"wrong Gate 5 profile in {report}")
    require({"application-lifecycle", "shutdown-hooks", "application-diagnostics"}
            <= set(data["capabilities"]), f"missing Gate 5 capabilities in {report}")
    require({"monitorenter", "monitorexit"} <= set(data["opcodes"]),
            f"Gate 5 must inherit monitor opcodes in {report}")
    methods = {method_id(record): (record["engine"], record["reason"]) for record in data["methods"]}
    main = (case.main.replace(".", "/"), "main", MAIN_DESC)
    expected_main = "interp" if mode == "interp" else "aot"
    require(methods[main][0] == expected_main,
            f"{case.main}.main ran as {methods[main][0]} in {mode} mode: {report}")
    if case.mixed:
        expected_mixed = "interp" if mode != "aot" else "aot"
        require(methods[case.mixed][0] == expected_mixed,
                f"{case.mixed} ran as {methods[case.mixed][0]} in {mode} mode: {report}")
        if mode == "mixed":
            require("policy" in methods[case.mixed][1],
                    f"mixed policy evidence missing in {report}")
    if label == "late-hook":
        require(all(record["owner"] != "ext/HookFactory" and record["owner"] != "ext/PluginHook"
                    for record in data["methods"]),
                f"late plugin leaked into {report}")
    return tag, binary, report


# -------------------------------------------------------------- trace checks

def check_lifecycle_trace(path, result, args, *, timeout=False):
    require(path.is_file(), f"missing trace {path}. See {result.log}")
    events = [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines()]
    require(events, f"empty trace {path}")
    require(events[0]["event"] == "app_start", f"trace does not start with app_start: {path}")
    for event in events:
        require(type(event.get("pid")) is int and event["pid"] == result.pid,
                f"wrong process id in {path}: {event}")
    require(sum(1 for event in events if event["event"] == "app_start") == 1,
            f"duplicate app_start in {path}")
    start = events[0]
    require(start["profile"] == PROFILE and start["imageABI"] == ABI,
            f"wrong app_start profile in {path}: {start}")
    require(start["argc"] == len(args), f"wrong app_start argc in {path}: {start}")
    exits = [event for event in events if event["event"] == "process_exit"]
    if timeout:
        require(not exits, f"stalled shutdown published process_exit in {path}")
    else:
        require(len(exits) == 1 and events[-1] is exits[0],
                f"process_exit must be the one final event in {path}")
        terminal = exits[0]
        require(posix_status(terminal["status"]) == result.status,
                f"trace status {terminal['status']} does not encode observed {result.status} in {path}")

    starts, ends, stacks = {}, {}, {}
    for event in events:
        kind = event["event"]
        if kind == "thread_start":
            thread = event["thread"]
            require(type(thread) is int and thread > 0 and thread not in starts,
                    f"invalid or duplicate thread_start in {path}: {event}")
            if event.get("shutdownHook"):
                require("parentThread" not in event,
                        f"VM-started hook must not have an ordinary parent: {event}")
                require(type(event.get("daemon")) is bool, f"hook start lacks daemon: {event}")
            else:
                parent = event.get("parentThread")
                require(parent is None or (parent in starts and parent not in ends),
                        f"invalid parent thread in {path}: {event}")
            starts[thread] = event
            stacks[thread] = []
            continue
        if kind in ("app_start", "shutdown_begin", "process_exit"):
            continue
        require("thread" in event, f"event without thread in {path}: {event}")
        thread = event["thread"]
        require(thread in starts and thread not in ends,
                f"event outside a thread lifetime in {path}: {event}")
        if kind == "thread_end":
            require(not stacks[thread], f"thread_end with open frames in {path}: {event}")
            require(event["result"] in ("terminated", "uncaught"),
                    f"invalid thread_end result in {path}: {event}")
            ends[thread] = event["result"]
        elif kind == "enter":
            require(event["depth"] == len(stacks[thread]),
                    f"wrong enter depth in {path}: {event}")
            stacks[thread].append(method_id(event))
        elif kind in ("return", "unwind"):
            require(stacks[thread] and stacks[thread].pop() == method_id(event)
                    and event["depth"] == len(stacks[thread]),
                    f"unbalanced {kind} in {path}: {event}")
        elif kind == "catch":
            require(stacks[thread] and stacks[thread][-1] == method_id(event),
                    f"catch outside the active method in {path}: {event}")

    main = None
    for thread, event in starts.items():
        if not event.get("shutdownHook") and event.get("name") == "main" and not event["daemon"]:
            require(main is None, f"multiple main threads in {path}")
            main = thread
    require(main is not None, f"missing main thread in {path}")
    require(start["mainThread"] == main, f"app_start mainThread disagrees in {path}")
    if timeout:
        return events

    terminal = exits[0]
    reason = terminal["reason"]
    require(reason in ("natural", "exit", "halt"), f"invalid terminal reason in {path}")
    begins = [event for event in events if event["event"] == "shutdown_begin"]
    if reason == "natural":
        require(main in ends, f"natural shutdown left main unretired in {path}")
    if reason in ("natural", "exit"):
        require(len(begins) == 1, f"{reason} shutdown needs one shutdown_begin in {path}")
        require(begins[0]["reason"] == reason, f"shutdown_begin reason in {path}")
    else:
        require(len(begins) <= 1, f"halt published multiple shutdown_begin in {path}")
    if begins:
        begin = begins[0]
        if reason in ("natural", "exit"):
            require(begin["status"] == terminal["status"],
                    f"shutdown status and terminal status disagree in {path}")
        require(type(begin["hookCount"]) is int and begin["hookCount"] >= 0,
                f"shutdown_begin hookCount in {path}")
        hook_starts = [event for event in starts.values() if event.get("shutdownHook")]
        if hook_starts:
            first_hook = min(events.index(event) for event in hook_starts)
            require(events.index(begin) < first_hook,
                    f"shutdown_begin must precede hook starts in {path}")
        if reason == "exit":
            require("requesterThread" in begin,
                    f"exit shutdown_begin lacks requesterThread in {path}")
    if reason == "natural":
        require("requesterThread" not in terminal, f"natural process_exit has a requester in {path}")
    else:
        require("requesterThread" in terminal, f"{reason} process_exit lacks requesterThread in {path}")
        require(terminal["requesterThread"] in starts,
                f"unknown requester thread in {path}: {terminal}")
    if reason != "halt":
        for thread, event in starts.items():
            if event.get("shutdownHook"):
                require(thread in ends,
                        f"hook thread {thread} did not retire before process_exit in {path}")
    return events


def method_engine(events, key, engine):
    records = [event for event in events if event["event"] == "enter" and method_id(event) == key]
    require(records, f"missing enter event for {key}")
    require(all(event["engine"] == engine for event in records),
            f"{key} engine mismatch: {[(event['engine']) for event in records]}")
    return records


# ------------------------------------------------------------------- cases

def expect_timeout(label, args):
    return ((label == "exit-contention" and args and args[0] == "hook-exit")
            or (label == "natural" and args and args[0] == "wait"))


def run_specs(label, jar):
    if label == "startup":
        return [([], None), ([""], None), (["   "], None), (["-x"], None),
                (["héllo→"], None)]
    if label == "natural":
        return [([], None), (["fail"], 1), (["wait"], None), (["wait-halt"], 31)]
    if label == "hooks":
        return [([], 7)]
    if label == "registration":
        return [([], 0)]
    if label == "exit":
        return [(["worker0"], 0), (["worker7"], 7), (["workerneg"], -3), (["workerwide"], 300),
                (["reflect0"], 0), (["reflect7"], 7), (["reflectneg"], -3),
                (["reflectwide"], 300), (["monitor"], 23), (["clinit"], 7)]
    if label == "exit-contention":
        return [(["status"], 3), (["hook-exit"], None), (["hook-exit-halt"], 21)]
    if label == "halt":
        return [(["before"], 19), (["during"], 23)]
    if label == "late-hook":
        return [([str(jar)], 5)]
    return [([], 0)]


def expected_stderr(label, name):
    if label == "natural" and name == "fail":
        return (rb'Exception in thread "main" java\.lang\.RuntimeException',
                rb"Natural\.main")
    if label == "hooks":
        return (rb'Exception in thread "bad-hook" java\.lang\.RuntimeException',
                rb"Hooks\$Throwing\.run", rb"hook-boom")
    return ()


def forbidden_lines(label, name):
    if label == "exit":
        return (b"after-exit", b"finally", b"main-end")
    if label == "exit-contention":
        return (b"post-exit", b"second-after", b"main-after")
    if label == "halt":
        return (b"after",)
    return ()


def compare_run(suite, label, case, binary, mode, args, expected_status, oracle,
                *, timeout=False, sorted_stdout=False):
    if timeout:
        require(oracle.timed_out, f"oracle did not stall for {label} {args}")
        trace = suite.work / "traces" / f"launch-{label}-{mode}-{run_suffix(args)}.jsonl"
        trace.parent.mkdir(parents=True, exist_ok=True)
        result = run_allow_timeout(suite, f"launch-{label}-{mode}-{run_suffix(args)}",
                                   [binary, *args],
                                   timeout=5 if label == "natural" else 20, cwd=binary.parent,
                                   env={"PATH": "", "JAVA_HOME": "", "JDK_HOME": "",
                                        "CATTY_TRACE": str(trace),
                                        "GORACE": "atexit_sleep_ms=0 halt_on_error=1"})
        require(result.timed_out, f"Catty did not stall for {label} {args}. See {result.log}")
        require(oracle.stdout == result.stdout,
                f"stalled output differs for {label} {args}. See {result.log}")
        require(not any(marker in result.stdout.splitlines()
                        for marker in forbidden_lines(label, args[0] if args else "")),
                f"forbidden marker after stalled exit. See {result.log}")
        events = check_lifecycle_trace(trace, result, args, timeout=True)
        return result, events

    result, trace = launch(suite, f"launch-{label}-{mode}-{run_suffix(args)}",
                           binary, args, timeout=30)
    if expected_status is None:
        expected_status = oracle.status
    require(result.status == posix_status(expected_status),
            f"{label} {args} {mode}: status {result.status}, want {posix_status(expected_status)}. See {result.log}")
    if expected_status != oracle.status:
        require(posix_status(expected_status) == oracle.status,
                f"oracle status {oracle.status} disagrees with expectation {expected_status}")
    if sorted_stdout:
        require(Counter(oracle.stdout.splitlines()) == Counter(result.stdout.splitlines()),
                f"{label} {args} {mode}: stdout lines differ.\nOracle:\n{oracle.stdout!r}\nCatty:\n{result.stdout!r}\nSee {result.log}")
    else:
        require(result.stdout == oracle.stdout,
                f"{label} {args} {mode}: stdout {result.stdout!r}, want {oracle.stdout!r}. See {result.log}")
    patterns = expected_stderr(label, args[0] if args else "")
    if patterns:
        require(stderr_has(result, patterns),
                f"{label} {args} {mode}: missing diagnostic {patterns}. See {result.log}")
    else:
        require(not result.stderr, f"{label} {args} {mode}: unexpected stderr {result.stderr!r}. See {result.log}")
    require(not any(marker in result.stdout.splitlines()
                    for marker in forbidden_lines(label, args[0] if args else "")),
            f"{label} {args} {mode}: forbidden marker present. See {result.log}")
    require(b"DATA RACE" not in result.stderr,
            f"{label} {args} {mode}: Go data race. See {result.log}")
    events = check_lifecycle_trace(trace, result, args)
    return result, events


def case_checks(label, mode, events, args):
    if label == "startup":
        method_engine(events, ("Startup", "main", MAIN_DESC), "interp" if mode == "interp" else "aot")
        method_engine(events, ("Startup", "step", "(LStartup$Box;)V"),
                      "interp" if mode in ("interp", "mixed") else "aot")
        method_engine(events, ("Startup$Box", "bump", "()V"),
                      "interp" if mode == "interp" else "aot")
        step = [event for event in events if event["event"] == "enter"
                and method_id(event) == ("Startup", "step", "(LStartup$Box;)V")]
        bump = [event for event in events if event["event"] == "enter"
                and method_id(event) == ("Startup$Box", "bump", "()V")]
        require(step and bump and step[0]["thread"] == bump[0]["thread"],
                "mixed engine crossing lost the logical thread")
        if mode == "mixed":
            require(step[0]["depth"] == 1 and bump[0]["depth"] == 2,
                    "mixed crossing lost its nested Java frames")
    if label == "hooks":
        hook_starts = [event for event in events if event["event"] == "thread_start" and event.get("shutdownHook")]
        require(len(hook_starts) == 3 and {event["name"] for event in hook_starts}
                == {"hook-a", "hook-b", "bad-hook"}, f"wrong hook starts: {hook_starts}")
        ends = {event["thread"]: event["result"] for event in events if event["event"] == "thread_end"}
        bad = [event for event in hook_starts if event["name"] == "bad-hook"][0]
        require(ends.get(bad["thread"]) == "uncaught", "throwing hook was not diagnosed as uncaught")
        require(any(event["event"] == "thread_start" and event.get("name") == "hook-worker"
                    and not event["daemon"] for event in events),
                "hook-started ordinary worker is missing")
    if label == "natural" and not (args and (args[0] == "wait" or args[0] == "wait-halt")):
        terminal = [event for event in events if event["event"] == "process_exit"][0]
        require(terminal["reason"] == "natural", "natural case did not record natural shutdown")
    if label == "natural" and args and args[0] == "wait-halt":
        terminal = [event for event in events if event["event"] == "process_exit"][0]
        require(terminal["reason"] == "halt" and terminal["status"] == 31,
                f"halt after waiting hook: {terminal}")
    if label == "exit":
        variant = args[0]
        process = [event for event in events if event["event"] == "process_exit"][0]
        if variant == "monitor":
            require(process["reason"] == "halt" and process["status"] == 23,
                    f"monitor exit was not resolved by halt: {process}")
        elif variant == "clinit":
            require(process["reason"] == "exit" and process["status"] == 7,
                    f"clinit exit result: {process}")
        else:
            require(process["reason"] == "exit", f"explicit exit reason: {process}")
    if label == "late-hook":
        starts = [event for event in events if event["event"] == "thread_start" and event.get("shutdownHook")]
        require(len(starts) == 1 and starts[0]["name"] == "plugin-hook",
                f"late hook thread start: {starts}")
        enters = [event for event in events if event["event"] == "enter"
                  and method_id(event) == ("ext/PluginHook", "run", "()V")]
        require(len(enters) == 1 and enters[0]["engine"] == "interp" and enters[0]["loaderId"] != 1,
                f"plugin hook did not execute interpreted under a plugin loader: {enters}")
        require(enters[0]["thread"] == starts[0]["thread"],
                "plugin hook execution is not on its logical thread")
    if label == "exit-contention" and args[0] == "hook-exit-halt":
        process = [event for event in events if event["event"] == "process_exit"][0]
        require(process["reason"] == "halt" and process["status"] == 21,
                f"coordinated halt result: {process}")


def oracle_only(suite, label):
    case = CASES[label]
    paths, adapter, jar = prepare(suite, label, case)
    if label == "late-hook":
        classpath = os.pathsep.join(str(p) for p in (adapter, *paths))
    else:
        classpath = os.pathsep.join(str(p) for p in paths)
    for args, expected_status in run_specs(label, jar):
        timeout = expect_timeout(label, args)
        command = ["java", "-Xverify:all", "-Dstdout.encoding=UTF-8", "-Dstderr.encoding=UTF-8",
                 "-cp", classpath, case.main, *args]
        if timeout:
            result = run_allow_timeout(suite, f"oracle-{label}-{args[0]}", command,
                                       timeout=5 if label == "natural" else 20)
            require(result.timed_out, f"oracle did not stall for {label} {args}")
        else:
            result = suite.run(f"oracle-{label}-{'-'.join(args) or 'plain'}", command, timeout=30)
            if expected_status is not None:
                require(result.status == posix_status(expected_status),
                        f"oracle {label} {args}: status {result.status}, want {posix_status(expected_status)}. See {result.log}")
    print(f"{label} OpenJDK reference passed", flush=True)


def exercise(suite, label, race):
    case = CASES[label]
    paths, adapter, jar = prepare(suite, label, case)
    case.paths = paths
    images = {}
    for mode in MODES:
        images[mode] = build_image(suite, label, case, paths, mode)
    digests = {mode: hashlib.sha256(tag_binary[1].read_bytes()).hexdigest()
               for mode, tag_binary in images.items()}
    race_images = {}
    if race:
        for mode in MODES:
            race_images[mode] = build_image(suite, label, case, paths, mode, race=True)
    if label == "late-hook":
        classpath = os.pathsep.join(str(p) for p in (adapter, *paths))
    else:
        classpath = os.pathsep.join(str(p) for p in paths)
    for args, expected_status in run_specs(label, jar):
        timeout = expect_timeout(label, args)
        oracle_command = ["java", "-Xverify:all", "-Dstdout.encoding=UTF-8", "-Dstderr.encoding=UTF-8",
                 "-cp", classpath, case.main, *args]
        if timeout:
            oracle = run_allow_timeout(suite, f"oracle-{label}-{args[0]}", oracle_command,
                                       timeout=20)
            require(oracle.timed_out, f"oracle did not stall for {label} {args}")
        else:
            oracle = suite.run(f"oracle-{label}-{'-'.join(args) or 'plain'}", oracle_command,
                               timeout=30)
        if not timeout:
            if expected_status is None:
                expected_status = oracle.status
            require(oracle.status == posix_status(expected_status),
                    f"oracle status {oracle.status} disagrees with expectation {expected_status}. See {oracle.log}")
        sorted_stdout = label in ("registration", "hooks")
        for mode in MODES:
            tag, binary, _ = images[mode]
            result, events = compare_run(suite, label, case, binary, mode, args,
                                         expected_status, oracle,
                                         timeout=timeout, sorted_stdout=sorted_stdout)
            case_checks(label, mode, events, args)
        for mode in MODES:
            if mode not in race_images:
                continue
            tag, binary, _ = race_images[mode]
            result, events = compare_run(suite, label, case, binary, mode, args,
                                         expected_status, oracle,
                                         timeout=timeout, sorted_stdout=sorted_stdout)
            case_checks(label, mode, events, args)
    for mode, tag_binary in images.items():
        after = hashlib.sha256(tag_binary[1].read_bytes()).hexdigest()
        require(after == digests[mode],
                f"{label} {mode}: executable changed during isolated runs")
    print(f"{label} passed", flush=True)


# ----------------------------------------------------- generated code checks

def generated_checks(suite):
    """Inspect representative mixed/AOT generated bodies with the shared tool."""
    startup = CASES["startup"]
    paths, _, _ = prepare(suite, "startup", startup)
    _, binary, report = build_image(suite, "startup-generated", startup, paths, "mixed",
                                    generated=True)
    data = json.loads(report.read_text(encoding="utf-8"))
    records = {method_id(record): record for record in data["methods"]}
    source = binary.parent / "source"
    for key, check in (((("Startup", "main", MAIN_DESC)), "control-flow"),
                       (("Startup$Box", "bump", "()V"), "mutation")):
        record = records.get(key)
        require(record and record.get("generated"),
                f"missing generated record for {key} in {report}")
        generated = record["generated"]
        suite.run(f"inspect-{key[1]}", [
            "go", "run", str(ROOT / "tests" / "inspect_generated.go"),
            generated["symbol"], check, str(source / generated["file"])]).succeeded()
    require(records[("Startup", "step", "(LStartup$Box;)V")]["engine"] == "interp",
            f"mixed crossing target is not interpreted in {report}")

    hooks = CASES["hooks"]
    paths, _, _ = prepare(suite, "hooks", hooks)
    _, binary, report = build_image(suite, "hooks-generated", hooks, paths, "mixed",
                                    generated=True)
    data = json.loads(report.read_text(encoding="utf-8"))
    records = {method_id(record): record for record in data["methods"]}
    source = binary.parent / "source"
    for key, check in ((("Hooks$HookA", "run", "()V"), "monitors"),
                       (("Hooks$HookB", "run", "()V"), "control-flow")):
        record = records.get(key)
        require(record and record.get("generated"),
                f"missing generated record for {key} in {report}")
        generated = record["generated"]
        suite.run(f"inspect-{key[1]}-{key[0].split('/')[-1]}", [
            "go", "run", str(ROOT / "tests" / "inspect_generated.go"),
            generated["symbol"], check, str(source / generated["file"])]).succeeded()
    print("generated-code checks passed", flush=True)


# -------------------------------------------------------------- telemetry

TELEMETRY_CASE = Case("Telemetry", ("Telemetry", "crossing", "(LTelemetry$Box;)V"))
TELEMETRY_STDOUT = b"plugin\nrejected\n1\n2\nlookup-miss\ncaught\ncontended\ndone\n"


def prepare_telemetry(suite):
    stub = compile_sources(suite, "javac-telemetry-stub", suite.work / "telemetry-stub",
                           sorted((GATE3 / "stub").rglob("*.java")))
    adapter = compile_sources(suite, "javac-telemetry-adapter", suite.work / "telemetry-adapter",
                              sorted((GATE3 / "adapter").rglob("*.java")))
    app = compile_sources(suite, "javac-telemetry", suite.classes("telemetry"),
                          [FIXTURES / "telemetry" / "Telemetry.java"], [stub])
    plugin = compile_sources(suite, "javac-telemetry-plugin", suite.classes("telemetry-plugin"),
                             sorted((FIXTURES / "telemetry" / "plugin").rglob("*.java")))
    jar = make_jar(suite, "jar-telemetry-plugin", plugin, suite.work / "telemetry-plugin.jar")
    invalid = suite.work / "telemetry-not-a-jar"
    invalid.mkdir(exist_ok=True)
    return app, adapter, jar, invalid


def check_metrics(document, result, mode):
    require(document.get("schemaVersion") == 1, f"metrics schema: {document}")
    require(document.get("pid") == result.pid, f"metrics pid: {document.get('pid')} != {result.pid}")
    require(document.get("profile") == PROFILE and document.get("terminalReason") == "natural"
            and document.get("status") == 0, f"metrics identity: {document}")
    counters = document["counters"]
    require(counters["classesLoaded"] > 0, "metrics classesLoaded")
    engines = counters["methodsByEngine"]
    require(engines["interp"] > 0 and engines["intrinsic"] > 0,
            f"metrics methodsByEngine: {engines}")
    if mode != "interp":
        require(engines["aot"] > 0, f"metrics methodsByEngine: {engines}")
    require(counters["reflectiveInvocations"] == 2,
            f"reflectiveInvocations: {counters['reflectiveInvocations']}")
    require(counters["reflectionLookupMisses"] == 1,
            f"reflectionLookupMisses: {counters['reflectionLookupMisses']}")
    require(counters["javaExceptions"] >= 2, f"javaExceptions: {counters['javaExceptions']}")
    require(counters["monitorContentions"] == 1,
            f"monitorContentions: {counters['monitorContentions']}")
    require(counters["liveJavaThreads"] == 0, f"liveJavaThreads: {counters['liveJavaThreads']}")
    require(counters["peakJavaThreads"] >= 2, f"peakJavaThreads: {counters['peakJavaThreads']}")
    if mode == "mixed":
        require(counters["crossEngineCalls"] >= 1,
                f"mixed telemetry missed the engine crossing: {counters['crossEngineCalls']}")
    require(document["goRuntime"]["heapAllocBytes"] > 0, "metrics heap allocation")
    require("postGCHeapBytes" in document["goRuntime"] and document["goRuntime"]["postGCHeapBytes"] > 0,
            f"metrics post-GC heap observation: {document['goRuntime']}")
    require("Go-managed" in document["goRuntime"]["heapAllocScope"],
            f"metrics heap scope label: {document['goRuntime']}")


def telemetry_checks(suite, race):
    app, adapter, jar, invalid = prepare_telemetry(suite)
    classpath = os.pathsep.join(str(p) for p in (adapter, app))
    oracle = suite.run("oracle-telemetry", [
        "java", "-Xverify:all", "-Dstdout.encoding=UTF-8", "-Dstderr.encoding=UTF-8",
        "-cp", classpath, "Telemetry", str(jar), str(invalid)], timeout=30).succeeded()
    require(oracle.stdout == TELEMETRY_STDOUT,
            f"telemetry oracle output: {oracle.stdout!r}. See {oracle.log}")
    images = {(mode, variant): build_image(suite, "telemetry", TELEMETRY_CASE, [app], mode,
                                           race=(variant == "race"))
              for mode in MODES for variant in (("plain", "race") if race else ("plain",))}
    skip_stdout = TELEMETRY_STDOUT
    reference_classes = {}
    for (mode, variant), (tag, binary, _) in images.items():
        metrics = suite.work / "metrics" / f"{tag}.json"
        metrics.parent.mkdir(parents=True, exist_ok=True)
        result, _ = launch(suite, f"telemetry-{tag}", binary, [str(jar), str(invalid)],
                           env={"CATTY_METRICS": str(metrics)}, trace=False)
        require(result.status == 0 and result.stdout == skip_stdout,
                f"telemetry {tag}: status {result.status}, stdout {result.stdout!r}. See {result.log}")
        require(not result.stderr, f"telemetry {tag}: stderr {result.stderr!r}. See {result.log}")
        document = json.loads(metrics.read_text(encoding="utf-8"))
        check_metrics(document, result, mode)
        # A rejected admission must leave publication counters intact.
        skip_metrics = suite.work / "metrics" / f"{tag}-skip.json"
        skip = launch(suite, f"telemetry-skip-{tag}", binary, [str(jar), "skip"],
                      env={"CATTY_METRICS": str(skip_metrics)}, trace=False)[0]
        require(skip.status == 0 and skip.stdout == b"plugin\n1\n2\nlookup-miss\ncaught\ncontended\ndone\n",
                f"telemetry skip {tag}: {skip.stdout!r}. See {skip.log}")
        skip_document = json.loads(skip_metrics.read_text(encoding="utf-8"))
        require(skip_document["counters"]["classesLoaded"] == document["counters"]["classesLoaded"],
                f"rejected admission changed classesLoaded for {tag}")
        require(skip_document["counters"]["methodsByEngine"] == document["counters"]["methodsByEngine"],
                f"rejected admission changed methodsByEngine for {tag}")
        reference_classes[(mode, variant)] = document["counters"]["classesLoaded"]
    # Instrumentation must not change externally visible behavior.
    plain = images[("aot", "plain")][1]
    disabled = launch(suite, "telemetry-uninstrumented", plain, [str(jar), str(invalid)],
                      trace=False)[0]
    require(disabled.status == 0 and disabled.stdout == skip_stdout and not disabled.stderr,
            f"uninstrumented telemetry run changed behavior. See {disabled.log}")
    print("telemetry passed", flush=True)

# --------------------------------------------------------- profile checks

def profile_checks(suite, race):
    stub = compile_sources(suite, "javac-stub", suite.work / "stub",
                           sorted((GATE3 / "stub").rglob("*.java")))
    probe = compile_sources(suite, "javac-gate5-probe", suite.classes("gate5-probe"),
                            [FIXTURES / "profile" / "ExitProbe.java"])
    for flag, tag in ((("--enable-plugins",), "gate3"), (("--enable-threads",), "gate4"), ((), "gate2")):
        folder = suite.work / "profile-refusal" / tag
        folder.mkdir(parents=True)
        args = [suite.bin, "build", *flag, "--classpath", str(probe), "--main", "ExitProbe",
                "--mode", "auto", "--report", str(folder / "report.json"),
                "--output", str(folder / "app")]
        suite.run(f"build-{tag}-refusal", args).failed(
            r"unavailable bootstrap method", r"System\.")
        require(not (folder / "report.json").exists() and not (folder / "app").exists(),
                f"{tag} refusal published artifacts")

    # A late JAR that references a Gate 5-only bootstrap member is refused by an
    # earlier plugin profile at the admission boundary.
    app = compile_sources(suite, "javac-gate4-late-app", suite.classes("gate4-late-app"),
                          sorted((FIXTURES / "profile" / "gate4app").rglob("*.java")), [stub])
    bad = compile_sources(suite, "javac-gate4-late-bad", suite.classes("gate4-late-bad"),
                          sorted((FIXTURES / "profile" / "plugin").rglob("*.java")))
    bad_jar = make_jar(suite, "jar-gate4-late-bad", bad, suite.work / "gate4-late-bad.jar")
    folder = suite.work / "profile-late" / "gate4"
    folder.mkdir(parents=True)
    report, binary = folder / "report.json", folder / "app"
    suite.run("build-gate4-late-refusal", [
        suite.bin, "build", "--enable-threads", "--classpath", str(app), "--main", "app.Main",
        "--mode", "auto", "--report", str(report), "--output", str(binary)]).succeeded()
    trace = suite.work / "traces" / "gate4-late-refusal.jsonl"
    trace.parent.mkdir(parents=True, exist_ok=True)
    result = suite.run("run-gate4-late-refusal", [binary, str(bad_jar)],
                       cwd=binary.parent, env={"CATTY_TRACE": str(trace),
                                               "PATH": "", "JAVA_HOME": "", "JDK_HOME": ""})
    require(result.status != 0, f"Gate 4 admitted a Gate 5-only late class. See {result.log}")
    events = [json.loads(line) for line in trace.read_text(encoding="utf-8").splitlines()]
    rejected = [event for event in events if event["event"] == "plugin_admission"
                and event["result"] == "rejected"]
    require(rejected, f"Gate 4 did not record an admission rejection: {result.log}")
    print("profile checks passed", flush=True)


# ----------------------------------------------------------- setup failures

def setup_failure_checks(suite):
    # A gate5 executable with a hand-patched lifecycle flag must fail before Java.
    case = CASES["startup"]
    paths, _, _ = prepare(suite, "startup", case)
    _, binary, _ = build_image(suite, "setup", case, paths, "aot", generated=True)
    source = binary.parent / "source"
    patched = suite.work / "setup" / "patched"
    shutil.copytree(source, patched)
    main_go = patched / "main.go"
    text = main_go.read_text(encoding="utf-8")
    marker = "\tos.Exit(rt.RunApp(image, os.Args[1:]))"
    require(marker in text, "generated main.go shape changed")
    main_go.write_text(text.replace(marker, "\timage.Lifecycle = false\n" + marker), encoding="utf-8")
    bad_binary = patched / "app"
    env = os.environ.copy()
    env.setdefault("GOCACHE", str(suite.work / "go-cache"))
    env["GOWORK"] = "off"
    suite.run("build-patched-image", ["go", "build", "-trimpath", "-buildvcs=false",
                                      "-o", str(bad_binary), "."], cwd=patched, env=env).succeeded()
    trace = suite.work / "traces" / "setup-bad-flag.jsonl"
    result = suite.run("run-bad-image-flag", [bad_binary],
                       cwd=bad_binary.parent,
                       env={"CATTY_TRACE": str(trace), "PATH": "", "JAVA_HOME": "", "JDK_HOME": ""})
    require(result.status == 2, f"bad image flag status {result.status}. See {result.log}")
    require(result.stdout == b"", f"bad image flag ran Java: {result.stdout!r}")
    require(b"lifecycle flag" in result.stderr, f"bad image flag diagnostic: {result.stderr!r}")
    if trace.exists():
        require(b"app_start" not in trace.read_bytes(), "bad image flag emitted app_start")

    # Diagnostics output setup failures are status 2 with no Java entry.
    _, binary, _ = build_image(suite, "setup-diagnostics", case, paths, "aot")
    for env_key in ("CATTY_TRACE", "CATTY_METRICS"):
        result = suite.run(f"run-bad-{env_key}", [binary],
                           cwd=binary.parent,
                           env={env_key: str(suite.work / "missing" / "nested" / "out"),
                                "PATH": "", "JAVA_HOME": "", "JDK_HOME": ""})
        require(result.status == 2, f"{env_key} setup status {result.status}. See {result.log}")
        require(result.stdout == b"", f"{env_key} setup ran Java: {result.stdout!r}")
        require(b"cannot open" in result.stderr, f"{env_key} diagnostic: {result.stderr!r}")

    # A Go host panic while Java frames are active reports the Java frame
    # snapshot and the original Go panic stack, and never runs hooks.
    fault = run_allow_timeout(suite, "host-fault",
                              ["go", "run", "./testdata/gate5/setup-failure/fault"],
                              timeout=120, cwd=ROOT)
    require(fault.status != 0 and not fault.timed_out, f"host fault exited cleanly: {fault.log}")
    text = fault.stderr.decode("utf-8", errors="replace")
    for pattern in (r"Catty host fault on Java thread \"main\"",
                    r"at Fault\.main\(",
                    r"panic: injected host fault",
                    r"goroutine "):
        require(re.search(pattern, text), f"host fault stderr missing /{pattern}/: {fault.log}")
    print("setup-failure checks passed", flush=True)


# ------------------------------------------------------------------- main

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--java-home", type=Path, default=os.environ.get("JAVA_HOME"))
    parser.add_argument("--prepare-only", action="store_true", help="run only the OpenJDK 21 fixtures")
    parser.add_argument("--case", choices=sorted(CASES))
    parser.add_argument("--race", action="store_true", help="also run race-instrumented Catty images")
    parser.add_argument("--keep", action="store_true", help="keep artifacts on success")
    args = parser.parse_args()
    if args.prepare_only and args.race:
        parser.error("--race requires Catty execution")

    work = Path(tempfile.mkdtemp(prefix="catty-gate5-"))
    print(f"Gate 5 artifacts in {work}", flush=True)
    try:
        suite = Suite(work, env=java_environment(args.java_home))
        check_toolchains(suite)
        if not args.prepare_only:
            suite.run("build-catty", ["go", "build", "-o", suite.bin, "./cmd/catty"]).succeeded()
        if args.prepare_only:
            if not args.case:
                app, adapter, jar, invalid = prepare_telemetry(suite)
                oracle = suite.run("oracle-telemetry", [
                    "java", "-Xverify:all", "-Dstdout.encoding=UTF-8",
                    "-Dstderr.encoding=UTF-8", "-cp",
                    os.pathsep.join(str(p) for p in (adapter, app)),
                    "Telemetry", str(jar), str(invalid)], timeout=30).succeeded()
                require(oracle.stdout == TELEMETRY_STDOUT,
                        f"telemetry oracle output: {oracle.stdout!r}. See {oracle.log}")
            for label in ((args.case,) if args.case else CASES):
                oracle_only(suite, label)
        else:
            if not args.case:
                profile_checks(suite, args.race)
                setup_failure_checks(suite)
                generated_checks(suite)
                telemetry_checks(suite, args.race)
            for label in ((args.case,) if args.case else CASES):
                exercise(suite, label, args.race)
    except (AssertionError, KeyError, OSError, ValueError) as error:
        print(f"Gate 5 failed. {error}\nArtifacts retained in {work}", file=sys.stderr)
        return 1
    if not args.keep:
        shutil.rmtree(work)
    print("Gate 5 OpenJDK 21 reference prepared." if args.prepare_only else "Gate 5 differential tests passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
