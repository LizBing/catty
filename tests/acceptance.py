#!/usr/bin/env python3
"""Catty Gate 1 acceptance tests. See ACCEPTANCE.md for the wire contracts."""

import argparse
from collections import Counter
from dataclasses import dataclass
import json
import os
from pathlib import Path
import re
import shutil
import struct
import subprocess
import sys
import tempfile
import unittest

from classfile import ClassFile, u2


ROOT = Path(__file__).resolve().parents[1]
DATA = ROOT / "testdata"
PROBE = ("Crossing", "probe", "()I")
APPLY = ("Crossing$Worker", "apply", "(LCrossing$Box;)V")
STEP_APPLY = ("Crossing$Step", *APPLY[1:])
BUMP = ("Crossing$Box", "bump", "()V")
REM = ("Remainder", "rem", "(II)I")
INCREMENT = ("State$Counter", "increment", "()I")
MAIN_DESC = "([Ljava/lang/String;)V"

MODES = ("auto", "aot", "interp")
OPCODES = frozenset("""
    aload aload_0 aload_1 aload_2 aload_3 astore astore_0 astore_1 astore_2 astore_3
    dup pop iconst_m1 iconst_0 iconst_1 iconst_2 iconst_3 iconst_4 iconst_5 bipush
    ldc ldc_w new getfield putfield getstatic anewarray aastore invokestatic invokespecial
    invokevirtual invokeinterface return ireturn athrow iadd if_acmpeq if_acmpne if_icmpeq
    iload iload_0 iload_1 iload_2 iload_3 istore istore_0 istore_1 istore_2 istore_3
    sipush aconst_null putstatic aaload arraylength areturn goto ifeq ifne if_icmpne irem
""".split())


def require(condition, message):
    if not condition:
        raise AssertionError(message)


@dataclass(frozen=True)
class Result:
    pid: int
    status: int
    stdout: bytes
    stderr: bytes
    log: Path

    def succeeded(self):
        require(self.status == 0, f"expected success, got status {self.status}. See {self.log}")
        return self

    def failed(self, *patterns):
        require(self.status > 0, f"expected a controlled failure, got status {self.status}. See {self.log}")
        stderr = self.stderr.decode("utf-8", errors="replace")
        for pattern in patterns:
            require(re.search(pattern, stderr, re.IGNORECASE),
                    f"stderr missing /{pattern}/. See {self.log}")
        return self


@dataclass(frozen=True)
class Build:
    binary: Path
    report: Path
    generated: Path


