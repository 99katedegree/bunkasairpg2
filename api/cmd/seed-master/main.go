// Command seed-master は本番用のマスターデータ（武器・アイテム・モンスター）を投入する。
//
//	go run ./cmd/seed-master              # 検証 → 投入
//	go run ./cmd/seed-master -dry-run     # DB に触らず検証とバランス表だけ出す
//	go run ./cmd/seed-master -report      # 投入したうえでバランス表も出す
//
// 素手は DB に入れない。装備なしのとき handler/me.go がハードコード値を返す仕様なので、
// このコマンドは data_weapons.go の bareHands をバランス計算の出発点にするだけ。
//
// 投入は ID 固定の upsert なので何度実行しても安全。既存の図鑑登録・所持品・
// バトル履歴は壊さず、数値の調整だけを上書きできる。
//
// 実データは data_weapons.go / data_items.go / data_monsters.go にある。
package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"os"
	"regexp"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"github.com/99katedegree/bunkasairpg2/api/internal/domain/entity"
)

func main() {
	var (
		envPath = flag.String("env", ".env", "読み込む .env のパス")
		dryRun  = flag.Bool("dry-run", false, "DB に書き込まず検証とバランス表だけ実行する")
		report  = flag.Bool("report", false, "投入後にバランス表を出力する")
	)
	flag.Parse()

	if problems := validate(); len(problems) > 0 {
		for _, p := range problems {
			slog.Error("マスターデータの検証に失敗", "detail", p)
		}
		os.Exit(1)
	}
	slog.Info("マスターデータの検証OK",
		"weapons", len(weapons), "items", len(items), "monsters", len(monsters))

	if *dryRun || *report {
		printBalanceReport()
	}
	if *dryRun {
		slog.Info("dry-run のため DB には書き込んでいない")
		return
	}

	if err := loadEnv(*envPath); err != nil {
		slog.Warn("env ファイルを読めなかった", "path", *envPath, "err", err)
	}

	db, err := connectDB()
	if err != nil {
		slog.Error("DB に接続できなかった", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		slog.Error("DB に疎通できなかった", "err", err)
		os.Exit(1)
	}

	if err := seed(ctx, db); err != nil {
		slog.Error("投入に失敗", "err", err)
		os.Exit(1)
	}
	slog.Info("投入完了",
		"weapons", len(weapons), "items", len(items), "monsters", len(monsters))
}

// ============================================================================
// 投入
// ============================================================================

func seed(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// モンスターが武器・アイテムを外部キーで参照するので先にこちらを入れる。
	if err := insertWeapons(ctx, tx); err != nil {
		return err
	}
	if err := insertItems(ctx, tx); err != nil {
		return err
	}
	if err := insertMonsters(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}

// ============================================================================
// 一括投入
//
// 1 行ずつ INSERT すると 1 行につき 1 往復かかる。リモートの RDS が相手だと
// 往復 10ms でも 670 行で 7 秒、100ms なら 4 分近くになる。複数行を 1 文に
// まとめて送ることで往復回数を数回まで落とす。
//
// バッチサイズはプレースホルダ数で決める。MySQL のプリペアドステートメントは
// パラメータ 65535 個が上限なので、1 行あたりの列数から逆算して余裕を持たせる。
// ============================================================================

// batchSize は 1 文に詰める行数を、1 行あたりの列数から決める。
func batchSize(columns int) int {
	const maxParams = 60000
	n := maxParams / columns
	if n > 500 {
		n = 500
	}
	if n < 1 {
		n = 1
	}
	return n
}

// placeholders は "(?, ?, ?), (?, ?, ?)" のような VALUES 部分を組み立てる。
func placeholders(rows, columns int, tail string) string {
	one := "(" + strings.Repeat("?, ", columns-1) + "?" + tail + ")"
	return strings.Repeat(one+", ", rows-1) + one
}

func insertWeapons(ctx context.Context, tx *sql.Tx) error {
	const cols = 7
	size := batchSize(cols)
	for i := 0; i < len(weapons); i += size {
		chunk := weapons[i:min(i+size, len(weapons))]
		args := make([]any, 0, len(chunk)*cols)
		for _, w := range chunk {
			var elementAttack any
			if w.ElementAttack != nil {
				elementAttack = *w.ElementAttack
			}
			args = append(args, w.ID, w.Name, w.IndexNumber, w.PhysicsAttack,
				elementAttack, w.PhysicsType, w.ElementType)
		}
		q := `INSERT INTO weapons
				(id, name, index_number, physics_attack, element_attack, physics_type, element_type, created_at, updated_at)
			VALUES ` + placeholders(len(chunk), cols, ", NOW(), NOW()") + `
			ON DUPLICATE KEY UPDATE
				name = VALUES(name),
				index_number = VALUES(index_number),
				physics_attack = VALUES(physics_attack),
				element_attack = VALUES(element_attack),
				physics_type = VALUES(physics_type),
				element_type = VALUES(element_type),
				updated_at = NOW()`
		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("武器の投入 (%d〜%d 件目): %w", i+1, i+len(chunk), err)
		}
	}
	return nil
}

func insertItems(ctx context.Context, tx *sql.Tx) error {
	const cols = 4
	size := batchSize(cols)
	for i := 0; i < len(items); i += size {
		chunk := items[i:min(i+size, len(items))]
		args := make([]any, 0, len(chunk)*cols)
		for _, it := range chunk {
			args = append(args, it.ID, it.Name, it.IndexNumber, it.EffectType)
		}
		q := `INSERT INTO items (id, name, index_number, effect_type, created_at, updated_at)
			VALUES ` + placeholders(len(chunk), cols, ", NOW(), NOW()") + `
			ON DUPLICATE KEY UPDATE
				name = VALUES(name),
				index_number = VALUES(index_number),
				effect_type = VALUES(effect_type),
				updated_at = NOW()`
		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("アイテムの投入 (%d〜%d 件目): %w", i+1, i+len(chunk), err)
		}
	}

	// 効果種別を作り替えたときに古い行が残らないよう、3 テーブルとも一度空にしてから入れ直す。
	// 投入するアイテム ID の集合で消すので、seed が管理していない行には触れない。
	ids := make([]any, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ID)
	}
	if len(ids) > 0 {
		in := "(" + strings.Repeat("?, ", len(ids)-1) + "?)"
		for _, table := range []string{"heal_items", "buff_items", "debuff_items"} {
			if _, err := tx.ExecContext(ctx,
				"DELETE FROM "+table+" WHERE item_id IN "+in, ids...); err != nil {
				return fmt.Errorf("%s の削除: %w", table, err)
			}
		}
	}

	// 効果種別ごとに分けて、それぞれ 1 文で入れる。
	var heal, buff, debuff []ItemSeed
	for _, it := range items {
		switch it.EffectType {
		case EffectHeal:
			heal = append(heal, it)
		case EffectBuff:
			buff = append(buff, it)
		case EffectDebuff:
			debuff = append(debuff, it)
		}
	}

	if len(heal) > 0 {
		args := make([]any, 0, len(heal)*2)
		for _, it := range heal {
			args = append(args, it.ID, it.Amount)
		}
		q := `INSERT INTO heal_items (item_id, amount, created_at, updated_at) VALUES ` +
			placeholders(len(heal), 2, ", NOW(), NOW()")
		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("回復アイテムの効果行: %w", err)
		}
	}
	for _, x := range []struct {
		table string
		rows  []ItemSeed
	}{{"buff_items", buff}, {"debuff_items", debuff}} {
		if len(x.rows) == 0 {
			continue
		}
		args := make([]any, 0, len(x.rows)*3)
		for _, it := range x.rows {
			args = append(args, it.ID, it.Rate, it.Target)
		}
		q := `INSERT INTO ` + x.table + ` (item_id, rate, target, created_at, updated_at) VALUES ` +
			placeholders(len(x.rows), 3, ", NOW(), NOW()")
		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("%s の効果行: %w", x.table, err)
		}
	}
	return nil
}

