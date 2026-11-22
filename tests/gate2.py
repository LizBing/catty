#!/usr/bin/env python3
"""Run the Gate 2 differential suite against an OpenJDK 21 installation."""

import argparse
from dataclasses import dataclass
import json
import struct
import os
from pathlib import Path
import re
import shutil
import sys
import tempfile

from acceptance import DATA, MAIN_DESC, MODES, ROOT, Suite, method_id, read_report, require, verify_trace
from classfile import ClassFile, u2


PROFILE = "catty-jvm21-gate2"
IMAGE_ABI = 2
ENABLE_PLUGINS = False
REQUIRED_CASES = (
    "long-crossing", "initialization", "numeric-edge-cases", "switch-and-wide",
    "arrays", "dispatch-and-constructors", "constant-values-and-reentry",
    "reflection", "deferred-linkage", "profile-and-image",
)

LONG_BRIDGE = ("LongCrossing", "bridge", "(JI)J")
LONG_LEAF = ("LongCrossing", "leaf", "(JI)J")
DEFAULT_INIT = ("InitSequence$WithDefault", "<clinit>", "()V")

# javap renders wide-prefixed instructions with a _w suffix; the profile
# inventory records the canonical opcode.
WIDE_ALIASES = {
    "iinc_w": "iinc", "iload_w": "iload", "lload_w": "lload",
    "fload_w": "fload", "dload_w": "dload", "aload_w": "aload",
    "istore_w": "istore", "lstore_w": "lstore", "fstore_w": "fstore",
    "dstore_w": "dstore", "astore_w": "astore",
}


@dataclass(frozen=True)
class Case:
    name: str
    main: str
    opcodes: tuple
    runs: tuple  # (argv, expected stdout) pairs
    interpreted: tuple  # Method ID forced to interpretation in the mixed executable
    ordered_entries: tuple
    forbidden_entries: tuple = ()


def canonical_opcode(name):
    return WIDE_ALIASES.get(name, name)


