# コントリビューション

[English](CONTRIBUTING.md)

readme-genへのコントリビューションに興味を持っていただきありがとうございます！

## 開発環境のセットアップ

### 前提条件

- Go 1.23以上
- [mise](https://mise.jdx.dev/)（オプション、推奨）

### 始め方

```bash
# リポジトリをクローン
git clone https://github.com/hulk510/readme-gen.git
cd readme-gen

# 依存関係をインストール
go mod download

# セットアップを確認
go build ./...
go test ./...
```

## 開発コマンド

### miseを使う場合（推奨）

```bash
mise run build      # bin/readme-genにビルド
mise run dev        # 開発モードで実行
mise run test       # 全テスト実行
mise run lint       # リンター実行
mise run install    # $GOPATH/binにインストール
mise run clean      # ビルド成果物を削除
```

### Goを直接使う場合

```bash
# ビルド
go build -o bin/readme-gen ./cmd/readme-gen

# 実行
go run ./cmd/readme-gen

# テスト
go test ./...

# 詳細出力でテスト
go test -v ./...

# 特定パッケージのテスト
go test ./internal/scanner/...

# インストール
go install ./cmd/readme-gen
```

## プロジェクト構造

```
cmd/readme-gen/     # CLIエントリーポイント
internal/
├── cmd/            # Cobraコマンド定義
├── i18n/           # 国際化
├── marker/         # マーカーベース更新
├── scanner/        # ディレクトリスキャン
├── template/       # READMEテンプレート
└── ui/             # ターミナルUIスタイル
```

## 変更を加える

### 新機能の追加

1. フィーチャーブランチを作成: `git checkout -b feature/your-feature`
2. まずテストを書く
3. 機能を実装
4. 全テストが通ることを確認: `go test ./...`
5. 必要に応じてドキュメントを更新
6. プルリクエストを送信

### 翻訳の追加

1. `internal/i18n/i18n.go` に新しいメッセージを追加
2. 英語と日本語の両方の翻訳を追加
3. 必要に応じてテンプレートを更新

### 新しいテンプレートの追加

1. `internal/template/templates/` にテンプレートファイルを作成
2. 命名形式: `{name}.md.tmpl`（英語）または `{name}_ja.md.tmpl`（日本語）
3. `internal/template/template_test.go` にテストを追加

## コードスタイル

- 標準的なGoの規約に従う
- コミット前に `gofmt` を実行
- 意味のある変数名・関数名を使用
- エクスポートされる関数にはコメントを追加

## プルリクエストの流れ

1. リポジトリをフォーク
2. フィーチャーブランチを作成
3. 明確なメッセージで変更をコミット
4. フォークにプッシュ
5. プルリクエストを作成

### コミットメッセージの形式

[Conventional Commits](https://www.conventionalcommits.org/) を使用:

```
feat: 新しいテンプレートオプションを追加
fix: 構造解析を修正
docs: READMEを更新
test: スキャナーのテストを追加
refactor: マーカーロジックを簡素化
```

## テスト

### テストの実行

```bash
# 全テスト
go test ./...

# カバレッジ付き
go test -cover ./...

# カバレッジレポート生成
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### テストの書き方

- テストは `*_test.go` ファイルに配置
- 適切な場合はテーブル駆動テストを使用
- 成功ケースとエラーケースの両方をテスト

## リリースプロセス

リリースは [release-please](https://github.com/googleapis/release-please) で自動化されています。

1. PRを `main` にマージ
2. release-pleaseがリリースPRを作成/更新
3. リリースPRをマージすると新しいリリースがトリガー
4. GoReleaserがバイナリをビルド・公開

## 質問がありますか？

質問や提案があればお気軽にissueを作成してください！
