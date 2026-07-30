package entity

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Monster struct {
	ID               uuid.UUID
	WeaponID         *int64
	ItemID           *int64
	IndexNumber      string
	Name             string
	Attack           int
	HitPoint         int
	ExperiencePoint  int
	RecommendedLevel int
	ImageURL         *string
	Slash            float64
	Blow             float64
	Shoot            float64
	Neutral          float64
	Flame            float64
	Water            float64
	Wood             float64
	Shine            float64
	Dark             float64
	Weapon           *Weapon
	Item             *Item
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// MonsterCatalogEntry はユーザーの図鑑エントリ（閲覧済みのみ）
type MonsterCatalogEntry struct {
	MonsterID *uuid.UUID // nil = 未登録スロット
}

// ============================================================================
// 耐性値のルール
//
// 耐性は「軽減率」で、ダメージ計算では 1 - 耐性 の形で係数になる（simulator.go）。
//
//	-1.0 → 係数  2.0  2倍弱点
//	 0.0 → 係数  1.0  通常
//	 0.5 → 係数  0.5  50%耐性
//	 1.0 → 係数  0.0  無効
//	 1.5 → 係数 -0.5  50%吸収
//	 2.0 → 係数 -1.0  100%吸収（与えたダメージがそのまま相手の回復になる）
//
// 上限下限は攻撃種別と属性で異なる。物理は無効までで吸収できない。
// 無属性は武器の属性タイプとして使われるが、この上限下限に関しては物理側に属し、
// 吸収されない。属性吸収を持つ相手に対する唯一の安全牌という位置づけになっている。
//
// DB のカラムは DECIMAL(4,2) でこれよりずっと広い値を格納できてしまうので、
// 範囲の保証はここが唯一の拠り所になる。管理画面のスライダーや投入コマンドは
// 独自に範囲を持たず、必ず ResistanceBounds / ValidateResistances を使うこと。
// ============================================================================

const (
	ResistanceMin        = -1.0 // 2倍弱点
	ResistanceMaxPhysics = 1.0  // 無効。攻撃種別と無属性はここまで
	ResistanceMaxElement = 2.0  // 100%吸収。無属性を除く5属性はここまで
	ResistanceStep       = 0.1  // 刻み幅
)

// 上限下限の適用先。無属性が物理側にいるのは上記のとおり意図的。
var (
	physicsAxes = []string{"slash", "blow", "shoot", "neutral"}
	elementAxes = []string{"flame", "water", "wood", "shine", "dark"}
)

// ResistanceBounds は軸名に対する耐性の許容範囲を返す。
// 未知の軸名なら ok が false になる。
func ResistanceBounds(axis string) (minimum, maximum float64, ok bool) {
	if slices.Contains(physicsAxes, axis) {
		return ResistanceMin, ResistanceMaxPhysics, true
	}
	if slices.Contains(elementAxes, axis) {
		return ResistanceMin, ResistanceMaxElement, true
	}
	return 0, 0, false
}

// Resistances は9軸の耐性を軸名で引ける形にして返す。
func (m *Monster) Resistances() map[string]float64 {
	return map[string]float64{
		"slash": m.Slash, "blow": m.Blow, "shoot": m.Shoot,
		"neutral": m.Neutral, "flame": m.Flame, "water": m.Water,
		"wood": m.Wood, "shine": m.Shine, "dark": m.Dark,
	}
}

// ValidateResistances は9軸すべてが許容範囲かつ刻み幅どおりかを検証する。
func (m *Monster) ValidateResistances() error {
	for _, axis := range append(append([]string{}, physicsAxes...), elementAxes...) {
		if err := ValidateResistance(axis, m.Resistances()[axis]); err != nil {
			return err
		}
	}
	return nil
}

// ValidateResistance は 1 軸分を検証する。
func ValidateResistance(axis string, value float64) error {
	minimum, maximum, ok := ResistanceBounds(axis)
	if !ok {
		return fmt.Errorf("%w: 未知の耐性軸 %q", ErrInvalidResistance, axis)
	}
	if value < minimum || value > maximum {
		return fmt.Errorf("%w: %s の耐性 %.2f が範囲 [%.1f, %.1f] の外",
			ErrInvalidResistance, axis, value, minimum, maximum)
	}
	// 0.1 刻み以外はスライダーで作れない値なので弾く。
	if math.Abs(value/ResistanceStep-math.Round(value/ResistanceStep)) > 1e-9 {
		return fmt.Errorf("%w: %s の耐性 %.2f が刻み幅 %.1f に乗っていない",
			ErrInvalidResistance, axis, value, ResistanceStep)
	}
	return nil
}
