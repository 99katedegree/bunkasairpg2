package main

// ============================================================================
// アイテムマスタ — 全 44 個
//
// buff は  (1 + バフ率)      を攻撃側に掛ける。素直な火力アップ＝自分の型を濃くする。
// debuff は(1 - 耐性 + デバフ率) として耐性倍率に加算される＝相手の型を薄くする。
//
// この非対称が攻略の肝になっている。
//   ・耐性 0（等倍）の相手にデバフを使っても倍率は 1.0 のままで意味がない
//   ・耐性 1.0（無効）の相手でも、デバフ 0.5 を入れれば係数 0.5 まで戻る
//   ・耐性 1.3（吸収）の相手でも、デバフ 1.0 を入れれば係数 0.7 の正ダメージになる
// 「全軸 1.0 でデバフ必須」のモンスターはここを踏まえて配置している。
//
// rate はシミュレータ側で 0.1 刻みに切り捨てられる（math.Floor(rate*10)/10）ので、
// 0.1 刻みで指定すること。
// ============================================================================

var items = []ItemSeed{

	// --- 回復 ---
	{ID: 1, IndexNumber: "0001", Name: "タールの雫", EffectType: EffectHeal, Amount: 40},
	{ID: 2, IndexNumber: "0002", Name: "ホルの実", EffectType: EffectHeal, Amount: 90},
	{ID: 3, IndexNumber: "0003", Name: "ガスの蜜", EffectType: EffectHeal, Amount: 150},
	{ID: 4, IndexNumber: "0004", Name: "ネムの血", EffectType: EffectHeal, Amount: 230},
	{ID: 5, IndexNumber: "0005", Name: "マロの涙", EffectType: EffectHeal, Amount: 330},
	{ID: 6, IndexNumber: "0006", Name: "ホスの髄", EffectType: EffectHeal, Amount: 450},
	{ID: 7, IndexNumber: "0007", Name: "タルの心", EffectType: EffectHeal, Amount: 620},
	{ID: 8, IndexNumber: "0008", Name: "ダムの魂", EffectType: EffectHeal, Amount: 810},

	// --- バフ（火力アップ） ---
	{ID: 9, IndexNumber: "0009", Name: "モコの牙", EffectType: EffectBuff, Rate: 0.3, Target: Slash},
	{ID: 10, IndexNumber: "0010", Name: "ケルの大牙", EffectType: EffectBuff, Rate: 0.8, Target: Slash},
	{ID: 11, IndexNumber: "0011", Name: "バムの拳", EffectType: EffectBuff, Rate: 0.3, Target: Blow},
	{ID: 12, IndexNumber: "0012", Name: "ジャルの大拳", EffectType: EffectBuff, Rate: 0.8, Target: Blow},
	{ID: 13, IndexNumber: "0013", Name: "グルの爪", EffectType: EffectBuff, Rate: 0.3, Target: Shoot},
	{ID: 14, IndexNumber: "0014", Name: "ビルの大爪", EffectType: EffectBuff, Rate: 0.8, Target: Shoot},
	{ID: 15, IndexNumber: "0015", Name: "ポメの石", EffectType: EffectBuff, Rate: 0.3, Target: Neutral},
	{ID: 16, IndexNumber: "0016", Name: "ナルの大石", EffectType: EffectBuff, Rate: 0.8, Target: Neutral},
	{ID: 17, IndexNumber: "0017", Name: "シムの焔", EffectType: EffectBuff, Rate: 0.3, Target: Flame},
	{ID: 18, IndexNumber: "0018", Name: "ラスの大焔", EffectType: EffectBuff, Rate: 0.8, Target: Flame},
	{ID: 19, IndexNumber: "0019", Name: "ハネの鱗", EffectType: EffectBuff, Rate: 0.3, Target: Water},
	{ID: 20, IndexNumber: "0020", Name: "クルの大鱗", EffectType: EffectBuff, Rate: 0.8, Target: Water},
	{ID: 21, IndexNumber: "0021", Name: "シムの蔓", EffectType: EffectBuff, Rate: 0.3, Target: Wood},
	{ID: 22, IndexNumber: "0022", Name: "レドの大蔓", EffectType: EffectBuff, Rate: 0.8, Target: Wood},
	{ID: 23, IndexNumber: "0023", Name: "マスの羽", EffectType: EffectBuff, Rate: 0.3, Target: Shine},
	{ID: 24, IndexNumber: "0024", Name: "ミルの大羽", EffectType: EffectBuff, Rate: 0.8, Target: Shine},
	{ID: 25, IndexNumber: "0025", Name: "シムの影", EffectType: EffectBuff, Rate: 0.3, Target: Dark},
	{ID: 26, IndexNumber: "0026", Name: "レトの大影", EffectType: EffectBuff, Rate: 0.8, Target: Dark},

	// --- デバフ（耐性削り） ---
	{ID: 27, IndexNumber: "0027", Name: "ヌノの牙砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Slash},
	{ID: 28, IndexNumber: "0028", Name: "リガの牙断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Slash},
	{ID: 29, IndexNumber: "0029", Name: "ネスの拳砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Blow},
	{ID: 30, IndexNumber: "0030", Name: "ガネの拳断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Blow},
	{ID: 31, IndexNumber: "0031", Name: "ゴルの爪砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Shoot},
	{ID: 32, IndexNumber: "0032", Name: "ラムの爪断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Shoot},
	{ID: 33, IndexNumber: "0033", Name: "ビルの石砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Neutral},
	{ID: 34, IndexNumber: "0034", Name: "ミルの石断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Neutral},
	{ID: 35, IndexNumber: "0035", Name: "ロガの焔砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Flame},
	{ID: 36, IndexNumber: "0036", Name: "ズケの焔断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Flame},
	{ID: 37, IndexNumber: "0037", Name: "リドの鱗砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Water},
	{ID: 38, IndexNumber: "0038", Name: "キムの鱗断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Water},
	{ID: 39, IndexNumber: "0039", Name: "キスの蔓砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Wood},
	{ID: 40, IndexNumber: "0040", Name: "ズラの蔓断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Wood},
	{ID: 41, IndexNumber: "0041", Name: "リメの羽砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Shine},
	{ID: 42, IndexNumber: "0042", Name: "バミの羽断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Shine},
	{ID: 43, IndexNumber: "0043", Name: "ホルの影砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Dark},
	{ID: 44, IndexNumber: "0044", Name: "ジャルの影断ち", EffectType: EffectDebuff, Rate: 1.0, Target: Dark},
}