CASES = (
    Case(
        "long-crossing", "LongCrossing",
        ("ldc2_w", "lxor", "lstore_3", "lload_3", "lreturn", "lushr", "l2i"),
        (((), b"270544960\n1348497537\n"),
         (("a", "b", "c"), b"270544960\n1348497538\n")),
        LONG_BRIDGE,
        (("LongCrossing", "main", MAIN_DESC), LONG_BRIDGE, LONG_LEAF),
    ),
    Case(
        "initialization", "InitSequence", ("ldc", "getstatic", "putstatic", "athrow"),
        (((), b"1\n10\n20\n30\n30\n40\n50\n1\n"),), DEFAULT_INIT,
        (("InitSequence", "main", MAIN_DESC), ("InitSequence$Parent", "<clinit>", "()V"),
         DEFAULT_INIT, ("InitSequence$Child", "<clinit>", "()V"),
         ("InitSequence$Child", "read", "()I"), ("InitSequence$Broken", "<clinit>", "()V")),
        (("InitSequence$WithoutDefault", "<clinit>", "()V"),
         ("InitSequence$Broken", "read", "()I")),
    ),
    Case(
        "numeric-edge-cases", "NumericEdgeCases",
        ("ladd", "lsub", "lmul", "ldiv", "lrem", "lneg", "land", "lor", "lxor", "lshl", "lshr", "lushr", "ishl", "ishr", "iushr", "lcmp", "fcmpl", "fcmpg", "dcmpl", "dcmpg", "frem", "drem", "iand", "ixor", "ifnonnull", "if_acmpeq", "i2b", "i2c", "i2s", "l2f", "l2d", "d2l", "ldc2_w", "ldc_w", "ldc", "dup_x1", "dup_x2", "dup2", "dup2_x1", "dup2_x2", "pop2"),
        (((), b'-1161981756646125697\n-202646972097526600\n-9223372036854775808\n-9223372036854775808\n0\n2147483645\n1515489821\n1515489821\n7001\n-22801902\n-22801902\n-136\n-136\n-1401480802\n-1401480802\n1763500739\n1763500739\n4609434218613702656\n-1072693248\n-1572864\n0\n0\n0\n0\n2147483647\n9223372036854775807\n2147483647\n9223372036854775807\n-2147483648\n-9223372036854775808\n-2147483648\n-9223372036854775808\n2147483647\n3000000000\n-2147483648\n-2147483649\n2147483647\n2147483648\n2147483647\n9223372036854775807\n-2147483648\n-2147483648\n-2147483648\n-9223372036854775808\n-1082130432\n1082130432\n-1066401792\n-1088841318\n-1077936128\n236\n236\n1065353216\n-1082130432\n-2147483648\n-2147483648\n0\n227\n227\n1065353216\n1065353216\n-2147483648\n-8388608\n-1082130432\n236\n236\n-1082130432\n-1082130432\n-2147483648\n-8388608\n1065353216\n227\n227\n32\n32\n160\n160\n96\n96\n218\n218\n236\n236\n218\n218\n-4616189618054758400\n4616189618054758400\n-4607745368753438720\n-4619792497756654797\n-4613937818241073152\n236\n236\n4607182418800017408\n-4616189618054758400\n-9223372036854775808\n-9223372036854775808\n0\n227\n227\n4607182418800017408\n4607182418800017408\n-9223372036854775808\n-4503599627370496\n-4616189618054758400\n236\n236\n9214871658872686752\n9214871658872686752\n4607182418800017407\n9218868437227405312\n-8500377982089056\n236\n236\n32\n32\n160\n160\n96\n96\n218\n218\n236\n236\n218\n218\n255\n2143289345\n9221120237041090561\n-2147483648\n-9223372036854775808\n2\n1\n4\n1638\n2409\n26\n41\n21\n534984\n275\n0\n5\n1069547520\n1\n-4611123068473966592\n1\n31337\n31338\n5\n64\n-1429410030\n65\n-1429410029\n1\n3368603\n64\n1010\n2001\n2002\n3\n'),
         (("a", "b", "c"), b'-1161981756646125697\n-202646972097526600\n-9223372036854775808\n-9223372036854775808\n3\n-2147483645\n1515489827\n1515489827\n7001\n-22801902\n-22801902\n-133\n-133\n-1401480802\n-1401480802\n1763500739\n1763500739\n4609434218613702656\n-1072693248\n-1572864\n0\n0\n0\n0\n2147483647\n9223372036854775807\n2147483647\n9223372036854775807\n-2147483648\n-9223372036854775808\n-2147483648\n-9223372036854775808\n2147483647\n3000000000\n-2147483648\n-2147483649\n2147483647\n2147483648\n2147483647\n9223372036854775807\n-2147483648\n-2147483648\n-2147483648\n-9223372036854775808\n-1082130432\n1082130432\n-1066401792\n-1088841318\n-1077936128\n236\n236\n1065353216\n-1082130432\n-2147483648\n-2147483648\n0\n227\n227\n1065353216\n1065353216\n-2147483648\n-8388608\n-1082130432\n236\n236\n-1082130432\n-1082130432\n-2147483648\n-8388608\n1065353216\n227\n227\n32\n32\n160\n160\n96\n96\n218\n218\n236\n236\n218\n218\n-4616189618054758400\n4616189618054758400\n-4607745368753438720\n-4619792497756654797\n-4613937818241073152\n236\n236\n4607182418800017408\n-4616189618054758400\n-9223372036854775808\n-9223372036854775808\n0\n227\n227\n4607182418800017408\n4607182418800017408\n-9223372036854775808\n-4503599627370496\n-4616189618054758400\n236\n236\n9214871658872686752\n9214871658872686752\n4607182418800017407\n9218868437227405312\n-8500377982089056\n236\n236\n32\n32\n160\n160\n96\n96\n218\n218\n236\n236\n218\n218\n255\n2143289345\n9221120237041090561\n-2147483648\n-9223372036854775808\n2\n1\n4\n1638\n2409\n26\n41\n21\n534984\n563\n0\n5\n1069547520\n1\n-4611123068473966592\n1\n31337\n31338\n5\n64\n-1429410030\n65\n-1429410029\n1\n3368603\n64\n1010\n2001\n2002\n3\n1437226410\n')),
        ("NumericEdgeCases", "leaf", "(I)I"),
        (("NumericEdgeCases", "main", "([Ljava/lang/String;)V"),
         ("NumericEdgeCases", "longOps", "(JI)I"),
         ("NumericEdgeCases", "longShiftCheck", "(JI)I"),
         ("NumericEdgeCases", "intShiftCheck", "(II)I")),
        (("NumericEdgeCases", "unusedPath", "(I)I"),),
    ),
    Case(
        "switch-and-wide", "SwitchWide",
        ("tableswitch", "lookupswitch", "iinc_w", "iload_w", "istore_w", "lload_w", "lstore_w", "goto", "iadd", "if_icmple", "if_icmpge", "ifle", "ifne", "ldc2_w", "ldc_w", "bipush", "sipush", "iconst_m1", "invokestatic", "invokevirtual", "ireturn", "lreturn", "ladd", "istore", "lstore", "iconst_0"),
        (((), b'41\n81\n-1\n206\n-10\n11\n15\n66\n31\n34\n41\n44\n51\n54\n7\n6\n4\n8\n101\n102\n100\n4\n2\n7\n6\n5\n9\n2138\n2138\n25681\n25681\n'),),
        ('SwitchWide', 'dense', '(I)I'),
        (('SwitchWide', 'main', '([Ljava/lang/String;)V'), ('SwitchWide', 'sparse', '(I)I'), ('SwitchWide', 'nested', '(II)I'), ('SwitchWide', 'loopSwitch', '(I)I')),
        (('SwitchWide', 'neverSwitch', '(I)I')),
    ),
    Case(
        "arrays", "Arrays",
        ("newarray", "multianewarray", "arraylength", "iaload", "laload", "faload", "daload", "baload", "caload", "saload", "iastore", "lastore", "fastore", "dastore", "bastore", "castore", "sastore", "aaload", "aastore", "anewarray", "checkcast", "instanceof", "invokeinterface", "invokespecial"),
        (((), b'10\n3\n0\n2\n0\n2\n0\n2\n0\n2\n0\n2\n0\n2\n0\n2\n0\n3\n0\n20\n33\n1\n5\n1\n65\n90\n7\n-3\n1000\n-1000\n1\n2\n1\n2\n30\n2\n1\n41\n42\n2\n1\n31\n1\n32\n1\n1\n1\n1\n1\n2\n1\n1\n1\n1\n1\n33\n50\n2\n40\n2\n3\n3\n1\n0\n9\n2\n1\n1\n3\n1\n2\n3\n4\n0\n7\n1\n1\n1\n2\n3\n0\n1\n5\n2\n1\n0\n50\n51\n52\n53\n54\n55\n56\n57\n58\n59\n1\n60\n61\n62\n1\n63\n64\n65\n66\n67\n68\n91\n69\n80\n55\n55\n7\n1\n3\n8\n1\n1\n'), (('x',), b'10\n3\n0\n2\n0\n2\n0\n2\n0\n2\n0\n2\n0\n2\n0\n2\n0\n3\n0\n20\n33\n1\n5\n1\n65\n90\n7\n-3\n1000\n-1000\n1\n2\n1\n2\n30\n2\n1\n41\n42\n2\n1\n31\n1\n32\n1\n1\n1\n1\n1\n2\n1\n1\n1\n1\n1\n33\n50\n2\n40\n2\n3\n3\n1\n0\n9\n2\n1\n1\n3\n1\n2\n3\n4\n0\n7\n1\n1\n1\n2\n3\n0\n1\n5\n2\n1\n0\n50\n51\n52\n53\n54\n55\n56\n57\n58\n59\n1\n60\n61\n62\n1\n63\n64\n65\n66\n67\n68\n91\n69\n80\n55\n55\n7\n1\n3\n8\n1\n1\n85\n2\n1\n1\n81\n1\n1\n82\n3\n')),
        ('Arrays', 'lengthOf', '([I)I'),
        (('Arrays', 'main', '([Ljava/lang/String;)V'), ('Arrays', 'defaults', '()V'), ('Arrays', 'primitives', '()V'), ('Arrays', 'references', '()V')),
        (('Arrays$Sized', 'size', '()I')),
    ),
    Case(
        "dispatch-and-constructors", "Constructors",
        ("invokevirtual", "invokeinterface", "invokespecial", "invokestatic", "new", "getfield", "putfield", "instanceof", "checkcast", "anewarray", "aaload", "aastore", "aconst_null", "bipush", "sipush"),
        (((), b'100\n62\n4\n62\n63\n63\n63\n61\n200\n51\n52\n52\n53\n52\n52\n52\n52\n300\n9\n18\n16\n50\n400\n10\n11\n12\n13\n22\n44\n23\n40\n41\n42\n43\n40\n30\n9\n22\n11\n45\n40\n21\n500\n22\n11\n22\n11\n21\n21\n600\n1\n16\n16\n1\n70\n1\n1\n1\n1\n700\n8\n9\n10\n7\n'),),
        ('Constructors', 'note', '(I)I'),
        (('Constructors', 'main', '([Ljava/lang/String;)V'), ('Constructors', 'dispatch', '()V'), ('Constructors', 'defaults', '()V'), ('Constructors', 'abstractShapes', '()V')),
        (('Constructors$Shape', 'area', '()I')),
    ),
    Case(
        "constant-values-and-reentry", "InitOrder",
        ("ldc", "ldc2_w", "getstatic", "putstatic", "invokestatic", "invokevirtual", "new", "athrow", "l2i", "lushr", "idiv", "iadd"),
        (((), b'1\n90\n91\n100\n101\n200\n201\n300\n301\n400\n401\n10\n20\n30\n30\n40\n50\n180\n60\n61\n660\n700\n64\n65\n70\n71\n70\n1\n70\n0\n45\n46\n293709995\n501\n502\n1\n602\n801\n802\n'),),
        ('InitOrder', 'mark', '(I)I'),
        (('InitOrder', 'main', '([Ljava/lang/String;)V'), ('InitOrder$Base', '<clinit>', '()V'), ('InitOrder$Derived', '<clinit>', '()V'), ('InitOrder$Derived', 'readDerived', '()I'), ('InitOrder$Base', 'readBase', '()I')),
        (('InitOrder$NoDefault', '<clinit>', '()V'), ('InitOrder$ChildFail', '<clinit>', '()V')),
    ),
    Case(
        "reflection", "Reflect",
        ("new", "dup", "invokespecial", "invokestatic", "invokevirtual", "checkcast", "instanceof", "getstatic", "ldc", "anewarray", "aastore", "astore", "aload_0", "iadd", "ireturn", "goto"),
        (((), b'1\n11\n2\n22\n12\n3\n1\n123\n1\n-1\n4\n5\n30\n14\n14\n6\n1\n105\n106\n165\n207\n308\n301\n409\n310\n401\n566\n1\n7\n8\n9\n10\n11\n12\n13\n1\n31\n1\n32\n1\n33\n1\n34\n1\n1\n1\n67\n1\n35\n1\n36\n1\n14\n15\n16\n41\n17\n18\n19\n20\n21\n1\n1\n1\n22\n23\n24\n2\n3\n'),),
        ('Reflect', 'leaf', '(I)I'),
        (('Reflect', 'main', '([Ljava/lang/String;)V'), ('Base', 'inherited', '()I'), ('Child', 'overridden', '()I'), ('Base', 'overridden', '()I'), ('Child', 'value', '()Ljava/lang/Integer;')),
        (('Base', 'packagePrivate', '()I'), ('Base', 'unused', '()I')),
    ),
    Case(
        "profile-and-image", "ProfileImage",
        ("ldc2_w", "lushr", "l2i", "lstore_2", "lload_2", "iadd", "isub", "imul", "idiv", "invokestatic", "invokevirtual", "getstatic", "if_icmple", "new"),
        (((), b'287454020\n301767\n905330\n226331\n8\n'),),
        ('ProfileImage', 'mix', '(I)I'),
        (('ProfileImage', 'main', '([Ljava/lang/String;)V'), ('ProfileImage', 'helper', '(I)I'), ('ProfileImage', 'scale', '(I)I')),
        (),
    ),
    Case(
        "deferred-linkage", "Link",
        ("invokestatic", "getstatic", "invokevirtual", "invokespecial", "iconst_1", "ireturn", "return", "iadd", "pop", "istore_0", "astore_1", "goto"),
        (((), b'1\n2\n3\n4\n5\n'),),
        ('Link', 'callGood', '()I'),
        (('Link', 'main', '([Ljava/lang/String;)V'), ('Link', 'startup', '()I'), ('Link', 'readValue', '()I'), ('Link', 'finish', '()I')),
        (('Link', 'unusedCall', '()V'), ('Link', 'unusedField', '()V')),
    ),
)


