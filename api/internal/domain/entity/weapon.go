package entity

import "time"

type Weapon struct {
	ID            int64
	Name          string
	IndexNumber   string
	PhysicsAttack float64
	ElementAttack *float64
	PhysicsType   string // slash | blow | shoot
	ElementType   string // neutral | flame | water | wood | shine | dark
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// BareHands は武器を 1 本も装備していないときに使われる初期装備。
// DB には保存しないので、この変数が唯一の定義になる。
// API のレスポンスとバトルの再計算の両方がここを見るため、片方だけずれることはない。
var BareHands = Weapon{
	ID:            0,
	Name:          "素手",
	IndexNumber:   "W000",
	PhysicsAttack: 10,
	ElementAttack: nil, // nil は計算上 1 として扱われる
	PhysicsType:   "blow",
	ElementType:   "neutral",
}
