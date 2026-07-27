package entity

import "github.com/google/uuid"

// 管理画面が「マスタに何が存在するか」を把握するための要約。
//
// 画像は持たない。どの画像を割り当てるかはフロント側の関心で、
// web/src/constants/*-images.ts が図鑑番号をキーに対応付けている。
// バックエンドは存在するものを列挙するだけに留める。

type MonsterSummary struct {
	ID          uuid.UUID
	Name        string
	IndexNumber string
}

type WeaponSummary struct {
	ID          int64
	Name        string
	IndexNumber string
}

type ItemSummary struct {
	ID          int64
	Name        string
	IndexNumber string
}
