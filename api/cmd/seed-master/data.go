package main

import "github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"

// ============================================================================
// マスターデータの型定義と、バランス調整で使う共通の目盛り。
//
// 実データは data_weapons.go / data_items.go / data_monsters.go にある。
// ここには「数値の意味」だけを置き、実際の値は置かない。
// ============================================================================

// --- 攻撃種別（物理タイプ） --------------------------------------------------
//
//	slash 斬撃: 剣 / 双剣 / 斧 / 短剣 / 爪
//	blow  打撃: ハンマー / 拳 / メイス / モーニングスター
//	shoot 射撃: 銃 / 弓 / 重銃 / 両手銃
const (
	Slash = "slash"
	Blow  = "blow"
	Shoot = "shoot"
)

// --- 属性 -------------------------------------------------------------------
const (
	Neutral = "neutral" // 無
	Flame   = "flame"   // 火
	Water   = "water"   // 水
	Wood    = "wood"    // 木
	Shine   = "shine"   // 光
	Dark    = "dark"    // 闇
)

// --- 耐性値の目盛り ----------------------------------------------------------
//
// ダメージ計算（internal/domain/battle/simulator.go）は
//
//	物理 = 物理攻撃力 * (1 + 物理バフ) * (1 - 物理耐性 + 物理デバフ)
//	属性 = 属性攻撃力 * (1 + 属性バフ) * (1 - 属性耐性 + 属性デバフ)
//	ダメージ = √|物理 * 属性| * レベル補正 * 乱数 * 符号
//
// なので耐性値は「軽減率」であり、次のように効く。
//
//	負の値 → 係数が 1 を超える     = 弱点
//	0.0    → 等倍
//	0〜1   → 軽減
//	1.0    → 係数 0。その軸だけで積が 0 になるので完全無効
//	1.0 超 → 係数が負。積が負になりモンスターが回復する = 吸収
//
// 設計時に必ず意識すること。
//  1. 物理・属性のどちらか一方でも 1.0 なら、もう一方が何であってもダメージ 0。
//  2. 攻撃種別と無属性は 1.0（無効）が上限で吸収できない。吸収を置けるのは
//     無属性を除く 5 属性だけ。したがって物理と属性の係数が同時に負になることはなく、
//     「両方吸収させて符号を反転させる」型のギミックは原理的に作れない。
//
// 上限下限は entity.ResistanceBounds が唯一の定義。ここではその範囲内の目盛りに
// 名前を付けているだけで、範囲そのものは持たない。
//
//	攻撃種別と無属性  -1.0 〜 1.0  吸収できない
//	無属性を除く5属性 -1.0 〜 2.0  吸収できる
const (
	MaxWeak   = entity.ResistanceMin        // 2倍弱点  (係数 2.0)
	VeryWeak  = -0.5                        // 大弱点   (係数 1.5)
	Weak      = -0.3                        // 弱点     (係数 1.3)
	Even      = 0.0                         // 等倍
	Guard     = 0.4                         // 軽減     (係数 0.6)
	HardGuard = 0.7                         // 強耐性   (係数 0.3)
	Nullify   = entity.ResistanceMaxPhysics // 無効     (係数 0。この軸に触れた時点でダメージ 0)
	Drain     = 1.3                         // 吸収     (係数 -0.3。相手を回復させる)
	BigDrain  = 1.6                         // 大吸収   (係数 -0.6)
	FullDrain = entity.ResistanceMaxElement // 100%吸収 (係数 -1.0。与えたダメージがそのまま回復に回る)
)

// --- アイテム効果種別 --------------------------------------------------------
const (
	EffectHeal   = "heal"
	EffectBuff   = "buff"
	EffectDebuff = "debuff"
)

// ============================================================================
// 型
// ============================================================================

// WeaponSeed は weapons テーブル 1 行。
// ID は固定値。web/src/constants/weapon-images.ts は ID をキーに画像を引くので、
// 一度配ったら ID は絶対に変えないこと。
//
// 素手は例外で DB に入れない。装備なし（equipped_weapon_id = NULL）のときに
// handler/me.go が id=0 のハードコード値を返す仕様なので、こちらは data_weapons.go の
// bareHands としてバランス計算にだけ使う。
type WeaponSeed struct {
	ID            int64
	IndexNumber   string
	Name          string
	Category      string // 剣 / ハンマー / 弓 … 表示には使わないがデータの意図を残す
	PhysicsAttack int
	ElementAttack *int // nil = 属性攻撃力なし。計算上 1 として扱われるので実質ただの物理武器
	PhysicsType   string
	ElementType   string
	Note          string
}

// ItemSeed は items + heal_items / buff_items / debuff_items の 1 セット。
// EffectType に応じて Amount / Rate + Target のどちらかだけを使う。
type ItemSeed struct {
	ID          int64
	IndexNumber string
	Name        string
	EffectType  string
	Amount      int     // heal のみ
	Rate        float64 // buff / debuff のみ。0.1 刻みでないと切り捨てられる
	Target      string  // buff / debuff のみ。攻撃種別または属性
	Note        string
}

// Resistance は 9 軸の耐性。ゼロ値が Even（等倍）なので、
// 特徴のある軸だけをフィールド名指定で書けばよい。
type Resistance struct {
	Slash   float64
	Blow    float64
	Shoot   float64
	Neutral float64
	Flame   float64
	Water   float64
	Wood    float64
	Shine   float64
	Dark    float64
}

// MonsterSeed は monsters テーブル 1 行。
// ID は固定 UUID。QR コードのトークンと web/src/constants/monster-images.ts の
// キーを兼ねるので、一度印刷したら変えないこと。
type MonsterSeed struct {
	ID               string
	IndexNumber      string
	Name             string
	Area             string
	RecommendedLevel int
	Attack           int
	HitPoint         int
	ExperiencePoint  int
	Res              Resistance
	DropWeaponID     *int64 // 初回撃破時のみ入手。既に持っていれば何も落とさない
	DropItemID       *int64 // 撃破するたび 1 個入手。周回できる
	// RequiresItem を true にすると「素の武器だけでは倒せない」検証エラーを黙らせる。
	// デバフアイテムで耐性を削ることが前提の謎解きモンスターに付ける。
	RequiresItem bool
	Note         string
}

func wid(id int64) *int64 { return &id }
func iid(id int64) *int64 { return &id }
func ea(v int) *int       { return &v }
