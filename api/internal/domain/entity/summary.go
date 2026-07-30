package entity

import "github.com/google/uuid"

// 管理画面が「マスタに何が存在するか」を把握するための要約。
//
// ImageURL は cmd/seed-master が図鑑番号から機械的に組み立てて書き込む
// （internal/adapter/handler の画像配信ルートを指す相対パス）。
// 画像がまだ用意されていない個体・武器・アイテムは nil。

type MonsterSummary struct {
	ID          uuid.UUID
	Name        string
	IndexNumber string
	ImageURL    *string
}

type WeaponSummary struct {
	ID          int64
	Name        string
	IndexNumber string
	ImageURL    *string
}

type ItemSummary struct {
	ID          int64
	Name        string
	IndexNumber string
	ImageURL    *string
}