func insertMonsters(ctx context.Context, tx *sql.Tx) error {
	const cols = 18
	size := batchSize(cols)
	for i := 0; i < len(monsters); i += size {
		chunk := monsters[i:min(i+size, len(monsters))]
		args := make([]any, 0, len(chunk)*cols)
		for _, m := range chunk {
			var weaponID, itemID any
			if m.DropWeaponID != nil {
				weaponID = *m.DropWeaponID
			}
			if m.DropItemID != nil {
				itemID = *m.DropItemID
			}
			args = append(args,
				m.ID, weaponID, itemID, m.IndexNumber, m.Name, m.Attack, m.HitPoint,
				m.ExperiencePoint, m.RecommendedLevel,
				m.Res.Slash, m.Res.Blow, m.Res.Shoot,
				m.Res.Neutral, m.Res.Flame, m.Res.Water, m.Res.Wood, m.Res.Shine, m.Res.Dark)
		}
		q := `INSERT INTO monsters
				(id, weapon_id, item_id, index_number, name, attack, hit_point, experience_point,
				 recommended_level, slash, blow, shoot, neutral, flame, water, wood, shine, dark,
				 created_at, updated_at)
			VALUES ` + placeholders(len(chunk), cols, ", NOW(), NOW()") + `
			ON DUPLICATE KEY UPDATE
				weapon_id = VALUES(weapon_id),
				item_id = VALUES(item_id),
				index_number = VALUES(index_number),
				name = VALUES(name),
				attack = VALUES(attack),
				hit_point = VALUES(hit_point),
				experience_point = VALUES(experience_point),
				recommended_level = VALUES(recommended_level),
				slash = VALUES(slash), blow = VALUES(blow), shoot = VALUES(shoot),
				neutral = VALUES(neutral), flame = VALUES(flame), water = VALUES(water),
				wood = VALUES(wood), shine = VALUES(shine), dark = VALUES(dark),
				updated_at = NOW()`
		if _, err := tx.ExecContext(ctx, q, args...); err != nil {
			return fmt.Errorf("モンスターの投入 (%d〜%d 件目): %w", i+1, i+len(chunk), err)
		}
	}
	return nil
}

