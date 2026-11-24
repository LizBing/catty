#!/usr/bin/env python3
"""Gate 5 application benchmark workload runner.

Compiles the workload under ``testdata/gate5/benchmark`` once with
``javac --release 21 -g``, digests every class file, and uses those identical
class files for the OpenJDK 21 oracle and for the Catty interpreter, mixed and
forced-AOT images.  The runner:

* verifies engine placement and workload checksums against the OpenJDK oracle
  before any timing,
* measures cold launch latency, steady-state batches (``System.nanoTime()``
  around fixed-operation batches), peak RSS per child process, and the
  retained-graph explicit-GC observation,
* records environment/build metadata and raw per-process/per-batch samples,
* writes ``out/gate5/bench/SUMMARY.json`` plus raw evidence for the full run,
* uses reduced repeats across the full matrix with ``--smoke``.

See GATE5-TESTING.md section "Application measurement".  Only the workload
sources under ``testdata/gate5/benchmark`` and this runner belong to the
benchmark; no Go source is modified.

Usage:
  python3 -B tests/gate5_bench.py --java-home /path/to/jdk-21 --smoke
  python3 -B tests/gate5_bench.py --java-home /path/to/jdk-21
"""

import argparse
import csv
import hashlib
import json
import math
import os
from pathlib import Path
import platform
import re
import select
import shutil
import statistics
import subprocess
import sys
import tempfile
import time

from acceptance import DATA, ROOT, Suite, require
from gate2 import java_environment
from gate3 import check_toolchains, compile_sources


BENCH = DATA / "gate5" / "benchmark"
MAIN = "gate5bench.Main"
PROFILE = "catty-jvm21-gate5"
IMAGE_ABI = 5
ALL_CONFIGS = ("java", "interp", "mixed", "aot")
CATTY_CONFIGS = ("interp", "mixed", "aot")
MODE_FOR = {"interp": "interp", "mixed": "auto", "aot": "aot"}
MAIN_KEY = ("gate5bench/Main", "main", "([Ljava/lang/String;)V")
HOT = ("gate5bench/Hot", "mix", "(II)I")
CALLER = ("gate5bench/Workload", "process", "(Lgate5bench/Node;II)I")
CALLBACK = ("gate5bench/Callback", "mix", "(II)I")
MIXED_POLICY = {"methods": [{"owner": HOT[0], "name": HOT[1],
                             "descriptor": HOT[2], "engine": "interp"}]}
JVM_FLAGS = ("-Xms256m", "-Xmx256m", "-XX:+UseSerialGC")
KERNELS = (
    "app",
    "calls-static",
    "calls-virtual",
    "calls-mixed",
    "reflect-lookup",
    "reflect-invoke",
    "memory-alloc",
    "shared-1",
    "shared-4",
    "monitor",
)
KERNEL_OPS = {
    "app": 1024,
    "calls-static": 20000,
    "calls-virtual": 20000,
    "calls-mixed": 4000,
    "reflect-lookup": 300,
    "reflect-invoke": 2000,
    "memory-alloc": 4096,
    "shared-1": 8192,
    "shared-4": 8192,
    "monitor": 4096,
}
OP_DEFINITIONS = {
    "app": "one visit over one node of the deterministic graph: field/array reads, "
           "transient allocation, AOT caller -> interpreted Hot.mix -> AOT Callback",
    "calls-static": "one direct static invocation of Calls.directTarget(int, int)",
    "calls-virtual": "one virtual invocation through Calls.Base, alternating Calls.Alpha "
                     "and Calls.Beta receivers",
    "calls-mixed": "one AOT Calls.crossingKernel -> interpreted Hot.mix -> AOT Callback "
                   "crossing chain",
    "reflect-lookup": "one public Class.getMethod(name, parameter types) lookup by name",
    "reflect-invoke": "one invocation of one cached Method via Method.invoke, including "
                      "boxing; the lookup happened before warmup and outside the timed window",
    "memory-alloc": "one escaping Escape object plus its int[4] payload retained in a ring",
    "shared-1": "one read-modify-write of a shared int[] cell plus a field update, one thread",
    "shared-4": "the same shared-array/field operation, ops summed over four fixed threads",
    "monitor": "one synchronized Counter.bump() increment, ops summed over four fixed threads",
}
WORKLOAD_CONSTANTS = {
    "seed": 12345,
    "appNodes": 1024,
    "appRoundsPerBatch": 1,
    "callsStaticOps": 20000,
    "callsVirtualOps": 20000,
    "callsCrossingOps": 4000,
    "reflectLookupOps": 300,
    "reflectInvokeOps": 2000,
    "allocOpsPerBatch": 4096,
    "allocRing": 1024,
    "allocRssOps": 32768,
    "sharedOpsPerBatch": 8192,
    "sharedThreads": 4,
    "sharedSlots": 512,
    "monitorOpsPerBatch": 4096,
    "monitorThreads": 4,
    "retainedNodes": 12000,
}
SELECTED_METHODS = {
    "main": MAIN_KEY,
    "appProcess": CALLER,
    "appCallback": CALLBACK,
    "hotMethod": HOT,
    "callsStaticTarget": ("gate5bench/Calls", "directTarget", "(II)I"),
    "callsVirtualAlpha": ("gate5bench/Calls$Alpha", "pick", "(II)I"),
    "callsVirtualBeta": ("gate5bench/Calls$Beta", "pick", "(II)I"),
    "callsCrossingKernel": ("gate5bench/Calls", "crossingKernel", "(II)I"),
    "reflectTarget": ("gate5bench/ReflectTarget", "combine", "(II)I"),
    "monitorBump": ("gate5bench/Memory$Counter", "bump", "()V"),
    "memorySharedKernel": ("gate5bench/Memory", "sharedOps",
                           "([ILgate5bench/Cell;IIII)I"),
}
FULL_BATCHES = 30
FULL_WARMUP = 10
SMOKE_BATCHES = 5
SMOKE_WARMUP = 2
VERIFY = "verify"
COLD = "cold"
STEADY = "steady"
STEADY_SMOKE = "steady-smoke"
RSS = "rss"
RETAINED = "retained"
GC_HEAP_PATTERN = re.compile(r"\)\s+(\d+)M->(\d+)M\((\d+)M\)")
FORBIDDEN_BYTECODE = (
    "invokedynamic",
    "java/lang/StringBuilder",
    "java/lang/Math",
    "java/lang/System.arraycopy",
    "currentTimeMillis",
    "Integer.parseInt",
    "Object.hashCode",
    "LambdaMetafactory",
)
QUANTILE_METHOD = ("nearest-rank over the ascending pooled sample "
                   "(index = ceil(p * n) - 1, clamped to [0, n-1])")
TIMED_ENV = {"CATTY_TRACE": "", "CATTY_METRICS": "", "GORACE": ""}


def sha256_bytes(data):
    return hashlib.sha256(data).hexdigest()


def sha256_file(path):
    return sha256_bytes(Path(path).read_bytes())


def write_json(path, value):
    path = Path(path)
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def write_csv(path, header, rows):
    path = Path(path)
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", newline="", encoding="utf-8") as handle:
        writer = csv.writer(handle)
        writer.writerow(header)
        writer.writerows(rows)


def quantile(values, probability):
    """Nearest-rank quantile over an ascending copy of ``values``."""
    if not values:
        raise ValueError("quantile of empty sample")
    ordered = sorted(values)
    rank = int(math.ceil(probability * len(ordered)))
    if rank < 1:
        rank = 1
    if rank > len(ordered):
        rank = len(ordered)
    return ordered[rank - 1]


def describe(values):
    if not values:
        return None
    ordered = sorted(values)
    mean = sum(ordered) / len(ordered)
    deviation = statistics.stdev(ordered) if len(ordered) > 1 else 0.0
    return {
        "count": len(ordered),
        "min": min(ordered),
        "median": statistics.median(ordered),
        "mean": mean,
        "max": max(ordered),
        "stdev": deviation,
        "coefficientOfVariation": (deviation / mean) if mean else 0.0,
        "p50": quantile(ordered, 0.50),
        "p99": quantile(ordered, 0.99),
    }


def relative_files(root):
    root = Path(root)
    return sorted((path for path in root.rglob("*") if path.is_file()),
                  key=lambda path: path.as_posix())


def dir_digest_go(root):
    """Reproduce the Catty builder's class-directory digest."""
    root = Path(root)
    files = sorted(path for path in root.rglob("*") if path.is_file())
    digest = hashlib.sha256()
    for path in files:
        rel = path.relative_to(root).as_posix()
        digest.update(rel.encode("utf-8"))
        digest.update(b"\x00")
        digest.update(path.read_bytes())
    return digest.hexdigest()


def decode_go_string_literal(text, marker):
    """Decode one Go interpreted string literal assigned after ``marker``."""
    index = text.find(marker)
    require(index >= 0, f"missing {marker!r} in generated main.go")
    pos = index + len(marker)
    require(pos < len(text) and text[pos] == '"', "generated image constant is not a Go string")
    pos += 1
    simple = {
        '"': '"', "\\": "\\", "a": "\a", "b": "\b", "f": "\f",
        "n": "\n", "r": "\r", "t": "\t", "v": "\v", "'": "'",
    }
    out = bytearray()
    while pos < len(text):
        char = text[pos]
        pos += 1
        if char == '"':
            return bytes(out)
        if char != "\\":
            out.extend(char.encode("utf-8"))
            continue
        require(pos < len(text), "truncated Go string escape")
        escape = text[pos]
        pos += 1
        if escape in simple:
            out.extend(simple[escape].encode("utf-8"))
        elif escape == "x":
            out.append(int(text[pos:pos + 2], 16))
            pos += 2
        elif escape == "u":
            out.extend(chr(int(text[pos:pos + 4], 16)).encode("utf-8"))
            pos += 4
        elif escape == "U":
            out.extend(chr(int(text[pos:pos + 8], 16)).encode("utf-8"))
            pos += 8
        elif escape.isdigit():
            digits = escape
            while len(digits) < 3 and pos < len(text) and text[pos].isdigit():
                digits += text[pos]
                pos += 1
            out.append(int(digits, 8))
        else:
            raise ValueError(f"unsupported Go string escape \\{escape}")
    raise ValueError("unterminated Go string literal")


def parse_header(lines, command, context):
    require(len(lines) >= 4 and lines[0] == "CATTY-GATE5-BENCH" and lines[1] == "1",
            f"{context}: missing benchmark protocol header: {lines[:3]!r}")
    require(lines[2] == "COMMAND", f"{context}: missing COMMAND record")
    require(lines[3] == command, f"{context}: expected command {command!r}, got {lines[3]!r}")
    return 4


def parse_verify(text, context="verify"):
    lines = text.splitlines()
    index = parse_header(lines, VERIFY, context)
    checks = {}
    order = []
    while index < len(lines) and lines[index] == "CHECK":
        require(index + 2 < len(lines), f"{context}: truncated CHECK record")
        name, value = lines[index + 1], lines[index + 2]
        require(name not in checks, f"{context}: duplicate CHECK {name}")
        checks[name] = int(value)
        order.append(name)
        index += 3
    work = None
    if index < len(lines) and lines[index] == "WORK":
        require(index + 1 < len(lines), f"{context}: truncated WORK record")
        work = int(lines[index + 1])
        index += 2
    require(lines[index:index + 2] == ["END", "0"], f"{context}: missing END record")
    require(index + 2 == len(lines), f"{context}: trailing output after END")
    return {"checks": checks, "order": order, "work": work}


