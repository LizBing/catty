#!/usr/bin/env python3
"""Generate minimal class-file fixtures for exception-handler frame verification.

The classes are hand-assembled so the StackMapTable can declare the handler's
local types independently of the actual throwing path.  OpenJDK 21 rejects the
bad fixture with VerifyError; Catty admission must reject it with VERIFY.
"""
import io
import struct
import zipfile
from pathlib import Path

HERE = Path(__file__).resolve().parent


def u1(v): return struct.pack(">B", v)
def u2(v): return struct.pack(">H", v)
def u4(v): return struct.pack(">I", v)


def utf8(s):
    b = s.encode("utf-8")
    return u1(1) + u2(len(b)) + b


def cls(i):
    return u1(7) + u2(i)


def name_and_type(n, d):
    return u1(12) + u2(n) + u2(d)


def methodref(c, nt):
    return u1(10) + u2(c) + u2(nt)


def cpool(entries):
    return u2(len(entries) + 1) + b"".join(entries)


def class_file(name, local_kind, stack_kind=9):
    # local_kind: 4 = java/lang/Object (valid), 16 = java/lang/String (invalid)
    entries = [
        utf8(name),                # 1
        cls(1),                    # 2
        utf8("java/lang/Object"),  # 3
        cls(3),                    # 4
        utf8("f"),                 # 5
        utf8("(I)I"),              # 6
        utf8("Code"),              # 7
        utf8("java/lang/RuntimeException"),  # 8
        cls(8),                    # 9
        utf8("StackMapTable"),     # 10
        utf8("java/lang/String"),  # 11
        utf8("<init>"),            # 12
        utf8("()V"),               # 13
        name_and_type(12, 13),     # 14
        methodref(9, 14),          # 15
        cls(11),                   # 16, String class for the bad frame
        methodref(4, 14),          # 17, Object.<init>
    ]
    code = (
        u1(0xBB) + u2(4) +       # new java/lang/Object
        u1(0x59) +               # dup
        u1(0xB7) + u2(17) +      # invokespecial Object.<init>
        u1(0x4C) +               # astore_1: local 1 is a real Object
        u1(0xBB) + u2(9) +       # new RuntimeException
        u1(0x59) +               # dup
        u1(0xB7) + u2(15) +      # invokespecial RuntimeException.<init>
        u1(0xBF) +               # athrow
        u1(0x4C) +               # astore_1 (caught throwable)
        u1(0x05) +               # iconst_2
        u1(0xAC)                 # ireturn
    )
    # Handler frame at pc 16: locals [int, local_kind], stack [RuntimeException].
    stack_entry = u1(1) if stack_kind == 1 else u1(7) + u2(stack_kind)
    smap = (
        u2(1) +
        u1(255) + u2(16) +
        u2(2) + u1(1) + u1(7) + u2(local_kind) +
        u2(1) + stack_entry
    )
    code_attr = (
        u2(2) + u2(2) + u4(len(code)) + code +
        u2(1) + u2(8) + u2(16) + u2(16) + u2(9) +
        u2(1) + u2(10) + u4(len(smap)) + smap
    )
    method = (
        u2(0x0009) + u2(5) + u2(6) + u2(1) +
        u2(7) + u4(len(code_attr)) + code_attr
    )
    return (
        b"\xca\xfe\xba\xbe" + u2(0) + u2(65) +
        cpool(entries) +
        u2(0x0021) + u2(2) + u2(4) +
        u2(0) + u2(0) +
        u2(1) + method +
        u2(0)
    )


def write_jar(path, class_name, local_kind, stack_kind=9):
    with zipfile.ZipFile(path, "w", zipfile.ZIP_DEFLATED) as z:
        z.writestr(class_name + ".class", class_file(class_name, local_kind, stack_kind))


if __name__ == "__main__":
    write_jar(HERE / "handler-good.jar", "good/Handler", 4)
    write_jar(HERE / "handler-bad.jar", "bad/Handler", 16)
    write_jar(HERE / "handler-bad-stack.jar", "badstack/Handler", 4, 1)