// ============================================================================
// 検証
// ============================================================================

func validate() []string {
	var problems []string
	add := func(format string, args ...any) {
		problems = append(problems, fmt.Sprintf(format, args...))
	}

	physicsTypes := map[string]bool{Slash: true, Blow: true, Shoot: true}
	elementTypes := map[string]bool{
		Neutral: true, Flame: true, Water: true, Wood: true, Shine: true, Dark: true,
	}

	// --- 武器 ---------------------------------------------------------------
	weaponIDs := map[int64]bool{}
	seenIndex := map[string]string{}
	for _, w := range weapons {
		if weaponIDs[w.ID] {
			add("武器 ID %d が重複している", w.ID)
		}
		weaponIDs[w.ID] = true
		if prev, ok := seenIndex[w.IndexNumber]; ok {
			add("武器 index_number %q が %s と %s で重複している", w.IndexNumber, prev, w.Name)
		}
		seenIndex[w.IndexNumber] = w.Name
		if !indexNumberPattern.MatchString(w.IndexNumber) {
			add("武器 %s の index_number %q が 4 桁の数字でない", w.Name, w.IndexNumber)
		}

		if !physicsTypes[w.PhysicsType] {
			add("武器 %s の physics_type %q が不正", w.Name, w.PhysicsType)
		}
		if !elementTypes[w.ElementType] {
			add("武器 %s の element_type %q が不正", w.Name, w.ElementType)
		}
		if w.PhysicsAttack <= 0 {
			add("武器 %s の physics_attack が 0 以下", w.Name)
		}
		if w.ElementAttack != nil && *w.ElementAttack <= 0 {
			add("武器 %s の element_attack が 0 以下（不要なら nil にする）", w.Name)
		}
	}

	// --- アイテム -----------------------------------------------------------
	itemIDs := map[int64]bool{}
	seenIndex = map[string]string{}
	for _, it := range items {
		if itemIDs[it.ID] {
			add("アイテム ID %d が重複している", it.ID)
		}
		itemIDs[it.ID] = true
		if prev, ok := seenIndex[it.IndexNumber]; ok {
			add("アイテム index_number %q が %s と %s で重複している", it.IndexNumber, prev, it.Name)
		}
		seenIndex[it.IndexNumber] = it.Name
		if !indexNumberPattern.MatchString(it.IndexNumber) {
			add("アイテム %s の index_number %q が 4 桁の数字でない", it.Name, it.IndexNumber)
		}

		switch it.EffectType {
		case EffectHeal:
			if it.Amount <= 0 {
				add("回復アイテム %s の amount が 0 以下", it.Name)
			}
		case EffectBuff, EffectDebuff:
			if it.Rate <= 0 {
				add("%s の rate が 0 以下", it.Name)
			}
			// シミュレータが math.Floor(rate*10)/10 で 0.1 刻みに落とすので、
			// 中途半端な値は意図した効果にならない。
			if math.Abs(it.Rate*10-math.Round(it.Rate*10)) > 1e-9 {
				add("%s の rate %.2f は 0.1 刻みでないため切り捨てられる", it.Name, it.Rate)
			}
			if !physicsTypes[it.Target] && !elementTypes[it.Target] {
				add("%s の target %q が不正", it.Name, it.Target)
			}
		default:
			add("%s の effect_type %q が不正", it.Name, it.EffectType)
		}
	}

	// --- モンスター ---------------------------------------------------------
	monsterIDs := map[string]bool{}
	seenIndex = map[string]string{}
	droppedWeapons := map[int64]bool{}
	droppedItems := map[int64]bool{}
	for _, m := range monsters {
		if monsterIDs[m.ID] {
			add("モンスター UUID %s が重複している", m.ID)
		}
		monsterIDs[m.ID] = true
		if prev, ok := seenIndex[m.IndexNumber]; ok {
			add("モンスター index_number %q が %s と %s で重複している", m.IndexNumber, prev, m.Name)
		}
		seenIndex[m.IndexNumber] = m.Name

		if m.Attack <= 0 || m.HitPoint <= 0 || m.ExperiencePoint <= 0 {
			add("%s の attack / hit_point / experience_point に 0 以下がある", m.Name)
		}
		if m.RecommendedLevel < 1 {
			add("%s の recommended_level が 1 未満", m.Name)
		}
		// 耐性の上限下限はドメインが唯一の定義。ここでは独自の範囲を持たない。
		for axis, v := range resistanceMap(m.Res) {
			if err := entity.ValidateResistance(axis, v); err != nil {
				add("%s: %v", m.Name, err)
			}
		}

		if m.DropWeaponID != nil && m.DropItemID != nil {
			// battle.go の報酬付与は weapon を優先し item を無視するので、両立させない。
			add("%s が武器とアイテムの両方をドロップ指定している（武器が優先されアイテムは死ぬ）", m.Name)
		}
		if m.DropWeaponID != nil {
			if !weaponIDs[*m.DropWeaponID] {
				add("%s のドロップ武器 ID %d が存在しない", m.Name, *m.DropWeaponID)
			}
			droppedWeapons[*m.DropWeaponID] = true
		}
		if m.DropItemID != nil {
			if !itemIDs[*m.DropItemID] {
				add("%s のドロップアイテム ID %d が存在しない", m.Name, *m.DropItemID)
			}
			droppedItems[*m.DropItemID] = true
		}

		// 図鑑番号は 4 桁の数字。上位 3 桁が種族、末尾 1 桁が個体番号。
		if !indexNumberPattern.MatchString(m.IndexNumber) {
			add("%s の index_number %q が 4 桁の数字でない", m.Name, m.IndexNumber)
		}

		// 武器だけで削れるか。削れないなら RequiresItem を明示させる。
		if !m.RequiresItem && bestWeaponAgainst(m, append([]WeaponSeed{bareHands}, weapons...)) == nil {
			add("%s はどの武器でもダメージが通らない（意図的なら RequiresItem: true を付ける）", m.Name)
		}
		if m.RequiresItem && bestWeaponAgainst(m, append([]WeaponSeed{bareHands}, weapons...)) != nil {
			add("%s は RequiresItem だが武器だけで倒せてしまう", m.Name)
		}
	}

	// --- 入手経路 -----------------------------------------------------------
	for _, w := range weapons {
		if !droppedWeapons[w.ID] {
			add("武器 %s (id=%d) を落とすモンスターがいない＝入手不可能", w.Name, w.ID)
		}
	}
	for _, it := range items {
		if !droppedItems[it.ID] {
			add("アイテム %s (id=%d) を落とすモンスターがいない＝入手不可能", it.Name, it.ID)
		}
	}

	// --- 最初の土地は素手で全部倒せること -----------------------------------
	problems = append(problems, firstAreaMustBeBeatableBareHanded()...)

	// --- 攻略順序と「存在しない組み合わせ」 -----------------------------------
	//
	// 上から順に潰していったとき、その時点の手持ちで本当に倒せるか。
	// 落とす武器がないと自分自身を倒せない循環と、その時点の武器に実在しない
	// 「攻撃種別 × 属性」を要求しているモンスターを、まとめてここで弾く。
	//
	// 武器表全体に対して同じ判定をしても意味がない。最終ティア T12 が 18 通りを
	// 全て揃えている以上、どの組み合わせも「いつかは」存在してしまうため。
	// 効くのは「そのモンスターに辿り着いた時点で」という進行順の条件だけ。
	for _, row := range simulateProgression() {
		if row.Damage > 0 {
			continue
		}
		m := row.Monster
		add("%s はここまでに入手できる装備ではダメージが通らない（通る組み合わせ: %s / その時点で持っている組み合わせ: %s）",
			m.Name, describeCombos(openCombos(m)), describeCombos(row.OwnedCombos))
	}

	sort.Strings(problems)
	return problems
}