def parse_cold(text, context="cold"):
    lines = text.splitlines()
    index = parse_header(lines, COLD, context)
    require(index + 6 == len(lines), f"{context}: unexpected record count")
    require(lines[index] == "RESULT", f"{context}: missing RESULT record")
    checksum = int(lines[index + 1])
    ops = int(lines[index + 2])
    require(lines[index + 3] == "app", f"{context}: unexpected cold workload name")
    require(lines[index + 4:index + 6] == ["END", "0"], f"{context}: missing END record")
    return {"checksum": checksum, "ops": ops}


def parse_steady(text, context="steady"):
    lines = text.splitlines()
    # The workload's steady command always reports the internal command name
    # "steady"; the batch count in the KERNEL records distinguishes the smoke
    # and full variants.
    index = parse_header(lines, "steady", context)
    require(index < len(lines) and lines[index] == "10", f"{context}: expected 10 kernel records")
    index += 1
    kernels = []
    for _ in range(10):
        require(index + 4 < len(lines) and lines[index] == "KERNEL",
                f"{context}: missing KERNEL record")
        name = lines[index + 1]
        ops = int(lines[index + 2])
        batches = int(lines[index + 3])
        warmup = int(lines[index + 4])
        index += 5
        samples = []
        while index < len(lines) and lines[index] == "NS":
            require(index + 1 < len(lines), f"{context}: truncated NS record")
            samples.append(int(lines[index + 1]))
            index += 2
        require(index + 2 < len(lines) and lines[index] == "CHECK",
                f"{context}: missing CHECK for {name}")
        require(lines[index + 1] == name, f"{context}: CHECK name mismatch for {name}")
        kernels.append({
            "name": name,
            "opsPerBatch": ops,
            "batches": batches,
            "warmup": warmup,
            "samples": samples,
            "checksum": int(lines[index + 2]),
        })
        index += 3
    require(lines[index:index + 2] == ["END", "0"], f"{context}: missing END record")
    require(index + 2 == len(lines), f"{context}: trailing output after END")
    return kernels


def parse_rss(text, context="rss"):
    lines = text.splitlines()
    index = parse_header(lines, RSS, context)
    checks = {}
    while index < len(lines) and lines[index] == "CHECK":
        require(index + 2 < len(lines), f"{context}: truncated CHECK record")
        checks[lines[index + 1]] = int(lines[index + 2])
        index += 3
    require(index + 1 < len(lines) and lines[index] == "WORK", f"{context}: missing WORK record")
    work = int(lines[index + 1])
    index += 2
    while index < len(lines) and lines[index] == "CHECK":
        require(index + 2 < len(lines), f"{context}: truncated CHECK record")
        checks[lines[index + 1]] = int(lines[index + 2])
        index += 3
    require(lines[index:index + 2] == ["END", "0"], f"{context}: missing END record")
    require(index + 2 == len(lines), f"{context}: trailing output after END")
    return {"checks": checks, "work": work}


def parse_retained(text, context="retained"):
    lines = text.splitlines()
    index = parse_header(lines, RETAINED, context)
    values = {}
    for tag in ("RETAINED", "CHECK-BEFORE", "GC-CALLS", "CHECK-AFTER"):
        require(index + 1 < len(lines) and lines[index] == tag,
                f"{context}: expected {tag} record, got {lines[index]!r}")
        values[tag] = int(lines[index + 1])
        index += 2
    require(lines[index:index + 2] == ["END", "0"], f"{context}: missing END record")
    require(index + 2 == len(lines), f"{context}: trailing output after END")
    require(values["GC-CALLS"] == 1, f"{context}: expected one explicit GC call")
    require(values["CHECK-BEFORE"] == values["CHECK-AFTER"],
            f"{context}: retained graph checksum changed across collection")
    return values


def steady_checksums(kernels):
    return [(kernel["name"], kernel["checksum"]) for kernel in kernels]


def parse_gc_log(path):
    """Return the explicit full-collection heap observation from -Xlog:gc output."""
    if not Path(path).is_file():
        return {"status": "failed", "reason": f"missing GC log {path}"}
    text = Path(path).read_text(encoding="utf-8", errors="replace")
    explicit = []
    heap_lines = []
    for line in text.splitlines():
        if "Pause Full" in line and "System.gc()" in line:
            explicit.append(line)
            match = GC_HEAP_PATTERN.search(line)
            if match:
                heap_lines.append((line, match))
    if not explicit:
        return {
            "status": "failed",
            "reason": "no explicit full collection (Pause Full (System.gc())) in -Xlog:gc output",
            "log": text,
        }
    if not heap_lines:
        return {
            "status": "failed",
            "reason": "explicit collection line has no heap-after field",
            "lines": explicit,
        }
    line, match = heap_lines[-1]
    after_mib = int(match.group(2))
    return {
        "status": "ok",
        "source": "OpenJDK -Xlog:gc:file, Pause Full (System.gc()) heap-after field",
        "units": "MiB in the log; heapAfterBytes converted *1024*1024",
        "heapAfterBytes": after_mib * 1024 * 1024,
        "heapBeforeMiB": int(match.group(1)),
        "heapAfterMiB": after_mib,
        "heapTotalMiB": int(match.group(3)),
        "collectionConfirmed": True,
        "rawLine": line,
    }


def read_metrics(path):
    return json.loads(Path(path).read_text(encoding="utf-8"))


def signal_exit_code():
    return -9


def wait4_measure(command, cwd, env, timeout, out_path, err_path):
    """Run one child, returning per-child peak RSS from wait4 rusage."""
    with open(out_path, "wb") as stdout_file, open(err_path, "wb") as stderr_file:
        started = time.monotonic()
        process = subprocess.Popen(
            [str(part) for part in command], cwd=str(cwd), env=env,
            stdin=subprocess.DEVNULL, stdout=stdout_file, stderr=stderr_file)
        timed_out = False
        status = None
        usage = None
        while True:
            try:
                finished, status, usage = os.wait4(process.pid, os.WNOHANG)
            except InterruptedError:
                continue
            if finished == process.pid:
                break
            if time.monotonic() - started > timeout:
                timed_out = True
                process.kill()
                try:
                    finished, status, usage = os.wait4(process.pid, 0)
                except ChildProcessError:
                    status = None
                    usage = None
                break
            time.sleep(0.005)
        process.returncode = (os.waitstatus_to_exitcode(status)
                              if status is not None else signal_exit_code())
        wall = time.monotonic() - started
    if sys.platform.startswith("darwin"):
        units = "bytes"
    elif sys.platform.startswith("linux"):
        units = "KiB"
    else:
        units = "unknown"
    peak = usage.ru_maxrss if usage is not None else None
    return {
        "pid": process.pid,
        "returncode": process.returncode,
        "timedOut": timed_out,
        "wallSeconds": wall,
        "peakRssRaw": peak,
        "peakRssUnits": units,
        "rssMethod": "os.wait4 ru_maxrss (per-child high-water mark)",
        "rssSource": f"POSIX wait4 rusage (Darwin units bytes, Linux units KiB); platform={sys.platform}",
        "stdout": Path(out_path).read_bytes(),
        "stderr": Path(err_path).read_bytes(),
        "stderrSharedWithRssTool": False,
    }


def time_wrapper_measure(command, cwd, env, timeout, out_path, err_path):
    """Fallback peak-RSS method using /usr/bin/time -l (Darwin) or -v (Linux)."""
    candidate = Path("/usr/bin/time")
    if not candidate.is_file():
        return None
    flag = "-l" if sys.platform.startswith("darwin") else "-v"
    with open(out_path, "wb") as stdout_file, open(err_path, "wb") as stderr_file:
        started = time.monotonic()
        process = subprocess.run(
            [str(candidate), flag] + [str(part) for part in command],
            cwd=str(cwd), env=env, stdin=subprocess.DEVNULL,
            stdout=stdout_file, stderr=stderr_file, timeout=timeout)
        wall = time.monotonic() - started
    text = Path(err_path).read_text(encoding="utf-8", errors="replace")
    if sys.platform.startswith("darwin"):
        match = re.search(r"(\d+)\s+maximum resident set size", text)
        units = "bytes"
    else:
        match = re.search(r"Maximum resident set size \(kbytes\):\s*(\d+)", text)
        units = "KiB"
    if not match:
        return None
    return {
        "pid": process.pid,
        "returncode": process.returncode,
        "timedOut": False,
        "wallSeconds": wall,
        "peakRssRaw": int(match.group(1)),
        "peakRssUnits": units,
        "rssMethod": f"/usr/bin/time {flag}",
        "rssSource": f"/usr/bin/time {flag} maximum resident set size; platform={sys.platform}",
        "stdout": Path(out_path).read_bytes(),
        "stderr": Path(err_path).read_bytes(),
        "stderrSharedWithRssTool": True,
    }


def direct_measure(command, cwd, env, timeout, out_path, err_path):
    """Run one child without a peak-RSS reading (last-resort fallback)."""
    with open(out_path, "wb") as stdout_file, open(err_path, "wb") as stderr_file:
        started = time.monotonic()
        process = subprocess.run(
            [str(part) for part in command], cwd=str(cwd), env=env,
            stdin=subprocess.DEVNULL, stdout=stdout_file, stderr=stderr_file,
            timeout=timeout)
        wall = time.monotonic() - started
    return {
        "pid": process.pid,
        "returncode": process.returncode,
        "timedOut": False,
        "wallSeconds": wall,
        "peakRssRaw": None,
        "peakRssUnits": "unavailable",
        "rssMethod": "unavailable",
        "rssSource": ("no per-child peak-RSS method: os.wait4 and a recognized "
                      "/usr/bin/time -l/-v output are unavailable on this platform"),
        "stdout": Path(out_path).read_bytes(),
        "stderr": Path(err_path).read_bytes(),
        "stderrSharedWithRssTool": False,
    }


def rotated(configs, index):
    count = len(configs)
    start = index % count
    return [configs[(start + offset) % count] for offset in range(count)]


def rss_to_bytes(raw, units):
    if raw is None:
        return None
    if units == "bytes":
        return raw
    if units == "KiB":
        return raw * 1024
    return None


