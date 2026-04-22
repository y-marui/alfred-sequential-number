# File Map

AI が作業のたびに依存関係を追記していくファイルレベルのマップ。
情報が足りない・古い場合は探索して追記・更新する。

## Entry Points

| ファイル | 用途 |
|---|---|
| `workflow/filter.py` | Alfred Script Filter から直接実行（キーワード入力中のリアルタイム表示） |
| `workflow/run.py` | Alfred Run Script から直接実行（選択アイテムの確定・クリップボードペースト） |
| `workflow/scripts/entry.py` | 開発用シミュレータ（`make run`）経由のエントリーポイント |

## Dependencies

```
workflow/filter.py
    └── workflow/seq.py (loader, formatter, multi_range, load_range)

workflow/run.py
    └── workflow/seq.py (loader)

workflow/seq.py          ← 数列生成コア（依存なし）
    ├── loader()         ← パターン文字列を解析して数列を生成
    ├── formatter        ← %a/%A（アルファベット変換）対応カスタム Formatter
    ├── multi_range()    ← 多次元範囲の直積を生成
    ├── load_range()     ← "N" or "S-E" 形式を [start, end] に変換
    ├── convert_alphabet() ← 1-indexed 整数 → 小文字アルファベット列
    └── convert_Alphabet() ← 1-indexed 整数 → 大文字アルファベット列

src/app/commands/seq_cmd.py
    └── workflow/seq.py (loader)   ← テスト用ラッパー

src/app/core.py
    └── src/app/commands/seq_cmd.py (generate)

tests/test_commands.py
    └── src/app/core.py (generate)

tests/test_services.py
    └── (サービス層のテスト)

tests/test_alfred.py
    └── src/alfred/ (response, router, cache, config, logger, safe_run)
```

## Key Files

| ファイル | 役割 |
|---|---|
| `workflow/info.plist` | Alfred ワークフロー定義・キーワードトリガー・Config Builder |
| `workflow/seq.py` | 数列生成コア（Alfred から直接参照） |
| `workflow/filter.py` | Script Filter ロジック（リアルタイム表示） |
| `workflow/run.py` | Run Script ロジック（クリップボード出力） |
| `src/alfred/` | Alfred SDK ヘルパー（response / router / cache / config / logger / safe_run） |
| `src/app/commands/seq_cmd.py` | seq_cmd テスト用ラッパー |
| `pyproject.toml` | プロジェクト設定・依存関係・バージョン |
| `Makefile` | 開発タスク統一インターフェース |