class Suite:
    def __init__(self, work, *, env=None):
        self.work = work
        self.env = env or {}
        self.bin = work / ("catty.exe" if os.name == "nt" else "catty")

    def run(self, label, args, *, cwd=ROOT, env=None, timeout=120):
        log = self.work / "logs" / f"{label}.json"
        log.parent.mkdir(parents=True, exist_ok=True)
        relevant_env = {"PATH", "CATTY_TRACE", "LC_ALL", "JAVA_HOME", "JDK_HOME", "CLASSPATH"}
        options = os.environ.copy()
        for key in ("CATTY_TRACE", "JAVA_TOOL_OPTIONS", "JDK_JAVA_OPTIONS", "_JAVA_OPTIONS", "CLASSPATH"):
            options.pop(key, None)
        options["LC_ALL"] = "C"
        options.setdefault("GOCACHE", str(self.work / "go-cache"))
        options.update(self.env)
        options.update(env or {})
        command = [str(arg) for arg in args]
        details = {"command": command, "cwd": str(cwd),
                   "environment": {k: options[k] for k in relevant_env if k in options}}
        try:
            with subprocess.Popen(command, cwd=cwd, env=options, stdin=subprocess.DEVNULL,
                                  stdout=subprocess.PIPE, stderr=subprocess.PIPE) as proc:
                try:
                    stdout, stderr = proc.communicate(timeout=timeout)
                except subprocess.TimeoutExpired:
                    proc.kill()
                    stdout, stderr = proc.communicate()
                    raise AssertionError(f"timed out after {timeout}s. See {log}")
                finally:
                    details.update(pid=proc.pid, status=proc.returncode,
                                   stdout=stdout.decode("utf-8", errors="replace"),
                                   stderr=stderr.decode("utf-8", errors="replace"))
        except OSError as exc:
            details["error"] = str(exc)
            raise AssertionError(f"cannot launch {command[0]}. See {log}") from exc
        finally:
            log.write_text(json.dumps(details, indent=2) + "\n", encoding="utf-8")
        return Result(proc.pid, proc.returncode, stdout, stderr, log)

    def classes(self, name):
        return self.work / "classes" / name

    def disassemble(self, name, classes):
        files = sorted(classes.rglob("*.class"))
        require(files, f"no classes in {classes}")
        result = self.run(f"javap-{name}", ["javap", "-p", "-c", "-s", *files]).succeeded()
        return result.stdout.decode("utf-8")

    def prepare(self):
        for tool, flag in (("go", "version"), ("java", "-version"), ("javac", "-version"),
                           ("javap", "-version"), ("jar", "--version")):
            result = self.run(f"version-{tool}", [tool, flag]).succeeded()
            print((result.stdout + result.stderr).decode("utf-8").strip(), flush=True)
        fixtures = {
            "crossing": "crossing", "rem": "remainder", "dispatch": "dispatch",
            "state": "state", "faults": "exceptions", "first": "packaging/first",
            "second": "packaging/second", "unsupported": "packaging/unsupported",
            "dynamic": "packaging/dynamic", "reserved-source": "packaging/reserved",
        }
        disassembly = {}
        for name, source in fixtures.items():
            classes = self.classes(name)
            classes.mkdir(parents=True)
            sources = sorted((DATA / source).rglob("*.java"))
            self.run(f"javac-{name}", ["javac", "--release", "21", "-g", "-d", classes,
                                      *sources]).succeeded()
            disassembly[name] = self.disassemble(name, classes)
            opcodes = set(re.findall(r"^\s+\d+:\s+([a-z]\w*)\b", disassembly[name], re.MULTILINE))
            allowed = OPCODES if name == "rem" else OPCODES - {"irem"}
            if name == "dynamic":
                require("invokedynamic" in opcodes, "concat fixture must use dynamic bootstrap")
                allowed |= {"invokedynamic"}
            require(opcodes <= allowed, f"{name} generated unsupported bytecode: {sorted(opcodes - allowed)}")
        self.jar = self.work / "entry.jar"
        self.run("jar-first", ["jar", "--create", "--file", self.jar, "-C",
                               self.classes("first"), "."]).succeeded()
        self.patch_fixtures(disassembly)
        self.oracles = {}
        self.reference("crossing", "Main", b"2\n")
        for n, expected in enumerate((b"1\n", b"-1\n", None, b"1\n", b"0\n")):
            self.reference("rem", "Remainder", expected, ("arg",) * n)
        self.reference("dispatch", "Dispatch", b"11\n22\n31\n32\n41\n42\n12\n43\n51\n", ("reveal",))
        self.reference("dispatch", "Dispatch", b"22\n11\n31\n32\n41\n42\n12\n43\n51\n", ("reveal", "switch"))
        self.reference("state", "State", b"60\n70\n11\n17\n21\n32\n81\n82\n")
        self.reference("faults", "Faults", b"11\n12\n13\n14\n15\n1\n")
        self.reference("exclusive", "Faults", b"11\n12\n13\n14\n16\n1\n")
        for label, paths, expected in (
            ("directory", [self.classes("first")], b"41\n"),
            ("jar", [self.jar], b"41\n"),
            ("first-wins", [self.jar, self.classes("second")], b"41\n"),
            ("second-wins", [self.classes("second"), self.jar], b"52\n"),
        ):
            self.reference(label, "example.Entry", expected, paths=paths)
        for failure in ("bad-branch", "bad-stackmap"):
            args = ["java", "-Xverify:all", "-cp", self.classes(failure), "Main"]
            self.run(f"oracle-{failure}", args).failed(r"VerifyError|ClassFormatError")

    def patch_fixtures(self, disassembly):
        source = self.classes("crossing") / "Crossing.class"
        parsed = ClassFile(source)
        code = parsed.methods[PROBE[1:]].code
        probe = disassembly["crossing"].split("static int probe()", 1)[1]
        branch = re.search(r"^\s+(\d+):\s+if_acmp(?:eq|ne)\s+", probe, re.MULTILINE)
        require(branch, "javap did not find a reference branch in Crossing.probe")
        calls = re.findall(
            r"^\s+(\d+):\s+invoke(?:interface|virtual)\s+.*"
            r"(?:Crossing\$Step\.apply:|java/lang/reflect/Method\.invoke:)", probe, re.MULTILINE)
        require(len(calls) == len(code.handlers) == 2, "unexpected Crossing.probe calls or handlers")
        self.crossing_sites = [(int(call), u2(parsed.data, handler + 4))
                               for call, handler in zip(calls, code.handlers)]
        pc = int(branch[1])
        require(parsed.data[code.start + pc] in (0xA5, 0xA6), "branch offset disagrees with javap")

        def patch_copy(label, offset, replacement):
            dest = self.classes(label)
            shutil.copytree(source.parent, dest)
            data = bytearray(parsed.data)
            data[offset:offset + len(replacement)] = replacement
            (dest / source.name).write_bytes(data)

        patch_copy("bad-branch", code.start + pc + 1, b"\x00\x01")
        stackmap = code.attributes["StackMapTable"]
        require(u2(parsed.data, stackmap) > 0, "Crossing.probe needs a StackMapTable entry")
        patch_copy("bad-stackmap", stackmap + 2, b"\x80")
        patch_copy("bad-version", 6, b"\x00\x42")  # Major version 66.

        reserved = self.classes("reserved") / "java" / "lang" / "Intruder.class"
        reserved.parent.mkdir(parents=True)
        original = self.classes("reserved-source") / "demo" / "lang" / "Intruder.class"
        data = original.read_bytes()
        require(b"demo/lang/Intruder" in data, "unexpected Intruder fixture layout")
        reserved.write_bytes(data.replace(b"demo/lang/Intruder", b"java/lang/Intruder"))

        exclusive = self.classes("exclusive")
        shutil.copytree(self.classes("faults"), exclusive)
        path = exclusive / "Faults.class"
        parsed = ClassFile(path)
        code = parsed.methods[("boundary", "(LFaults$Bomb;)V")].code
        boundary = disassembly["faults"].split("static void boundary(", 1)[1]
        boundary = boundary.split("public static void main(", 1)[0]
        fire = re.search(r"^\s+(\d+):\s+invokevirtual\s+.*Faults\$Bomb\.fire:",
                         boundary, re.MULTILINE)
        require(fire and len(code.handlers) == 1, "unexpected boundary fixture layout")
        pc = int(fire[1])
        handler = code.handlers[0]
        require(u2(parsed.data, handler) < pc < u2(parsed.data, handler + 2),
                "fire must be strictly inside the original try range")
        struct.pack_into(">H", parsed.data, handler + 2, pc)
        path.write_bytes(parsed.data)
        self.disassemble("exclusive", exclusive)

    def policy(self, label, selections):
        folder = self.work / "policies"
        folder.mkdir(exist_ok=True)
        path = folder / f"{label}.json"
        methods = [dict(owner=owner, name=name, descriptor=descriptor, engine=engine)
                   for (owner, name, descriptor), engine in selections]
        path.write_text(json.dumps({"methods": methods}, indent=2) + "\n", encoding="utf-8")
        return path

    def inspect_generated(self, build, choices):
        root = build.generated.resolve()
        require(root.is_dir(), f"missing generated source directory: {root}")
        for method, check in ((BUMP, "mutation"), (PROBE, "control-flow")):
            generated = choices[method]["generated"]
            source = (root / generated["file"]).resolve()
            require(root in source.parents and source.is_file(), f"missing generated source: {source}")
            token = generated["symbol"]
            self.run(f"source-{method[1]}", ["go", "run", ROOT / "tests" / "inspect_generated.go",
                                             token, check, source]).succeeded()

    def reference(self, label, main, expected, args=(), paths=None):
        classpath = paths or [self.classes(label)]
        suffix = f"-{len(args)}" if args else ""
        result = self.run(f"oracle-{label}{suffix}", ["java", "-Xverify:all", "-cp",
                           os.pathsep.join(map(str, classpath)), main, *args])
        if expected is None:
            result.failed(r"java\.lang\.ArithmeticException", r"Remainder\.rem\(", r"Remainder\.main\(")
            expected = b""
        else:
            result.succeeded()
        require(result.stdout == expected,
                f"OpenJDK {label}{suffix}: expected {expected!r}, got {result.stdout!r}. See {result.log}")
        self.oracles[(label, args)] = result

    def build(self, label, paths, main, mode="auto", *, policy=None, expected_errors=(),
              generated=False, output=None, enable_plugins=False):
        folder = self.work / "builds" / label
        folder.mkdir(parents=True)
        result = Build(output or folder / ("app.exe" if os.name == "nt" else "app"),
                       folder / "report.json", folder / "generated")
        command = [self.bin, "build", "--classpath", os.pathsep.join(map(str, paths)),
                   "--main", main, "--mode", mode, "--report", result.report, "--output", result.binary]
        if enable_plugins:
            command.append("--enable-plugins")
        if policy is not None:
            command += ["--policy", policy]
        if generated:
            command += ["--generated-source", result.generated]
        run = self.run(f"build-{label}", command)
        if expected_errors:
            run.failed(*expected_errors)
            require(not result.report.exists(), f"failed build released a report: {result.report}")
            if output is None:
                require(not result.binary.exists(), f"failed build released an executable: {result.binary}")
        else:
            run.succeeded()
            require(result.binary.is_file(), f"missing {result.binary}. See {run.log}")
            require(result.report.is_file(), f"missing {result.report}. See {run.log}")
        return result

    def traced_run(self, label, binary, args=(), *, isolated=False):
        trace = self.work / "traces" / f"{label}.jsonl"
        trace.parent.mkdir(exist_ok=True)
        env = {"CATTY_TRACE": str(trace)}
        if isolated:
            env.update(PATH="", JAVA_HOME="", JDK_HOME="")
        result = self.run(f"launch-{label}", [binary, *args], cwd=binary.parent, env=env)
        return result, trace

    def launch(self, label, binary, oracle, args=(), *, isolated=False):
        result, trace = self.traced_run(label, binary, args, isolated=isolated)
        if oracle.status == 0:
            result.succeeded()
        else:
            result.failed(r"java[./]lang[./]ArithmeticException", r"Remainder[./]rem\b", r"Remainder[./]main\b")
        require(result.stdout == oracle.stdout,
                f"expected {oracle.stdout!r}, got {result.stdout!r}. See {result.log}")
        return verify_trace(trace, result)