class BenchRun:
    def __init__(self, suite, args, work, configs):
        self.suite = suite
        self.args = args
        self.work = work
        self.configs = list(configs)
        self.raw = work / "raw"
        self.raw.mkdir(parents=True, exist_ok=True)
        self.classes = suite.classes("benchmark")
        self.batch_command = STEADY_SMOKE if args.smoke else STEADY
        self.batches = SMOKE_BATCHES if args.smoke else FULL_BATCHES
        self.warmup = SMOKE_WARMUP if args.smoke else FULL_WARMUP
        self.repeat = args.repeat if args.repeat is not None else (2 if args.smoke else 7)
        self.build_repeats = (args.build_repeats if args.build_repeats is not None
                              else (1 if args.smoke else 3))
        self.order = []
        self.failures = []
        self.builds = {}
        self.reports = {}
        self.oracle = {}
        self.correctness = {}
        self.cold = {}
        self.steady = {}
        self.memory = {}
        self.retained = {}
        self.environment = {}
        self.workload = {}
        self.steady_csv_rows = []
        self.started_at = None
        self.finished_at = None
        self.duration = 0.0

    # ------------------------------------------------------------------
    # process helpers
    # ------------------------------------------------------------------
    def run_env(self, extra=None):
        options = os.environ.copy()
        for key in ("CATTY_TRACE", "CATTY_METRICS", "JAVA_TOOL_OPTIONS",
                    "JDK_JAVA_OPTIONS", "_JAVA_OPTIONS", "CLASSPATH", "GORACE"):
            options.pop(key, None)
        options["LC_ALL"] = "C"
        options.setdefault("GOCACHE", str(self.work / "go-cache"))
        options.update(self.suite.env)
        options.update(extra or {})
        return options

    def record_order(self, phase, config, sample, command, duration):
        self.order.append({
            "phase": phase,
            "config": config,
            "sample": sample,
            "command": command,
            "startedAt": time.strftime("%Y-%m-%dT%H:%M:%S%z"),
            "durationSeconds": duration,
        })

    def run_logged(self, phase, config, sample, label, command, *,
                   timeout=120, env=None, cwd=ROOT):
        started = time.monotonic()
        result = self.suite.run(label, command, timeout=timeout, env=env, cwd=cwd)
        duration = time.monotonic() - started
        self.record_order(phase, config, sample, [str(part) for part in command], duration)
        return result, duration

    def catty_command(self, config, command):
        return [str(self.builds[config]["binary"]), command]

    def java_command(self, command, extra=()):
        return ["java", *JVM_FLAGS, *extra, "-cp", str(self.classes), MAIN, command]

    def command_for(self, config, command, extra_java=()):
        if config == "java":
            return self.java_command(command, extra_java)
        return self.catty_command(config, command)

    def add_failure(self, phase, config, detail):
        self.failures.append({"phase": phase, "config": config, "detail": detail})
        print(f"FAIL {phase} [{config}]: {detail}", file=sys.stderr, flush=True)

    # ------------------------------------------------------------------
    # environment
    # ------------------------------------------------------------------
    def environment_phase(self):
        def git(*arguments):
            label = "git-" + "-".join(argument.replace("--", "").replace("/", "-")
                                      for argument in arguments)
            result = self.suite.run(label, ["git", "-C", str(ROOT), *arguments], timeout=60)
            return None if result.status != 0 else result.stdout.decode("utf-8", "replace").strip()

        head = git("rev-parse", "HEAD")
        branch = git("rev-parse", "--abbrev-ref", "HEAD")
        status = git("status", "--porcelain") or ""
        go_version = self.suite.run("go-version", ["go", "version"]).stdout.decode(
            "utf-8", "replace").strip()
        go_env = {}
        for key in ("GOOS", "GOARCH", "GOVERSION", "GOCACHE", "GOFLAGS", "GOEXPERIMENT",
                    "CGO_ENABLED", "GODEBUG", "GOMAXPROCS"):
            result = self.suite.run(f"go-env-{key}", ["go", "env", key], timeout=60)
            if result.status == 0:
                go_env[key] = result.stdout.decode("utf-8", "replace").strip()
        java_version_result = self.suite.run("java-version-env", ["java", "-version"])
        javac_version_result = self.suite.run("javac-version-env", ["javac", "-version"])
        cpu = {}
        if sys.platform.startswith("darwin"):
            for key in ("machdep.cpu.brand_string", "hw.model", "hw.ncpu", "hw.memsize"):
                result = self.suite.run("sysctl-" + key.replace(".", "-"), ["sysctl", "-n", key])
                if result.status == 0:
                    cpu[key] = result.stdout.decode("utf-8", "replace").strip()
        elif sys.platform.startswith("linux"):
            cpuinfo = Path("/proc/cpuinfo")
            if cpuinfo.is_file():
                for line in cpuinfo.read_text(encoding="utf-8", errors="replace").splitlines():
                    if line.startswith("model name"):
                        cpu["model name"] = line.split(":", 1)[1].strip()
                        break
            cpu["ncpu"] = os.cpu_count()
        uname = platform.uname()
        go_cache = self.run_env().get("GOCACHE", "")
        self.environment = {
            "repository": {
                "root": str(ROOT),
                "revision": head,
                "branch": branch,
                "dirty": bool(status),
                "statusPorcelain": status.splitlines(),
                "statusSha256": sha256_bytes(status.encode("utf-8")),
            },
            "toolchains": {
                "go": go_version,
                "goEnv": go_env,
                "java": (java_version_result.stdout + java_version_result.stderr).decode(
                    "utf-8", "replace").strip(),
                "javac": (javac_version_result.stdout + javac_version_result.stderr).decode(
                    "utf-8", "replace").strip(),
            },
            "platform": {
                "system": uname.system,
                "release": uname.release,
                "version": uname.version,
                "machine": uname.machine,
                "platform": platform.platform(),
                "cpu": cpu,
                "logicalCpuCount": os.cpu_count(),
                "python": platform.python_version(),
                "executable": sys.executable,
            },
            "processEnvironment": {
                "GOMAXPROCS": os.environ.get("GOMAXPROCS", ""),
                "GOMAXPROCSMeaning": "empty means the Go runtime default (logical CPU count)",
                "GOGC": os.environ.get("GOGC", ""),
                "GOGCMeaning": "empty means the Go default (100)",
                "GOMEMLIMIT": os.environ.get("GOMEMLIMIT", ""),
                "GOMEMLIMITMeaning": "empty means unlimited",
                "GOFLAGS": os.environ.get("GOFLAGS", ""),
                "CGO_ENABLED": os.environ.get("CGO_ENABLED", ""),
                "GOCACHE": go_cache,
                "GOCACHEPolicy": ("suite work-directory cache; the first builder run compiles the "
                                  "Go runtime and dependencies, later runs reuse the cache"),
                "LC_ALL": "C",
                "timedRunExclusions": ["CATTY_TRACE", "CATTY_METRICS", "GORACE", "GC logging"],
            },
            "jvm": {
                "flags": list(JVM_FLAGS),
                "retainedGcLogging": "-Xlog:gc:file=<path> (retained-graph runs only)",
                "verifyClasses": "java -Xverify:all (correctness oracle only, never a timed run)",
            },
            "benchmark": {
                "mode": "smoke" if self.args.smoke else "full",
                "configs": self.configs,
                "processesPerConfig": self.repeat,
                "batchesPerProcess": self.batches,
                "warmupBatches": self.warmup,
                "buildRepeatsPerMode": self.build_repeats,
                "quantileMethod": QUANTILE_METHOD,
                "normalizedMetric": ("batch elapsed System.nanoTime() divided by the batch's "
                                     "documented operation count; not individual operation latency"),
                "executionOrderPolicy": "configurations interleaved by rotating order per sample",
            },
        }

    # ------------------------------------------------------------------
    # build phase
    # ------------------------------------------------------------------
    def build_phase(self):
        sources = sorted(BENCH.rglob("*.java"), key=lambda path: path.as_posix())
        require(sources, f"no benchmark sources under {BENCH}")
        started = time.monotonic()
        compile_sources(self.suite, "javac-benchmark", self.classes, sources)
        javac_seconds = time.monotonic() - started
        self.record_order("build", "javac", None,
                          ["javac", "--release", "21", "-g"], javac_seconds)

        class_files = relative_files(self.classes)
        require(class_files, f"javac produced no classes in {self.classes}")
        digests = []
        for path in class_files:
            data = path.read_bytes()
            require(data[6:8] == b"\x00\x41",
                    f"{path.name}: class file is not version 65.0 (Java 21)")
            digests.append({
                "path": path.relative_to(self.classes).as_posix(),
                "bytes": len(data),
                "sha256": sha256_bytes(data),
            })
        digest_manifest = sha256_bytes(json.dumps(digests, sort_keys=True).encode("utf-8"))
        source_digests = [{"path": path.relative_to(ROOT).as_posix(),
                           "bytes": path.stat().st_size,
                           "sha256": sha256_file(path)} for path in sources]

        disassembly = self.suite.disassemble("benchmark", self.classes)
        for token in FORBIDDEN_BYTECODE:
            require(token not in disassembly,
                    f"benchmark bytecode contains forbidden implementation token {token!r}")
        require("invokedynamic" not in disassembly,
                "benchmark bytecode contains invokedynamic (no string concatenation or lambdas)")

        policy_path = self.raw / "mixed-policy.json"
        write_json(policy_path, MIXED_POLICY)

        self.workload = {
            "sources": [entry["path"] for entry in source_digests],
            "sourceDigests": source_digests,
            "classDirectory": str(self.classes),
            "classCount": len(digests),
            "classDirectoryBytes": sum(entry["bytes"] for entry in digests),
            "classes": digests,
            "classManifestSha256": digest_manifest,
            "classDirectoryDigest": dir_digest_go(self.classes),
            "javac": {
                "command": ["javac", "--release", "21", "-g", "-d", str(self.classes),
                            *[str(path) for path in sources]],
                "durationSeconds": javac_seconds,
                "classFileVersion": "65.0",
            },
            "bytecodeAudit": {
                "noInvokedynamic": True,
                "forbiddenTokensAbsent": list(FORBIDDEN_BYTECODE),
                "javapOutputSha256": sha256_bytes(disassembly.encode("utf-8")),
                "javapOutputBytes": len(disassembly.encode("utf-8")),
            },
            "main": MAIN,
            "mixedPolicy": {
                "path": str(policy_path),
                "sha256": sha256_file(policy_path),
                "json": MIXED_POLICY,
                "forcedMethod": list(HOT),
                "caller": list(CALLER),
                "callback": list(CALLBACK),
            },
            "jvmFlags": list(JVM_FLAGS),
            "kernelOrder": list(KERNELS),
            "kernelOpsPerBatch": dict(KERNEL_OPS),
            "opDefinitions": dict(OP_DEFINITIONS),
            "workloadConstants": dict(WORKLOAD_CONSTANTS),
            "timingStructure": {
                "lookupOutsideCachedInvoke": True,
                "threadKernelsStartAndJoinInsideTheTimedBatch": True,
                "clock": "System.nanoTime() around each batch only",
                "checksumReporting": "outside the measured batch window",
            },
            "buildBoundaries": {
                "compilation": "javac --release 21 -g over the checked-in sources, once",
                "packaging": "catty build reads the class directory, generates Go source and "
                             "links the standalone executable; its wall time includes image "
                             "generation, the generated Go package compile and the link",
                "inputArtifactForOpenJDK": str(self.classes),
                "artifactPerCattyMode": "standalone executable plus embedded image JSON",
            },
            "steadyCommand": self.batch_command,
            "batchesPerProcess": self.batches,
            "warmupBatches": self.warmup,
        }
        print(f"compiled {len(digests)} classes in {javac_seconds:.3f}s "
              f"(manifest {digest_manifest[:16]}...)", flush=True)

        for config in self.configs:
            if config == "java":
                continue
            self.build_config(config)

        input_digests = {config: report["inputs"][0]["digest"]
                         for config, report in self.reports.items()}
        if input_digests:
            require(len(set(input_digests.values())) == 1,
                    f"builder input digests disagree across modes: {input_digests}")
            require(next(iter(input_digests.values()))
                    == self.workload["classDirectoryDigest"],
                    "builder input digest does not match the runner's class-directory digest")
        self.workload["inputDigestByMode"] = input_digests

    def build_config(self, config):
        mode = MODE_FOR[config]
        folder = self.work / "builds" / config
        folder.mkdir(parents=True, exist_ok=True)
        policy = self.raw / "mixed-policy.json" if config == "mixed" else None
        repeats = []
        binary = None
        for index in range(self.build_repeats):
            attempt = folder / f"attempt-{index}"
            attempt.mkdir(parents=True, exist_ok=True)
            report = attempt / "report.json"
            generated = attempt / "generated"
            output = attempt / "app"
            command = [self.suite.bin, "build", "--enable-lifecycle",
                       "--classpath", self.classes, "--main", MAIN,
                       "--mode", mode, "--report", report,
                       "--generated-source", generated, "--output", output]
            if policy is not None:
                command += ["--policy", policy]
            result, duration = self.run_logged(
                "build", config, index, f"build-{config}-{index}",
                command, timeout=300, env={"GOFLAGS": ""})
            result.succeeded()
            require(output.is_file() and report.is_file(),
                    f"builder published no artifacts for {config}")
            report_data = json.loads(report.read_text(encoding="utf-8"))
            engines = self.check_report(report_data, config, report)
            image_bytes, image_data = self.extract_image(generated)
            self.check_image(image_data, report_data, config)
            executable = output.read_bytes()
            entry = {
                "attempt": index,
                "durationSeconds": duration,
                "command": [str(part) for part in command],
                "environment": {"GOFLAGS": ""},
                "executableBytes": len(executable),
                "executableSha256": sha256_bytes(executable),
                "imageJsonBytes": len(image_bytes),
                "imageJsonSha256": sha256_bytes(image_bytes),
                "reportSha256": sha256_file(report),
                "enginePlacement": engines,
                "reportInputDigest": report_data["inputs"][0]["digest"],
                "goVersionMetadata": self.go_version_metadata(output, config),
            }
            repeats.append(entry)
            if index == 0:
                binary = output
                self.reports[config] = report_data
            else:
                require(entry["executableSha256"] == repeats[0]["executableSha256"],
                        f"{config}: builder is not reproducible across identical inputs")
                require(entry["imageJsonSha256"] == repeats[0]["imageJsonSha256"],
                        f"{config}: embedded image JSON differs across identical builds")
        durations = [entry["durationSeconds"] for entry in repeats]
        self.builds[config] = {
            "binary": binary,
            "report": folder / "attempt-0" / "report.json",
            "generated": folder / "attempt-0" / "generated",
            "mode": mode,
            "repeats": repeats,
            "medianDurationSeconds": statistics.median(durations),
            "executableBytes": repeats[0]["executableBytes"],
            "executableSha256": repeats[0]["executableSha256"],
            "imageJsonBytes": repeats[0]["imageJsonBytes"],
            "imageJsonSha256": repeats[0]["imageJsonSha256"],
            "enginePlacement": repeats[0]["enginePlacement"],
        }
        print(f"built {config}: median {statistics.median(durations):.3f}s, "
              f"{repeats[0]['executableBytes']} bytes, image {repeats[0]['imageJsonBytes']} bytes",
              flush=True)

    def check_report(self, report, config, path):
        require(report["profile"] == PROFILE,
                f"{config}: wrong profile {report['profile']!r} in {path}")
        require(report["imageABI"] == IMAGE_ABI, f"{config}: wrong image ABI in {path}")
        require(report["lifecycle"] is True and report["threads"] is True
                and report["plugins"] is True,
                f"{config}: wrong Gate 5 feature flags in {path}")
        capabilities = set(report.get("capabilities") or ())
        require({"application-lifecycle", "shutdown-hooks", "application-diagnostics"}
                <= capabilities, f"{config}: missing Gate 5 capabilities in {path}")
        require({"monitorenter", "monitorexit"} <= set(report["opcodes"]),
                f"{config}: missing monitor opcodes in {path}")
        require(report["main"] == MAIN, f"{config}: wrong main class in {path}")
        require(report["toolchains"]["go"].strip(), f"{config}: missing Go toolchain in {path}")
        require(report["inputs"]
                and Path(report["inputs"][0]["origin"]).resolve() == self.classes.resolve(),
                f"{config}: wrong builder input origin in {path}")
        require(not report.get("unresolved"),
                f"{config}: unresolved references in {path}: {report.get('unresolved')}")
        bootstrap = report["bootstrap"]
        runtime_methods = {(method["name"], method["descriptor"])
                           for method in bootstrap["methods"]
                           if method["owner"] == "java/lang/Runtime"}
        require({("gc", "()V"), ("getRuntime", "()Ljava/lang/Runtime;"),
                 ("exit", "(I)V"), ("halt", "(I)V")} <= runtime_methods,
                f"{config}: bootstrap inventory is missing Runtime members")
        system_methods = {(method["name"], method["descriptor"])
                          for method in bootstrap["methods"]
                          if method["owner"] == "java/lang/System"}
        require({("exit", "(I)V"), ("nanoTime", "()J")} <= system_methods,
                f"{config}: bootstrap inventory is missing System members")

        placements = {}
        for record in report["methods"]:
            if not record["owner"].startswith("gate5bench/"):
                continue
            key = (record["owner"], record["name"], record["descriptor"])
            placements[key] = record["engine"]
        require(placements, f"{config}: report contains no application methods")
        for key, engine in placements.items():
            require(engine in ("aot", "interp", "abstract"),
                    f"{config}: invalid engine {engine!r} for {key}")
        expected = self.expected_placement(config, placements)
        mismatches = {key: (placements.get(key), engine) for key, engine in expected.items()
                      if placements.get(key) != engine}
        require(not mismatches,
                f"{config}: wrong engine placement (actual, expected): {mismatches}")
        return {"/".join(key[:2]) + key[2]: engine for key, engine in placements.items()}

    def expected_placement(self, config, placements):
        expected = {}
        for key in placements:
            if placements[key] == "abstract":
                expected[key] = "abstract"
            elif config == "interp":
                expected[key] = "interp"
            elif config == "aot":
                expected[key] = "aot"
            elif config == "mixed":
                expected[key] = "interp" if key == HOT else "aot"
            else:
                raise AssertionError("engine placement is only checked for Catty modes")
        for key in (HOT, CALLER, CALLBACK, MAIN_KEY):
            require(key in placements, f"{config}: report is missing {key}")
        if config == "mixed":
            require(placements[HOT] == "interp", "mixed: forced method did not run interpreted")
            require(placements[CALLER] == "aot", "mixed: caller did not remain AOT")
            require(placements[CALLBACK] == "aot", "mixed: callback did not remain AOT")
        return expected

    def extract_image(self, generated):
        source = (generated / "main.go").read_text(encoding="utf-8")
        encoded = decode_go_string_literal(source, "const cattyImageJSON = ")
        return encoded, json.loads(encoded.decode("utf-8"))

    def check_image(self, image, report, config):
        require(image["profile"] == PROFILE and image["abi"] == IMAGE_ABI,
                f"{config}: embedded image profile/ABI mismatch")
        require(image["lifecycle"] is True and image["threads"] is True
                and image["plugins"] is True,
                f"{config}: embedded image feature flags mismatch")
        main = image["main"]
        require((main["owner"], main["name"], main["descriptor"]) == MAIN_KEY,
                f"{config}: embedded image main mismatch: {main}")
        image_classes = {entry["name"] for entry in image["classes"]}
        report_classes = {entry["name"] for entry in report["classes"]}
        require(image_classes == report_classes,
                f"{config}: embedded image classes differ from report classes")
        report_engines = {(method["owner"], method["name"], method["descriptor"]): method["engine"]
                          for method in report["methods"]}
        for entry in image["classes"]:
            for method in entry.get("methods") or ():
                key = (method["owner"], method["name"], method["descriptor"])
                require(report_engines.get(key) == method["engine"],
                        f"{config}: embedded image engine for {key} disagrees with the report")

    def go_version_metadata(self, binary, config):
        result = self.suite.run(f"go-version-{config}", ["go", "version", "-m", binary])
        if result.status != 0:
            return {"status": "failed",
                    "stderr": result.stderr.decode("utf-8", "replace")}
        text = result.stdout.decode("utf-8", "replace")
        return {"status": "ok", "text": text, "sha256": sha256_bytes(text.encode("utf-8"))}

    # ------------------------------------------------------------------
    # correctness
    # ------------------------------------------------------------------
    def correctness_phase(self):
        oracle_verify = self.suite.run(
            "oracle-verify", self.java_command(VERIFY, ("-Xverify:all",)), timeout=120)
        oracle_verify.succeeded()
        require(not oracle_verify.stderr, "OpenJDK oracle verify wrote to stderr")
        oracle_records = parse_verify(oracle_verify.stdout.decode("utf-8"), "oracle verify")
        require(oracle_records["order"] == list(KERNELS),
                f"oracle verify kernel order is unexpected: {oracle_records['order']}")
        self.store_oracle_output(
            "verify", oracle_verify.stdout,
            extra={"jvmFlags": list(JVM_FLAGS) + ["-Xverify:all"]})
        self.record_order("correctness", "java", None,
                          [str(part) for part in self.java_command(
                              VERIFY, ("-Xverify:all",))], 0.0)

        java_verify = self.suite.run("java-verify", self.java_command(VERIFY),
                                     timeout=120).succeeded()
        require(java_verify.stdout == oracle_verify.stdout,
                "benchmark-flag Java verify disagrees with the -Xverify:all oracle")
        require(not java_verify.stderr, "benchmark-flag Java verify wrote to stderr")

        metrics = {}
        for config in self.configs:
            if config == "java":
                continue
            result = self.suite.run(f"verify-{config}", self.catty_command(config, VERIFY),
                                    env=TIMED_ENV, timeout=120).succeeded()
            require(result.stdout == oracle_verify.stdout,
                    f"{config}: verify stdout disagrees with the OpenJDK oracle. See {result.log}")
            require(not result.stderr, f"{config}: verify wrote to stderr")
            metrics_path = self.raw / f"metrics-verify-{config}.json"
            if metrics_path.exists():
                metrics_path.unlink()
            metric_result = self.suite.run(
                f"verify-metrics-{config}", self.catty_command(config, VERIFY),
                env={"CATTY_METRICS": str(metrics_path)}, timeout=120).succeeded()
            require(metric_result.stdout == oracle_verify.stdout,
                    f"{config}: instrumented verify stdout changed with metrics enabled")
            require(metrics_path.is_file(), f"{config}: metrics file was not written")
            document = read_metrics(metrics_path)
            require(document["profile"] == PROFILE
                    and document["terminalReason"] == "natural"
                    and document["status"] == 0,
                    f"{config}: unexpected metrics identity: {document['profile']} "
                    f"{document['terminalReason']} {document['status']}")
            counters = document["counters"]
            require(counters["reflectiveInvocations"] >= 1,
                    f"{config}: no reflective invocations counted")
            if config == "mixed":
                require(counters["crossEngineCalls"] > 0,
                        "mixed: metrics observed no engine-crossing calls")
            else:
                require(counters["crossEngineCalls"] == 0,
                        f"{config}: unexpected engine-crossing calls in a single-engine configuration")
            engines = counters["methodsByEngine"]
            if config == "interp":
                require(engines.get("aot", 0) == 0 and engines.get("interp", 0) > 0,
                        "interp: metrics method inventory is not exclusively interpreted")
            elif config == "aot":
                require(engines.get("interp", 0) == 0 and engines.get("aot", 0) > 0,
                        "aot: metrics method inventory is not exclusively AOT")
            else:
                require(engines.get("interp", 0) >= 1 and engines.get("aot", 0) > 0,
                        "mixed: metrics method inventory does not contain both engines")
            metrics[config] = {
                "counters": counters,
                "goRuntime": document["goRuntime"],
                "peakJavaThreads": counters["peakJavaThreads"],
            }

        cold_oracle = self.suite.run("oracle-cold", self.java_command(COLD),
                                     timeout=120).succeeded()
        require(not cold_oracle.stderr, "OpenJDK cold oracle wrote to stderr")
        parse_cold(cold_oracle.stdout.decode("utf-8"), "oracle cold")
        self.store_oracle_output("cold", cold_oracle.stdout)
        for config in self.configs:
            if config == "java":
                result = self.suite.run("java-cold", self.java_command(COLD),
                                        timeout=120).succeeded()
            else:
                result = self.suite.run(f"cold-{config}", self.catty_command(config, COLD),
                                        env=TIMED_ENV, timeout=120).succeeded()
            require(result.stdout == cold_oracle.stdout,
                    f"{config}: cold stdout disagrees with the OpenJDK oracle")
            require(not result.stderr, f"{config}: cold wrote to stderr")

        steady_oracle = self.suite.run("oracle-steady", self.java_command(self.batch_command),
                                       timeout=300).succeeded()
        require(not steady_oracle.stderr, "OpenJDK steady oracle wrote to stderr")
        steady_oracle_kernels = parse_steady(
            steady_oracle.stdout.decode("utf-8"), "oracle steady")
        require([kernel["name"] for kernel in steady_oracle_kernels] == list(KERNELS),
                "oracle steady kernel order is unexpected")
        self.check_steady_metadata(steady_oracle_kernels, "oracle steady")
        self.store_oracle_output("steady", steady_oracle.stdout)
        steady_expected = steady_checksums(steady_oracle_kernels)
        for config in self.configs:
            if config == "java":
                result = self.suite.run("java-steady", self.java_command(self.batch_command),
                                        timeout=300).succeeded()
            else:
                result = self.suite.run(f"steady-{config}",
                                        self.catty_command(config, self.batch_command),
                                        env=TIMED_ENV, timeout=300).succeeded()
            kernels = parse_steady(result.stdout.decode("utf-8"),
                                   f"{config} steady")
            require(steady_checksums(kernels) == steady_expected,
                    f"{config}: steady checksums disagree with the OpenJDK oracle")
            require(not result.stderr, f"{config}: steady correctness run wrote to stderr")

        rss_oracle = self.suite.run("oracle-rss", self.java_command(RSS),
                                    timeout=180).succeeded()
        require(not rss_oracle.stderr, "OpenJDK rss oracle wrote to stderr")
        parse_rss(rss_oracle.stdout.decode("utf-8"), "oracle rss")
        self.store_oracle_output("rss", rss_oracle.stdout)
        for config in self.configs:
            if config == "java":
                result = self.suite.run("java-rss", self.java_command(RSS),
                                        timeout=180).succeeded()
            else:
                result = self.suite.run(f"rss-{config}", self.catty_command(config, RSS),
                                        env=TIMED_ENV, timeout=180).succeeded()
            require(result.stdout == rss_oracle.stdout,
                    f"{config}: rss stdout disagrees with the OpenJDK oracle")
            require(not result.stderr, f"{config}: rss wrote to stderr")

        gc_log = self.raw / "oracle-retained-gc.log"
        if gc_log.exists():
            gc_log.unlink()
        retained_oracle = self.suite.run(
            "oracle-retained",
            self.java_command(RETAINED, (f"-Xlog:gc:file={gc_log}",)),
            timeout=120).succeeded()
        require(not retained_oracle.stderr, "OpenJDK retained oracle wrote to stderr")
        retained_values = parse_retained(retained_oracle.stdout.decode("utf-8"),
                                         "oracle retained")
        oracle_gc = parse_gc_log(gc_log)
        require(oracle_gc["status"] == "ok",
                f"OpenJDK retained oracle did not confirm the explicit collection: {oracle_gc}")
        self.store_oracle_output("retained", retained_oracle.stdout)
        retained_correctness = {}
        for config in self.configs:
            if config == "java":
                retained_correctness[config] = {"stdoutMatches": True,
                                                "collectionConfirmed": True}
                continue
            metrics_path = self.raw / f"metrics-retained-correctness-{config}.json"
            if metrics_path.exists():
                metrics_path.unlink()
            result = self.suite.run(
                f"retained-correctness-{config}", self.catty_command(config, RETAINED),
                env={"CATTY_METRICS": str(metrics_path)}, timeout=120).succeeded()
            require(result.stdout == retained_oracle.stdout,
                    f"{config}: retained stdout disagrees with the OpenJDK oracle")
            require(not result.stderr, f"{config}: retained correctness run wrote to stderr")
            require(metrics_path.is_file(), f"{config}: retained metrics file missing")
            document = read_metrics(metrics_path)
            runtime = document["goRuntime"]
            require(runtime.get("postGCHeapBytes") is not None,
                    f"{config}: retained run recorded no post-Runtime.gc() heap observation")
            require(runtime["numGC"] >= 1,
                    f"{config}: retained run recorded no completed GC")
            retained_correctness[config] = {
                "stdoutMatches": True,
                "collectionConfirmed": runtime.get("postGCHeapBytes") is not None,
                "postGCHeapBytes": runtime.get("postGCHeapBytes"),
            }

        self.correctness = {
            "status": "ok",
            "oracle": {
                "verify": {"kernelOrder": oracle_records["order"],
                           "work": oracle_records["work"]},
                "cold": {"stdoutSha256": sha256_bytes(cold_oracle.stdout)},
                "steady": {"checksums": steady_expected,
                           "batchCommand": self.batch_command,
                           "batches": self.batches,
                           "warmup": self.warmup},
                "rss": {"stdoutSha256": sha256_bytes(rss_oracle.stdout)},
                "retained": {"values": retained_values, "gc": oracle_gc},
            },
            "enginePlacement": {
                "expected": {
                    "interp": "all application methods interpreted",
                    "mixed": {"forced": list(HOT), "forcedEngine": "interp",
                              "caller": list(CALLER), "callerEngine": "aot",
                              "callback": list(CALLBACK), "callbackEngine": "aot"},
                    "aot": "all application methods AOT",
                },
                "reports": {config: self.builds[config]["enginePlacement"]
                            for config in CATTY_CONFIGS if config in self.builds},
                "selectedTargets": {
                    config: {
                        name: self.builds[config]["enginePlacement"].get(
                            "/".join(key[:2]) + key[2])
                        for name, key in SELECTED_METHODS.items()
                    } for config in CATTY_CONFIGS if config in self.builds
                },
                "virtualReceiverTypes": ["gate5bench/Calls$Alpha", "gate5bench/Calls$Beta"],
                "reflectionLookupName": "combine",
            },
            "runtimeMetrics": metrics,
            "retained": retained_correctness,
            "openjdkVerifyAll": {
                "command": [str(part) for part in self.java_command(VERIFY, ("-Xverify:all",))],
                "status": oracle_verify.status,
                "stdoutSha256": sha256_bytes(oracle_verify.stdout),
            },
        }
        built = [config for config in CATTY_CONFIGS if config in self.builds]
        if "mixed" in built:
            detail = "mixed engine placement (forced interp, AOT caller/callback) and crossing verified"
        else:
            detail = "engine placement verified for " + ", ".join(built)
        print(f"correctness: verify/cold/steady/rss/retained match OpenJDK; {detail}",
              flush=True)

    def check_steady_metadata(self, kernels, context):
        require(len(kernels) == len(KERNELS), f"{context}: expected {len(KERNELS)} kernels")
        for kernel in kernels:
            require(kernel["batches"] == self.batches,
                    f"{context}: expected {self.batches} batches for {kernel['name']}")
            require(kernel["warmup"] == self.warmup,
                    f"{context}: expected warmup {self.warmup} for {kernel['name']}")
            require(len(kernel["samples"]) == self.batches,
                    f"{context}: {kernel['name']} has {len(kernel['samples'])} batch samples")
            require(kernel["opsPerBatch"] == KERNEL_OPS[kernel["name"]],
                    f"{context}: {kernel['name']} reports {kernel['opsPerBatch']} ops, "
                    f"expected {KERNEL_OPS[kernel['name']]}")
            require(all(sample > 0 for sample in kernel["samples"]),
                    f"{context}: {kernel['name']} has a non-positive nanoTime batch")

    def store_oracle_output(self, name, stdout, extra=None):
        path = self.raw / "oracle" / f"{name}.stdout"
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_bytes(stdout)
        value = {"path": str(path), "bytes": len(stdout), "sha256": sha256_bytes(stdout)}
        if extra:
            value.update(extra)
        self.oracle[name] = value

    # ------------------------------------------------------------------
    # cold samples
    # ------------------------------------------------------------------
    def cold_phase(self):
        expected = (self.raw / "oracle" / "cold.stdout").read_bytes()
        expected_record = parse_cold(expected.decode("utf-8"), "cold oracle")
        samples_by_config = {config: [] for config in self.configs}
        for index in range(self.repeat):
            for config in rotated(self.configs, index):
                samples_by_config[config].append(
                    self.cold_sample(config, index, expected, expected_record))
        for config in self.configs:
            successful = [sample for sample in samples_by_config[config]
                          if sample["status"] == "ok"]
            if len(successful) != self.repeat:
                self.add_failure("cold", config,
                                 f"{len(successful)}/{self.repeat} successful cold samples")
            self.cold[config] = {
                "samples": samples_by_config[config],
                "successfulSamples": len(successful),
                "requestedSamples": self.repeat,
                "elapsedToResult": describe(
                    [sample["elapsedToResultSeconds"] for sample in successful]),
                "totalProcessSeconds": describe(
                    [sample["totalSeconds"] for sample in successful]),
                "checksum": expected_record["checksum"],
                "ops": expected_record["ops"],
            }

    def cold_sample(self, config, index, expected_stdout, expected_record):
        command = self.command_for(config, COLD)
        env = self.run_env(TIMED_ENV)
        out_path = self.work / "child" / f"cold-{config}-{index}.stdout"
        err_path = self.work / "child" / f"cold-{config}-{index}.stderr"
        out_path.parent.mkdir(parents=True, exist_ok=True)
        started = time.monotonic()
        measurement = self.launch_cold(command, env, started, out_path, err_path, timeout=120)
        total = time.monotonic() - started
        self.record_order("cold", config, index, [str(part) for part in command], total)
        sample = {
            "config": config,
            "sample": index,
            "pid": measurement["pid"],
            "status": "failed",
            "command": [str(part) for part in command],
            "launchEpochSeconds": measurement["launchEpochSeconds"],
            "elapsedToResultSeconds": measurement["elapsedToResultSeconds"],
            "totalSeconds": total,
            "timedOut": measurement["timedOut"],
            "stdoutMatchesOracle": False,
            "checksum": None,
            "ops": None,
            "failure": None,
        }
        if measurement["timedOut"]:
            sample["failure"] = "cold sample timed out"
        elif measurement["returncode"] != 0:
            sample["failure"] = f"cold sample exited with status {measurement['returncode']}"
        elif measurement["stderr"]:
            sample["failure"] = "cold sample wrote to stderr"
        elif measurement["stdout"] != expected_stdout:
            sample["failure"] = "cold sample stdout disagrees with the OpenJDK oracle"
        elif measurement["elapsedToResultSeconds"] is None:
            sample["failure"] = "cold sample never printed a RESULT record"
        else:
            try:
                record = parse_cold(measurement["stdout"].decode("utf-8"), "cold sample")
            except AssertionError as error:
                sample["failure"] = str(error)
            else:
                if record["checksum"] != expected_record["checksum"]:
                    sample["failure"] = "cold checksum mismatch"
                else:
                    sample.update(status="ok", stdoutMatchesOracle=True,
                                  checksum=record["checksum"], ops=record["ops"])
        if sample["status"] != "ok":
            self.add_failure("cold", config, f"sample {index}: {sample['failure']}")
        return sample

    def launch_cold(self, command, env, started, out_path, err_path, timeout):
        """Launch one process and timestamp the first complete RESULT record."""
        with open(err_path, "wb") as stderr_file:
            process = subprocess.Popen(
                [str(part) for part in command], cwd=str(ROOT), env=env,
                stdin=subprocess.DEVNULL, stdout=subprocess.PIPE, stderr=stderr_file)
            fd = process.stdout.fileno()
            raw = bytearray()
            pending = b""
            stage = 0
            result_time = None
            deadline = started + timeout
            timed_out = False
            while True:
                remaining = deadline - time.monotonic()
                if remaining <= 0:
                    timed_out = True
                    process.kill()
                    break
                ready, _, _ = select.select([fd], [], [], min(0.05, remaining))
                if not ready:
                    if process.poll() is not None:
                        while True:
                            chunk = os.read(fd, 65536)
                            if not chunk:
                                break
                            raw.extend(chunk)
                            pending += chunk
                        break
                    continue
                chunk = os.read(fd, 65536)
                if not chunk:
                    break
                raw.extend(chunk)
                pending += chunk
                lines = pending.split(b"\n")
                pending = lines.pop()
                for line in lines:
                    text = line.decode("utf-8", errors="replace")
                    if stage == 0 and text == "RESULT":
                        stage = 1
                    elif stage == 1:
                        result_time = time.monotonic()
                        stage = 2
            if process.poll() is None:
                try:
                    process.wait(timeout=10)
                except subprocess.TimeoutExpired:
                    process.kill()
                    process.wait()
            returncode = process.returncode
        return {
            "pid": process.pid,
            "returncode": returncode,
            "timedOut": timed_out,
            "launchEpochSeconds": started,
            "elapsedToResultSeconds": (result_time - started) if result_time is not None else None,
            "stdout": bytes(raw),
            "stderr": Path(err_path).read_bytes(),
        }

    # ------------------------------------------------------------------
    # steady samples
    # ------------------------------------------------------------------
    def steady_phase(self):
        expected = self.correctness["oracle"]["steady"]["checksums"]
        for config in self.configs:
            processes = [self.steady_sample(config, index, expected)
                         for index in range(self.repeat)]
            self.steady[config] = self.summarize_steady(config, processes, expected)

    def steady_sample(self, config, index, expected_checksums):
        command = self.command_for(config, self.batch_command)
        started = time.monotonic()
        result = self.suite.run(
            f"timed-steady-{config}-{index}", command, env=TIMED_ENV, timeout=600)
        duration = time.monotonic() - started
        self.record_order("steady", config, index,
                          [str(part) for part in command], duration)
        entry = {
            "config": config,
            "sample": index,
            "pid": result.pid,
            "status": "failed",
            "command": [str(part) for part in command],
            "wallSeconds": duration,
            "kernels": [],
            "failure": None,
            "stdoutSha256": sha256_bytes(result.stdout),
        }
        if result.status != 0:
            entry["failure"] = f"steady process exited with status {result.status}"
        elif result.stderr:
            entry["failure"] = "steady process wrote to stderr"
        else:
            try:
                kernels = parse_steady(result.stdout.decode("utf-8"),
                                       f"{config} steady sample")
                self.check_steady_metadata(kernels, f"{config} steady sample")
            except AssertionError as error:
                entry["failure"] = str(error)
            else:
                if steady_checksums(kernels) != expected_checksums:
                    entry["failure"] = "steady checksums disagree with the OpenJDK oracle"
                else:
                    entry.update(status="ok", kernels=kernels)
        if entry["status"] != "ok":
            self.add_failure("steady", config, f"sample {index}: {entry['failure']}")
        return entry

    def summarize_steady(self, config, processes, expected_checksums):
        successful = [process for process in processes if process["status"] == "ok"]
        if len(successful) != self.repeat:
            self.add_failure("steady", config,
                             f"{len(successful)}/{self.repeat} successful steady processes")
        summary = {
            "processes": len(successful),
            "requestedProcesses": self.repeat,
            "batchesPerProcess": self.batches,
            "warmupBatches": self.warmup,
            "command": self.batch_command,
            "kernels": {},
        }
        expected = dict(expected_checksums)
        for kernel_name in KERNELS:
            per_process = []
            pooled_ns_per_op = []
            total_ops = 0
            total_seconds = 0.0
            ops_per_batch = None
            for process in successful:
                kernel = next((entry for entry in process["kernels"]
                               if entry["name"] == kernel_name), None)
                require(kernel is not None,
                        f"{config}: successful steady process lacks kernel {kernel_name}")
                ops_per_batch = kernel["opsPerBatch"]
                process_seconds = sum(kernel["samples"]) / 1e9
                process_ops = kernel["batches"] * ops_per_batch
                per_process.append({
                    "process": process["sample"],
                    "pid": process["pid"],
                    "batches": kernel["batches"],
                    "ops": process_ops,
                    "seconds": process_seconds,
                    "opsPerSecond": (process_ops / process_seconds) if process_seconds else None,
                    "normalizedNsPerOpP50": quantile(
                        [sample / ops_per_batch for sample in kernel["samples"]], 0.50),
                })
                total_ops += process_ops
                total_seconds += process_seconds
                for batch_index, nanos in enumerate(kernel["samples"]):
                    normalized = nanos / ops_per_batch
                    pooled_ns_per_op.append(normalized)
                    self.steady_csv_rows.append(
                        [config, process["sample"], kernel_name, batch_index,
                         nanos, ops_per_batch, normalized])
            summary["kernels"][kernel_name] = {
                "opsPerBatch": ops_per_batch,
                "measuredBatches": len(pooled_ns_per_op),
                "totalCompletedOps": total_ops,
                "totalMeasuredSeconds": total_seconds,
                "aggregateOpsPerSecond": (total_ops / total_seconds) if total_seconds else None,
                "normalizedNsPerOp": describe(pooled_ns_per_op),
                "perProcess": per_process,
                "expectedChecksum": expected[kernel_name],
            }
        return summary

    # ------------------------------------------------------------------
    # memory samples
    # ------------------------------------------------------------------
    def memory_phase(self):
        expected_rss = (self.raw / "oracle" / "rss.stdout").read_bytes()
        expected_retained = (self.raw / "oracle" / "retained.stdout").read_bytes()
        rss_by_config = {config: [] for config in self.configs}
        retained_by_config = {config: [] for config in self.configs}
        for index in range(self.repeat):
            for config in rotated(self.configs, index):
                rss_by_config[config].append(
                    self.rss_sample(config, index, expected_rss))
            for config in rotated(self.configs, index):
                retained_by_config[config].append(
                    self.retained_sample(config, index, expected_retained))
        for config in self.configs:
            successful = [sample for sample in rss_by_config[config]
                          if sample["status"] == "ok"]
            if len(successful) != self.repeat:
                self.add_failure("memory", config,
                                 f"{len(successful)}/{self.repeat} successful RSS samples")
            rss_values = [sample["peakRssBytes"] for sample in successful]
            self.memory[config] = {
                "samples": rss_by_config[config],
                "successfulSamples": len(successful),
                "requestedSamples": self.repeat,
                "peakRss": describe(rss_values),
                "peakRssMiB": describe([value / (1024.0 * 1024.0) for value in rss_values]),
                "rssMethod": successful[0]["rssMethod"] if successful else None,
                "rssSource": successful[0]["rssSource"] if successful else None,
                "rssUnits": successful[0]["peakRssUnits"] if successful else None,
            }
            retained_samples = retained_by_config[config]
            retained_successful = [sample for sample in retained_samples
                                   if sample["status"] == "ok"]
            if len(retained_successful) != self.repeat:
                self.add_failure("retained", config,
                                 f"{len(retained_successful)}/{self.repeat} "
                                 "successful retained samples")
            heap_after = [sample["heapAfterBytes"] for sample in retained_successful
                          if sample.get("heapAfterBytes") is not None]
            rss_retained = [sample["peakRssBytes"] for sample in retained_successful
                            if sample.get("peakRssBytes") is not None]
            self.retained[config] = {
                "samples": retained_samples,
                "successfulSamples": len(retained_successful),
                "requestedSamples": self.repeat,
                "retainedItems": (retained_successful[0]["retainedItems"]
                                  if retained_successful else None),
                "collectionConfirmed": (all(sample.get("collectionConfirmed")
                                            for sample in retained_successful)
                                        if retained_successful else False),
                "heapMetric": (retained_successful[0]["heapMetric"]
                               if retained_successful else None),
                "heapAfterBytes": describe(heap_after),
                "peakRssBytes": describe(rss_retained),
            }

    def rss_sample(self, config, index, expected_stdout):
        command = self.command_for(config, RSS)
        env = self.run_env(TIMED_ENV)
        out_path = self.work / "child" / f"rss-{config}-{index}.stdout"
        err_path = self.work / "child" / f"rss-{config}-{index}.stderr"
        out_path.parent.mkdir(parents=True, exist_ok=True)
        started = time.monotonic()
        measurement = self.measure_child(command, env, 180, out_path, err_path)
        wall = time.monotonic() - started
        self.record_order("memory", config, index, [str(part) for part in command], wall)
        sample = {
            "config": config,
            "sample": index,
            "pid": measurement["pid"],
            "status": "failed",
            "command": [str(part) for part in command],
            "wallSeconds": wall,
            "peakRssRaw": measurement["peakRssRaw"],
            "peakRssUnits": measurement["peakRssUnits"],
            "peakRssBytes": rss_to_bytes(measurement["peakRssRaw"],
                                         measurement["peakRssUnits"]),
            "rssMethod": measurement["rssMethod"],
            "rssSource": measurement["rssSource"],
            "failure": None,
        }
        failure = self.check_child(measurement, expected_stdout, "rss sample")
        if failure is not None:
            sample["failure"] = failure
        else:
            try:
                parse_rss(measurement["stdout"].decode("utf-8"), "rss sample")
            except AssertionError as error:
                sample["failure"] = str(error)
            else:
                if sample["peakRssBytes"] is None:
                    sample["failure"] = "peak RSS metric unavailable for this platform"
                else:
                    sample["status"] = "ok"
        if sample["status"] != "ok":
            self.add_failure("memory", config, f"sample {index}: {sample['failure']}")
        return sample

    def retained_sample(self, config, index, expected_stdout):
        extra = ()
        env = self.run_env(TIMED_ENV)
        metrics_path = None
        gc_log = None
        if config == "java":
            gc_log = self.work / "child" / f"retained-{config}-{index}.gc.log"
            gc_log.parent.mkdir(parents=True, exist_ok=True)
            if gc_log.exists():
                gc_log.unlink()
            extra = (f"-Xlog:gc:file={gc_log}",)
        else:
            metrics_path = self.work / "child" / f"retained-{config}-{index}.metrics.json"
            metrics_path.parent.mkdir(parents=True, exist_ok=True)
            if metrics_path.exists():
                metrics_path.unlink()
            env["CATTY_METRICS"] = str(metrics_path)
        command = self.command_for(config, RETAINED, extra)
        out_path = self.work / "child" / f"retained-{config}-{index}.stdout"
        err_path = self.work / "child" / f"retained-{config}-{index}.stderr"
        started = time.monotonic()
        measurement = self.measure_child(command, env, 180, out_path, err_path)
        wall = time.monotonic() - started
        self.record_order("retained", config, index,
                          [str(part) for part in command], wall)
        sample = {
            "config": config,
            "sample": index,
            "pid": measurement["pid"],
            "status": "failed",
            "command": [str(part) for part in command],
            "wallSeconds": wall,
            "peakRssRaw": measurement["peakRssRaw"],
            "peakRssUnits": measurement["peakRssUnits"],
            "peakRssBytes": rss_to_bytes(measurement["peakRssRaw"],
                                         measurement["peakRssUnits"]),
            "rssMethod": measurement["rssMethod"],
            "rssSource": measurement["rssSource"],
            "retainedItems": None,
            "heapAfterBytes": None,
            "heapMetric": None,
            "collectionConfirmed": False,
            "failure": None,
        }
        failure = self.check_child(measurement, expected_stdout, "retained sample")
        if failure is not None:
            sample["failure"] = failure
        else:
            try:
                values = parse_retained(measurement["stdout"].decode("utf-8"),
                                        "retained sample")
            except AssertionError as error:
                sample["failure"] = str(error)
            else:
                sample["retainedItems"] = values["RETAINED"]
                if config == "java":
                    observation = parse_gc_log(gc_log)
                    if observation["status"] != "ok":
                        sample["failure"] = observation["reason"]
                    else:
                        sample["heapAfterBytes"] = observation["heapAfterBytes"]
                        sample["heapMetric"] = {
                            "scope": "OpenJDK Java heap occupancy after the explicit collection",
                            "source": observation["source"],
                            "units": observation["units"],
                            "rawLine": observation["rawLine"],
                            "collectionConfirmed": True,
                        }
                        sample["collectionConfirmed"] = True
                elif not metrics_path.is_file():
                    sample["failure"] = "Catty metrics file missing for retained sample"
                else:
                    document = read_metrics(metrics_path)
                    runtime = document["goRuntime"]
                    if runtime.get("postGCHeapBytes") is None:
                        sample["failure"] = "no post-Runtime.gc() Go heap observation in metrics"
                    elif runtime["numGC"] < 1:
                        sample["failure"] = "metrics recorded no completed GC"
                    else:
                        sample["heapAfterBytes"] = runtime["postGCHeapBytes"]
                        sample["heapMetric"] = {
                            "scope": "Catty live Go-managed heap after Runtime.gc(), "
                                     "including runtime structures",
                            "source": "CATTY_METRICS goRuntime.postGCHeapBytes",
                            "units": "bytes",
                            "scopeLabel": runtime.get("postGCHeapScope"),
                            "heapAllocBytes": runtime.get("heapAllocBytes"),
                            "totalAllocBytes": runtime.get("totalAllocBytes"),
                            "numGC": runtime.get("numGC"),
                            "collectionConfirmed": True,
                        }
                        sample["collectionConfirmed"] = True
                if sample["failure"] is None:
                    sample["status"] = "ok"
        if sample["status"] != "ok":
            self.add_failure("retained", config, f"sample {index}: {sample['failure']}")
        return sample

    def measure_child(self, command, env, timeout, out_path, err_path):
        measurement = None
        if hasattr(os, "wait4"):
            try:
                measurement = wait4_measure(command, ROOT, env, timeout, out_path, err_path)
            except (AttributeError, OSError, ChildProcessError):
                measurement = None
        if measurement is None:
            try:
                measurement = time_wrapper_measure(command, ROOT, env, timeout,
                                                   out_path, err_path)
            except subprocess.TimeoutExpired:
                measurement = {"pid": -1, "returncode": -9, "timedOut": True,
                               "wallSeconds": timeout, "peakRssRaw": None,
                               "peakRssUnits": "unavailable",
                               "rssMethod": "unavailable",
                               "rssSource": "peak RSS measurement timed out",
                               "stdout": b"", "stderr": b""}
            except OSError:
                measurement = None
        if measurement is None:
            measurement = direct_measure(command, ROOT, env, timeout, out_path, err_path)
            self.add_failure(
                "memory", "all",
                "per-child peak RSS unavailable (os.wait4 and /usr/bin/time -l/-v "
                "both failed); the sample output is still captured and the missing "
                "metric is recorded rather than fabricated")
        return measurement

    def check_child(self, measurement, expected_stdout, context):
        if measurement["timedOut"]:
            return f"{context} timed out"
        if measurement["returncode"] != 0:
            return f"{context} exited with status {measurement['returncode']}"
        if measurement["stderr"] and not measurement.get("stderrSharedWithRssTool"):
            return f"{context} wrote to stderr"
        if measurement["stdout"] != expected_stdout:
            return f"{context} stdout disagrees with the OpenJDK oracle"
        return None

    # ------------------------------------------------------------------
    # summary and evidence
    # ------------------------------------------------------------------
    def observations(self):
        notes = []
        for config in self.configs:
            kernels = self.steady.get(config, {}).get("kernels", {})
            if kernels.get("shared-1") and kernels.get("shared-4"):
                one = kernels["shared-1"]["aggregateOpsPerSecond"]
                four = kernels["shared-4"]["aggregateOpsPerSecond"]
                if one and four:
                    notes.append(
                        f"{config}: shared-4 reaches {four / one * 100:.1f}% of shared-1 "
                        "aggregate ops/s for the same total operation count over four fixed "
                        "threads (thread start/join and shared-memory synchronization included)")
            if kernels.get("calls-static") and kernels.get("calls-mixed"):
                direct = kernels["calls-static"]["aggregateOpsPerSecond"]
                crossing = kernels["calls-mixed"]["aggregateOpsPerSecond"]
                if direct and crossing:
                    notes.append(
                        f"{config}: calls-mixed reaches {crossing / direct * 100:.1f}% of "
                        "calls-static aggregate ops/s (the forced interpretation crossing "
                        "versus a direct call)")
            if kernels.get("app"):
                app = kernels["app"]["normalizedNsPerOp"]
                if app:
                    notes.append(
                        f"{config}: app batch latency p50 {app['p50']:.1f} ns/op, "
                        f"p99 {app['p99']:.1f} ns/op (batch elapsed divided by 1024 node "
                        "visits; not per-node latency)")
        for config in self.configs:
            retained = self.retained.get(config, {})
            metric = retained.get("heapMetric") or {}
            if retained.get("heapAfterBytes") and metric:
                notes.append(
                    f"{config}: retained graph ({retained.get('retainedItems')} items) "
                    f"median post-GC heap {retained['heapAfterBytes']['median'] / (1024 * 1024):.2f} "
                    f"MiB; scope: {metric.get('scope')}")
        notes.append(
            "Scope: OpenJDK 21 executes JIT-compiled bytecode with a fixed 256 MiB SerialGC "
            "heap; Catty runs the same class files as standalone Go images with no JIT. "
            "Measurements are single-host, unpinned, and exclude trace, metrics and GC "
            "probes from every timed run except the dedicated retained-graph samples.")
        return notes

    def summary_value(self, status):
        return {
            "schemaVersion": 1,
            "kind": "gate5-application-benchmark",
            "status": status,
            "mode": "smoke" if self.args.smoke else "full",
            "startedAt": self.started_at,
            "finishedAt": self.finished_at,
            "durationSeconds": self.duration,
            "workDir": str(self.work),
            "workDirRetained": self.workDirRetained,
            "evidenceDir": str(self.args.out) if self.args.out else None,
            "environment": self.environment,
            "workload": self.workload,
            "build": {
                "javacSeconds": self.workload["javac"]["durationSeconds"],
                "modes": {config: self.build_summary(config)
                          for config in CATTY_CONFIGS if config in self.builds},
            },
            "correctness": self.correctness,
            "cold": self.cold,
            "steady": self.steady,
            "memory": self.memory,
            "retained": self.retained,
            "executionOrder": self.order,
            "failures": self.failures,
            "observations": self.observations(),
            "quantileMethod": QUANTILE_METHOD,
            "limitations": [
                "Cold launch samples measure process launch to the first verified RESULT "
                "record with a monotonic clock; the OS filesystem cache is warm after the "
                "correctness phase and is not dropped (no root cache purge).",
                "Steady per-batch times include fixed platform-thread start/join cost for the "
                "shared-4 and monitor kernels; the thread count is fixed at 4 per batch.",
                "Normalized ns/op is batch elapsed time divided by the documented operation "
                "count; it must not be read as individual operation latency.",
                "The Catty retained heap observation is process-wide live Go-managed heap "
                "(including runtime structures); the OpenJDK value is post-collection Java "
                "heap occupancy. The two heap metrics have different runtime scopes.",
                "OpenJDK steady numbers come from fresh JVM processes with a fixed warmup; "
                "Catty interp/AOT runs have no JIT. The comparison is the runtimes as "
                "shipped, not equivalent generated machine code.",
                "Peak RSS on Darwin is reported by wait4 ru_maxrss in bytes; on Linux "
                "ru_maxrss is in KiB. A /usr/bin/time -l/-v fallback exists when os.wait4 is "
                "unavailable.",
                "Only one host/CPU was measured; no CPU pinning or cgroup isolation was "
                "applied.",
            ],
        }

    def build_summary(self, config):
        build = self.builds[config]
        return {
            "mode": build["mode"],
            "medianBuildSeconds": build["medianDurationSeconds"],
            "buildSeconds": [entry["durationSeconds"] for entry in build["repeats"]],
            "executableBytes": build["executableBytes"],
            "executableSha256": build["executableSha256"],
            "imageJsonBytes": build["imageJsonBytes"],
            "imageJsonSha256": build["imageJsonSha256"],
            "reportSha256": build["repeats"][0]["reportSha256"],
            "reproducible": True,
            "enginePlacement": build["enginePlacement"],
            "repeats": build["repeats"],
        }

    def publish(self, out_dir):
        out_dir = Path(out_dir)
        out_dir.mkdir(parents=True, exist_ok=True)
        raw_out = out_dir / "raw"
        if raw_out.exists():
            shutil.rmtree(raw_out)
        shutil.copytree(self.raw, raw_out)
        logs_out = out_dir / "logs"
        if logs_out.exists():
            shutil.rmtree(logs_out)
        if (self.work / "logs").is_dir():
            shutil.copytree(self.work / "logs", logs_out)
        build_out = out_dir / "build"
        if build_out.exists():
            shutil.rmtree(build_out)
        build_out.mkdir(parents=True, exist_ok=True)
        # The retained main.go is inspection evidence, not a buildable package:
        # a nested go.mod keeps `go vet ./...` and `go test ./...` out of it.
        (build_out / "go.mod").write_text(
            "module catty-gate5-benchmark-evidence\n\ngo 1.21\n", encoding="utf-8")
        for config in CATTY_CONFIGS:
            if config not in self.builds:
                continue
            folder = build_out / config
            folder.mkdir(parents=True, exist_ok=True)
            shutil.copy2(self.builds[config]["report"], folder / "report.json")
            shutil.copy2(self.builds[config]["generated"] / "main.go",
                         folder / "main.go")
        if (self.raw / "mixed-policy.json").is_file():
            shutil.copy2(self.raw / "mixed-policy.json", build_out / "mixed-policy.json")

    def write_raw(self):
        write_json(self.raw / "class-digests.json", {
            "classDirectory": str(self.classes),
            "classManifestSha256": self.workload["classManifestSha256"],
            "classDirectoryDigest": self.workload["classDirectoryDigest"],
            "classes": self.workload["classes"],
            "sources": self.workload["sourceDigests"],
        })
        write_json(self.raw / "build.json", {
            "javac": self.workload["javac"],
            "goCachePolicy": self.environment["processEnvironment"]["GOCACHEPolicy"],
            "modes": {config: self.build_summary(config)
                      for config in CATTY_CONFIGS if config in self.builds},
        })
        write_json(self.raw / "environment.json", self.environment)
        write_json(self.raw / "correctness.json", self.correctness)
        write_json(self.raw / "cold.json", self.cold)
        write_csv(self.raw / "cold.csv",
                  ["config", "sample", "pid", "status", "elapsedToResultSeconds",
                   "totalSeconds", "checksum", "ops", "failure"],
                  [[sample["config"], sample["sample"], sample["pid"], sample["status"],
                    sample["elapsedToResultSeconds"], sample["totalSeconds"],
                    sample["checksum"], sample["ops"], sample["failure"]]
                   for config in self.configs
                   for sample in self.cold.get(config, {}).get("samples", [])])
        write_json(self.raw / "steady.json", self.steady)
        write_csv(self.raw / "steady-batches.csv",
                  ["config", "process", "kernel", "batch", "nanos", "opsPerBatch",
                   "nsPerOp"],
                  self.steady_csv_rows)
        write_json(self.raw / "memory.json", self.memory)
        write_csv(self.raw / "memory.csv",
                  ["config", "sample", "pid", "status", "peakRssRaw", "peakRssUnits",
                   "peakRssBytes", "rssMethod", "failure"],
                  [[sample["config"], sample["sample"], sample["pid"], sample["status"],
                    sample["peakRssRaw"], sample["peakRssUnits"],
                    sample["peakRssBytes"], sample["rssMethod"], sample["failure"]]
                   for config in self.configs
                   for sample in self.memory.get(config, {}).get("samples", [])])
        write_json(self.raw / "retained.json", self.retained)
        write_csv(self.raw / "retained-samples.csv",
                  ["config", "sample", "pid", "status", "retainedItems",
                   "heapAfterBytes", "heapMetricScope", "collectionConfirmed",
                   "peakRssBytes", "failure"],
                  [[sample["config"], sample["sample"], sample["pid"], sample["status"],
                    sample["retainedItems"], sample["heapAfterBytes"],
                    (sample["heapMetric"] or {}).get("scope"),
                    sample["collectionConfirmed"], sample["peakRssBytes"],
                    sample["failure"]]
                   for config in self.configs
                   for sample in self.retained.get(config, {}).get("samples", [])])

    def print_summary(self, status):
        print("", flush=True)
        print(f"Gate 5 application benchmark ({'smoke' if self.args.smoke else 'full'}): "
              f"{status.upper()}", flush=True)
        print(f"artifacts: {self.work}", flush=True)
        if self.args.out:
            print(f"evidence: {self.args.out}", flush=True)
        if self.workload:
            print(f"javac: {self.workload['javac']['durationSeconds']:.3f}s, "
                  f"{self.workload['classCount']} classes, "
                  f"manifest {self.workload['classManifestSha256'][:16]}", flush=True)
        for config in CATTY_CONFIGS:
            if config in self.builds:
                build = self.builds[config]
                print(f"build {config:6s}: median {build['medianDurationSeconds']:.3f}s, "
                      f"{build['executableBytes']} bytes, "
                      f"image {build['imageJsonBytes']} bytes, "
                      f"sha256 {build['executableSha256'][:16]}", flush=True)
        if self.cold:
            parts = []
            for config in self.configs:
                data = self.cold[config]["elapsedToResult"]
                if data:
                    parts.append(f"{config}={data['median'] * 1000:.2f}ms")
            print("cold launch->verified result (median): " + ", ".join(parts), flush=True)
        if self.steady:
            print("steady aggregate ops/s (normalized ns/op p50/p99):", flush=True)
            header = "kernel".ljust(14) + "".join(config.rjust(24) for config in self.configs)
            print(header, flush=True)
            for kernel in KERNELS:
                row = kernel.ljust(14)
                for config in self.configs:
                    data = self.steady.get(config, {}).get("kernels", {}).get(kernel)
                    if not data or not data["aggregateOpsPerSecond"]:
                        row += "n/a".rjust(24)
                        continue
                    ns = data["normalizedNsPerOp"]
                    row += (f"{data['aggregateOpsPerSecond']:.0f} "
                            f"{ns['p50']:.1f}/{ns['p99']:.1f}").rjust(24)
                print(row, flush=True)
        if self.memory:
            parts = []
            for config in self.configs:
                data = self.memory[config]["peakRss"]
                if data:
                    parts.append(f"{config}={data['median'] / (1024 * 1024):.1f}MiB")
            print("peak RSS (median, memory workload): " + ", ".join(parts), flush=True)
        if self.retained:
            parts = []
            for config in self.configs:
                data = self.retained[config]["heapAfterBytes"]
                if data:
                    parts.append(f"{config}={data['median'] / (1024 * 1024):.2f}MiB")
            print("retained graph post-GC heap (median): " + ", ".join(parts), flush=True)
        if self.failures:
            print(f"failures: {len(self.failures)}", flush=True)
            for failure in self.failures[:12]:
                print(f"  {failure['phase']} [{failure['config']}]: {failure['detail']}",
                      flush=True)
            if len(self.failures) > 12:
                print(f"  ... {len(self.failures) - 12} more in SUMMARY.json", flush=True)


