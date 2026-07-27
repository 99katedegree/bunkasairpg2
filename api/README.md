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

- 武器 126 本 / アイテム 44 個 / 魔物 500 体（341 種族）。
  推奨レベルは Lv1-49 が 334 体、Lv50-90 が 166 体でおよそ 2 : 1。
- 最終ティア T12 の 18 本は 3 攻撃種別 × 6 属性を全て揃えてあり、積も
  2,000,000 で統一してある。実効火力は √積 に比例するので 18 本は同じ強さ。
- 図鑑番号は武器・アイテム・魔物とも 4 桁の数字。魔物だけは上位 3 桁が種族、
  末尾 1 桁が歪みの深さという意味を持ち、同じ種族はより深く歪んだ個体と
  上位 3 桁を共有する（例: 0010 / 0011 / 0012）。
  図鑑は番号順、攻略順はスライスの並び順で、この 2 つは独立している。
- 世界設定（六輝と姿定めぬ竜、12 の土地の由来）は `cmd/seed-master/story.go`。
  エリア名・武器名・アイテム名・魔物名はすべてそこから引かれている。
- 最初の土地「陽だまりの草原」の 24 体は、素手・Lv1 でどの 1 体からでも
  勝てることが不変条件。文化祭で QR をばらばらに配置するため、ここが崩れると
  最初の 1 体で詰む参加者が出る。`validate()` が機械的に検査している。
- 実データは `cmd/seed-master/data_weapons.go` / `data_items.go` / `data_monsters.go`。
  `data.go` には型と耐性値の目盛りだけが入っている。
- 素手は DB に入れない。装備なしのとき `handler/me.go` が id=0 のハードコード値を返す仕様で、
  seed 側は `data_weapons.go` の `bareHands` をバランス計算の出発点に使うだけ。
- ID 固定の upsert なので何度実行しても安全。図鑑・所持品・バトル履歴は壊さず、
  数値の調整だけを上書きできる。
- 画像は武器・アイテム・魔物とも **図鑑番号**（4 桁）をキーに
  `web/src/constants/*-images.ts` から引かれる。内部 ID ではなく図鑑番号を使うのは、
  QR や図鑑に印刷されて人が目にする番号がこちらだから。管理画面の
  `/monsters/index-numbers` などもこの番号を返し、登録漏れを突き合わせる。
- 投入前に、入手不可能な武器・ドロップの循環・倒せないモンスターがないかを検証する。
  問題があれば 1 件も書き込まずに終了する。
