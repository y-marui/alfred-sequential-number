# Alfred Sequential Number

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-sequential-number/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-sequential-number/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-sequential-number/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-sequential-number/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

An Alfred workflow to generate sequential numbers and paste them to the clipboard.

## Requirements

- Alfred 5 (Powerpack required)
- Python 3.9+

## Installation

Download the latest `.alfredworkflow` from [GitHub Releases](https://github.com/y-marui/alfred-sequential-number/releases) and double-click to install.

## Usage

### `seq`
```
seq <length or range>

seq 10
>>> 1 2 3 ... 8 9 10
seq 3-10
>>> 3 4 5 ... 8 9 10
```
Generate sequential numbers in decimal. Same as `seq fmt %d`.

### `seq bin`
```
seq bin <length or range>

seq bin 10
>>> 1 10 11 ... 1000 1001 1010
```
Generate sequential numbers in binary.

### `seq oct`
```
seq oct <length or range>

seq oct 8
>>> 1 2 3 ... 6 7 10
```
Generate sequential numbers in octal.

### `seq hex`
```
seq hex <length or range>

seq hex 16
>>> 1 2 3 ... e f 10
```
Generate sequential numbers in hexadecimal (lower case).

### `seq Hex`
```
seq Hex <length or range>

seq Hex 8
>>> 1 2 3 ... E F 10
```
Generate sequential numbers in hexadecimal (upper case).

### `seq alf`
```
seq alf <length or range>

seq alf 27
>>> a b c ... y z aa
```
Generate sequential numbers alphabetically (lower case).

### `seq Alf`
```
seq Alf <length or range>

seq Alf 27
>>> A B C ... Y Z AA
```
Generate sequential numbers alphabetically (upper case).

### `seq fmt`
```
seq fmt <format> <length or range> [<length or range> ...]

seq fmt Sample-#a-# 3 2
>>> Sample-a-1 Sample-a-2 Sample-b-1 Sample-b-2 Sample-c-1 Sample-c-2
```
Generate multi-dimensional sequential values with a custom format string.

#### Format specifiers

| Specifier | Meaning |
|---|---|
| `%b` | binary |
| `%o` | octal |
| `%d`, `#` | decimal |
| `%x` | hexadecimal lower case |
| `%X` | hexadecimal upper case |
| `%a`, `#a` | alphabetic lower case: a, b, c, ... |
| `%A`, `#A` | alphabetic upper case: A, B, C, ... |

## Development

```bash
make install    # install dev dependencies
make test       # run tests
make lint       # ruff check
make build      # build dist/*.alfredworkflow
```

## License

MIT — see [LICENSE](LICENSE)

---

*This is the reference (English) version. The canonical Japanese version is [README-jp.md](README-jp.md). Update both files in the same commit.*