CASE_MAINS = {case.name: case.main for case in CASES}


def java_environment(home):
    if home is None:
        return {}
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
        first_line = version.partition("\n")[0]
        require(re.search(pattern, version),
                f"Gate 2 requires OpenJDK 21 tools. {tool} reported {first_line!r}. "
                f"Use --java-home to select a JDK 21 installation. See {result.log}")
        print(version, flush=True)


def case_paths(suite, case):
    if case.name == "deferred-linkage":
        return [suite.classes(case.name + "-link"), suite.classes(case.name + "-variant")]
    return [suite.classes(case.name)]


def prepare_case(suite, case):
    classes = suite.classes(case.name)
    classes.mkdir(parents=True)
    case_dir = DATA / "gate2" / case.name
    # The deferred-linkage variant source declares the same class as the main
    # dependency and is compiled separately.
    sources = sorted(p for p in case_dir.rglob("*.java") if "variant" not in p.parts)
    require(sources, f"no Java sources for {case.name}")
    suite.run(f"javac-{case.name}", ["javac", "--release", "21", "-g", "-d", classes, *sources]).succeeded()
    if case.name == "deferred-linkage":
        variant = suite.classes(case.name + "-variant")
        variant.mkdir(parents=True)
        suite.run(f"javac-{case.name}-variant",
                  ["javac", "--release", "21", "-g", "-d", variant,
                   case_dir / "variant" / "DepChanged.java"]).succeeded()
        link = suite.classes(case.name + "-link")
        shutil.copytree(classes, link)
        (link / "Dep.class").unlink()
        classpath = case_paths(suite, case)
    else:
        classpath = [classes]
    disassembly = suite.disassemble(case.name, classes)
    opcodes = set(re.findall(r"^\s+\d+:\s+([a-z]\w*)\b", disassembly, re.MULTILINE))
    missing = set(case.opcodes) - opcodes
    require(not missing, f"{case.name} missing bytecodes {sorted(missing)}")
    results = []
    for index, (args, expected) in enumerate(case.runs):
        result = suite.run(f"oracle-{case.name}-{index}",
                           ["java", "-Xverify:all", "-cp", os.pathsep.join(map(str, classpath)),
                            case.main, *args]).succeeded()
        require((result.stdout, result.stderr) == (expected, b""),
                f"OpenJDK {case.name} {args}: expected {expected!r} and no stderr. See {result.log}")
        results.append((args, result))
    if case.name == "long-crossing":
        prepare_invalid_long(suite, classes, disassembly)
    return results