// ============================================================================
// バランス表
// ============================================================================

// progressRow は 1 体分の想定値。
type progressRow struct {
	Monster     MonsterSeed
	Loadout     string // 想定した武器（+ 併用するデバフアイテム）
	Damage      int    // こちらの 1 ターンの与ダメージ
	TakenDamage int    // 相手の 1 ターンの与ダメージ
	// OwnedCombos はその時点で入手済みの武器が持つ「攻撃種別 × 属性」。
	// ダメージが通らなかったときに、何が足りないのかを示すために使う。
	OwnedCombos [][2]string
}

// simulateProgression は index_number 順に上から潰していく想定で、
// 「そのモンスターに辿り着いた時点で持っているはずの武器・アイテム」だけを候補に
// 推奨レベルで戦った場合の想定値を並べる。
//
// 前のモンスターが落とす武器でないと倒せない、といった順序の破綻をここで検出できる。
func simulateProgression() []progressRow {
	ownedWeapons := []WeaponSeed{bareHands}
	var ownedItems []ItemSeed

	rows := make([]progressRow, 0, len(monsters))
	for _, m := range monsters {
		level := m.RecommendedLevel
		row := progressRow{
			Monster:     m,
			Loadout:     "—",
			TakenDamage: monsterDamage(m, level),
			OwnedCombos: combosOf(ownedWeapons),
		}

		if best := bestWeaponAgainst(m, ownedWeapons); best != nil {
			row.Loadout = best.Name
			row.Damage = playerDamage(*best, m, level)
		} else if m.RequiresItem {
			// 武器だけでは通らない相手は、手持ちのデバフアイテムを併用した想定で出す。
			if w, used, d := bestDebuffComboAgainst(m, ownedWeapons, ownedItems, level); w != nil {
				names := make([]string, 0, len(used))
				for _, it := range used {
					names = append(names, it.Name)
				}
				row.Loadout = w.Name + " + " + strings.Join(names, " + ")
				row.Damage = d
			}
		}
		rows = append(rows, row)

		if m.DropWeaponID != nil {
			ownedWeapons = append(ownedWeapons, *weaponByID(*m.DropWeaponID))
		}
		if m.DropItemID != nil {
			ownedItems = append(ownedItems, *itemByID(*m.DropItemID))
		}
	}
	return rows
}

