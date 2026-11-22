#!/usr/bin/env python3
"""Run the Gate 3 plugin-admission suite against OpenJDK 21.

The runner builds every Catty executable before any plugin JAR exists, then
packages the JARs and runs the same executable against late plugin definitions.
"""

import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import struct
import sys
import tempfile
import zipfile
from concurrent.futures import ThreadPoolExecutor

from acceptance import DATA, ROOT, Suite, require, verify_trace


PROFILE = "catty-jvm21-gate3"
IMAGE_ABI = 3
API_STUB = DATA / "gate3" / "stub"
ADAPTER = DATA / "gate3" / "adapter"
SHARED = DATA / "gate3" / "shared"
GATE3 = DATA / "gate3"

# name -> dict(main, modes, jars (labels in arg order), plugins)
CASES = {
    "late-override": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "late-override" / "plugin-a"),
                 ("b", GATE3 / "late-override" / "plugin-b")],
    ),
    "name-lookup": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "name-lookup" / "plugin")],
    ),
    "parent-first": dict(
        main="app.Main",
        plugins=[("shadow", GATE3 / "parent-first" / "plugin-shadow"),
                 ("grand", GATE3 / "parent-first" / "plugin-grand"),
                 ("independent", GATE3 / "parent-first" / "plugin-null")],
    ),
    "loader-identity": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "loader-identity" / "plugin")],
    ),
    "replacement": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "replacement" / "plugin-a"),
                 ("b", GATE3 / "replacement" / "plugin-b")],
    ),
    "crossing-throw": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "crossing-throw" / "plugin")],
    ),
    "initialization": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "initialization" / "plugin")],
    ),
    "deferred-linkage": dict(
        main="app.Main",
        plugins=[("full", GATE3 / "deferred-linkage" / "plugin-full"),
                 ("missing", GATE3 / "deferred-linkage" / "plugin-missing")],
    ),
    "uncaught-plugin": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "uncaught-plugin" / "plugin")],
    ),
    "reflection-identity": dict(
        main="app.Main",
        plugins=[("a", GATE3 / "reflection-identity" / "plugin-a"),
                 ("b", GATE3 / "reflection-identity" / "plugin-b")],
    ),
    "access": dict(
        main="app.Main",
        plugins=[("valid", GATE3 / "access" / "plugin-valid"),
                 ("invalid", GATE3 / "access" / "plugin-invalid")],
    ),
    "rollback": dict(
        main="app.Main",
        plugins=[("good", GATE3 / "rollback" / "plugin-good"),
                 ("bad", GATE3 / "rollback" / "plugin-bad-src"),
                 ("corrected", GATE3 / "rollback" / "plugin-corrected")],
    ),
}

MODES = ("interp", "mixed", "aot")


def java_env(home):
    home = home.expanduser().resolve()
    binaries = home / "bin"
    for tool in ("java", "javac", "javap", "jar"):
        require(shutil.which(tool, path=str(binaries)), f"{home} does not provide {tool}")
    return {"JAVA_HOME": str(home), "PATH": str(binaries) + os.pathsep + os.environ.get("PATH", "")}


def check_toolchains(suite):
    checks = (
        ("go", "version", r"^go version "),
        ("java", "-version", r'^openjdk version "21(?:[."+\-])'),
        ("javac", "-version", r"^javac 21(?:[.\s+\-]|$)"),
        ("javap", "-version", r"^21(?:[.\s+\-]|$)"),
        ("jar", "--version", r"^jar 21(?:[.\s+\-]|$)"),
    )
    for tool, flag, pattern in checks:
        result = suite.run(f"version-{tool}", [tool, flag]).succeeded()
        version = (result.stdout + result.stderr).decode("utf-8").strip()
        require(re.search(pattern, version),
                f"Gate 3 requires OpenJDK 21 tools. {tool} reported {version.partition(chr(10))[0]!r}. See {result.log}")
        print(version, flush=True)


