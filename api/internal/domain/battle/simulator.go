package battle

import (
	"errors"
	"math"
)

type ActionType string

const (
	ActionAttack       ActionType = "attack"
	ActionUseItem      ActionType = "use-item"
	ActionChangeWeapon ActionType = "change-weapon"
)

type Action struct {
	Turn     int
	Type     ActionType
	ItemID   *int64
	WeaponID *int64
}

type WeaponParams struct {
	ID            int64
	PhysicsAttack float64
	ElementAttack *float64
	PhysicsType   string
	ElementType   string
}

type ItemParams struct {
	ID         int64
	EffectType string   // heal | buff | debuff
	Amount     *int     // heal のみ
	Rate       *float64 // buff/debuff のみ
	Target     *string  // buff/debuff のみ
}

type MonsterParams struct {
	Attack                      int
	HitPoint                    int
	MaxHitPoint                 int
	ExperiencePoint             int
	Slash, Blow, Shoot          float64
	Neutral, Flame, Water, Wood float64
	Shine, Dark                 float64
}

type MonsterState struct {
	MonsterParams
	CurrentHitPoint int
}

type UserState struct {
	Level           int
	HitPoint        int
	MaxHitPoint     int
	ExperiencePoint int
	Weapon          WeaponParams
}

type SimulatorInput struct {
	Seed    int64
	Actions []Action
	Monster MonsterParams
	IsBoss  bool
	User    UserState
	Items   map[int64]ItemParams
	Weapons map[int64]WeaponParams
}

type SimulatorResult struct {
	Won             bool
	MonsterDefeated bool
	ExperiencePoint int
	Level           int
	HitPoint        int
	DropWeaponID    *int64
	DropItemID      *int64
}

var ErrBattleLost = errors.New("battle lost")
var ErrBattleInvalid = errors.New("battle result invalid")

func Simulate(input SimulatorInput) (SimulatorResult, error) {
	rng := NewRNG(input.Seed)

	buffs := map[string]float64{}
	debuffs := map[string]float64{}

	ms := MonsterState{
		MonsterParams:   input.Monster,
		CurrentHitPoint: input.Monster.HitPoint,
	}
	us := input.User

	// アイテムは 1 ターンに 1 個まで。攻撃するたびに解除される。
	// リファクタ前のフロント（isItemUsed）と同じ制約をサーバー側でも保証する。
	itemUsedThisTurn := false

	for _, action := range input.Actions {
		switch action.Type {
		case ActionAttack:
			itemUsedThisTurn = false
			damage := calcPlayerDamage(rng, us, ms, buffs, debuffs)
			if damage < 0 {
				ms.CurrentHitPoint = min(ms.CurrentHitPoint-damage, ms.MaxHitPoint)
			} else {
				ms.CurrentHitPoint = max(ms.CurrentHitPoint-damage, 0)
			}

		case ActionUseItem:
			if action.ItemID == nil {
				return SimulatorResult{}, ErrBattleInvalid
			}
			item, ok := input.Items[*action.ItemID]
			if !ok {
				return SimulatorResult{}, ErrBattleInvalid
			}
			// アイテムは 1 ターンに 1 個まで。攻撃を挟むまで次は使えない。
			if itemUsedThisTurn {
				return SimulatorResult{}, ErrBattleInvalid
			}
			itemUsedThisTurn = true
			applyItem(&us, buffs, debuffs, item)
			continue // アイテム使用はターンを消費せず、モンスター攻撃も発生しない

		case ActionChangeWeapon:
			if action.WeaponID == nil {
				return SimulatorResult{}, ErrBattleInvalid
			}
			w, ok := input.Weapons[*action.WeaponID]
			if !ok {
				return SimulatorResult{}, ErrBattleInvalid
			}
			us.Weapon = w
			continue // 武器変更はモンスター攻撃を発生させない
		}

		// モンスターが生きていればモンスター攻撃
		if ms.CurrentHitPoint > 0 {
			monDmg := calcMonsterDamage(rng, ms, us)
			us.HitPoint = max(us.HitPoint-monDmg, 0)

			// ボスのみ弱点シフト
			if input.IsBoss {
				ms.ShiftWeakness(rng)
			}
		}

		// プレイヤー死亡チェック
		if us.HitPoint == 0 {
			return SimulatorResult{Won: false}, ErrBattleLost
		}

		// モンスター撃破
		if ms.CurrentHitPoint == 0 {
			break
		}
	}

	if ms.CurrentHitPoint > 0 {
		return SimulatorResult{}, ErrBattleInvalid
	}

	return SimulatorResult{
		Won:             true,
		MonsterDefeated: true,
		ExperiencePoint: ms.ExperiencePoint,
	}, nil
}

