# UI Design

Alfred Sequential Number は Alfred 5 の Script Filter / Run Script を使用するワークフローであり、
独自の UI コンポーネントを持たない。Alfred 自身が UI を制御する。

## Alfred Script Filter Items

Script Filter が返す JSON アイテムの仕様。

### 通常アイテム

```json
{
  "title": "1, 2, 3, ..., 10",
  "subtitle": "Paste Sequential Numbers (decimal)",
  "arg": "1\r\n2\r\n3\r\n...\r\n10",
  "valid": true
}
```

| フィールド | 内容 |
|---|---|
| `title` | 数列プレビュー（5要素以下はそのまま、超過時は "1, 2, 3, ..., N" 形式） |
| `subtitle` | フォーマット説明文 |
| `arg` | `\r\n` 区切りの全数列（Run Script がクリップボードへペースト、Alfred の仕様により CRLF を使用） |

### エラー・プレースホルダーアイテム

```json
{
  "title": "No results",
  "subtitle": "Paste Sequential Numbers (decimal)",
  "valid": false
}
```

## UI ガイドライン

- `title` / `subtitle` に Unicode 絵文字を使用しない（AI_CONTEXT.md の UI ガイドライン参照）
- アイコンは `workflow/icon.png` で制御する
- 外観モード（ライト/ダーク）は Alfred が制御するため、ワークフロー側での対応は不要
