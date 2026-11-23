#!/usr/bin/env python3
"""Generate class chains that exercise the hierarchy-depth ceiling.

The archive entry names sort ancestors before descendants, which is the order
that exposed the discarded-remaining-path-length bug.
"""

import sys
from pathlib import Path


def write_chain(root, package, prefix, first, count, base=None, base_package=None):
    package_dir = root / package
    package_dir.mkdir(parents=True, exist_ok=True)
    for i in range(count):
        number = first + i
        name = f"{prefix}{number:03d}"
        if number == 0:
            extends = ""
        elif number == first and base is not None:
            extends = f" extends {base_package}.{base}"
        else:
            extends = f" extends {prefix}{number - 1:03d}"
        source = package_dir / f"{name}.java"
        source.write_text(
            f"package {package};\n\n"
            f"public class {name}{extends} {{\n"
            f"    public static int depth() {{ return {number + 1}; }}\n"
            f"}}\n",
            encoding="utf-8",
        )


def main():
    out = Path(sys.argv[1])
    write_chain(out / "big128", "chain", "C", 0, 128)
    write_chain(out / "big129", "chain", "C", 0, 129)
    write_chain(out / "parent1", "span", "P", 0, 120)
    write_chain(out / "parent2", "span", "P", 120, 8, base="P119", base_package="span")
    write_chain(out / "parent3", "span", "P", 128, 1, base="P127", base_package="span")


if __name__ == "__main__":
    main()
