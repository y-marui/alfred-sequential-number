# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Breaking (implementation):** Reimplemented the workflow in Go
  (`cmd/sequential-number-alfred` + `internal/seq`, `internal/seqcmd`,
  `internal/scriptfilter`), replacing the Python `src/alfred`/`src/app` implementation.
  The `seq` keyword, bundle ID, and the behavior of `seq`/`bin`/`oct`/`hex`/`Hex`/`alf`/
  `Alf`/`fmt` (including the leading-space "preview every format" mode and the
  multi-range Cartesian-product behavior) are unchanged; results are byte-for-byte
  equivalent JSON.
- Alfred now invokes a compiled universal (amd64+arm64) binary directly instead of
  `python3`/`uv run python`; the `Use uv` Config Builder toggle is removed.
- Build/test tooling moved from `uv`/`ruff`/`mypy`/`pytest` to `go build`/`gofmt`/`go
  vet`/`go test`.

### Fixed

- `README.md`/`README-jp.md`'s `seq fmt` usage was documented as
  `seq fmt <format> <length or range>` with a multi-dimensional example that never
  matched the actual (and tested) argument order; corrected to
  `seq fmt <length or range> [<format>]`, matching runtime behavior since the original
  Python implementation.

## [1.0.0] - 2026-04-14

### Added

- Sequential number generation in decimal, binary, octal, hex (upper/lower), and
  alphabetic (upper/lower) formats via the `seq`/`bin`/`oct`/`hex`/`Hex`/`alf`/`Alf`
  subcommands
- Custom format strings via `seq fmt <length or range> [<format>]`
  (`%b`/`%o`/`%d`/`%x`/`%X`/`%a`/`%A` and `#`/`#a`/`#A` specifiers)
- Leading-space mode to preview every format for one length/range at once
- Multi-range Cartesian-product generation (e.g. `seq bin 3 5`)
- Copies the generated sequence to the clipboard and pastes it on selection
- Alfred SDK: `response`, `cache`, `config`, `logger`, `router`, `safe_run`
- Vendor packaging via `scripts/vendor.sh`
- Build pipeline via `scripts/build.sh`
- GitHub Actions CI (lint, test, build)
- GitHub Actions Release (tag → `.alfredworkflow` → GitHub Release)
- Full pytest test suite

Note: this repository has not yet had a tagged GitHub Release; `1.0.0` here matches the
version already recorded in `workflow/info.plist` at the time of the Go rewrite above.