def method_id(record):
    return tuple(record[key] for key in ("owner", "name", "descriptor"))


def read_report(build, main, paths, *, winning_paths=None, default="aot", overrides=None):
    report = json.loads(build.report.read_text(encoding="utf-8"))
    require(isinstance(report["profile"], str) and report["profile"], f"missing profile in {build.report}")
    abi = report["imageABI"]
    require(type(abi) in (int, str) and abi, f"missing image ABI in {build.report}")
    require(report["main"] == main, f"wrong entry point in {build.report}")
    require(isinstance(report["toolchains"]["go"], str) and report["toolchains"]["go"].strip(),
            f"missing Go version in {build.report}")
    require([Path(entry["origin"]) for entry in report["inputs"]] == paths,
            f"wrong classpath order in {build.report}")
    for entry in report["inputs"]:
        require(isinstance(entry["digest"], str) and entry["digest"], f"missing digest: {entry}")

    choices = {}
    for record in report["methods"]:
        key = method_id(record)
        require(key not in choices, f"duplicate method in {build.report}: {key}")
        require(record["engine"] in (*MODES[1:], "intrinsic", "abstract"),
                f"invalid entry kind: {record}")
        require(isinstance(record["reason"], str) and record["reason"].strip(),
                f"missing selection reason: {record}")
        length = record["retainedCodeBytes"]
        require(type(length) is int and (length > 0 if record["engine"] == "interp" else length == 0),
                f"wrong retained length for {key}: {length}")
        choices[key] = record

    # For a JAR, the source of truth for its methods is the directory used to create it.
    selected_paths = winning_paths if winning_paths is not None else paths
    seen = set()
    expect_overrides = overrides or {}
    for path in selected_paths:
        for class_path in sorted(path.rglob("*.class")):
            class_file = ClassFile(class_path)
            if class_file.name in seen:
                continue
            seen.add(class_file.name)
            for (name, descriptor), method in class_file.methods.items():
                key = (class_file.name, name, descriptor)
                record = choices.get(key)
                require(record is not None, f"missing method {key} in {build.report}")
                kind = "abstract" if method.flags & 0x0400 else expect_overrides.get(key, default)
                require(record["engine"] == kind,
                        f"{key} ran as {record['engine']}, expected {kind} in {build.report}")
                if kind == "interp":
                    require(record["retainedCodeBytes"] == method.code.size,
                            f"wrong retained byte count for {key} in {build.report}")
    require(bool(seen), f"no application classes checked in {build.report}")
    require(any(record["engine"] == "intrinsic" for record in choices.values()),
            f"missing bootstrap intrinsic records in {build.report}")
    return report, choices


