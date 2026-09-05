# Alfred Sequential Number

> **これは日本語版（正本）です。**
> 英語版（参照）は [README.md](README.md) を参照してください。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-sequential-number/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-sequential-number/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-sequential-number/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-sequential-number/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

Alfred で連番を生成してクリップボードに貼り付けるワークフロー。

## Requirements

- Alfred 5（Powerpack が必要）
- Go（ビルド時のみ、`go.mod` を参照）

## Installation

[GitHub Releases](https://github.com/y-marui/alfred-sequential-number/releases) から最新の `.alfredworkflow` をダウンロードしてダブルクリックでインストールします。

## Usage

`seq` キーワードで連番を生成します。

### `seq`
```
seq <長さ または 範囲>

seq 10
>>> 1 2 3 ... 8 9 10
seq 3-10
>>> 3 4 5 ... 8 9 10
```
10進数の連番を生成します。`seq fmt %d` と同じ動作です。

### `seq bin`
```
seq bin <長さ または 範囲>

seq bin 10
>>> 1 10 11 ... 1000 1001 1010
```
2進数の連番を生成します。

### `seq oct`
```
seq oct <長さ または 範囲>

seq oct 8
>>> 1 2 3 ... 6 7 10
```
8進数の連番を生成します。

### `seq hex`
```
seq hex <長さ または 範囲>

seq hex 16
>>> 1 2 3 ... e f 10
```
16進数（小文字）の連番を生成します。

### `seq Hex`
```
seq Hex <長さ または 範囲>

seq Hex 8
>>> 1 2 3 ... E F 10
```
16進数（大文字）の連番を生成します。

### `seq alf`
```
seq alf <長さ または 範囲>

seq alf 27
>>> a b c ... y z aa
```
アルファベット（小文字）の連番を生成します。

### `seq Alf`
```
seq Alf <長さ または 範囲>

seq Alf 27
>>> A B C ... Y Z AA
```
アルファベット（大文字）の連番を生成します。

### `seq fmt`
```
seq fmt <長さ または 範囲> [<フォーマット>]

seq fmt 3
>>> item-1 item-2 item-3
seq fmt 3 Sample-#a
>>> Sample-a Sample-b Sample-c
```
カスタムフォーマット文字列で連番を生成します（デフォルトフォーマット: `item-#`）。

#### Format specifiers

| 指定子 | 意味 |
|---|---|
| `%b` | 2進数 |
| `%o` | 8進数 |
| `%d`, `#` | 10進数 |
| `%x` | 16進数（小文字） |
| `%X` | 16進数（大文字） |
| `%a`, `#a` | アルファベット（小文字）: a, b, c, ... |
| `%A`, `#A` | アルファベット（大文字）: A, B, C, ... |

## Development

```bash
make test            # テスト実行
make lint             # gofmt -l + go vet
make build-workflow   # dist/*.alfredworkflow を生成
```

## Project Structure

```
alfred-sequential-number/
├── cmd/
│   └── sequential-number-alfred/  # Alfred が実行するバイナリ
├── internal/
│   ├── seq/             # 連番生成ロジック（コア）
│   ├── seqcmd/           # クエリディスパッチ・フォーマット/サブコマンドルーティング
│   └── scriptfilter/    # Alfred Script Filter JSON 型
└── workflow/             # Alfred パッケージ（info.plist, icon.png）
```

## License

MIT — [LICENSE](LICENSE) を参照

---

*この文書には英語版（参照版）[README.md](README.md) があります。編集時は同一コミットで更新してください。*