def prepare_invalid_long(suite, classes, disassembly):
    invalid = suite.classes("bad-long-slot")
    shutil.copytree(classes, invalid)
    path = invalid / "LongCrossing.class"
    parsed = ClassFile(path)
    code = parsed.methods[LONG_BRIDGE[1:]].code
    require("StackMapTable" in code.attributes, "LongCrossing.bridge must contain a stack map")
    _, declaration, body = disassembly.partition("static long bridge(long, int);")
    require(declaration, "javap did not find LongCrossing.bridge")
    bridge = body.partition("public static void main(")[0]
    match = re.search(r"^\s+(\d+):\s+lload_3\b", bridge, re.MULTILINE)
    require(match, "javap did not find the saved long load in LongCrossing.bridge")
    offset = code.start + int(match[1])
    require(parsed.data[offset] == 0x21, "expected lload_3 at the selected offset")
    parsed.data[offset] = 0x1F  # lload_1 reads the second slot of the incoming long.
    path.write_bytes(parsed.data)
    suite.run("oracle-bad-long-slot", ["java", "-Xverify:all", "-cp", invalid, "LongCrossing"]).failed(
        r"VerifyError", r"LongCrossing\.bridge", r"lload_1")


def check_run(suite, case, label, binary, args, oracle, choices, *, isolated=False):
    result, trace = suite.traced_run(label, binary, args, isolated=isolated)
    actual = (result.status, result.stdout, result.stderr)
    expected = (oracle.status, oracle.stdout, oracle.stderr)
    require(actual == expected, f"expected {expected!r}, got {actual!r}. See {result.log}")
    entries = [event for event in verify_trace(trace, result) if event["event"] == "enter"]
    observed = []
    for event in entries:
        method = method_id(event)
        require(method in choices and choices[method]["engine"] == event["engine"],
                f"engine trace disagrees with the report for {method}. See {trace}")
        require(method not in case.forbidden_entries, f"unexpected entry into {method}. See {trace}")
        if method in case.ordered_entries:
            observed.append(method)
    require(tuple(observed) == case.ordered_entries,
            f"expected ordered entries {case.ordered_entries}, got {observed}. See {trace}")
    if case.name == "long-crossing":
        depths = [event["depth"] for event in entries if method_id(event) in case.ordered_entries]
        require(depths == [0, 1, 2], f"expected nested long calls, got depths {depths}. See {trace}")