def verify_trace(path, result):
    require(path.is_file(), f"missing trace {path}. See {result.log}")
    events = [json.loads(line) for line in path.read_text(encoding="utf-8").splitlines()]
    require(events, f"empty trace {path}")
    require(all(type(event["pid"]) is int and event["pid"] == result.pid for event in events),
            f"trace did not come from the launched process ({result.pid}): {path}")
    threads = {event["thread"] for event in events}
    require(len(threads) == 1 and None not in threads, f"expected one logical Java thread: {path}")
    stack = []
    for event in events:
        kind = event["event"]
        if kind == "enter":
            require(type(event["depth"]) is int and event["depth"] == len(stack),
                    f"wrong enter depth: {event} in {path}")
            require(event["engine"] in ("aot", "interp", "intrinsic"),
                    f"invalid executing engine: {event} in {path}")
            stack.append(event)
        elif kind in ("return", "unwind"):
            require(stack, f"{kind} without an active frame: {event} in {path}")
            entry = stack.pop()
            require((method_id(event), event["depth"]) == (method_id(entry), entry["depth"]),
                    f"unbalanced {kind}: {event} vs {entry} in {path}")
            if kind == "unwind":
                require(event["throwable"] is not None, f"unwind lost the throwable: {event}")
        elif kind == "catch":
            require(stack and method_id(event) == method_id(stack[-1]),
                    f"catch outside the active method: {event} in {path}")
            require(type(event["throwPC"]) is int and event["throwPC"] >= 0
                    and type(event["handlerPC"]) is int and event["handlerPC"] >= 0,
                    f"catch missing original bytecode offsets: {event}")
            require(event["throwable"] is not None, f"catch lost the throwable: {event}")
    require(not stack and any(event["event"] == "enter" for event in events),
            f"missing or unclosed frames in {path}")
    return events


