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

## コマンド

| コマンド | 用途 |
|---|---|
| `go run ./cmd/server` | API サーバー起動 |
| `go run ./cmd/create-admin` | 管理者アカウント作成（`EMAIL` / `PASSWORD` を env で指定） |
| `go run ./cmd/seed-master` | 本番用マスターデータ（武器・アイテム・モンスター）の投入 |

### seed-master

```sh
go run ./cmd/seed-master -dry-run    # DB に触らず、検証とバランス表だけ出す
go run ./cmd/seed-master             # 投入
go run ./cmd/seed-master -report     # 投入したうえでバランス表も出す
```

- 武器 120 本 / アイテム 44 個 / モンスター 500 体。
  推奨レベルは Lv1-49 が 334 体、Lv50-90 が 166 体でおよそ 2 : 1。
- 実データは `cmd/seed-master/data_weapons.go` / `data_items.go` / `data_monsters.go`。
  `data.go` には型と耐性値の目盛りだけが入っている。
- 素手は DB に入れない。装備なしのとき `handler/me.go` が id=0 のハードコード値を返す仕様で、
  seed 側は `data_weapons.go` の `bareHands` をバランス計算の出発点に使うだけ。
- ID 固定の upsert なので何度実行しても安全。図鑑・所持品・バトル履歴は壊さず、
  数値の調整だけを上書きできる。
- 武器 / アイテムの画像は **DB の ID**、モンスターの画像は **UUID** をキーに
  `web/src/constants/*-images.ts` から引かれる。QR を印刷したあとに ID を変えないこと。
- 投入前に、入手不可能な武器・ドロップの循環・倒せないモンスターがないかを検証する。
  問題があれば 1 件も書き込まずに終了する。
