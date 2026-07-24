package battle

// BossParams はハードコードされたボスのパラメータ。
// フロントエンドの /front/src/const/boss.ts と同一値を維持すること。
var BossParams = MonsterParams{
	Attack:          1000,
	HitPoint:        4000,
	MaxHitPoint:     4000,
	ExperiencePoint: 16000,
	Slash:           1.0,
	Blow:            1.0,
	Shoot:           1.0,
	Neutral:         1.0,
	Flame:           1.4,
	Water:           1.4,
	Wood:            1.4,
	Shine:           1.4,
	Dark:            1.4,
}

var bossPhysicsTypes = []string{"slash", "blow", "shoot"}
var bossElementTypes = []string{"neutral", "flame", "water", "wood", "shine", "dark"}

func (m *MonsterState) ShiftWeakness(rng *RNG) {
	pIdx := int(rng.Next() * 3)
	eIdx := int(rng.Next() * 6)
	physicsType := bossPhysicsTypes[pIdx]
	elementType := bossElementTypes[eIdx]

	m.Slash = 1.0
	m.Blow = 1.0
	m.Shoot = 1.0
	m.Neutral = 1.0
	m.Flame = 1.4
	m.Water = 1.4
	m.Wood = 1.4
	m.Shine = 1.4
	m.Dark = 1.4

	switch physicsType {
	case "slash":
		m.Slash = 0.7
	case "blow":
		m.Blow = 0.7
	case "shoot":
		m.Shoot = 0.7
	}
	switch elementType {
	case "neutral":
		m.Neutral = 0.9
	case "flame":
		m.Flame = 0.9
	case "water":
		m.Water = 0.9
	case "wood":
		m.Wood = 0.9
	case "shine":
		m.Shine = 0.9
	case "dark":
		m.Dark = 0.9
	}
}
