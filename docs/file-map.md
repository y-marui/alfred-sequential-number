# File Map

> File-level dependency map for alfred-sequential-number.
> Update this as the codebase evolves.

## Entry Points

| File | Role |
|---|---|
| `cmd/sequential-number-alfred/main.go` | Alfred executes this binary — the sole entry point |

## Call Flow

```
cmd/sequential-number-alfred/main.go
  └─ dispatch(query)                     [recovers panics into an error item]
       └─ internal/seqcmd.Dispatch(query)
            ├─ dispatchAllFormats(query)  [leading-space "preview every format" mode]
            │    └─ internal/seq.Generate(pattern, narr)  (once per format)
            ├─ "fmt <range> [format]"     → internal/seq.Generate(format, narr)
            ├─ subcommand (bin/oct/hex/Hex/alf/Alf) → internal/seq.Generate(format, narr)
            └─ default (decimal)         → internal/seq.Generate("%d", narr)
```

## Package Dependency Table

| Package | Imports from | Notes |
|---|---|---|
| `internal/scriptfilter` | stdlib only | Alfred Script Filter JSON types (`Item`, `Response`) and `Write` |
| `internal/seq` | stdlib only | Pattern parsing, bijective-base-26 alphabet conversion, `Generate` |
| `internal/seqcmd` | `internal/seq`, `internal/scriptfilter` | Query dispatch, format/subcommand routing, preview/error item rendering |
| `cmd/sequential-number-alfred` | `internal/seqcmd`, `internal/scriptfilter` | Reads `os.Args[1]`, recovers panics, writes JSON to stdout |

## Tests

| File | Tests |
|---|---|
| `internal/seq/seq_test.go` | Alphabet conversion, `LoadRange`, `Generate` (decimal/binary/`#`/`#a`/`#A`/custom format), the multi-range "extra args ignored" quirk, invalid pattern/empty-range errors |
| `internal/seqcmd/seqcmd_test.go` | Command dispatch (decimal, bin, oct, hex, Hex, alf, Alf, fmt), leading-space preview-all-formats mode, "No results"/"Enter a length or range" placeholder items |