def matching(events, event_kind, method):
    return [(index, event) for index, event in enumerate(events)
            if event["event"] == event_kind and method_id(event) == method]


class Acceptance(unittest.TestCase):
    def crossing_trace(self, events, mode):
        base = "interp" if mode == "interp" else "aot"
        apply_engine = "aot" if mode == "aot" else "interp"
        main, = matching(events, "enter", ("Main", "main", MAIN_DESC))
        probe, = matching(events, "enter", PROBE)
        first_apply, second_apply = matching(events, "enter", APPLY)
        first_bump, second_bump = matching(events, "enter", BUMP)
        first_throw, second_throw = matching(events, "unwind", APPLY)
        first_catch, second_catch = matching(events, "catch", PROBE)
        probe_return, = matching(events, "return", PROBE)
        main_return, = matching(events, "return", ("Main", "main", MAIN_DESC))
        self.assertEqual([event["engine"] for _, event in (main, probe, first_bump, second_bump)],
                         [base] * 4)
        self.assertEqual([first_apply[1]["engine"], second_apply[1]["engine"]], [apply_engine] * 2)
        self.assertEqual(main[1]["depth"], 0)
        self.assertEqual(probe[1]["depth"], 1)
        for (_, apply), (_, bump) in ((first_apply, first_bump), (second_apply, second_bump)):
            self.assertEqual(bump["depth"], apply["depth"] + 1)
        self.assertEqual(first_apply[1]["depth"], probe[1]["depth"] + 1)

        first_site, second_site = SUITE.crossing_sites
        self.assertEqual((first_catch[1]["throwPC"], first_catch[1]["handlerPC"]), first_site)
        self.assertEqual((second_catch[1]["throwPC"], second_catch[1]["handlerPC"]), second_site)
        direct, = [i for i, event in enumerate(events)
                   if event["event"] == "dispatch" and event["opcode"] == "invokeinterface"
                   and method_id(event["declaration"]) == STEP_APPLY and method_id(event["target"]) == APPLY]
        reflect, = [i for i, event in enumerate(events)
                    if event["event"] == "reflect" and method_id(event["declaration"]) == STEP_APPLY
                    and method_id(event["target"]) == APPLY]
        wrapper, = [(i, event) for i, event in enumerate(events)
                    if event["event"] == "wrap" and event["wrapper"] == second_catch[1]["throwable"]]
        order = [main[0], probe[0], direct, first_apply[0], first_bump[0], first_throw[0],
                 first_catch[0], reflect, second_apply[0], second_bump[0], second_throw[0],
                 wrapper[0], second_catch[0], probe_return[0], main_return[0]]
        self.assertEqual(order, sorted(set(order)))
        signal = first_throw[1]["throwable"]
        self.assertEqual(first_catch[1]["throwable"], signal)
        self.assertEqual(second_throw[1]["throwable"], signal)
        self.assertEqual(wrapper[1]["cause"], signal)
        self.assertNotEqual(wrapper[1]["wrapper"], signal)

    def test_crossing_and_standalone(self):
        paths = [SUITE.classes("crossing")]
        oracle = SUITE.oracles[("crossing", ())]
        identities = []
        for mode in MODES:
            with self.subTest(mode=mode):
                mixed = mode == "auto"
                build = SUITE.build(f"crossing-{mode}", paths, "Main", mode,
                                    policy=DATA / "crossing" / "mixed.json" if mixed else None,
                                    generated=mixed)
                report, choices = read_report(build, "Main", paths, default="interp" if mode == "interp" else "aot",
                                              overrides={APPLY: "interp"} if mixed else None)
                identities.append((report["profile"], report["imageABI"]))
                trace = SUITE.launch(f"crossing-{mode}", build.binary, oracle)
                self.crossing_trace(trace, mode)
                if mixed:
                    SUITE.inspect_generated(build, choices)
                    standalone = SUITE.work / "standalone"
                    standalone.mkdir()
                    binary = standalone / build.binary.name
                    shutil.copy2(build.binary, binary)
                    self.assertEqual(list(standalone.iterdir()), [binary])
                    offline = paths[0].with_name("crossing-offline")
                    paths[0].rename(offline)
                    try:
                        self.crossing_trace(SUITE.launch("standalone", binary, oracle, isolated=True), mode)
                    finally:
                        offline.rename(paths[0])
        self.assertEqual(len(set(identities)), 1)

    def test_fallback_and_failed_replacement(self):
        paths = [SUITE.classes("rem")]
        build = SUITE.build("remainder-auto", paths, "Remainder")
        _, choices = read_report(build, "Remainder", paths, overrides={REM: "interp"})
        self.assertIn("irem", choices[REM]["reason"].lower())
        for count in range(5):
            with self.subTest(argument_count=count):
                args = ("arg",) * count
                trace = SUITE.launch(f"remainder-{count}", build.binary, SUITE.oracles[("rem", args)], args)
                rem, = matching(trace, "enter", REM)
                self.assertEqual(rem[1]["engine"], "interp")
                main, = matching(trace, "enter", ("Remainder", "main", MAIN_DESC))
                self.assertEqual(main[1]["engine"], "aot")
                if count == 2:
                    rem_unwind, = matching(trace, "unwind", REM)
                    main_unwind, = matching(trace, "unwind", ("Remainder", "main", MAIN_DESC))
                    self.assertEqual(rem_unwind[1]["throwable"], main_unwind[1]["throwable"])

        original = build.binary.read_bytes()
        SUITE.build("remainder-forced-aot", paths, "Remainder", "aot", output=build.binary,
                    expected_errors=(r"irem", r"Remainder", r"rem\s*\(II\)I", r"(?:pc|offset|bci)"))
        self.assertEqual(build.binary.read_bytes(), original)
        SUITE.launch("remainder-preserved", build.binary, SUITE.oracles[("rem", ())])

    def test_dispatch_and_reflection(self):
        paths = [SUITE.classes("dispatch")]
        for mode in ("aot", "interp"):
            with self.subTest(mode=mode):
                build = SUITE.build(f"dispatch-{mode}", paths, "Dispatch", mode)
                read_report(build, "Dispatch", paths, default=mode)
                for args in (("reveal",), ("reveal", "switch")):
                    trace = SUITE.launch(f"dispatch-{mode}-{len(args)}", build.binary,
                                         SUITE.oracles[("dispatch", args)], args)
                    first = ("Dispatch$First", "pick", "()I")
                    second = ("Dispatch$Second", "pick", "()I")
                    picks = [method_id(event) for event in trace if event["event"] == "enter"
                             and method_id(event) in (first, second)]
                    self.assertEqual(picks, [first, second] if len(args) == 1 else [second, first])
                    parent = ("Dispatch$Parent", "reveal", "(LDispatch$Token;)Ljava/lang/Object;")
                    override = ("Dispatch$First", *parent[1:])
                    targets = [method_id(event["target"]) for event in trace if event["event"] == "reflect"
                               and method_id(event["declaration"]) == parent]
                    # The wrong-type adapter error may be logged before or after target selection.
                    self.assertEqual(targets[:2], [parent, override])
                    self.assertEqual(len(matching(trace, "enter", parent)), 1)
                    self.assertEqual(len(matching(trace, "enter", override)), 1)

    def test_class_initialization_and_hidden_fields(self):
        paths = [SUITE.classes("state")]
        policy = SUITE.policy("state", [(INCREMENT, "interp")])
        build = SUITE.build("state-mixed", paths, "State", policy=policy)
        read_report(build, "State", paths, overrides={INCREMENT: "interp"})
        trace = SUITE.launch("state-mixed", build.binary, SUITE.oracles[("state", ())])
        initializer, = matching(trace, "enter", ("State$Counter", "<clinit>", "()V"))
        broken, = matching(trace, "enter", ("State$Broken", "<clinit>", "()V"))
        increments = matching(trace, "enter", INCREMENT)
        self.assertEqual(len(increments), 2)
        self.assertTrue(initializer[0] < increments[0][0] < increments[1][0] < broken[0])
        self.assertEqual([event["engine"] for _, event in increments], ["interp", "interp"])
        main, = matching(trace, "enter", ("State", "main", MAIN_DESC))
        self.assertEqual(main[1]["engine"], "aot")

    def test_packaging_and_classpath(self):
        first, second = SUITE.classes("first"), SUITE.classes("second")
        expected_choices = []
        for label, paths, winning, origin in (
            ("directory", [first], first, str(first / "example" / "Entry.class")),
            ("jar", [SUITE.jar], first, f"{SUITE.jar}!/example/Entry.class"),
            ("first-wins", [SUITE.jar, second], first, f"{SUITE.jar}!/example/Entry.class"),
            ("second-wins", [second, SUITE.jar], second, str(second / "example" / "Entry.class")),
        ):
            with self.subTest(packaging=label):
                build = SUITE.build(label, paths, "example.Entry")
                report, choices = read_report(build, "example.Entry", paths, winning_paths=[winning])
                origins = [entry["origin"] for entry in report["classes"] if entry["name"] == "example/Entry"]
                self.assertEqual(origins, [origin])
                expected_choices.append({key: value["engine"] for key, value in choices.items()})
                SUITE.launch(label, build.binary, SUITE.oracles[(label, ())])
        self.assertEqual(expected_choices[0], expected_choices[1])

    def test_invalid_inputs(self):
        crossing = SUITE.classes("crossing")
        for label, paths, main, evidence in (
            ("bad-version", [SUITE.classes("bad-version")], "Main", (r"Crossing\.class", r"version")),
            ("bad-branch", [SUITE.classes("bad-branch")], "Main",
             (r"Crossing", r"probe\s*\(\)I", r"branch|target", r"pc|offset|bci")),
            ("bad-stackmap", [SUITE.classes("bad-stackmap")], "Main",
             (r"Crossing", r"probe\s*\(\)I", r"stack.?map")),
            ("synchronized", [SUITE.classes("unsupported")], "Locked",
             (r"Locked", r"main", r"synchroniz")),
            ("dynamic", [SUITE.classes("dynamic")], "Concat",
             (r"Concat", r"join", r"invokedynamic|bootstrap")),
            ("reserved", [SUITE.classes("reserved"), crossing], "Main",
             (r"java[./]lang[./]Intruder", r"reserved|bootstrap")),
            ("invalid-main", [crossing], "Crossing", (r"Crossing", r"main")),
        ):
            with self.subTest(input=label):
                SUITE.build(label, paths, main, expected_errors=evidence)

        unknown = ("Crossing$Worker", "missing", APPLY[2])
        unlowerable = ("Remainder", "rem", "(II)I")
        for label, paths, main, selections, evidence in (
            ("unknown", [crossing], "Main", [(unknown, "interp")], (r"unknown|missing|found", r"missing")),
            ("duplicate", [crossing], "Main", [(APPLY, "interp"), (APPLY, "interp")], (r"duplicate", r"apply")),
            ("conflicting", [crossing], "Main", [(APPLY, "interp"), (APPLY, "aot")],
             (r"conflict|duplicate", r"apply")),
            ("abstract", [crossing], "Main", [(STEP_APPLY, "interp")], (r"abstract|concrete|impossible", r"apply")),
            ("unlowerable", [SUITE.classes("rem")], "Remainder", [(unlowerable, "aot")], (r"irem", r"rem")),
        ):
            with self.subTest(policy=label):
                SUITE.build(f"policy-{label}", paths, main, policy=SUITE.policy(label, selections),
                            expected_errors=evidence)

    def test_exceptions_and_exclusive_range(self):
        boundary = ("Faults", "boundary", "(LFaults$Bomb;)V")
        covered = ("Faults", "covered", "(LFaults$Bomb;)V")
        fire = ("Faults$Bomb", "fire", "()V")
        main = ("Faults", "main", MAIN_DESC)
        for label in ("faults", "exclusive"):
            paths = [SUITE.classes(label)]
            for mode in ("aot", "interp"):
                with self.subTest(fixture=label, mode=mode):
                    build = SUITE.build(f"{label}-{mode}", paths, "Faults", mode)
                    read_report(build, "Faults", paths, default=mode)
                    events = SUITE.launch(f"{label}-{mode}", build.binary, SUITE.oracles[(label, ())])
                    first, second = matching(events, "unwind", fire)
                    _, covered_catch = matching(events, "catch", covered)[0]
                    self.assertEqual(len(matching(events, "catch", covered)), 1)
                    self.assertEqual(covered_catch["throwable"], first[1]["throwable"])
                    catches = matching(events, "catch", boundary)
                    unwinds = matching(events, "unwind", boundary)
                    if label == "exclusive":
                        self.assertEqual(catches, [])
                        _, unwind = unwinds[0]
                        self.assertEqual(len(unwinds), 1)
                        _, caller_catch = matching(events, "catch", main)[-1]
                        self.assertEqual(unwind["throwable"], second[1]["throwable"])
                        self.assertEqual(caller_catch["throwable"], unwind["throwable"])
                    else:
                        self.assertEqual(unwinds, [])
                        _, catch = catches[0]
                        self.assertEqual(len(catches), 1)
                        self.assertEqual(catch["throwable"], second[1]["throwable"])

    def test_runtime_and_go_packages(self):
        result = SUITE.run("go-test", ["go", "test", "-count=1", "-json", "./..."], timeout=240).succeeded()
        events = [json.loads(line) for line in result.stdout.splitlines()]
        for name in ("TestCattyAcceptanceImageABI", "TestCattyAcceptanceInterpreterOperandGC"):
            actions = {event["Action"] for event in events if event.get("Test") == name}
            self.assertTrue({"run", "pass"} <= actions, f"{name} must run and pass. See {result.log}")
            self.assertFalse({"skip", "fail"} & actions, f"{name} skipped or failed. See {result.log}")


def main():
    global SUITE
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--prepare-only", action="store_true",
                        help="compile and disassemble fixtures and verify the OpenJDK results")
    args = parser.parse_args()
    work = Path(tempfile.mkdtemp(prefix="catty-acceptance-")).resolve()
    print(f"Acceptance artifacts in {work}", flush=True)
    SUITE = Suite(work)
    try:
        SUITE.prepare()
        if args.prepare_only:
            print(f"Verified {len(SUITE.oracles)} OpenJDK runs and 2 verifier rejections.")
            return 0
        SUITE.run("go-build-catty", ["go", "build", "-o", SUITE.bin, "./cmd/catty"], timeout=240).succeeded()
    except AssertionError as exc:
        print(exc, file=sys.stderr)
        return 1
    result = unittest.TextTestRunner(verbosity=2).run(unittest.defaultTestLoader.loadTestsFromTestCase(Acceptance))
    return 0 if result.wasSuccessful() else 1


if __name__ == "__main__":
    sys.exit(main())