def check_report_inventory(report, case):
    expected = (("catty-jvm21-gate3", 3) if ENABLE_PLUGINS else (PROFILE, IMAGE_ABI))
    require((report["profile"], report["imageABI"]) == expected,
            f"expected profile {expected[0]} and image ABI {expected[1]}")
    if ENABLE_PLUGINS:
        require(report.get("plugins") is True, f"{case.name}: Gate 3 build did not advertise plugins")
    inventory = set(report["opcodes"])
    missing = {canonical_opcode(op) for op in case.opcodes} - inventory
    require(not missing, f"{case.name}: report inventory misses opcodes {sorted(missing)}")
    require(isinstance(report["bootstrap"]["methods"], list) and report["bootstrap"]["methods"],
            "missing bootstrap method inventory")


def check_deferred_report(report):
    unresolved = {(entry["kind"], entry["owner"], entry["name"], entry["descriptor"], entry.get("pc"))
                  for entry in report.get("unresolved", [])}
    for kind, name, descriptor in (("method", "good", "()I"), ("field", "value", "I")):
        require(any(entry[0] == kind and entry[1] == "Dep" and entry[2] == name and entry[3] == descriptor
                    and isinstance(entry[4], int) and entry[4] >= 0 for entry in unresolved),
                f"deferred-linkage report did not retain Dep.{name}{descriptor} with a bytecode offset: {sorted(unresolved)}")


def check_mode(suite, case, mode, oracles):
    paths = case_paths(suite, case)
    label = f"{case.name}-{mode}"
    mixed = mode == "auto"
    overrides = {case.interpreted: "interp"} if mixed else {}
    policy = suite.policy(label, list(overrides.items())) if mixed else None
    build = suite.build(label, paths, case.main, mode, policy=policy, generated=mixed,
                        enable_plugins=ENABLE_PLUGINS)
    report, choices = read_report(build, case.main, paths, default="interp" if mode == "interp" else "aot",
                                  overrides=overrides)
    check_report_inventory(report, case)
    if case.name == "deferred-linkage":
        check_deferred_report(report)
    for index, (args, oracle) in enumerate(oracles):
        check_run(suite, case, f"{label}-{index}", build.binary, args, oracle, choices)
    if mixed and case.name == "long-crossing":
        standalone = suite.work / "standalone"
        standalone.mkdir()
        binary = standalone / build.binary.name
        shutil.copy2(build.binary, binary)
        offline = paths[0].with_name(f"{case.name}-offline")
        paths[0].rename(offline)
        try:
            args, oracle = oracles[0]
            check_run(suite, case, "standalone", binary, args, oracle, choices, isolated=True)
            original = [json.loads(line) for line in
                        (suite.work / "traces" / f"{label}-0.jsonl").read_text().splitlines()]
            isolated_trace = [json.loads(line) for line in
                              (suite.work / "traces" / "standalone.jsonl").read_text().splitlines()]
            strip_pid = lambda events: [{key: value for key, value in event.items() if key != "pid"}
                                        for event in events]
            require(strip_pid(original) == strip_pid(isolated_trace),
                    "the isolated executable produced a different engine trace")
        finally:
            offline.rename(paths[0])
    return build


def generated_function(source, symbol):
    marker = f"func {symbol}("
    start = source.index(marker)
    depth = 0
    for index in range(start, len(source)):
        if source[index] == "{":
            depth += 1
        elif source[index] == "}":
            depth -= 1
            if depth == 0:
                return source[start:index + 1]
    raise AssertionError(f"unterminated generated function {symbol}")


def generated_symbol(report, owner, name, descriptor):
    for record in report["methods"]:
        if (record["owner"], record["name"], record["descriptor"]) == (owner, name, descriptor):
            generated = record.get("generated") or {}
            require(generated.get("symbol"), f"{owner}.{name}{descriptor} was not lowered to Go")
            return generated["symbol"]
    raise AssertionError(f"report misses {owner}.{name}{descriptor}")


