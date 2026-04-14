"""seq command - generate sequential numbers.

Wraps workflow/seq.py logic for use in tests and development.
Alfred invokes seq.py directly via inline scripts in info.plist;
this module exists for programmatic access and testing.

Usage in Alfred keyword triggers:
    seq <length or range>         -- decimal
    seq bin <length or range>     -- binary
    seq oct <length or range>     -- octal
    seq hex <length or range>     -- hex lower
    seq Hex <length or range>     -- hex upper
    seq alf <length or range>     -- alphabetic lower
    seq Alf <length or range>     -- alphabetic upper
    seq fmt <format> <range...>   -- custom format
"""

from __future__ import annotations

import sys
from pathlib import Path

# workflow/seq.py is not on sys.path by default outside of Alfred; add it.
_workflow_root = Path(__file__).resolve().parents[3] / "workflow"
if str(_workflow_root) not in sys.path:
    sys.path.insert(0, str(_workflow_root))

from seq import loader  # noqa: E402  # type: ignore[import-not-found]

_FORMAT_MAP = {
    "bin": "%b",
    "oct": "%o",
    "hex": "%x",
    "Hex": "%X",
    "alf": "%a",
    "Alf": "%A",
}


def generate(subcommand: str, args: list[str]) -> list[str]:
    """Return a list of sequential strings.

    Args:
        subcommand: One of 'bin', 'oct', 'hex', 'Hex', 'alf', 'Alf', 'fmt', or '' for decimal.
        args: Range/length tokens (and format string for 'fmt').

    Returns:
        List of generated strings.
    """
    if subcommand == "fmt":
        if not args:
            return []
        narr = args[:1]
        fmt = args[1] if len(args) > 1 else "item-#"
        return list(loader(fmt, narr))

    fmt = _FORMAT_MAP.get(subcommand, "%d")
    return list(loader(fmt, args[:1]))
