# UI Design

Alfred Script Filter workflows present results as a list of items in the Alfred
launcher. This document defines the UI conventions for result items in this workflow.

## Result Item Structure

Alfred result items are JSON objects with the following fields used in this workflow:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | Primary text (large, always visible) — a preview of the generated sequence |
| `subtitle` | string | no | Secondary text (small, below title) — the format description |
| `arg` | string | no | The full sequence, `\r\n`-joined; copied to clipboard on Enter |
| `valid` | bool | yes | If false, Enter does not trigger an action |

## Text Guidelines

### No Unicode Emoji in `title` / `subtitle`

- **Prohibited:** `🔍 Search`, `✅ Done`, `📄 Document`
- **Allowed:** ASCII symbols — `>`, `*`, `[x]`, `(!)`, `--`
- **Reason:** Emoji rendering is inconsistent across Alfred versions and macOS
  updates. ASCII symbols are universally stable.

### Preview Truncation (`title`)

- 5 or fewer values: joined with `", "` (e.g. `1, 2, 3, 4, 5`).
- More than 5 values: first 3 plus the last, joined with `", ..., "`
  (e.g. `1, 2, 3, ..., 10`).

### Empty / Error States

- `seq fmt` with no range → a placeholder item, `title: "Enter a length or range"`,
  `valid: false`, to guide the user.
- No results (invalid pattern/range) → `title: "No results"`, `subtitle: <format
  description>`, `valid: false`.
- Unexpected internal error → `cmd/sequential-number-alfred`'s panic recovery
  automatically shows a `"Workflow Error"` item; do not hide errors silently.

## Icon

- Workflow icon: `workflow/icon.png` (PNG, any size — Alfred scales it).
- Alfred controls light/dark mode; do not ship separate light/dark icons.
- No per-item icons are used in this workflow.

## Keyboard Shortcuts

These are standard Alfred behaviors — do not override them in the workflow:

| Key | Action |
|---|---|
| ↩ Enter | Copy `arg` to clipboard and paste |
| ⌘C | Copy `arg` to clipboard |
| ⌘L | Show `title` in Large Type |

## Layout Conventions

### Normal result (decimal/bin/oct/hex/Hex/alf/Alf/fmt)

```
title:    <preview, e.g. "1, 2, 3, ..., 10">
subtitle: <format description, e.g. "Paste Sequential Numbers (decimal)">
arg:      <full sequence, \r\n-joined>
valid:    true
```

### Leading-space "preview every format" mode

One item per registered format, each labeled with its subcommand keyword:

```
title:    <preview for that format>
subtitle: <subcommand>: <format description>   (e.g. "bin: Paste Sequential Numbers (binary)")
arg:      <full sequence for that format>
valid:    true
```

### Placeholder / error items

```
title:    "Enter a length or range" | "No results" | "Workflow Error"
subtitle: <format description, or the error message>
valid:    false
```
