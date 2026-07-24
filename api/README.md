# bunkasairpg API

## 技術スタック

| カテゴリ | ライブラリ |
|---|---|
| フレームワーク | [Echo](https://echo.labstack.com/) |
| DI | [uber-go/dig](https://github.com/uber-go/dig) |
| コード生成 | [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) |
| DB クライアント | [sqlc](https://sqlc.dev/) |
| マイグレーション | [golang-migrate](https://github.com/golang-migrate/migrate) |
| ストレージ | [aws-sdk-go-v2](https://github.com/aws/aws-sdk-go-v2)（R2 / S3互換） |
| 設定 / 環境変数 | [envconfig](https://github.com/kelseyhightower/envconfig) |
| ロギング | slog（Go標準） |
| バリデーション | [go-playground/validator](https://github.com/go-playground/validator) |
| 認証 | [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |

## アーキテクチャ

- スキーマ駆動開発（OpenAPI → oapi-codegen でサーバーインターフェース生成）
- クリーンアーキテクチャ

```
internal/
├── adapter/         # Echo ハンドラー（oapi-codegen インターフェース実装）
├── domain/          # エンティティ・リポジトリインターフェース
├── infrastructure/  # RDS（sqlc）・R2 実装
└── usecase/         # ユースケース
```

## ストレージ

- データベース: Amazon RDS
- オブジェクトストレージ: Cloudflare R2