def check_generated_go(generated_builds):
    required = {
        "numeric-edge-cases": ("NumericEdgeCases", "remainders", "(JIFD)I"),
        "switch-and-wide": ("SwitchWide", "offsetDense", "(I)I"),
        "initialization": ("InitSequence", "main", "([Ljava/lang/String;)V"),
    }
    cases_by_name = {case.name: case for case in CASES}
    for case_name, method in required.items():
        build, paths = generated_builds[case_name]
        case = cases_by_name[case_name]
        report, _ = read_report(build, CASE_MAINS[case_name], paths,
                                overrides={case.interpreted: "interp"})
        symbol = generated_symbol(report, *method)
        source = (build.generated / "generated.go").read_text(encoding="utf-8")
        require("catty: interpreter" not in source and "rt.Execute" not in source,
                f"{case_name}: generated Go contains an interpreter handoff")
        body = generated_function(source, symbol)
        if case_name == "numeric-edge-cases":
            require("rt.FRem(" in body and "rt.DRem(" in body,
                    "numeric generated Go is missing lowered float/double remainder")
            require("rt.LongValue(rt.AsLong(" in body and " % " in body,
                    "numeric generated Go is missing lowered long arithmetic/remainder")
        elif case_name == "switch-and-wide":
            require("switch rt.AsInt(stack[" in body and "case " in body and "default:" in body
                    and "block =" in body,
                    "switch generated Go does not contain real switch control flow")
        else:
            class_file = ClassFile(paths[0] / "InitSequence.class")
            code = class_file.methods[("main", "([Ljava/lang/String;)V")].code
            handlers = [u2(class_file.data, offset + 4) for offset in code.handlers]
            require(handlers, "initialization fixture lost its exception handlers")
            for handler in handlers:
                require(f"case {handler}:" in body,
                        f"exception generated Go misses original handler offset {handler}")
            ranges = [(u2(class_file.data, offset), u2(class_file.data, offset + 2))
                      for offset in code.handlers]
            transfer_pcs = [int(value) for value in re.findall(r"t\.HandleThrow\(m, (\d+)", body)]
            require(transfer_pcs, "exception generated Go has no throwing-instruction transfer")
            protected = [pc for pc in transfer_pcs if any(start <= pc < end for start, end in ranges)]
            require(protected,
                    f"exception generated Go has no transfer from a protected bytecode offset: {transfer_pcs}")
    print("PASS generated-go checks", flush=True)


def malformed_copy(suite, source, name):
    target = suite.classes(name)
    shutil.copytree(suite.classes(source), target)
    return target


def method_section(disassembly, header):
    start = disassembly.index(header)
    return disassembly[start:].split("\n\n", 1)[0]


def check_malformed_bytecode(suite):
    if suite.classes("numeric-edge-cases").is_dir():
        classes = malformed_copy(suite, "numeric-edge-cases", "bad-stack-shape")
        path = classes / "NumericEdgeCases.class"
        parsed = ClassFile(path)
        code = parsed.methods[("intArith", "(II)I")].code
        section = method_section(suite.disassemble("bad-stack-shape", classes), "static int intArith(int, int);")
        match = re.search(r"^\s+(\d+):\s+iadd\b", section, re.MULTILINE)
        require(match, "javap did not find intArith.iadd")
        pc = int(match[1])
        require(parsed.data[code.start + pc] == 0x60, "expected iadd at the selected offset")
        parsed.data[code.start + pc] = 0x61  # ladd consumes two longs from an int stack.
        path.write_bytes(parsed.data)
        suite.build("bad-stack-shape", [classes], "NumericEdgeCases",
                    expected_errors=(r"NumericEdgeCases", r"intArith", r"ladd|long|stack", r"\bpc\b|offset"))
        print("PASS bad-stack-shape", flush=True)
    if suite.classes("switch-and-wide").is_dir():
        classes = malformed_copy(suite, "switch-and-wide", "bad-control-flow")
        path = classes / "SwitchWide.class"
        parsed = ClassFile(path)
        code = parsed.methods[("dense", "(I)I")].code
        section = method_section(suite.disassemble("bad-control-flow", classes), "static int dense(int);")
        pcs = {int(value) for value in re.findall(r"^\s+(\d+):", section, re.MULTILINE)}
        match = re.search(r"^\s+(\d+):\s+tableswitch\b", section, re.MULTILINE)
        require(match, "javap did not find the dense tableswitch")
        pc = int(match[1])
        require(parsed.data[code.start + pc] == 0xAA, "expected tableswitch at the selected offset")
        pad = (4 - ((pc + 1) % 4)) % 4
        default_pos = code.start + pc + 1 + pad
        first_target = struct.unpack_from(">i", parsed.data, default_pos + 12)[0]
        broken = first_target + 1
        while broken in pcs:
            broken += 1
        require(broken < code.size, "no non-instruction target for the malformed switch")
        struct.pack_into(">i", parsed.data, default_pos, broken)
        path.write_bytes(parsed.data)
        suite.build("bad-control-flow", [classes], "SwitchWide",
                    expected_errors=(r"SwitchWide", r"branch target|instruction boundary|stack map",
                                     r"\bpc\b|offset|target"))
        print("PASS bad-control-flow", flush=True)


