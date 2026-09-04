# Specification

> Functional specification, behavior definition, and data flow for alfred-sequential-number.

## Overview

This workflow is an Alfred 5 Script Filter that accepts a keyword + query, generates a
sequence of decimal/binary/octal/hex/alphabetic/custom-format values, and returns a JSON
result to Alfred. Selecting the result copies the sequence to the clipboard (one value per
line, `\r\n`-separated) and pastes it.

## Commands

### `seq` (default, decimal)

**Trigger:** `seq <length or range>`

**Behavior:**
1. If no range is given, defaults to length `12`.
2. Generate the sequence via `internal/seq.Generate("%d", narr)`.
3. Return one result item (`arg` = the full sequence, `\r\n`-joined).

Same output as `seq fmt <length or range> %d`.

---

### `seq bin` / `seq oct` / `seq hex` / `seq Hex` / `seq alf` / `seq Alf`

**Trigger:** `seq <subcommand> <length or range>`

**Behavior:** Same as the default, but with the format fixed to that subcommand's pattern:

| Subcommand | Pattern | Output example (`seq X 4`) |
|---|---|---|
| `bin` | `%b` | 1, 10, 11, 100 |
| `oct` | `%o` | 1, 2, 3, 4 |
| `hex` | `%x` | 1, 2, 3, 4 |
| `Hex` | `%X` | 1, 2, 3, 4 |
| `alf` | `%a` | a, b, c, d |
| `Alf` | `%A` | A, B, C, D |

If more than one range argument follows the subcommand (e.g. `seq bin 3 5`), every range
contributes a dimension to the Cartesian product (see Multi-Range Behavior below), even
though the format string only has one placeholder.

---

### `seq fmt`

**Trigger:** `seq fmt <length or range> [<format>]`

**Behavior:**
1. `<length or range>` is required; if omitted, returns an "Enter a length or range"
   placeholder item (`valid: false`).
2. `<format>` is optional, default `item-#`.
3. Generate via `internal/seq.Generate(format, [range])` — a single dimension only (unlike
   the subcommand paths above, `fmt` only ever takes its first range argument).

Format specifiers usable in `<format>`:

| Specifier | Meaning |
|---|---|
| `%d`, `#` | decimal |
| `%b` | binary |
| `%o` | octal |
| `%x` | hex lower |
| `%X` | hex upper |
| `%a`, `#a` | alphabetic lower |
| `%A`, `#A` | alphabetic upper |

---

### Leading-space mode ("double-space after `seq`")

**Trigger:** a query starting with a space, e.g. typing `seq` then a space then another
space (`" 12"`, `" "`).

**Behavior:** Preview every format (`ALL_FORMATS`: decimal/bin/oct/hex/Hex/alf/Alf/fmt) for
one length/range (default `12` if the query is only whitespace) as separate result items.

---

## Multi-Range Behavior

`internal/seq.Generate` computes the Cartesian product of every range passed in `narr`
(one dimension per range, outermost = first range, innermost = last), then renders each
combination against the pattern's format fields in order. If a pattern has fewer fields
than `narr` has ranges, the extra trailing values are accepted but unused — this makes the
outer dimension's value repeat once per inner-dimension value (e.g. `seq bin 3 5` produces
`1,1,1,1,1,10,10,10,10,10,11,11,11,11,11` — 3 values × 5 repeats). This matches the
original Python implementation's `str.format()`, which silently ignores unused positional
arguments; it is preserved for behavior parity rather than treated as a bug to fix.

## Data Flow

```
Alfred input (keyword + query string)
  │
  ▼
cmd/sequential-number-alfred/main.go   reads os.Args[1]
  │
  ▼
dispatch(query)                        recovers any panic → error item
  │
  ▼
internal/seqcmd.Dispatch(query)        leading-space → preview all formats;
  │                                    otherwise routes fmt/subcommand/default
  ▼
internal/seq.Generate(pattern, narr)   Cartesian product of ranges → formatted values
  │
  ▼
scriptfilter.Response.Write(os.Stdout) prints JSON to stdout → Alfred renders result(s)
  │
  ▼
User selects item → arg (\r\n-joined sequence) → Alfred's native Clipboard output node
  → copied to clipboard and pasted
```

## Error Handling

- Any panic in `dispatch()` is recovered in `cmd/sequential-number-alfred/main.go`, which
  emits a single "Workflow Error" result item.
- Invalid pattern (unsupported `%` verb) or invalid/empty range → generation returns an
  error, surfaced as a `{"title": "No results", "subtitle": <desc>, "valid": false}` item.
- `seq fmt` with no range argument → an "Enter a length or range" placeholder item.

## Configuration Variables

None. `userconfigurationconfig` is an empty array (see `docs/configuration-builder.md`).

## Constraints

- Script Filter response time target: **< 100 ms** (compiled Go binary, no I/O).
- All output must go through `scriptfilter.Response.Write()` — never `fmt.Print*` directly.
- `cmd/sequential-number-alfred` contains no business logic; it only reads `os.Args[1]`,
  recovers panics, and writes the response `internal/seqcmd.Dispatch` returns.