def render_text_summary(summary):
    lines = [
        "Gate 5 application benchmark summary",
        f"mode: {summary['mode']}",
        f"status: {summary['status']}",
        f"started: {summary['startedAt']}",
        f"finished: {summary['finishedAt']}",
        f"durationSeconds: {summary['durationSeconds']}",
        f"git revision: {summary['environment']['repository']['revision']}",
        f"git dirty: {summary['environment']['repository']['dirty']}",
        f"go: {summary['environment']['toolchains']['go']}",
        f"java: {summary['environment']['toolchains']['java'].splitlines()[0]}",
        f"class manifest sha256: {summary['workload']['classManifestSha256']}",
        "failures: " + (json.dumps(summary["failures"]) if summary["failures"] else "none"),
        "",
        "cold launch->result (median seconds): " + json.dumps(
            {config: (summary["cold"].get(config, {}).get("elapsedToResult") or {}).get("median")
             for config in ALL_CONFIGS}),
        "",
        "steady aggregate ops/s by kernel:",
    ]
    for kernel in KERNELS:
        values = {}
        for config in ALL_CONFIGS:
            data = summary["steady"].get(config, {}).get("kernels", {}).get(kernel)
            values[config] = data.get("aggregateOpsPerSecond") if data else None
        lines.append(f"  {kernel}: " + json.dumps(values))
    lines.append("")
    lines.append("peak RSS bytes (median): " + json.dumps(
        {config: (summary["memory"].get(config, {}).get("peakRss") or {}).get("median")
         for config in ALL_CONFIGS}))
    lines.append("retained heap after GC bytes (median): " + json.dumps(
        {config: (summary["retained"].get(config, {}).get("heapAfterBytes") or {}).get("median")
         for config in ALL_CONFIGS}))
    lines.append("")
    lines.append("Quantile method: " + summary["quantileMethod"])
    return "\n".join(lines) + "\n"