def check_invalid_long(suite):
    suite.build("bad-long-slot", [suite.classes("bad-long-slot")], "LongCrossing",
                expected_errors=(r"LongCrossing", r"bridge\s*\(JI\)J", r"lload|local|category.?two|slot",
                                 r"\b(?:pc|offset|bci)\b"))


def generated_image(generated):
    source = generated / "main.go"
    require(source.is_file(), f"missing generated main.go in {generated}")
    text = source.read_text(encoding="utf-8")
    match = re.search(r'const cattyImageJSON = (".*")\n', text)
    require(match, "generated main.go does not contain the image literal")
    return text, json.loads(json.loads(match.group(1)))


def write_image(source, text, image):
    quoted = json.dumps(json.dumps(image, separators=(",", ":")))
    updated = re.sub(r'const cattyImageJSON = (".*")\n',
                     lambda _match: f"const cattyImageJSON = {quoted}\n", text)
    source.write_text(updated, encoding="utf-8")


def image_variant(suite, build, label, mutate):
    folder = suite.work / "image-variants" / label
    shutil.copytree(build.generated, folder)
    source = folder / "main.go"
    text, image = generated_image(folder)
    mutate(image)
    write_image(source, text, image)
    binary = folder / "app"
    suite.run(f"image-build-{label}", ["go", "build", "-trimpath", "-buildvcs=false", "-o", binary, "."],
              cwd=folder, env={"GOWORK": "off"}, timeout=240).succeeded()
    return folder, binary