// printBalanceReport は想定ターン数を一覧で出す。
// 与ターンが 3〜7、被ターンがそれより十分大きければ狙いどおり。
func printBalanceReport() {
	fmt.Println()
	fmt.Println("その時点で入手済みの装備・推奨レベルで戦った場合の想定値（乱数は 1.0 固定）")
	fmt.Println("  与ターン = こちらが倒しきるまで / 被ターン = 倒されるまで")
	fmt.Println(strings.Repeat("-", 126))
	fmt.Printf("%-5s %-20s %3s %6s %5s %6s  %-34s %7s %6s %6s %5s\n",
		"No", "名前", "Lv", "HP", "攻撃", "EXP", "想定装備", "与ダメ", "与T", "被ダメ", "被T")
	fmt.Println(strings.Repeat("-", 126))

	totalExp := 0
	area := ""
	for _, row := range simulateProgression() {
		m := row.Monster
		if m.Area != area {
			area = m.Area
			fmt.Printf("\n[%s]\n", area)
		}

		dmgText, killText := "—", "—"
		if row.Damage > 0 {
			dmgText = fmt.Sprintf("%d", row.Damage)
			killText = fmt.Sprintf("%.1f", float64(m.HitPoint)/float64(row.Damage))
		}
		deathText := "—"
		if row.TakenDamage > 0 {
			deathText = fmt.Sprintf("%.1f",
				float64(playerMaxHP(m.RecommendedLevel))/float64(row.TakenDamage))
		}

		fmt.Printf("%-5s %-20s %3d %6d %5d %6d  %-34s %7s %6s %6d %5s\n",
			m.IndexNumber, m.Name, m.RecommendedLevel, m.HitPoint, m.Attack, m.ExperiencePoint,
			row.Loadout, dmgText, killText, row.TakenDamage, deathText)

		totalExp += m.ExperiencePoint
	}

	fmt.Println(strings.Repeat("-", 126))
	fmt.Printf("全 %d 体を 1 周した場合の累計 EXP: %d → Lv%d（最大HP %d）\n",
		len(monsters), totalExp, levelFromExp(totalExp), playerMaxHP(levelFromExp(totalExp)))
	fmt.Println()
}

