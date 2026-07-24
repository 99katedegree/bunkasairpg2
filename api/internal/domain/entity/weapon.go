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
