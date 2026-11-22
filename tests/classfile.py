"""Locate fields to patch in the javac-produced acceptance fixtures."""

import struct
from dataclasses import dataclass


def u2(data, offset):
    return struct.unpack_from(">H", data, offset)[0]


def u4(data, offset):
    return struct.unpack_from(">I", data, offset)[0]


@dataclass(frozen=True)
class Code:
    start: int
    size: int
    handlers: tuple
    attributes: dict


@dataclass(frozen=True)
class Method:
    flags: int
    code: object


class ClassFile:
    def __init__(self, path):
        self.data = bytearray(path.read_bytes())
        assert self.data[:8] == bytes.fromhex("cafebabe00000041"), path
        self.offset = 8
        self.pool = self.constant_pool()
        self.offset += 2  # access_flags
        self.name = self.pool[self.pool[self.read_u2()]].decode("ascii")
        self.offset += 2  # super_class
        interface_count = self.read_u2()
        self.offset += 2 * interface_count
        for _ in range(self.read_u2()):
            self.member()
        self.methods = {}
        for _ in range(self.read_u2()):
            flags, name, descriptor, attributes = self.member()
            code = self.read_code(attributes["Code"]) if "Code" in attributes else None
            self.methods[(name, descriptor)] = Method(flags, code)

    def read_u2(self):
        value = u2(self.data, self.offset)
        self.offset += 2
        return value

    def constant_pool(self):
        sizes = {3: 4, 4: 4, 5: 8, 6: 8, 7: 2, 8: 2, 9: 4, 10: 4,
                 11: 4, 12: 4, 15: 3, 16: 2, 17: 4, 18: 4, 19: 2, 20: 2}
        count = self.read_u2()
        pool = [None]
        while len(pool) < count:
            tag = self.data[self.offset]
            self.offset += 1
            if tag == 1:
                length = self.read_u2()
                pool.append(bytes(self.data[self.offset:self.offset + length]))
                self.offset += length
                continue
            pool.append(u2(self.data, self.offset) if tag == 7 else None)
            self.offset += sizes[tag]
            if tag in (5, 6):
                pool.append(None)
        return pool

    def attributes(self, offset):
        result = {}
        count = u2(self.data, offset)
        offset += 2
        for _ in range(count):
            name = self.pool[u2(self.data, offset)].decode("ascii")
            length = u4(self.data, offset + 2)
            result[name] = offset + 6
            offset += 6 + length
        return result, offset

    def member(self):
        flags = self.read_u2()
        name = self.pool[self.read_u2()].decode("ascii")
        descriptor = self.pool[self.read_u2()].decode("ascii")
        attributes, self.offset = self.attributes(self.offset)
        return flags, name, descriptor, attributes

    def read_code(self, offset):
        size = u4(self.data, offset + 4)
        start = offset + 8
        exception_count = u2(self.data, start + size)
        table_start = start + size + 2
        attributes, _ = self.attributes(table_start + 8 * exception_count)
        return Code(start, size, tuple(range(table_start, table_start + 8 * exception_count, 8)),
                    attributes)