// ============================================================================
// バトル計算のミラー（internal/domain/battle と同じ式、乱数は 1.0 固定）
// ============================================================================

func playerDamage(w WeaponSeed, m MonsterSeed, level int) int {
	return playerDamageWithDebuffs(w, m, level, nil)
}

func playerDamageWithDebuffs(w WeaponSeed, m MonsterSeed, level int, debuffs map[string]float64) int {
	elementAttack := 1.0
	if w.ElementAttack != nil {
		elementAttack = float64(*w.ElementAttack)
	}
	// デバフは耐性倍率への加算。simulator.go の calcPlayerDamage と同じ式にすること。
	physics := float64(w.PhysicsAttack) *
		(1 - resistanceOf(m.Res, w.PhysicsType) + debuffs[w.PhysicsType])
	element := elementAttack *
		(1 - resistanceOf(m.Res, w.ElementType) + debuffs[w.ElementType])
	base := physics * element

	sign := 1.0
	if base < 0 {
		sign = -1.0
	}
	levelFactor := 0.8 + math.Sqrt(float64(level))/5.0
	return int(math.Floor(math.Sqrt(math.Abs(base)) * levelFactor * sign))
}

func monsterDamage(m MonsterSeed, level int) int {
	return int(math.Floor(float64(m.Attack) / (1.0 + math.Sqrt(float64(level))/1.7)))
}

func playerMaxHP(level int) int { return 100 + 8*(level-1) }

func levelFromExp(exp int) int {
	if exp <= 0 {
		return 1
	}
	const baseExp = 19.0
	const rateOfIncrease = 0.067
	level := int(math.Floor(1 + math.Log(1+float64(exp)*rateOfIncrease/baseExp)/math.Log(1+rateOfIncrease)))
	if level < 1 {
		return 1
	}
	return level
}

