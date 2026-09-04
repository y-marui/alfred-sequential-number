# Architecture

## Overview

An Alfred Workflow (Go): `cmd/sequential-number-alfred` is the single universal
(amd64+arm64) binary `workflow/info.plist` invokes. Its Script Filter node passes the
query following the `seq` keyword as `$1`; the binary parses it, generates a sequence
via `internal/seq`, and prints Alfred Script Filter JSON via `internal/scriptfilter`.
Selecting the result copies the sequence to the clipboard and pastes it through
Alfred's own native Clipboard Output node — no script is involved in that step; see
[docs/specification.md](specification.md) for the full data flow.
`scripts/build-workflow.sh` packages the binary with `workflow/info.plist` and
`workflow/icon.png` into a `.alfredworkflow`.

This structure — a thin `cmd/` entry point over independently testable `internal/`
packages, no generic command-router abstraction, Script Filter JSON via a small
`scriptfilter` package — deliberately matches
[y-marui/alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text),
[y-marui/alfred-markdown-ref](https://github.com/y-marui/alfred-markdown-ref), and
[y-marui/alfred-password-generator](https://github.com/y-marui/alfred-password-generator),
this author's other Alfred Workflows already implemented in Go. This workflow itself
was originally a Python implementation (`src/alfred`/`src/app`); see `CHANGELOG.md`'s
`[Unreleased]` entry for what changed and why in that rewrite.

## Entry Points

- `cmd/sequential-number-alfred` — a single command, no subcommands. The query it
  receives (e.g. `"10"`, `"bin 10"`, `"fmt 3 item-%d"`, `" 12"`) determines behavior; see
  [docs/specification.md](specification.md#commands).

One Alfred trigger reaches it: the `seq` keyword, wired in `workflow/info.plist`.

## Directory Structure

| Directory | Role |
|---|---|
| `cmd/sequential-number-alfred/` | The binary Alfred invokes; recovers panics into a Script Filter error item and writes the response |
| `internal/seqcmd/` | Query dispatch, argument parsing, preview/error item rendering — the Alfred result row(s) |
| `internal/seq/` | Pattern expansion and sequence generation, unit tested independently of Alfred |
| `internal/scriptfilter/` | Alfred Script Filter JSON response types |
| `workflow/` | `info.plist` (the Alfred object graph), `icon.png` |
| `scripts/build-workflow.sh` | Builds the universal binary and packages `workflow/` into `dist/*.alfredworkflow` |
| `scripts/extract-changelog.sh` | Extracts one version's notes from `CHANGELOG.md` for GitHub Releases |
| `docs/` | Specification, notes |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`) |

## Key Dependencies

None. `internal/seq` and `internal/seqcmd` use only the Go standard library
(`strconv`, `strings`).

## Alfred Configuration Builder (`userconfigurationconfig`)

Alfred 5 の Configuration Builder は `info.plist` の `userconfigurationconfig` キーで定義する。
利用可能な全型・各キーの詳細は [`docs/configuration-builder.md`](configuration-builder.md) を参照。

This workflow declares no Config Builder variables — `userconfigurationconfig` is an
empty array. (The Python implementation had a `use_uv` checkbox that selected between
`uv run python` and `python3`; a compiled binary needs no such interpreter toggle, so
it was removed in the Go rewrite.)

### Passing variables

Alfred はスクリプト実行時に各 `variable` を環境変数として渡す。
インストール直後は `prefs.plist` が存在しないため変数は未セットになる場合がある。
将来スクリプト側で変数を読む場合は常にデフォルト値を持たせること。

~~~go
// Go
value := os.Getenv("my_variable")
if value == "" {
	value = "fallback"
}
~~~

**注意:** `checkbox` 型の unchecked 値は `"0"` ではなく空文字 `""` になる。
`[ "$var" = "1" ]` で判定し、`"0"` との比較は避けること。

### Relationship between `variables` / `prefs.plist` / `default`

| 場所 | 役割 |
|---|---|
| `userconfigurationconfig[].config.default` | Configuration Builder UI の初期表示のみ。変数への書き込みは行わない。 |
| `prefs.plist`（同ディレクトリ） | ユーザーが Configuration Builder で保存した値。Alfred が自動生成・更新する。 |
| `info.plist` の `variables` | スクリプトに常に渡したい固定の環境変数。Configuration Builder で管理する変数はここに入れない。 |