def ensure_go_on_path():
    if shutil.which("go"):
        return
    candidate = Path("/usr/local/go/bin")
    if (candidate / "go").is_file():
        os.environ["PATH"] = str(candidate) + os.pathsep + os.environ.get("PATH", "")


def execute(args, work):
    suite = Suite(work, env=java_environment(args.java_home))
    check_toolchains(suite)
    configs = args.configs or list(ALL_CONFIGS)
    run = BenchRun(suite, args, work, configs)
    run.started_at = time.strftime("%Y-%m-%dT%H:%M:%S%z")
    started = time.monotonic()

    setup = suite.run("build-catty", ["go", "build", "-o", suite.bin, "./cmd/catty"],
                      timeout=300)
    setup.succeeded()

    run.environment_phase()
    run.build_phase()
    run.correctness_phase()
    run.cold_phase()
    run.steady_phase()
    run.memory_phase()

    run.finished_at = time.strftime("%Y-%m-%dT%H:%M:%S%z")
    run.duration = time.monotonic() - started
    status = "ok" if not run.failures else "failed"
    run.workDirRetained = bool(args.keep or args.smoke or not args.out or status != "ok")
    run.write_raw()
    summary = run.summary_value(status)
    write_json(work / "raw" / "SUMMARY.json", summary)
    if args.out:
        run.publish(args.out)
        write_json(Path(args.out) / "SUMMARY.json", summary)
        (Path(args.out) / "SUMMARY.txt").write_text(render_text_summary(summary),
                                                    encoding="utf-8")
    run.print_summary(status)
    if not run.workDirRetained:
        shutil.rmtree(work, ignore_errors=True)
        print(f"temporary work directory removed; evidence retained in {args.out}",
              flush=True)
    return 0 if status == "ok" else 1