func calcPlayerDamage(rng *RNG, us UserState, ms MonsterState, buffs, debuffs map[string]float64) int {
	pt := us.Weapon.PhysicsType
	et := us.Weapon.ElementType

	// フロント側は (elementAttack || 1) なので 0 も 1 に読み替えられる。
	// nil だけを見ると 0 のときに結果がずれ、再計算が食い違うので同じ扱いにする。
	ea := 1.0
	if us.Weapon.ElementAttack != nil && *us.Weapon.ElementAttack != 0 {
		ea = *us.Weapon.ElementAttack
	}

	// デバフは耐性倍率への加算。以前は (1 - 耐性*(1-デバフ)) と乗算していたが、
	// あれは耐性を符号ごと 0 に近づける式なので、弱点(耐性が負)にかけると
	// ダメージが下がり、デバフが 1 を超えると耐性と弱点が入れ替わっていた。
	// 等倍(耐性 0)の相手には何の効果もなかった。加算なら常に単調増加する。
	physics := us.Weapon.PhysicsAttack * (1 + buffs[pt]) * (1 - resistance(ms, pt) + debuffs[pt])
	element := ea * (1 + buffs[et]) * (1 - resistance(ms, et) + debuffs[et])

	// 物理と属性は加算ではなく乗算。属性を攻撃全体にかかる係数として効かせるための
	// 構造で、加算にすると炎吸収の相手に炎属性武器で殴ったとき属性分だけ吸収されて
	// 物理分はダメージが通ってしまう。攻撃全体が吸収される方が納得感があるためこの形。
	base := physics * element

	lf := 0.8 + math.Sqrt(float64(us.Level))/5.0
	r := 0.95 + rng.Next()*0.1
	sign := 1.0
	if base < 0 {
		sign = -1.0
	}
	return int(math.Floor(math.Sqrt(math.Abs(base)) * lf * r * sign))
}

func calcMonsterDamage(rng *RNG, ms MonsterState, us UserState) int {
	r := 0.95 + rng.Next()*0.1
	lf := 1.0 + math.Sqrt(float64(us.Level))/1.7
	return int(math.Floor(float64(ms.Attack) * r / lf))
}

func applyItem(us *UserState, buffs, debuffs map[string]float64, item ItemParams) {
	switch item.EffectType {
	case "heal":
		if item.Amount != nil {
			us.HitPoint = min(us.HitPoint+*item.Amount, us.MaxHitPoint)
		}
	case "buff":
		if item.Rate != nil && item.Target != nil {
			buffs[*item.Target] += math.Floor(*item.Rate*10) / 10
		}
	case "debuff":
		if item.Rate != nil && item.Target != nil {
			debuffs[*item.Target] += math.Floor(*item.Rate*10) / 10
		}
	}
}

func resistance(ms MonsterState, typ string) float64 {
	switch typ {
	case "slash":
		return ms.Slash
	case "blow":
		return ms.Blow
	case "shoot":
		return ms.Shoot
	case "neutral":
		return ms.Neutral
	case "flame":
		return ms.Flame
	case "water":
		return ms.Water
	case "wood":
		return ms.Wood
	case "shine":
		return ms.Shine
	case "dark":
		return ms.Dark
	}
	return 1.0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
