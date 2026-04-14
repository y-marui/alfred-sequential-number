import json
import sys

from seq import loader

DEFAULT_LENGTH = "12"
DEFAULT_FMT = "item-#"

# (fmt, desc, subcommand label for double-space listing)
ALL_FORMATS = [
    ("%d", "Paste Sequential Numbers (decimal)", ""),
    ("%b", "Paste Sequential Numbers (binary)", "bin"),
    ("%o", "Paste Sequential Numbers (octal)", "oct"),
    ("%x", "Paste Sequential Numbers (hex lower)", "hex"),
    ("%X", "Paste Sequential Numbers (hex upper)", "Hex"),
    ("%a", "Paste Sequential Numbers (alpha lower)", "alf"),
    ("%A", "Paste Sequential Numbers (alpha upper)", "Alf"),
    (DEFAULT_FMT, "Paste Sequential Numbers (custom format)", "fmt"),
]

SUBCMD_TO_FMT = {label: fmt for fmt, _, label in ALL_FORMATS if label}


def make_preview(seq: list[str]) -> str:
    if len(seq) <= 5:
        return ", ".join(seq)
    return ", ".join(seq[:3]) + ", ..., " + seq[-1]


def make_item(fmt: str, narr: list[str], desc: str) -> dict | None:
    try:
        seq = list(loader(fmt, narr))
    except Exception:
        return None
    if not seq:
        return None
    return {
        "title": make_preview(seq),
        "subtitle": desc,
        "arg": "\r\n".join(seq),
        "valid": True,
    }


def output(items: list) -> None:
    sys.stdout.write(json.dumps({"items": items}))


raw_query = sys.argv[1] if len(sys.argv) > 1 else ""

# Double-space after "seq": query starts with a space → show all format variants
if raw_query.startswith(" "):
    narr = [raw_query.strip() or DEFAULT_LENGTH]
    items = []
    for fmt, desc, label in ALL_FORMATS:
        display_desc = f"{label}: {desc}" if label else desc
        item = make_item(fmt, narr, display_desc)
        if item:
            items.append(item)
    output(items if items else [{"title": "No results", "valid": False}])
    sys.exit(0)

# Alfred passes the whole query as a single $1 (scriptargtype=1); split on spaces.
args = raw_query.split() if raw_query else []

if not args:
    fmt, narr, desc = "%d", [DEFAULT_LENGTH], "Paste Sequential Numbers (decimal)"
elif args[0] == "fmt":
    # seq fmt <range> [format]  (format is optional, default: item-#)
    rest = args[1:]
    if not rest:
        output(
            [
                {
                    "title": "Enter a length or range",
                    "subtitle": "Paste Sequential Numbers (custom format)",
                    "valid": False,
                }
            ]
        )
        sys.exit(0)
    narr = [rest[0]]
    fmt = rest[1] if len(rest) > 1 else DEFAULT_FMT
    desc = "Paste Sequential Numbers (custom format)"
elif args[0] in SUBCMD_TO_FMT:
    fmt = SUBCMD_TO_FMT[args[0]]
    narr = args[1:] if args[1:] else [DEFAULT_LENGTH]
    desc = next(d for _, d, lbl in ALL_FORMATS if lbl == args[0])
else:
    fmt, narr, desc = "%d", args, "Paste Sequential Numbers (decimal)"

try:
    seq = list(loader(fmt, narr))
except Exception:
    seq = []

if not seq:
    output([{"title": "No results", "subtitle": desc, "valid": False}])
    sys.exit(0)

output(
    [
        {
            "title": make_preview(seq),
            "subtitle": desc,
            "arg": "\r\n".join(seq),
            "valid": True,
        }
    ]
)
