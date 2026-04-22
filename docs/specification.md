# Specification

## Overview

Alfred Sequential Number は連続した数列を生成してクリップボードにペーストする Alfred 5 ワークフロー。
decimal / binary / octal / hex / alphabetic / カスタムフォーマットの 8 種類の出力形式に対応。

## Commands

### デフォルト（decimal）

```
seq <length or range>
seq <N>        → 1 〜 N の整数列
seq <S>-<E>   → S 〜 E の整数列
```

### サブコマンド

| サブコマンド | フォーマット指定子 | 出力例（seq X 4） |
|---|---|---|
| `bin` | `%b` | 1, 10, 11, 100 |
| `oct` | `%o` | 1, 2, 3, 4 |
| `hex` | `%x` | 1, 2, 3, 4 |
| `Hex` | `%X` | 1, 2, 3, 4 |
| `alf` | `%a` | a, b, c, d |
| `Alf` | `%A` | A, B, C, D |
| `fmt` | 任意 | （カスタム） |

### ダブルスペースモード

`seq ` の後にスペースを入力すると全フォーマットのプレビューを同時表示する。

### カスタムフォーマット（`seq fmt`）

```
seq fmt <range> [<format>]
```

フォーマット文字列には以下の指定子が使用できる:

| 指定子 | 意味 |
|---|---|
| `%d`, `#` | decimal |
| `%b` | binary |
| `%o` | octal |
| `%x` | hex lower |
| `%X` | hex upper |
| `%a`, `#a` | alphabetic lower |
| `%A`, `#A` | alphabetic upper |

多次元の直積を生成する場合は range を複数指定する:

```
seq fmt Sample-#a-# 3 2
→ Sample-a-1 Sample-a-2 Sample-b-1 Sample-b-2 Sample-c-1 Sample-c-2
```

## Data Flow

```
ユーザー入力（Alfred キーワード）
    ↓
workflow/filter.py  ← Alfred Script Filter（入力中にリアルタイム更新）
    ↓
seq.loader(fmt, narr)  ← パターン解析・数列生成
    ↓
Alfred Script Filter JSON（{"items": [...]})
    ↓
Alfred UI（プレビュー表示）
    ↓ ↩ Enter
workflow/run.py  ← Alfred Run Script（選択確定時）
    ↓
seq.loader(fmt, narr)
    ↓
stdout（\r\n 区切りの数列文字列）→ クリップボードペースト
```

## Input Format

### Range

| 入力形式 | 解釈 |
|---|---|
| `N`（正整数） | 1 〜 N |
| `S-E`（ハイフン区切り） | S 〜 E |

### Pattern Parsing（seq.loader）

1. `%X` 形式（X は指定子文字）→ `{:X}` に変換
2. `#a`, `#A` 形式 → `{:a}`, `{:A}` に変換
3. `#` 単体 → `{:d}` に変換
4. その他の文字はそのまま出力に使用

## Output Format

- Alfred Script Filter JSON: `{"items": [{"title": "...", "subtitle": "...", "arg": "...", "valid": true}]}`
- `title`: プレビュー（最大5要素、それ以上は "1, 2, 3, ..., N" 形式に省略）
- `arg`: `\r\n` 区切りの全数列文字列（Run Script がクリップボードへペースト）※Alfred の仕様により CRLF を使用

## Error Handling

- 無効な入力（空・不正フォーマット）→ `{"title": "No results", "valid": false}` を表示
- `seq fmt` でフォーマット文字列のみ（range なし）→ 入力を促すプレースホルダーアイテムを表示