// bestDebuffComboAgainst は「武器 + その武器の攻撃種別/属性に噛み合うデバフアイテム」の
// 組み合わせのうち、最もダメージの出るものを返す。耐性を全部 1.0 で固めた
// 謎解きモンスターが本当に攻略可能かを確かめるために使う。
func bestDebuffComboAgainst(m MonsterSeed, candidates []WeaponSeed, pouch []ItemSeed, level int) (*WeaponSeed, []ItemSeed, int) {
	var (
		bestWeapon *WeaponSeed
		bestItems  []ItemSeed
		bestDamage int
	)
	for i := range candidates {
		w := candidates[i]
		// 攻撃種別・属性それぞれについて、最も耐性を削れるデバフを 1 つずつ選ぶ。
		physicsDebuff := strongestDebuff(pouch, w.PhysicsType)
		elementDebuff := strongestDebuff(pouch, w.ElementType)

		debuffs := map[string]float64{}
		var used []ItemSeed
		for _, it := range []*ItemSeed{physicsDebuff, elementDebuff} {
			if it != nil {
				debuffs[it.Target] += it.Rate
				used = append(used, *it)
			}
		}
		if d := playerDamageWithDebuffs(w, m, level, debuffs); d > bestDamage {
			bestWeapon, bestItems, bestDamage = &candidates[i], used, d
		}
	}
	return bestWeapon, bestItems, bestDamage
}

func strongestDebuff(pouch []ItemSeed, target string) *ItemSeed {
	var best *ItemSeed
	for i := range pouch {
		it := &pouch[i]
		if it.EffectType == EffectDebuff && it.Target == target {
			if best == nil || it.Rate > best.Rate {
				best = it
			}
		}
	}
	return best
}

// bestWeaponAgainst は候補のうち最もダメージの出る武器を返す。
// 全てダメージ 0 以下（＝無効か吸収しかしない）なら nil。
func bestWeaponAgainst(m MonsterSeed, candidates []WeaponSeed) *WeaponSeed {
	var best *WeaponSeed
	bestDamage := 0
	for i := range candidates {
		if d := playerDamage(candidates[i], m, m.RecommendedLevel); d > bestDamage {
			bestDamage = d
			best = &candidates[i]
		}
	}
	return best
}

// indexNumberPattern は図鑑番号の形式。武器・アイテム・魔物とも 4 桁の数字。
// 魔物だけは上位 3 桁が種族、末尾 1 桁が歪みの深さという意味を持つ。
var indexNumberPattern = regexp.MustCompile(`^\d{4}$`)

// openCombos はそのモンスターにダメージが通る「攻撃種別 × 属性」を返す。
// 物理耐性・属性耐性がどちらも 1.0 未満であることが条件。
func openCombos(m MonsterSeed) [][2]string {
	var out [][2]string
	for _, p := range []string{Slash, Blow, Shoot} {
		if resistanceOf(m.Res, p) >= 1.0 {
			continue
		}
		for _, e := range []string{Neutral, Flame, Water, Wood, Shine, Dark} {
			if resistanceOf(m.Res, e) >= 1.0 {
				continue
			}
			out = append(out, [2]string{p, e})
		}
	}
	return out
}

