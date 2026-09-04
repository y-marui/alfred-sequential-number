# Developing

This document covers the development workflow, conventions, and guidelines for contributors to this project.

## Development Setup

```bash
git clone https://github.com/y-marui/alfred-sequential-number
cd alfred-sequential-number
go build ./...
```

**Prerequisites:**
- macOS (required for Alfred)
- Go (see `go.mod` for the toolchain version)
- Alfred 5 with Powerpack
- `jq` (optional, for pretty-printed dev output): `brew install jq`
- `gh` CLI (required for releases): `brew install gh`

## Development Workflow

### Daily commands

```bash
go run ./cmd/sequential-number-alfred "10"          # simulate Alfred locally
go run ./cmd/sequential-number-alfred "bin 3-10"
go run ./cmd/sequential-number-alfred "fmt 3 item-%d"
go run ./cmd/sequential-number-alfred ""
make test                 # run test suite
make lint                 # gofmt -l + go vet
make fmt                  # gofmt -w (auto-fix)
make build                # go build ./...
make build-workflow       # build dist/*.alfredworkflow
```

Pipe through `jq` for pretty-printed JSON:

```bash
go run ./cmd/sequential-number-alfred "bin 10" | jq .
```

### Testing in Alfred

1. `make build-workflow` — generates `dist/*.alfredworkflow`
2. Double-click the `.alfredworkflow` file to install in Alfred
3. Open Alfred and type `seq` to verify behavior

During rapid iteration you can symlink `workflow/` to Alfred's workflow directory,
but `go run ./cmd/sequential-number-alfred "query"` is usually faster for logic changes.

## Adding a New Subcommand

1. Register the format in `allFormats` in `internal/seqcmd/seqcmd.go`:

```go
var allFormats = []formatEntry{
	// ...
	{"%z", "Paste Sequential Numbers (my format)", "myfmt"},
}
```

2. Add the corresponding verb to `internal/seq/seq.go`'s `parsePattern`/`render` if it's a
   new format specifier (not just a new subcommand alias for an existing one).
3. Add tests in `internal/seq/seq_test.go` and `internal/seqcmd/seqcmd_test.go`.
4. Update `README.md`/`README-jp.md`'s Usage section and `workflow/info.plist`'s Script
   Filter `subtext`.

## Naming Conventions

| Scope | Convention | Example |
|---|---|---|
| Go files | one file per package concern | `seq.go`, `seqcmd.go` |
| Go packages | short, lowercase, no underscores | `seq`, `seqcmd`, `scriptfilter` |
| Exported functions / types | `PascalCase` | `Generate`, `Response`, `Item` |
| Unexported functions / variables | `camelCase` | `parsePattern`, `defaultLength` |
| Alfred subcommand names | lowercase/mixed-case matching output case | `"bin"`, `"fmt"`, `"Hex"`, `"Alf"` |
| Commit messages | Conventional Commits | `feat:`, `fix:`, `docs:`, `chore:` |
| Branch names | `feat/`, `fix/`, `docs/`, `chore/` | `feat/add-base32-format` |

## Code Style

- **Formatter:** `gofmt`. CI enforces this (`make lint`).
- **Linter:** `go vet`.
- **Comments:** Write *why*, not *what*. Do not comment self-evident code.
- **No debug prints:** Remove all stray `fmt.Print*` statements before committing;
  the only writer to stdout is `scriptfilter.Response.Write`.
- **No third-party dependencies** unless clearly justified — keep `go.mod` dependency-free.

## Commit Guidelines

- Commit per **feature unit**, after confirming it works.
- **No WIP commits** — do not commit code that does not run.
- **No `--no-verify`** — never skip pre-commit hooks.

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add base32 format specifier
fix: off-by-one in alphabetic range wraparound
chore: update Go toolchain to 1.28
docs: update README usage section
refactor: simplify seqcmd dispatch logic
```

## Pull Request Checklist

- [ ] `make lint` passes
- [ ] `make test` passes
- [ ] `make build-workflow` succeeds
- [ ] New commands have tests
- [ ] `README.md` updated if user-facing changes
- [ ] `CHANGELOG.md` entry added under `[Unreleased]`

## Code Review Guidelines

**Reviewers check for:**
- Architectural constraints respected (no business logic in `cmd/sequential-number-alfred`)
- No hardcoded absolute paths (use `$HOME` / env vars)
- No debug prints in production code
- No Unicode emoji in Alfred result item `title` / `subtitle`
- Tests cover the new or changed behavior
- Alfred env variables managed via Config Builder, not `variables` key

**Security-sensitive changes** (auth, encryption, data access) require explicit
security review before merge.

**Self-review:** Individual contributors open a PR and self-review before merging
to `main`.
