# readme-gen

[English](README.md)

README.mdを構造自動同期で管理するCLIツール。

## 特徴

- テンプレートからREADME生成（oss / general）
- マーカーベースのディレクトリ構造自動更新
- Claude Code連携でAI説明文生成
- CI連携用のチェックコマンド
- 日本語/英語対応

## ユースケース

**こんな時に最適:**
- 小〜中規模プロジェクトのREADMEをサクッと作りたい時
- `create-next-app`や`go mod init`で作成後、READMEが放置されている時
- 個人プロジェクトでREADMEのメンテナンスが面倒な時
- まずは叩き台が欲しい時

**向いていない用途:**
- 大規模モノレポや複雑な構造のプロジェクト
- 詳細なドキュメントが必要な場合（専用ドキュメントツールを使用）
- 完璧な本番用ドキュメントの生成

readme-genは**構造の同期**と**初期雛形の生成**に特化しています。包括的なドキュメント生成ツールではありません。

## 構造

<!-- readme-gen:structure:start -->
```
├── .claude/           # Claude Code skills
│   └── skills/
├── .github/           # GitHub Actions
│   └── workflows/
├── cmd/               # CLIエントリーポイント
│   └── readme-gen/
├── extras/            # ユーザー配布用スキル
│   └── skills/
└── internal/          # 内部パッケージ
    ├── cmd/           # Cobraコマンド定義
    ├── i18n/          # 国際化（日/英）
    ├── marker/        # マーカー更新処理
    ├── scanner/       # ディレクトリスキャン
    ├── template/      # テンプレート処理
    │   └── templates/
    └── ui/            # Charm UIスタイル
```
<!-- readme-gen:structure:end -->

## インストール

```bash
# Go
go install github.com/hulk510/readme-gen@latest

# または curl
curl -fsSL https://raw.githubusercontent.com/hulk510/readme-gen/main/install.sh | bash
```

## 使い方

### 初期化

```bash
# 対話モード
readme-gen init

# テンプレート指定
readme-gen init --template oss

# 非対話モード（全てデフォルト）
readme-gen init --yes

# AI生成付き
readme-gen init --yes --with-ai
```

### 構造更新

```bash
# 現在の構造を表示
readme-gen structure

# README.mdの構造を更新
readme-gen structure --update
```

### 差分チェック

```bash
# 構造が最新かチェック（CI用）
readme-gen check
# → 差分があればexit 1
```

## コマンドオプション

### `readme-gen init`

| オプション | 説明 |
|-----------|------|
| `-t, --template` | テンプレート選択（oss, general） |
| `-y, --yes` | 非対話モード |
| `--with-skills` | Claude Code skillsを追加 |
| `--with-ai` | AIで説明を自動生成 |
| `--no-skills` | skills追加をスキップ |
| `--no-ai` | AI生成をスキップ |
| `--lang` | 言語指定（en, ja） |

### `readme-gen structure`

| オプション | 説明 |
|-----------|------|
| `--update` | README.mdの構造を更新 |

### `readme-gen check`

| オプション | 説明 |
|-----------|------|
| `--lang` | 言語指定（en, ja） |

## Claude Code連携

`readme-gen init` でClaude Code skillsを追加すると、`.claude/skills/readme-update.md` が作成されます。

Claude Codeでプロジェクトを開くと、構造変更時に自動でREADME更新を提案してくれます。

### AI説明生成

`--with-ai` オプションまたは対話で選択すると、Claude Codeが各ディレクトリの説明を自動生成します。

```bash
readme-gen init --yes --with-ai
```

## 開発

開発のセットアップとガイドラインは [CONTRIBUTING.ja.md](CONTRIBUTING.ja.md) を参照してください。

### クイックスタート

```bash
# クローン
git clone https://github.com/hulk510/readme-gen.git
cd readme-gen

# 依存関係インストール
go mod download

# ビルド
mise run build
# または
go build -o bin/readme-gen ./cmd/readme-gen

# テスト実行
mise run test
# または
go test ./...

# ローカル実行
mise run dev
# または
go run ./cmd/readme-gen
```

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) を参照してください。