// firstAreaMustBeBeatableBareHanded は最初のエリアの不変条件を検査する。
//
// 文化祭では QR コードを会場にばらばらに配置するので、プレイヤーがどの個体に
// 最初に出会うかを制御できない。最初の 1 体が素手で倒せないと、その参加者は
// 二度と先へ進めない。したがって最初のエリアの全個体について
//
//  1. 耐性が 0 以下（等倍か弱点）だけであること
//     → 素手（打 × 無）の係数が必ず正になり、ダメージが通る
//  2. 素手で削り切るターン数が、倒されるターン数より十分に短いこと
//
// を満たす必要がある。ここは遊べるかどうかに直結するので、
// バランス調整の都合で緩めてはいけない。
func firstAreaMustBeBeatableBareHanded() []string {
	var problems []string
	if len(monsters) == 0 {
		return problems
	}
	firstArea := monsters[0].Area

	for _, m := range monsters {
		if m.Area != firstArea {
			break
		}

		for axis, v := range resistanceMap(m.Res) {
			if v > 0 {
				problems = append(problems, fmt.Sprintf(
					"%s（%s）の %s 耐性が %+.1f。最初のエリアは素手が必ず通るよう耐性 0 以下だけにすること",
					m.Name, firstArea, axis, v))
			}
		}

		dmg := playerDamage(bareHands, m, m.RecommendedLevel)
		if dmg <= 0 {
			problems = append(problems, fmt.Sprintf(
				"%s（%s）に素手のダメージが通らない", m.Name, firstArea))
			continue
		}
		taken := monsterDamage(m, m.RecommendedLevel)
		if taken <= 0 {
			continue
		}
		killTurns := float64(m.HitPoint) / float64(dmg)
		deathTurns := float64(playerMaxHP(m.RecommendedLevel)) / float64(taken)
		// 2 倍の余裕を要求する。回復アイテムなしで確実に勝ち切れる水準。
		if killTurns*2 > deathTurns {
			problems = append(problems, fmt.Sprintf(
				"%s（%s）は素手 Lv%d で余裕がない（撃破 %.1f ターン / 被撃破 %.1f ターン）",
				m.Name, firstArea, m.RecommendedLevel, killTurns, deathTurns))
		}
	}
	return problems
}

// combosOf は武器の集合が持つ「攻撃種別 × 属性」を重複なく返す。
func combosOf(ws []WeaponSeed) [][2]string {
	seen := map[string]bool{}
	var out [][2]string
	for _, w := range ws {
		key := w.PhysicsType + "/" + w.ElementType
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, [2]string{w.PhysicsType, w.ElementType})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i][0] != out[j][0] {
			return out[i][0] < out[j][0]
		}
		return out[i][1] < out[j][1]
	})
	return out
}

func describeCombos(combos [][2]string) string {
	if len(combos) == 0 {
		return "なし"
	}
	parts := make([]string, 0, len(combos))
	for _, c := range combos {
		parts = append(parts, c[0]+"×"+c[1])
	}
	return strings.Join(parts, ", ")
}

func resistanceOf(r Resistance, typ string) float64 {
	switch typ {
	case Slash:
		return r.Slash
	case Blow:
		return r.Blow
	case Shoot:
		return r.Shoot
	case Neutral:
		return r.Neutral
	case Flame:
		return r.Flame
	case Water:
		return r.Water
	case Wood:
		return r.Wood
	case Shine:
		return r.Shine
	case Dark:
		return r.Dark
	}
	return 1.0
}

func resistanceMap(r Resistance) map[string]float64 {
	return map[string]float64{
		Slash: r.Slash, Blow: r.Blow, Shoot: r.Shoot,
		Neutral: r.Neutral, Flame: r.Flame, Water: r.Water,
		Wood: r.Wood, Shine: r.Shine, Dark: r.Dark,
	}
}

func weaponByID(id int64) *WeaponSeed {
	for i := range weapons {
		if weapons[i].ID == id {
			return &weapons[i]
		}
	}
	return nil
}

func itemByID(id int64) *ItemSeed {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
	}
	return nil
}

// ============================================================================
// 環境変数・DB接続（cmd/create-admin と同じ手順）
// ============================================================================

// loadEnv は KEY=VALUE 形式の .env を読んで os.Setenv する。
// 既に環境変数がセットされている場合は上書きしない。
func loadEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
	return scanner.Err()
}

func connectDB() (*sql.DB, error) {
	host := requireEnv("DB_HOST")
	port := envOrDefault("DB_PORT", "3306")
	user := requireEnv("DB_USER")
	pass := requireEnv("DB_PASSWORD")
	name := requireEnv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC", user, pass, host, port, name)
	return sql.Open("mysql", dsn)
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "env var %s is required\n", key)
		os.Exit(1)
	}
	return v
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
