package entity

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrAlreadyExists  = errors.New("already exists")
	ErrBattleLost     = errors.New("battle lost")
	ErrBattleInvalid  = errors.New("battle result invalid")
	ErrBattleExpired  = errors.New("battle expired")
	ErrItemNotFound   = errors.New("item not found")
	ErrItemStockEmpty = errors.New("item stock empty")
	ErrWeaponNotOwned = errors.New("weapon not owned")
	ErrNoMonsters     = errors.New("no monsters registered")
	ErrInvalidCount   = errors.New("count must be greater than 0")

	ErrInvalidResistance = errors.New("resistance out of range")
)