def main():
    parser = argparse.ArgumentParser(
        description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--java-home", type=Path, default=os.environ.get("JAVA_HOME"),
                        help="JDK 21 home (adds its bin directory to PATH)")
    parser.add_argument("--smoke", action="store_true",
                        help="reduced repeats across the full configuration matrix")
    parser.add_argument("--repeat", type=int, default=None,
                        help="fresh processes per configuration for "
                             "cold/steady/memory/retained")
    parser.add_argument("--build-repeats", type=int, default=None,
                        help="Catty builder repetitions per mode (median build time)")
    parser.add_argument("--configs", default=None,
                        help="comma-separated subset of java,interp,mixed,aot")
    parser.add_argument("--out", type=Path, default=None,
                        help="evidence directory (default out/gate5/bench for the full run)")
    parser.add_argument("--keep", action="store_true",
                        help="keep the temporary work directory even when evidence is published")
    args = parser.parse_args()
    if args.smoke and args.repeat is not None and args.repeat > 3:
        parser.error("--smoke with --repeat > 3 is not a smoke run")
    if args.repeat is not None and args.repeat < 1:
        parser.error("--repeat must be positive")
    if args.build_repeats is not None and args.build_repeats < 1:
        parser.error("--build-repeats must be positive")
    if args.configs:
        configs = [item.strip() for item in args.configs.split(",") if item.strip()]
        unknown = [item for item in configs if item not in ALL_CONFIGS]
        if unknown:
            parser.error(f"unknown configurations: {unknown}")
        if not configs:
            parser.error("--configs selected no configuration")
        args.configs = configs
    else:
        args.configs = list(ALL_CONFIGS)
    if args.out is None and not args.smoke:
        args.out = ROOT / "out" / "gate5" / "bench"
    ensure_go_on_path()

    work = Path(tempfile.mkdtemp(prefix="catty-gate5-bench-")).resolve()
    print(f"Gate 5 benchmark artifacts in {work}", flush=True)
    try:
        return execute(args, work)
    except (AssertionError, OSError, ValueError, KeyError) as error:
        print(f"Gate 5 benchmark failed. {error}", file=sys.stderr, flush=True)
        print(f"Artifacts retained in {work}", file=sys.stderr, flush=True)
        return 1


if __name__ == "__main__":
    sys.exit(main())