def check_profile_image(suite):
    paths = [suite.classes("profile-and-image")]
    build = suite.build("profile-image-inventory", paths, "ProfileImage", "auto", generated=True)
    report, _ = read_report(build, "ProfileImage", paths)
    check_report_inventory(report, next(case for case in CASES if case.name == "profile-and-image"))
    text, image = generated_image(build.generated)
    require(image["profile"] == report["profile"] == PROFILE,
            f"image and report disagree on the profile: {image['profile']} / {report['profile']}")
    require(int(image["abi"]) == report["imageABI"] == IMAGE_ABI,
            f"image and report disagree on the image ABI: {image['abi']} / {report['imageABI']}")
    require(image.get("opcodes") == report["opcodes"], "image and report opcode inventories disagree")
    for opcode in ("ldc2_w", "tableswitch", "lookupswitch", "wide", "multianewarray", "invokeinterface",
                   "ifnonnull", "l2f", "i2b", "fcmpl", "newarray", "athrow"):
        require(opcode in set(report["opcodes"]), f"profile inventory misses opcode {opcode}")

    methods = {(entry["owner"], entry["name"], entry["descriptor"])
               for entry in report["bootstrap"]["methods"]}
    for member in (
        ("java/lang/Object", "getClass", "()Ljava/lang/Class;"),
        ("java/lang/Class", "getMethod", "(Ljava/lang/String;[Ljava/lang/Class;)Ljava/lang/reflect/Method;"),
        ("java/lang/reflect/Method", "invoke", "(Ljava/lang/Object;[Ljava/lang/Object;)Ljava/lang/Object;"),
        ("java/lang/Integer", "valueOf", "(I)Ljava/lang/Integer;"),
        ("java/lang/Integer", "intValue", "()I"),
        ("java/lang/Long", "valueOf", "(J)Ljava/lang/Long;"),
        ("java/lang/Long", "longValue", "()J"),
        ("java/lang/Float", "intBitsToFloat", "(I)F"),
        ("java/lang/Float", "floatToRawIntBits", "(F)I"),
        ("java/lang/Double", "longBitsToDouble", "(J)D"),
        ("java/lang/Double", "doubleToRawLongBits", "(D)J"),
    ):
        require(member in methods, f"bootstrap inventory misses {member}")
    fields = {(entry["owner"], entry["name"], entry["descriptor"])
              for entry in report["bootstrap"]["fields"]}
    for member in (("java/lang/System", "out", "Ljava/io/PrintStream;"),
                   *((f"java/lang/{box}", "TYPE", "Ljava/lang/Class;") for box in
                     ("Boolean", "Byte", "Character", "Short", "Integer", "Long", "Float", "Double"))):
        require(member in fields, f"bootstrap inventory misses {member}")

    for label, mutate, pattern in (
        ("wrong-abi", lambda image: image.__setitem__("abi", IMAGE_ABI - 1), r"abi"),
        ("wrong-profile", lambda image: image.__setitem__("profile", "catty-jvm21-gate1"), r"profile"),
    ):
        folder, binary = image_variant(suite, build, label, mutate)
        trace = suite.work / "traces" / f"image-{label}.jsonl"
        result = suite.run(f"image-{label}", [binary], cwd=folder, env={"CATTY_TRACE": str(trace)})
        require(result.status == 2, f"{label}: expected a rejected image (status 2), got {result.status}. See {result.log}")
        require(re.search(pattern, result.stderr.decode("utf-8", errors="replace"), re.IGNORECASE),
                f"{label}: stderr does not name the rejected {pattern}. See {result.log}")
        require(result.stdout == b"", f"{label}: rejected image produced application output. See {result.log}")
        if trace.is_file():
            events = [json.loads(line) for line in trace.read_text(encoding="utf-8").splitlines() if line]
            require(not any(event.get("event") == "enter" for event in events),
                    f"{label}: rejected image entered a Java method before validation")

    missing = suite.work / "missing-bootstrap"
    sources = missing / "sources"
    classes = missing / "classes"
    sources.mkdir(parents=True)
    classes.mkdir()
    source = sources / "MissingBootstrap.java"
    source.write_text("public final class MissingBootstrap {\n"
                      "    public static void main(String[] args) {\n"
                      "        System.out.println(System.nanoTime());\n"
                      "    }\n"
                      "}\n", encoding="utf-8")
    suite.run("javac-missing-bootstrap",
              ["javac", "--release", "21", "-g", "-d", classes, source]).succeeded()
    prior = build.binary
    original = prior.read_bytes()
    suite.build("missing-bootstrap", [classes], "MissingBootstrap", output=prior,
                expected_errors=(r"java/lang/System", r"nanoTime", r"bootstrap"))
    require(prior.read_bytes() == original, "a failed rebuild replaced the previous executable")
    require(not (suite.work / "builds" / "missing-bootstrap" / "report.json").exists(),
            "a failed missing-bootstrap build released a report")
    print("PASS profile-and-image checks", flush=True)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--java-home", type=Path, default=os.environ.get("JAVA_HOME"),
                        help="JDK 21 home directory (defaults to JAVA_HOME or tools on PATH)")
    parser.add_argument("--prepare-only", action="store_true",
                        help="compile fixtures and check OpenJDK results without building Catty")
    parser.add_argument("--case", choices=[case.name for case in CASES],
                        help="run one fixture without the full gate checks")
    parser.add_argument("--enable-plugins", action="store_true",
                        help="run the supported positive corpus with the Gate 3 profile enabled")
    args = parser.parse_args()
    global ENABLE_PLUGINS
    ENABLE_PLUGINS = args.enable_plugins
    work = Path(tempfile.mkdtemp(prefix="catty-gate2-")).resolve()
    print(f"Gate 2 artifacts in {work}", flush=True)
    cases = [case for case in CASES if args.case is None or case.name == args.case]
    try:
        suite = Suite(work, env=java_environment(args.java_home))
        check_toolchains(suite)
        prepared = [(case, prepare_case(suite, case)) for case in cases]
        count = sum(len(oracles) for _, oracles in prepared)
        negative = " and one verifier rejection" if any(case.name == "long-crossing" for case in cases) else ""
        print(f"Verified {count} OpenJDK 21 runs{negative}.", flush=True)
        if args.prepare_only:
            return 0
        suite.run("go-build-catty", ["go", "build", "-o", suite.bin, "./cmd/catty"], timeout=240).succeeded()
    except AssertionError as exc:
        print(exc, file=sys.stderr)
        return 1

    failures = []
    generated_builds = {}
    for case, oracles in prepared:
        for mode in MODES:
            try:
                build = check_mode(suite, case, mode, oracles)
                if mode == "auto":
                    generated_builds[case.name] = (build, case_paths(suite, case))
                print(f"PASS {case.name}/{mode}", flush=True)
            except AssertionError as exc:
                failures.append(f"{case.name}/{mode}: {exc}")
    if args.case is None and not failures:
        try:
            check_generated_go(generated_builds)
        except AssertionError as exc:
            failures.append(f"generated-go: {exc}")
    if any(case.name == "long-crossing" for case in cases) and not ENABLE_PLUGINS:
        try:
            check_invalid_long(suite)
            print("PASS bad-long-slot", flush=True)
        except AssertionError as exc:
            failures.append(f"bad-long-slot: {exc}")
    if any(case.name in {"numeric-edge-cases", "switch-and-wide"} for case in cases) and not ENABLE_PLUGINS:
        try:
            check_malformed_bytecode(suite)
        except AssertionError as exc:
            failures.append(f"malformed-bytecode: {exc}")
    if any(case.name == "profile-and-image" for case in cases) and not ENABLE_PLUGINS:
        try:
            check_profile_image(suite)
        except AssertionError as exc:
            failures.append(f"profile-and-image checks: {exc}")
    if args.case is None:
        missing = [name for name in REQUIRED_CASES if name not in {case.name for case in CASES}]
        if missing:
            failures.append("Missing required corpus cases: " + ", ".join(missing))
        try:
            if ENABLE_PLUGINS:
                raise AssertionError("skipped for Gate 3-enabled positive corpus")
            result = suite.run("gate1-regression", [sys.executable, "-B", ROOT / "tests" / "acceptance.py"],
                               timeout=900).succeeded()
            print(f"PASS Gate 1 regression. See {result.log}", flush=True)
        except AssertionError as exc:
            if not (ENABLE_PLUGINS and str(exc) == "skipped for Gate 3-enabled positive corpus"):
                failures.append(f"Gate 1 regression: {exc}")
    for failure in failures:
        print(f"FAIL {failure}", file=sys.stderr)
    if failures:
        return 1
    if args.case is None:
        print("Gate 2 positive corpus passed with Gate 3 enabled." if ENABLE_PLUGINS else "Gate 2 passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
