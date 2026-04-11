# discord-messenger

個人向けの通知を Discord に集約するためのプロジェクト。\
Amazon SNS, AWS Lambda, Discord Webhook を組み合わせて通知基盤を構築しています。

## 主な機能

### cost-notification

定期的にAWSコスト情報を取得して通知します。

## 主要技術

- AWS CDK (TypeScript)
- Go

## ディレクトリ構成

```
discord-messenger/
|-- .github/workflows/
|   `-- workflows/      # GitHub Actions ワークフロー定義
|-- infra/              # AWS CDK プロジェクト (TypeScript)
|   |-- bin/            # CDK エントリポイント
|   `-- lib/            # スタック/Construct 定義
|-- lambda/             # Lambda 実装 (Go)
|   |-- cmd/            # Lambda エントリポイント
|   `-- internal/       # 共通処理やドメインロジック
`-- README.md
```