def sha256(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def compile_sources(suite, label, out, sources, classpath=()):
    out.mkdir(parents=True, exist_ok=True)
    cp = os.pathsep.join(str(p) for p in classpath)
    args = ["javac", "--release", "21", "-g", "-d", str(out)]
    if cp:
        args += ["-cp", cp]
    args += [str(p) for p in sources]
    suite.run(label, args).succeeded()
    return out


def make_jar(suite, label, src, out):
    out.parent.mkdir(parents=True, exist_ok=True)
    suite.run(label, ["jar", "--create", "--file", str(out), "-C", str(src), "."]).succeeded()
    return out


def class_digests(classes):
    return {p.name: sha256(p) for p in sorted(classes.rglob("*.class"))}


def prepare_corpus(suite):
    stub = compile_sources(suite, "javac-stub", suite.work / "stub",
                           sorted(API_STUB.rglob("*.java")))
    common = compile_sources(suite, "javac-shared", suite.work / "shared",
                             sorted(SHARED.glob("*.java")), [stub])
    adapter = compile_sources(suite, "javac-adapter", suite.work / "adapter",
                              sorted(ADAPTER.rglob("*.java")))
    return stub, common, adapter


def special_compile(suite, name, app, stub, common):
    """Compile plugin class directories for special layouts."""
    root = GATE3 / name
    if name == "deferred-linkage":
        full = compile_sources(suite, "javac-deferred-full", suite.work / "plugin-full",
                               sorted((root / "plugin-full").rglob("*.java")), [app, stub])
        missing = suite.work / "plugin-missing"
        compile_sources(suite, "javac-deferred-missing-dep", missing,
                        [root / "plugin-missing" / "ext" / "Dep.java"],
                        [app, stub])
        link_out = suite.work / "plugin-missing-link"
        compile_sources(suite, "javac-deferred-missing-link", link_out,
                        [root / "plugin-missing" / "ext" / "Link.java"], [full, app, stub])
        shutil.copy2(link_out / "ext" / "Link.class", missing / "ext" / "Link.class")
        return {"full": full, "missing": missing}
    if name == "rollback":
        good = compile_sources(suite, "javac-rollback-good", suite.work / "plugin-good",
                               sorted((root / "plugin-good").rglob("*.java")), [app, stub])
        bad = compile_sources(suite, "javac-rollback-bad", suite.work / "plugin-bad",
                              sorted((root / "plugin-bad-src").rglob("*.java")), [app, stub])
        corrected = compile_sources(suite, "javac-rollback-corrected", suite.work / "plugin-corrected",
                                    sorted((root / "plugin-corrected").rglob("*.java")), [app, stub])
        return {"good": good, "bad": bad, "corrected": corrected}
    if name == "access":
        valid = compile_sources(suite, "javac-access-valid", suite.work / "plugin-valid",
                                sorted((root / "plugin-valid").rglob("*.java")), [app, stub])
        invalid = compile_sources(suite, "javac-access-invalid", suite.work / "plugin-invalid",
                                  sorted((root / "plugin-invalid").rglob("*.java")), [app, stub])
        return {"valid": valid, "invalid": invalid}
    return None


def package_jars(suite, name, plugins, classdirs):
    root = GATE3 / name
    jars = []
    if name == "rollback":
        badsrc = classdirs["bad"]
        badjar_dir = suite.work / "badjar"
        if badjar_dir.exists():
            shutil.rmtree(badjar_dir)
        shutil.copytree(badsrc, badjar_dir)
        broken = badjar_dir / "ext" / "Broken.class"
        data = bytearray(broken.read_bytes())
        data[0:4] = b"\x00\x00\x00\x00"
        broken.write_bytes(data)
        good = make_jar(suite, "jar-rollback-good", classdirs["good"], suite.work / "rollback-good.jar")
        bad = make_jar(suite, "jar-rollback-bad", badjar_dir, suite.work / "rollback-bad.jar")
        corrected = make_jar(suite, "jar-rollback-corrected", classdirs["corrected"],
                             suite.work / "rollback-corrected.jar")
        return [good, bad, corrected]
    for index, (label, _src) in enumerate(plugins):
        classdir = classdirs[label] if isinstance(classdirs, dict) else classdirs[index]
        jars.append(make_jar(suite, f"jar-{name}-{label}", classdir, suite.work / f"{name}-{label}.jar"))
    return jars


def build_images(suite, name, app, report_dir):
    bins = {}
    policies = suite.work / "policies"
    policies.mkdir(exist_ok=True)
    policy = policies / f"{name}.json"
    policy.write_text(json.dumps({"methods": [{
        "owner": "app/Main", "name": "main",
        "descriptor": "([Ljava/lang/String;)V", "engine": "interp"}]}) + "\n")
    def build(mode):
        out = suite.work / "bin" / f"{name}-{mode}"
        report = report_dir / f"{name}-{mode}.json"
        generated = suite.work / "generated" / f"{name}-{mode}"
        args = [suite.bin, "build", "--enable-plugins", "--classpath", str(app),
                "--main", "app.Main", "--report", str(report), "--output", str(out),
                "--generated-source", str(generated)]
        if mode == "interp":
            args += ["--mode", "interp"]
        elif mode == "aot":
            args += ["--mode", "aot"]
        else:
            args += ["--mode", "auto", "--policy", str(policy)]
        suite.run(f"build-{name}-{mode}", args).succeeded()
        return out, report, generated

    with ThreadPoolExecutor(max_workers=len(MODES)) as pool:
        for mode, value in zip(MODES, pool.map(build, MODES)):
            bins[mode] = value
    return bins


def oracle_run(suite, label, adapter, app, main, args, extra=()):
    paths = [str(adapter), str(app), *[str(p) for p in extra]]
    cp = os.pathsep.join(paths)
    return suite.run(label, ["java", "-Xverify:all", "-cp", cp, main, *map(str, args)])


def trace_events(path):
    events = []
    if not Path(path).exists():
        return events
    for line in Path(path).read_text(encoding="utf-8").splitlines():
        if line.strip():
            events.append(json.loads(line))
    return events


def check_common_trace(events, run_name):
    admissions = [e for e in events if e.get("event") == "plugin_admission" and e.get("result") == "accepted"]
    require(admissions, f"{run_name}: no accepted plugin_admission event")
    defines = [e for e in events if e.get("event") == "class_define"]
    require(defines, f"{run_name}: no class_define event")
    first_admission = events.index(admissions[0])
    first_define = events.index(defines[0])
    require(first_admission < first_define or any(
        e.get("event") == "plugin_admission" and e.get("result") == "accepted" for e in events[:first_define]),
        f"{run_name}: plugin class defined before admission")
    threads = {e.get("thread") for e in events if "thread" in e}
    require(len(threads) <= 1, f"{run_name}: multiple logical threads {threads}")
    # Every plugin method body must enter the interpreter.
    for e in events:
        if e.get("event") == "enter" and e.get("loaderId", 0) > 1:
            require(e.get("engine") == "interp", f"{run_name}: plugin method not interpreted: {e}")
    return admissions


def check_dispatch_trace(name, events, run_name):
    dispatches = [e for e in events if e.get("event") == "dispatch" and e.get("opcode") in ("invokevirtual", "invokeinterface")]
    plugin_targets = [e for e in dispatches if e.get("target", {}).get("loaderId", 0) > 1]
    if name == "late-override":
        relevant = [e for e in dispatches
                    if e.get("caller", {}).get("name") in ("virt", "inter")
                    and e.get("declaration", {}).get("name") == "value"]
        require(relevant, f"{run_name}: no late-override dispatch events")
        for site, expected_plugin in (("virt", 4), ("inter", 2)):
            events_at = [e for e in relevant if e["caller"]["name"] == site]
            owners = [e["target"]["owner"] for e in events_at]
            require(owners.count("app/Base") >= 2, f"{run_name}: {site} image target missing")
            require(owners.count("ext/Worker") >= expected_plugin,
                    f"{run_name}: {site} plugin targets missing: {owners}")
            hits = [e.get("cache", {}).get("hit") for e in events_at]
            require(True in hits and False in hits, f"{run_name}: {site} cache never transitioned")
            target_loaders = [e["target"]["loaderId"] for e in events_at if e["target"]["owner"] == "ext/Worker"]
            if site == "virt":
                require(len(set(target_loaders)) >= 2,
                        f"{run_name}: {site} sibling targets share a loader")
    if name in ("loader-identity", "replacement", "reflection-identity"):
        require(plugin_targets or any(e.get("event") == "reflect" for e in events),
                f"{run_name}: no plugin dispatch or reflect event")


def check_generated_dispatch(suite, name, app, generated_dir):
    text = (Path(generated_dir) / "generated.go").read_text()
    require("t.InvokeVirtual(m," in text or "t.InvokeInterface(m," in text,
            f"{name}: generated Go has no shared dispatch call")
    require("t.InvokeStatic(m," in text or "t.InvokeSpecial(m," in text,
            f"{name}: generated Go has no static/special invocation")
    # Every engine call carries the original bytecode offset and constant-pool
    # entry. Reject an emitter that dropped either argument.
    for match in re.finditer(r"t\.Invoke(?:Virtual|Interface|Static|Special)\(m, ([^)]*)\)", text):
        args = match.group(1).split(",")
        require(len(args) >= 3, f"{name}: dispatch call missing CP/PC: {match.group(0)}")
    if name.startswith("late-override"):
        disassembly = suite.disassemble(f"dispatch-{name}", app)
        section = disassembly.split("static int virt(app.Base, app.Box)", 1)[1]
        virtual = re.search(r"^\s+(\d+): invokevirtual", section[:section.find("static int")], re.MULTILINE)
        section = disassembly.split("static int inter(app.Step, app.Box)", 1)[1]
        interface = re.search(r"^\s+(\d+): invokeinterface", section[:section.find("static int")], re.MULTILINE)
        require(virtual and interface, f"{name}: javap did not find the wrapper sites")
        require(f"t.InvokeVirtual(m, {virtual.group(1)}," in text,
                f"{name}: generated virtual call does not carry the original PC {virtual.group(1)}")
        require(f"t.InvokeInterface(m, {interface.group(1)}," in text,
                f"{name}: generated interface call does not carry the original PC {interface.group(1)}")


def check_report(name, report_path, app_classes, jar_paths):
    report = json.loads(Path(report_path).read_text())
    require(report["profile"] == PROFILE, f"{name}: wrong report profile {report['profile']}")
    require(report["imageABI"] == IMAGE_ABI, f"{name}: wrong report ABI {report['imageABI']}")
    require(report.get("plugins") is True, f"{name}: report does not advertise plugins")
    require("plugin-admission" in report.get("capabilities", []), f"{name}: missing capability")
    limits = report.get("pluginLimits") or {}
    require(limits.get("compressedBytes") == 16 << 20, f"{name}: compressed limit not advertised")
    require(limits.get("loaders") == 64, f"{name}: loader limit not advertised")
    for key in ("readJar", "register"):
        require(any(r["owner"] == "catty/runtime/Plugins" and r["name"] == key
                    for r in report["bootstrap"]["methods"]),
                f"{name}: API method {key} missing from inventory")
    require(any(r["owner"] == "java/lang/Class" and r["name"] == "forName"
                for r in report["bootstrap"]["methods"]),
            f"{name}: Class.forName missing from inventory")
    text = json.dumps(report)
    for path in jar_paths:
        require(str(path) not in text, f"{name}: report references plugin JAR {path}")
    require("ext/Worker" not in text or name in ("name-lookup",), f"{name}: report contains plugin class")
    return report


def check_rollback_trace(events, run_name):
    rejected = [e for e in events if e.get("event") == "plugin_admission" and e.get("result") == "rejected"]
    require(rejected, f"{run_name}: no rejected admission event")
    require(all("loaderId" not in e for e in rejected), f"{run_name}: rejected admission published a loader")
    accepted_after = False
    for index, event in enumerate(events):
        if event.get("event") == "plugin_admission" and event.get("result") == "accepted":
            if any(e.get("event") == "plugin_admission" and e.get("result") == "rejected" for e in events[:index]):
                accepted_after = True
    require(accepted_after, f"{run_name}: no accepted admission after rejection")
    rejected_index = events.index(rejected[0])
    require(not any(e.get("event") == "class_define" and e.get("loaderId") not in (1, 2)
                    for e in events[:rejected_index]),
            f"{run_name}: class defined before the rejected candidate was reported")


def run_case(suite, name, app, stub, common, adapter, bins, jars, expected_stdout):
    (suite.work / "traces").mkdir(parents=True, exist_ok=True)
    before = {}
    for mode in MODES:
        before[mode] = sha256(bins[mode][0])
    jar_before = {str(jar): sha256(jar) for jar in jars}
    results = {}
    for mode in MODES:
        binary, report, generated = bins[mode]
        trace = suite.work / "traces" / f"{name}-{mode}.jsonl"
        result = suite.run(f"run-{name}-{mode}", [binary, *jars],
                           env={"CATTY_TRACE": str(trace)})
        require(result.status == 0, f"{name}-{mode}: status {result.status}. See {result.log}")
        require(result.stdout == expected_stdout,
                f"{name}-{mode}: stdout mismatch\nexpected={expected_stdout!r}\nactual={result.stdout!r}\nSee {result.log}")
        require(result.stderr == b"", f"{name}-{mode}: stderr not empty: {result.stderr!r}")
        verify_trace(trace, result)
        events = trace_events(trace)
        check_common_trace(events, f"{name}-{mode}")
        check_dispatch_trace(name, events, f"{name}-{mode}")
        if name == "rollback":
            check_rollback_trace(events, f"{name}-{mode}")
        require(sha256(bins[mode][0]) == before[mode], f"{name}-{mode}: executable digest changed")
        results[mode] = result
    for jar in jars:
        require(sha256(jar) == jar_before[str(jar)], f"{name}: JAR digest changed during run")
    return results


def policy_checks(suite, binary):
    """Archive and JAR policy variants through the Java registration intrinsic."""
    plugin = suite.work / "policy-plugin"
    compile_sources(suite, "javac-policy", plugin,
                    sorted((GATE3 / "policy" / "plugin").rglob("*.java")), [suite.work / "stub"])
    plain = (plugin / "ext" / "Plain.class").read_bytes()
    base = suite.work / "policy-base.jar"
    with zipfile.ZipFile(base, "w") as z:
        z.writestr("ext/Plain.class", plain)
    accepted = suite.run("policy-accepted", [binary, base]).succeeded()
    require(accepted.stdout == b"accepted\n", f"policy: valid JAR not accepted: {accepted.stdout!r}")

    variants = {}

    def variant(label, data):
        path = suite.work / f"policy-{label}.jar"
        path.write_bytes(data)
        result = suite.run(f"policy-{label}", [binary, path])
        require(result.status == 0, f"policy {label}: status {result.status}. See {result.log}")
        code = result.stdout.decode().strip()
        variants[label] = code
        return code

    valid = base.read_bytes()
    variants["nonzip"] = variant("nonzip", b"not a zip archive")
    variants["truncated-central"] = variant("truncated-central", valid[:-30])
    variants["truncated-entry"] = variant("truncated-entry", valid[:len(valid) // 2])
    badcrc = bytearray(valid)
    # Flip a byte in the last local entry's compressed payload; the central
    # directory starts at the first EOCD signature.
    eocd = badcrc.rfind(b"PK\x05\x06")
    require(eocd > 0, "policy: no EOCD in base JAR")
    badcrc[eocd - 1] ^= 0xFF
    variants["bad-crc"] = variant("bad-crc", bytes(badcrc))

    def write_zip(path, entries, compression=zipfile.ZIP_DEFLATED):
        with zipfile.ZipFile(path, "w", compression=compression) as z:
            for entry_name, data in entries:
                z.writestr(entry_name, data)
        return path.read_bytes()

    variants["duplicate"] = variant("duplicate", write_zip(suite.work / "dup.jar",
                                                           [("ext/Plain.class", plain), ("ext/Plain.class", plain)]))
    variants["empty"] = variant("empty", write_zip(suite.work / "empty.jar",
                                                   [("META-INF/MANIFEST.MF", b"Manifest-Version: 1.0\r\n\r\n")]))
    variants["dotdot"] = variant("dotdot", write_zip(suite.work / "dotdot.jar", [("../ext/Plain.class", plain)]))
    variants["backslash"] = variant("backslash", write_zip(suite.work / "backslash.jar", [("ext\\Plain.class", plain)]))
    variants["absolute"] = variant("absolute", write_zip(suite.work / "absolute.jar", [("/ext/Plain.class", plain)]))
    variants["drive"] = variant("drive", write_zip(suite.work / "drive.jar", [("C:/ext/Plain.class", plain)]))
    variants["name-mismatch"] = variant("name-mismatch", write_zip(suite.work / "mismatch.jar",
                                                                  [("ext/Other.class", plain)]))
    badmagic = bytearray(plain)
    badmagic[0:4] = b"\x00\x00\x00\x00"
    variants["bad-class"] = variant("bad-class", write_zip(suite.work / "badclass.jar",
                                                           [("ext/Plain.class", bytes(badmagic))]))
    major = bytearray(plain)
    major[6:8] = struct.pack(">H", 66)
    variants["major-66"] = variant("major-66", write_zip(suite.work / "major66.jar",
                                                         [("ext/Plain.class", bytes(major))]))
    variants["meta-class"] = variant("meta-class", write_zip(suite.work / "metaclass.jar",
                                                             [("META-INF/Plain.class", plain)]))
    variants["module-info"] = variant("module-info", write_zip(suite.work / "moduleinfo.jar",
                                                               [("module-info.class", plain)]))
    variants["multi-release"] = variant("multi-release", write_zip(
        suite.work / "mr.jar",
        [("META-INF/MANIFEST.MF", b"Manifest-Version: 1.0\r\nMulti-Release: true\r\n\r\n"),
         ("ext/Plain.class", plain)]))
    variants["versioned"] = variant("versioned", write_zip(
        suite.work / "versioned.jar",
        [("META-INF/versions/21/ext/Plain.class", plain), ("ext/Plain.class", plain)]))
    variants["class-path"] = variant("class-path", write_zip(
        suite.work / "classpath.jar",
        [("META-INF/MANIFEST.MF", b"Manifest-Version: 1.0\r\nClass-Path: dep.jar\r\n\r\n"),
         ("ext/Plain.class", plain)]))
    variants["sealed"] = variant("sealed", write_zip(
        suite.work / "sealed.jar",
        [("META-INF/MANIFEST.MF", b"Manifest-Version: 1.0\r\nSealed: true\r\n\r\n"),
         ("ext/Plain.class", plain)]))
    variants["signature"] = variant("signature", write_zip(
        suite.work / "sig.jar",
        [("META-INF/PLAIN.SF", b"Signature-Version: 1.0\r\n"), ("ext/Plain.class", plain)]))
    def set_encrypted(data):
        out = bytearray(data)
        pos = 0
        while True:
            pos = out.find(b"PK\x03\x04", pos)
            if pos < 0:
                break
            out[pos + 6:pos + 8] = b"\x01\x00"
            pos += 4
        pos = 0
        while True:
            pos = out.find(b"PK\x01\x02", pos)
            if pos < 0:
                break
            out[pos + 8:pos + 10] = b"\x01\x00"
            pos += 4
        return bytes(out)

    variants["encrypted"] = variant("encrypted", set_encrypted(write_zip(
        suite.work / "encrypted.jar", [("ext/Plain.class", plain)])))
    variants["unsupported-compression"] = variant("unsupported-compression", write_zip(
        suite.work / "bzip2.jar", [("ext/Plain.class", plain)], zipfile.ZIP_BZIP2))
    class_limit = [(f"pkg/C{i}.class", plain) for i in range(1025)]
    variants["class-count"] = variant("class-count", write_zip(suite.work / "classcount.jar", class_limit))
    entry_limit = [("pad/%d" % i, b"x") for i in range(4096)]
    entry_limit.append(("ext/Plain.class", plain))
    variants["entry-count"] = variant("entry-count", write_zip(suite.work / "entrycount.jar", entry_limit))
    total_pad = [("pad%d.bin" % i, b"\x00" * (4 << 20)) for i in range(17)]
    total_pad.append(("ext/Plain.class", plain))
    variants["expanded-total"] = variant("expanded-total", write_zip(suite.work / "expanded.jar", total_pad))

    # A ZIP entry larger than the expanded entry ceiling.
    huge = b"\x00" * ((4 << 20) + 1)
    variants["entry-limit"] = variant("entry-limit", write_zip(
        suite.work / "huge.jar", [("ext/Plain.class", plain), ("pad.bin", huge)]))
    # Exactly at the entry ceiling is accepted.
    exact = b"\x00" * (4 << 20)
    atlimit = suite.work / "atlimit.jar"
    write_zip(atlimit, [("ext/Plain.class", plain), ("pad.bin", exact)])
    result = suite.run("policy-at-entry-limit", [binary, atlimit]).succeeded()
    require(result.stdout == b"accepted\n", f"policy: exact entry ceiling rejected: {result.stdout!r}")

    # Reserved Catty API namespace.
    reserved = suite.work / "reserved"
    compile_sources(suite, "javac-policy-reserved", reserved,
                    sorted((GATE3 / "policy" / "reserved").rglob("*.java")))
    evil = (reserved / "catty" / "runtime" / "Evil.class").read_bytes()
    variants["reserved"] = variant("reserved", write_zip(suite.work / "reserved.jar",
                                                         [("catty/runtime/Evil.class", evil)]))
    # Missing required superclass from the image.
    dependency = suite.work / "dependency"
    compile_sources(suite, "javac-policy-dependency", dependency,
                    sorted((GATE3 / "policy" / "dependency").rglob("*.java")),
                    [suite.work / "stub", suite.work / "shared"])
    needs = (dependency / "ext" / "NeedsBase.class").read_bytes()
    variants["dependency"] = variant("dependency", write_zip(suite.work / "dependency.jar",
                                                             [("ext/NeedsBase.class", needs)]))

    expected = {
        "nonzip": "ARCHIVE", "truncated-central": "ARCHIVE", "truncated-entry": "ARCHIVE",
        "bad-crc": "ARCHIVE", "duplicate": "ARCHIVE", "empty": "ARCHIVE",
        "dotdot": "ARCHIVE", "backslash": "ARCHIVE", "absolute": "ARCHIVE", "drive": "ARCHIVE",
        "name-mismatch": "CLASS_FORMAT", "bad-class": "CLASS_FORMAT",
        "major-66": "PROFILE", "meta-class": "PACKAGE", "module-info": "PACKAGE",
        "multi-release": "PACKAGE", "versioned": "PACKAGE", "class-path": "PACKAGE",
        "sealed": "PACKAGE", "signature": "PACKAGE", "reserved": "PACKAGE",
        "dependency": "DEPENDENCY", "entry-limit": "LIMIT",
        "encrypted": "ARCHIVE", "unsupported-compression": "ARCHIVE",
        "class-count": "LIMIT", "entry-count": "LIMIT", "expanded-total": "LIMIT",
    }
    for label, want in expected.items():
        require(variants.get(label) == want,
                f"policy {label}: expected {want}, got {variants.get(label)}")
    return variants


def standalone_check(suite, work, prepared, expected_stdout):
    cfg, app, plugin_dirs, compiled, bins = prepared
    empty = work / "standalone"
    empty.mkdir(parents=True, exist_ok=True)
    binary = empty / "app"
    shutil.copy2(bins["mixed"][0], binary)
    jar_a = empty / "a.jar"
    jar_b = empty / "b.jar"
    shutil.copy2(work / "late-override-a.jar", jar_a)
    shutil.copy2(work / "late-override-b.jar", jar_b)
    trace = work / "traces" / "standalone.jsonl"
    result = suite.run("standalone", [binary, jar_a, jar_b], cwd=empty,
                       env={"PATH": "", "JAVA_HOME": "", "CLASSPATH": "", "CATTY_TRACE": str(trace)})
    require(result.status == 0, f"standalone: status {result.status}. See {result.log}")
    require(result.stdout == expected_stdout, f"standalone: stdout mismatch {result.stdout!r}")
    require(result.stderr == b"", f"standalone: stderr {result.stderr!r}")
    verify_trace(trace, result)
    events = trace_events(trace)
    check_common_trace(events, "standalone")
    check_dispatch_trace("late-override", events, "standalone")
    require(sha256(binary) == sha256(bins["mixed"][0]), "standalone: executable changed")


def profile_checks(suite, work, policy_app, stub):
    # Default builds stay Gate 2 and cannot reference the Gate 3 API.
    default_bin = work / "bin" / "default-policy"
    default_report = work / "reports" / "default-policy.json"
    result = suite.run("build-default-plugin-api",
                       [suite.bin, "build", "--classpath", str(policy_app), "--main", "app.Main",
                        "--mode", "auto", "--report", str(default_report),
                        "--output", str(default_bin)])
    require(result.status != 0, "default build accepted a Gate 3 API reference")
    require(not default_bin.exists() and not default_report.exists(),
            "failed default build published output")
    stub_bin = work / "bin" / "stub-on-classpath"
    stub_report = work / "reports" / "stub-on-classpath.json"
    classpath = os.pathsep.join([str(policy_app), str(stub)])
    stub_result = suite.run("build-stub-on-classpath",
                            [suite.bin, "build", "--enable-plugins", "--classpath", classpath,
                             "--main", "app.Main", "--mode", "auto", "--report", str(stub_report),
                             "--output", str(stub_bin)])
    require(stub_result.status != 0, "Gate 3 build accepted the Catty API stub on its classpath")
    require(not stub_bin.exists() and not stub_report.exists(),
            "reserved-package rejection published output")
    # Gate 3 options are rejected by the published Gate 2 runtime pair when a
    # report/image pair is forged. The report must advertise the pair.
    report = json.loads((work / "reports" / "late-override-mixed.json").read_text())
    require(report["profile"] == PROFILE and report["imageABI"] == IMAGE_ABI,
            "Gate 3 report does not advertise its profile/ABI pair")

    # Gate 3 keeps the Gate 2 irem fallback and the strict forced-AOT rejection.
    irem_app = compile_sources(suite, "javac-app-irem", work / "apps" / "irem",
                               sorted((GATE3 / "irem" / "app").rglob("*.java")), [stub])
    irem_out = work / "bin" / "irem-auto"
    irem_report = work / "reports" / "irem-auto.json"
    suite.run("build-irem-auto", [suite.bin, "build", "--enable-plugins", "--classpath", str(irem_app),
                                  "--main", "app.Main", "--mode", "auto", "--report", str(irem_report),
                                  "--output", str(irem_out)]).succeeded()
    report = json.loads(irem_report.read_text())
    rem = [m for m in report["methods"] if m["owner"] == "app/Main" and m["name"] == "rem"]
    require(rem and rem[0]["engine"] == "interp" and "irem" in rem[0]["reason"],
            "Gate 3 auto build did not fall back to the interpreter for irem")
    run = suite.run("run-irem-auto", [irem_out]).succeeded()
    require(run.stdout == b"2\n", f"irem auto output {run.stdout!r}")
    strict = work / "bin" / "irem-aot"
    strict_result = suite.run("build-irem-aot", [suite.bin, "build", "--enable-plugins",
                                                 "--classpath", str(irem_app), "--main", "app.Main",
                                                 "--mode", "aot", "--report", str(work / "reports" / "irem-aot.json"),
                                                 "--output", str(strict)])
    require(strict_result.status != 0, "forced-AOT build accepted irem")
    require(not strict.exists(), "forced-AOT rejection published an executable")
    require(b"irem" in strict_result.stderr, "forced-AOT diagnostic does not name irem")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--java-home", type=Path)
    parser.add_argument("--prepare-only", action="store_true")
    parser.add_argument("--case", choices=sorted(CASES))
    parser.add_argument("--keep", action="store_true")
    args = parser.parse_args()
    env = java_env(args.java_home) if args.java_home else {}
    work = Path(tempfile.mkdtemp(prefix="catty-gate3-"))
    print(f"Gate 3 artifacts in {work}", flush=True)
    suite = Suite(work, env=env)
    suite.bin = work / "bin" / "catty"
    try:
        check_toolchains(suite)
        go = shutil.which("go", path=env.get("PATH"))
        require(go, "go executable not found; set --java-home so its PATH is preserved")
        if not args.prepare_only:
            suite.run("build-catty", [go, "build", "-o", str(suite.bin), "./cmd/catty"]).succeeded()

        stub, common, adapter = prepare_corpus(suite)
        selected = [args.case] if args.case else list(CASES)
        prepared = {}
        for name in selected:
            cfg = CASES[name]
            case_dir = GATE3 / name
            app_sources = sorted((case_dir / "app").rglob("*.java")) + sorted(SHARED.glob("*.java"))
            app = compile_sources(suite, f"javac-app-{name}", work / "apps" / name,
                                  app_sources, [stub])
            # Pin class-file version and at least one important opcode.
            for class_file in sorted(app.rglob("*.class")):
                data = class_file.read_bytes()
                require(data[6:8] == b"\x00\x41", f"{name}: {class_file.name} is not class-file version 65.0")
            disassembly = suite.disassemble(f"app-{name}", app)
            require("invokevirtual" in disassembly or "invokestatic" in disassembly,
                    f"{name}: disassembly missing calls")
            prepared[name] = (cfg, app)
            if name == "access":
                compile_sources(suite, "javac-access-oracle", work / "oracle-access",
                                [GATE3 / "access" / "oracle" / "OracleMain.java"], [app, stub])

        # Oracle outcomes (before any plugin JAR exists).
        expected = {}
        for name, (cfg, app) in prepared.items():
            compiled = special_compile(suite, name, app, stub, common)
            if name == "deferred-linkage":
                plugin_dirs = [compiled["full"], compiled["missing"]]
            elif name == "rollback":
                plugin_dirs = [compiled["good"], compiled["bad"], compiled["corrected"]]
            elif name == "access":
                plugin_dirs = [compiled["valid"], compiled["invalid"]]
            else:
                plugin_dirs = []
                for label, srcdir in cfg["plugins"]:
                    plugin_dirs.append(compile_sources(suite, f"javac-plugin-{name}-{label}",
                                                       work / "plugins" / name / label,
                                                       sorted(srcdir.rglob("*.java")), [app, stub]))
            prepared[name] = (cfg, app, plugin_dirs, compiled, None)

        for name, (cfg, app, plugin_dirs, compiled, _) in list(prepared.items()):
            jars_expected = [work / f"{name}-{label}.jar" for label, _ in cfg["plugins"]]
            if name == "rollback":
                jars_expected = [work / "rollback-good.jar", work / "rollback-bad.jar",
                                 work / "rollback-corrected.jar"]
            jardirs = compiled if name in ("rollback", "access") else plugin_dirs
            bins = None
            if not args.prepare_only:
                for jar in jars_expected:
                    require(not jar.exists(), f"{name}: plugin JAR existed before image build: {jar}")
                bins = build_images(suite, name, app, work / "reports")
                for mode, (binary, report, generated) in bins.items():
                    check_report(name, report, app, jars_expected)
                    require(binary.exists() and binary.stat().st_size > 0, f"{name}-{mode}: no executable")
                    generated_text = (generated / "main.go").read_text()
                    require("catty/plugin" in generated_text and "SetAdmission" in generated_text,
                            f"{name}-{mode}: generated main does not install the adapter")
                    if mode != "interp":
                        check_generated_dispatch(suite, f"{name}-{mode}", app, generated)
            jars = package_jars(suite, name, cfg["plugins"], jardirs)
            prepared[name] = (cfg, app, plugin_dirs, compiled, bins if not args.prepare_only else None)
            for jar, expected_jar in zip(jars, jars_expected):
                require(jar == expected_jar, f"{name}: jar order mismatch {jar} != {expected_jar}")
                require(jar.exists() and jar.stat().st_size > 0, f"{name}: missing JAR {jar}")
            if name == "access":
                oracle = oracle_run(suite, "oracle-access-verify", adapter, app,
                                    "app.OracleMain", [jars[1]], [work / "oracle-access"])
                require(oracle.stdout == b"VerifyError\n",
                        f"access: OpenJDK did not report VerifyError: {oracle.stdout!r}")
            if name == "uncaught-plugin":
                oracle = oracle_run(suite, "oracle-uncaught", adapter, app, cfg["main"], jars)
                require(oracle.status == 1, f"uncaught: oracle status {oracle.status}")
                require(b"java.lang.RuntimeException" in oracle.stderr,
                        f"uncaught: oracle did not report RuntimeException: {oracle.stderr!r}")
            elif name == "access":
                expected[name] = b"7\n5\n1\ntrue\n"
            elif name == "rollback":
                expected[name] = b"10\ntrue\n20\n1\n1\n50\n"
            else:
                oracle = oracle_run(suite, f"oracle-{name}", adapter, app, cfg["main"], jars)
                require(oracle.status == 0, f"{name}: oracle status {oracle.status}. See {oracle.log}")
                require(oracle.stderr == b"", f"{name}: oracle stderr: {oracle.stderr!r}")
                expected[name] = oracle.stdout

        if args.prepare_only:
            print(f"Prepared {len(prepared)} Gate 3 cases and checked OpenJDK outcomes.")
            return

        for name, (cfg, app, plugin_dirs, compiled, bins) in prepared.items():
            jars = []
            if name == "rollback":
                jars = [work / "rollback-good.jar", work / "rollback-bad.jar",
                        work / "rollback-corrected.jar"]
            else:
                for label, _src in cfg["plugins"]:
                    jars.append(work / f"{name}-{label}.jar")
            if name == "uncaught-plugin":
                (work / "traces").mkdir(parents=True, exist_ok=True)
                # The uncaught case intentionally exits 1 on Catty.
                for mode in MODES:
                    binary, report, generated = bins[mode]
                    trace = work / "traces" / f"{name}-{mode}.jsonl"
                    result = suite.run(f"run-{name}-{mode}", [binary, *jars],
                                       env={"CATTY_TRACE": str(trace)})
                    require(result.status == 1, f"{name}-{mode}: expected status 1, got {result.status}")
                    require(b"java.lang.RuntimeException" in result.stderr,
                            f"{name}-{mode}: missing RuntimeException: {result.stderr!r}")
                    require(b"ext.Die.die" in result.stderr or b"ext/Die" in result.stderr,
                            f"{name}-{mode}: missing plugin frame: {result.stderr!r}")
                    stderr = result.stderr.decode("utf-8", errors="replace")
                    plugin_at = stderr.find("ext.Die")
                    image_at = stderr.find("app.Main")
                    require(0 <= plugin_at < image_at,
                            f"{name}-{mode}: plugin/image frame order wrong: {stderr!r}")
                    verify_trace(trace, result)
                    events = trace_events(trace)
                    check_common_trace(events, f"{name}-{mode}")
                print(f"PASS {name}", flush=True)
                continue
            run_case(suite, name, app, stub, common, adapter, bins, jars, expected[name])
            print(f"PASS {name}", flush=True)

        if args.case:
            print(f"Gate 3 case {args.case} passed.")
            return

        policy_app = compile_sources(suite, "javac-app-policy", work / "apps" / "policy",
                                     sorted((GATE3 / "policy" / "app").rglob("*.java")), [stub])
        policy_bins = build_images(suite, "policy", policy_app, work / "reports")
        policy_checks(suite, policy_bins["mixed"][0])
        print("PASS admission-policy", flush=True)
        standalone_check(suite, work, prepared["late-override"], expected["late-override"])
        print("PASS standalone", flush=True)
        profile_checks(suite, work, policy_app, stub)
        print("PASS profile-and-packaging", flush=True)
        print("Gate 3 passed.")
    finally:
        if not args.keep:
            shutil.rmtree(work, ignore_errors=True)
        else:
            print(f"kept {work}")


if __name__ == "__main__":
    main()
