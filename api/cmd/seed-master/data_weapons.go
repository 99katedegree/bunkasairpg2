package main

// ============================================================================
// 武器マスタ — 全 126 本
//
// このファイルは cmd/seed-master の投入データ。手で書き換えてよい。
//
// ダメージは √(物理攻撃力 × 属性攻撃力) に比例するので、実効火力は積で決まる。
// 物理寄り／属性寄りの振り分けはカードの見た目とキャラ付けのためのもので、
// 火力そのものには影響しない。tier ごとの積は以下。
//
//   T1 3,600 / T2 9,500 / T3 25,000 / T4 50,000 / T5 95,000 / T6 160,000
//   T7 260,000 / T8 420,000 / T9 650,000 / T10 950,000 / T11 1,400,000 / T12 2,000,000
//
// 属性はわざと歯抜けにしてある。tier をまたいで初めて解禁される組み合わせがあり、
// 「その属性の武器を持っていないと手が出ない」モンスターの前提になっている。
// ただし最終ティア T12 だけは 3 攻撃種別 × 6 属性の 18 通りを全て揃えてあり、
// そこで歯抜けが解消する。モンスター側のギミックは「その時点で存在する
// 組み合わせ」からしか作らないこと（validate の組み合わせ検証が守っている）。
// ============================================================================

// bareHands は DB に保存しない初期装備。
// api/internal/adapter/handler/me.go が装備なしのときに返すハードコード値と同じ内容で、
// バランス検証の出発点として使う。ここを変えたら me.go も合わせること。
var bareHands = WeaponSeed{
	ID: 0, IndexNumber: "0000", Name: "素手", Category: "拳",
	PhysicsAttack: 10, ElementAttack: nil, PhysicsType: Blow, ElementType: Neutral,
	Note: "DB には入れない。装備なしのとき me.go が返す値のミラー",
}

var weapons = []WeaponSeed{

	// --- T1 陽だまりの草原 (Lv1-4) 積 3,600 ---
	{ID: 1, IndexNumber: "0001", Name: "モスソード", Category: "剣", PhysicsAttack: 105, ElementAttack: ea(35), PhysicsType: Slash, ElementType: Neutral},
	{ID: 2, IndexNumber: "0002", Name: "モスハンマー", Category: "ハンマー", PhysicsAttack: 170, ElementAttack: ea(21), PhysicsType: Blow, ElementType: Neutral},
	{ID: 3, IndexNumber: "0003", Name: "ロスボウ", Category: "弓", PhysicsAttack: 66, ElementAttack: ea(55), PhysicsType: Shoot, ElementType: Wood},
	{ID: 4, IndexNumber: "0004", Name: "ロスダガー", Category: "短剣", PhysicsAttack: 66, ElementAttack: ea(55), PhysicsType: Slash, ElementType: Flame},
	{ID: 5, IndexNumber: "0005", Name: "フロナックル", Category: "拳", PhysicsAttack: 105, ElementAttack: ea(35), PhysicsType: Blow, ElementType: Wood},
	{ID: 6, IndexNumber: "0006", Name: "ヘルガン", Category: "銃", PhysicsAttack: 105, ElementAttack: ea(35), PhysicsType: Shoot, ElementType: Neutral},

	// --- T2 囁きの森 (Lv5-9) 積 9,500 ---
	{ID: 7, IndexNumber: "0007", Name: "ネブソード", Category: "剣", PhysicsAttack: 170, ElementAttack: ea(56), PhysicsType: Slash, ElementType: Neutral},
	{ID: 8, IndexNumber: "0008", Name: "シルツインソード", Category: "双剣", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Slash, ElementType: Water},
	{ID: 9, IndexNumber: "0009", Name: "ウンアックス", Category: "斧", PhysicsAttack: 275, ElementAttack: ea(34), PhysicsType: Slash, ElementType: Flame},
	{ID: 10, IndexNumber: "0010", Name: "モスメイス", Category: "メイス", PhysicsAttack: 170, ElementAttack: ea(56), PhysicsType: Blow, ElementType: Neutral},
	{ID: 11, IndexNumber: "0011", Name: "プルナックル", Category: "拳", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Blow, ElementType: Shine},
	{ID: 12, IndexNumber: "0012", Name: "プルガン", Category: "銃", PhysicsAttack: 170, ElementAttack: ea(56), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 13, IndexNumber: "0013", Name: "ヴルボウ", Category: "弓", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Shoot, ElementType: Wood},
	{ID: 14, IndexNumber: "0014", Name: "ノククロー", Category: "爪", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Slash, ElementType: Dark},

	// --- T3 鉄錆の廃坑 (Lv10-14) 積 25,000 ---
	{ID: 15, IndexNumber: "0015", Name: "セムソード", Category: "剣", PhysicsAttack: 275, ElementAttack: ea(91), PhysicsType: Slash, ElementType: Flame},
	{ID: 16, IndexNumber: "0016", Name: "ケーボウ", Category: "弓", PhysicsAttack: 175, ElementAttack: ea(145), PhysicsType: Shoot, ElementType: Water},
	{ID: 17, IndexNumber: "0017", Name: "ソルアックス", Category: "斧", PhysicsAttack: 445, ElementAttack: ea(56), PhysicsType: Slash, ElementType: Wood},
	{ID: 18, IndexNumber: "0018", Name: "パルメイス", Category: "メイス", PhysicsAttack: 275, ElementAttack: ea(91), PhysicsType: Blow, ElementType: Shine},
	{ID: 19, IndexNumber: "0019", Name: "ムスハンマー", Category: "ハンマー", PhysicsAttack: 445, ElementAttack: ea(56), PhysicsType: Blow, ElementType: Neutral},
	{ID: 20, IndexNumber: "0020", Name: "ヘルキャノン", Category: "重銃", PhysicsAttack: 445, ElementAttack: ea(56), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 21, IndexNumber: "0021", Name: "ロスツインガン", Category: "両手銃", PhysicsAttack: 175, ElementAttack: ea(145), PhysicsType: Shoot, ElementType: Flame},
	{ID: 22, IndexNumber: "0022", Name: "モスダガー", Category: "短剣", PhysicsAttack: 175, ElementAttack: ea(145), PhysicsType: Slash, ElementType: Water},
	{ID: 23, IndexNumber: "0023", Name: "ミカフレイル", Category: "モーニングスター", PhysicsAttack: 275, ElementAttack: ea(91), PhysicsType: Blow, ElementType: Dark},

	// --- T4 銀鏡の湖 (Lv15-19) 積 50,000 ---
	{ID: 24, IndexNumber: "0024", Name: "デンスソード", Category: "剣", PhysicsAttack: 385, ElementAttack: ea(130), PhysicsType: Slash, ElementType: Neutral},
	{ID: 25, IndexNumber: "0025", Name: "フェツインソード", Category: "双剣", PhysicsAttack: 245, ElementAttack: ea(205), PhysicsType: Slash, ElementType: Water},
	{ID: 26, IndexNumber: "0026", Name: "ノクハンマー", Category: "ハンマー", PhysicsAttack: 630, ElementAttack: ea(79), PhysicsType: Blow, ElementType: Neutral},
	{ID: 27, IndexNumber: "0027", Name: "リムメイス", Category: "メイス", PhysicsAttack: 385, ElementAttack: ea(130), PhysicsType: Blow, ElementType: Wood},
	{ID: 28, IndexNumber: "0028", Name: "フミガン", Category: "銃", PhysicsAttack: 385, ElementAttack: ea(130), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 29, IndexNumber: "0029", Name: "リケクロー", Category: "爪", PhysicsAttack: 245, ElementAttack: ea(205), PhysicsType: Slash, ElementType: Dark},
	{ID: 30, IndexNumber: "0030", Name: "ムルボウ", Category: "弓", PhysicsAttack: 245, ElementAttack: ea(205), PhysicsType: Shoot, ElementType: Shine},
	{ID: 31, IndexNumber: "0031", Name: "ノドアックス", Category: "斧", PhysicsAttack: 630, ElementAttack: ea(79), PhysicsType: Slash, ElementType: Flame},
	{ID: 32, IndexNumber: "0032", Name: "ウンキャノン", Category: "重銃", PhysicsAttack: 630, ElementAttack: ea(79), PhysicsType: Shoot, ElementType: Neutral},

	// --- T5 朽ちた王城 (Lv20-26) 積 95,000 ---
	{ID: 33, IndexNumber: "0033", Name: "ヘルソード", Category: "剣", PhysicsAttack: 535, ElementAttack: ea(180), PhysicsType: Slash, ElementType: Neutral},
	{ID: 34, IndexNumber: "0034", Name: "グラツインソード", Category: "双剣", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Slash, ElementType: Flame},
	{ID: 35, IndexNumber: "0035", Name: "セムクロー", Category: "爪", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Slash, ElementType: Dark},
	{ID: 36, IndexNumber: "0036", Name: "スコハンマー", Category: "ハンマー", PhysicsAttack: 870, ElementAttack: ea(110), PhysicsType: Blow, ElementType: Neutral},
	{ID: 37, IndexNumber: "0037", Name: "ロスフレイル", Category: "モーニングスター", PhysicsAttack: 535, ElementAttack: ea(180), PhysicsType: Blow, ElementType: Shine},
	{ID: 38, IndexNumber: "0038", Name: "ケーメイス", Category: "メイス", PhysicsAttack: 535, ElementAttack: ea(180), PhysicsType: Blow, ElementType: Wood},
	{ID: 39, IndexNumber: "0039", Name: "アリキャノン", Category: "重銃", PhysicsAttack: 870, ElementAttack: ea(110), PhysicsType: Shoot, ElementType: Water},
	{ID: 40, IndexNumber: "0040", Name: "スルツインガン", Category: "両手銃", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Shoot, ElementType: Shine},
	{ID: 41, IndexNumber: "0041", Name: "ムスボウ", Category: "弓", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 42, IndexNumber: "0042", Name: "フォダガー", Category: "短剣", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Slash, ElementType: Water},

	// --- T6 煉獄の火山 (Lv27-33) 積 160,000 ---
	{ID: 43, IndexNumber: "0043", Name: "ノクソード", Category: "剣", PhysicsAttack: 695, ElementAttack: ea(230), PhysicsType: Slash, ElementType: Flame},
	{ID: 44, IndexNumber: "0044", Name: "アルハンマー", Category: "ハンマー", PhysicsAttack: 1130, ElementAttack: ea(140), PhysicsType: Blow, ElementType: Neutral},
	{ID: 45, IndexNumber: "0045", Name: "コルナックル", Category: "拳", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Blow, ElementType: Flame},
	{ID: 46, IndexNumber: "0046", Name: "ムルアックス", Category: "斧", PhysicsAttack: 1130, ElementAttack: ea(140), PhysicsType: Slash, ElementType: Flame},
	{ID: 47, IndexNumber: "0047", Name: "ネブツインソード", Category: "双剣", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Slash, ElementType: Flame},
	{ID: 48, IndexNumber: "0048", Name: "ステガン", Category: "銃", PhysicsAttack: 695, ElementAttack: ea(230), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 49, IndexNumber: "0049", Name: "マドツインガン", Category: "両手銃", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Shoot, ElementType: Flame},
	{ID: 50, IndexNumber: "0050", Name: "シルクロー", Category: "爪", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Slash, ElementType: Shine},
	{ID: 51, IndexNumber: "0051", Name: "ローフレイル", Category: "モーニングスター", PhysicsAttack: 695, ElementAttack: ea(230), PhysicsType: Blow, ElementType: Dark},
	{ID: 52, IndexNumber: "0052", Name: "フミキャノン", Category: "重銃", PhysicsAttack: 1130, ElementAttack: ea(140), PhysicsType: Shoot, ElementType: Flame},

	// --- T7 氷晶の尖塔 (Lv34-41) 積 260,000 ---
	{ID: 53, IndexNumber: "0053", Name: "ローソード", Category: "剣", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Slash, ElementType: Water},
	{ID: 54, IndexNumber: "0054", Name: "クラツインソード", Category: "双剣", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Water},
	{ID: 55, IndexNumber: "0055", Name: "フォハンマー", Category: "ハンマー", PhysicsAttack: 1440, ElementAttack: ea(180), PhysicsType: Blow, ElementType: Water},
	{ID: 56, IndexNumber: "0056", Name: "ポルメイス", Category: "メイス", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Blow, ElementType: Water},
	{ID: 57, IndexNumber: "0057", Name: "ルスボウ", Category: "弓", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Shoot, ElementType: Water},
	{ID: 58, IndexNumber: "0058", Name: "ロスアックス", Category: "斧", PhysicsAttack: 1440, ElementAttack: ea(180), PhysicsType: Slash, ElementType: Wood},
	{ID: 59, IndexNumber: "0059", Name: "ラクキャノン", Category: "重銃", PhysicsAttack: 1440, ElementAttack: ea(180), PhysicsType: Shoot, ElementType: Shine},
	{ID: 60, IndexNumber: "0060", Name: "スクダガー", Category: "短剣", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Water},
	{ID: 61, IndexNumber: "0061", Name: "カテクロー", Category: "爪", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Water},
	{ID: 62, IndexNumber: "0062", Name: "レムナックル", Category: "拳", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Blow, ElementType: Neutral},
	{ID: 63, IndexNumber: "0063", Name: "プルフレイル", Category: "モーニングスター", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Blow, ElementType: Dark},

	// --- T8 世界樹の樹海 (Lv42-49) 積 420,000 ---
	{ID: 64, IndexNumber: "0064", Name: "モスボウ", Category: "弓", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Shoot, ElementType: Wood},
	{ID: 65, IndexNumber: "0065", Name: "ニグアックス", Category: "斧", PhysicsAttack: 1830, ElementAttack: ea(230), PhysicsType: Slash, ElementType: Wood},
	{ID: 66, IndexNumber: "0066", Name: "ラヴァメイス", Category: "メイス", PhysicsAttack: 1120, ElementAttack: ea(375), PhysicsType: Blow, ElementType: Wood},
	{ID: 67, IndexNumber: "0067", Name: "フーソード", Category: "剣", PhysicsAttack: 1120, ElementAttack: ea(375), PhysicsType: Slash, ElementType: Dark},
	{ID: 68, IndexNumber: "0068", Name: "パルツインガン", Category: "両手銃", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Shoot, ElementType: Shine},
	{ID: 69, IndexNumber: "0069", Name: "ノドキャノン", Category: "重銃", PhysicsAttack: 1830, ElementAttack: ea(230), PhysicsType: Shoot, ElementType: Wood},
	{ID: 70, IndexNumber: "0070", Name: "スペクロー", Category: "爪", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Slash, ElementType: Dark},
	{ID: 71, IndexNumber: "0071", Name: "イグハンマー", Category: "ハンマー", PhysicsAttack: 1830, ElementAttack: ea(230), PhysicsType: Blow, ElementType: Wood},
	{ID: 72, IndexNumber: "0072", Name: "ルナガン", Category: "銃", PhysicsAttack: 1120, ElementAttack: ea(375), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 73, IndexNumber: "0073", Name: "カルツインソード", Category: "双剣", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Slash, ElementType: Wood},
	{ID: 74, IndexNumber: "0074", Name: "テネナックル", Category: "拳", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Blow, ElementType: Shine},

	// --- T9 深淵の迷宮 (Lv50-58) 積 650,000 ---
	{ID: 75, IndexNumber: "0075", Name: "フロソード", Category: "剣", PhysicsAttack: 1400, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Neutral},
	{ID: 76, IndexNumber: "0076", Name: "プルハンマー", Category: "ハンマー", PhysicsAttack: 2280, ElementAttack: ea(285), PhysicsType: Blow, ElementType: Neutral},
	{ID: 77, IndexNumber: "0077", Name: "アルツインソード", Category: "双剣", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Slash, ElementType: Neutral},
	{ID: 78, IndexNumber: "0078", Name: "ニヴァボウ", Category: "弓", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 79, IndexNumber: "0079", Name: "ムルキャノン", Category: "重銃", PhysicsAttack: 2280, ElementAttack: ea(285), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 80, IndexNumber: "0080", Name: "ゲルクロー", Category: "爪", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Slash, ElementType: Dark},
	{ID: 81, IndexNumber: "0081", Name: "ロスメイス", Category: "メイス", PhysicsAttack: 1400, ElementAttack: ea(465), PhysicsType: Blow, ElementType: Neutral},
	{ID: 82, IndexNumber: "0082", Name: "シネダガー", Category: "短剣", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Slash, ElementType: Neutral},
	{ID: 83, IndexNumber: "0083", Name: "フルツインガン", Category: "両手銃", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Shoot, ElementType: Dark},
	{ID: 84, IndexNumber: "0084", Name: "レクスナックル", Category: "拳", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Blow, ElementType: Dark},
	{ID: 85, IndexNumber: "0085", Name: "スコフレイル", Category: "モーニングスター", PhysicsAttack: 1400, ElementAttack: ea(465), PhysicsType: Blow, ElementType: Neutral},

	// --- T10 竜骨の墓所 (Lv59-68) 積 950,000 ---
	{ID: 86, IndexNumber: "0086", Name: "オッシソード", Category: "剣", PhysicsAttack: 1690, ElementAttack: ea(565), PhysicsType: Slash, ElementType: Shine},
	{ID: 87, IndexNumber: "0087", Name: "カタハンマー", Category: "ハンマー", PhysicsAttack: 2760, ElementAttack: ea(345), PhysicsType: Blow, ElementType: Dark},
	{ID: 88, IndexNumber: "0088", Name: "リトツインソード", Category: "双剣", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Slash, ElementType: Shine},
	{ID: 89, IndexNumber: "0089", Name: "ネクキャノン", Category: "重銃", PhysicsAttack: 2760, ElementAttack: ea(345), PhysicsType: Shoot, ElementType: Dark},
	{ID: 90, IndexNumber: "0090", Name: "ドラボウ", Category: "弓", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Shoot, ElementType: Shine},
	{ID: 91, IndexNumber: "0091", Name: "フネアックス", Category: "斧", PhysicsAttack: 2760, ElementAttack: ea(345), PhysicsType: Slash, ElementType: Shine},
	{ID: 92, IndexNumber: "0092", Name: "クラメイス", Category: "メイス", PhysicsAttack: 1690, ElementAttack: ea(565), PhysicsType: Blow, ElementType: Dark},
	{ID: 93, IndexNumber: "0093", Name: "コスダガー", Category: "短剣", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Slash, ElementType: Shine},
	{ID: 94, IndexNumber: "0094", Name: "トゥムツインガン", Category: "両手銃", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Shoot, ElementType: Shine},
	{ID: 95, IndexNumber: "0095", Name: "エクナックル", Category: "拳", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Blow, ElementType: Dark},
	{ID: 96, IndexNumber: "0096", Name: "オッサフレイル", Category: "モーニングスター", PhysicsAttack: 1690, ElementAttack: ea(565), PhysicsType: Blow, ElementType: Shine},

	// --- T11 六輝の神殿 (Lv69-79) 積 1,400,000 ---
	{ID: 97, IndexNumber: "0097", Name: "クスナックル", Category: "拳", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Blow, ElementType: Shine},
	{ID: 98, IndexNumber: "0098", Name: "ヌーソード", Category: "剣", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Slash, ElementType: Shine},
	{ID: 99, IndexNumber: "0099", Name: "オラハンマー", Category: "ハンマー", PhysicsAttack: 3350, ElementAttack: ea(420), PhysicsType: Blow, ElementType: Shine},
	{ID: 100, IndexNumber: "0100", Name: "サンボウ", Category: "弓", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Shoot, ElementType: Wood},
	{ID: 101, IndexNumber: "0101", Name: "アドアックス", Category: "斧", PhysicsAttack: 3350, ElementAttack: ea(420), PhysicsType: Slash, ElementType: Flame},
	{ID: 102, IndexNumber: "0102", Name: "ウェツインソード", Category: "双剣", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Slash, ElementType: Water},
	{ID: 103, IndexNumber: "0103", Name: "ルーキャノン", Category: "重銃", PhysicsAttack: 3350, ElementAttack: ea(420), PhysicsType: Shoot, ElementType: Shine},
	{ID: 104, IndexNumber: "0104", Name: "タリスメイス", Category: "メイス", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Blow, ElementType: Neutral},
	{ID: 105, IndexNumber: "0105", Name: "サケツインガン", Category: "両手銃", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 106, IndexNumber: "0106", Name: "ウォークロー", Category: "爪", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Slash, ElementType: Wood},
	{ID: 107, IndexNumber: "0107", Name: "カンガン", Category: "銃", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Shoot, ElementType: Shine},
	{ID: 108, IndexNumber: "0108", Name: "アシフレイル", Category: "モーニングスター", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Blow, ElementType: Wood},

	// --- T12 裂けた空 (Lv80-90) 積 2,000,000 ---
	//
	// 最終ティアだけは 3 攻撃種別 × 6 属性の 18 通りを全て揃えてある。
	// ここまで来ればどんなギミックにも必ず answer が存在する、という到達点。
	// 積は全て 2,000,000 ちょうどで揃えてあり、ダメージは √積 に比例するので
	// 18 本の実効火力は完全に同じ。物理寄り／属性寄りの振り分けは見た目だけの差。
	//
	// 積を崩すと「どれか 1 本だけが正解」になってこの帯の設計が壊れるので、
	// 数値を触るときは physics_attack × element_attack = 2,000,000 を必ず保つこと。

	// 無属性
	{ID: 109, IndexNumber: "0109", Name: "リマソード", Category: "剣", PhysicsAttack: 2000, ElementAttack: ea(1000), PhysicsType: Slash, ElementType: Neutral},
	{ID: 110, IndexNumber: "0110", Name: "ルイハンマー", Category: "ハンマー", PhysicsAttack: 4000, ElementAttack: ea(500), PhysicsType: Blow, ElementType: Neutral},
	{ID: 111, IndexNumber: "0111", Name: "フィキャノン", Category: "重銃", PhysicsAttack: 2500, ElementAttack: ea(800), PhysicsType: Shoot, ElementType: Neutral},

	// 火属性
	{ID: 112, IndexNumber: "0112", Name: "ピュアックス", Category: "斧", PhysicsAttack: 4000, ElementAttack: ea(500), PhysicsType: Slash, ElementType: Flame},
	{ID: 113, IndexNumber: "0113", Name: "デレナックル", Category: "拳", PhysicsAttack: 1600, ElementAttack: ea(1250), PhysicsType: Blow, ElementType: Flame},
	{ID: 114, IndexNumber: "0114", Name: "バラツインガン", Category: "両手銃", PhysicsAttack: 1250, ElementAttack: ea(1600), PhysicsType: Shoot, ElementType: Flame},

	// 水属性
	{ID: 115, IndexNumber: "0115", Name: "マルダガー", Category: "短剣", PhysicsAttack: 1250, ElementAttack: ea(1600), PhysicsType: Slash, ElementType: Water},
	{ID: 116, IndexNumber: "0116", Name: "フィッメイス", Category: "メイス", PhysicsAttack: 2500, ElementAttack: ea(800), PhysicsType: Blow, ElementType: Water},
	{ID: 117, IndexNumber: "0117", Name: "カラボウ", Category: "弓", PhysicsAttack: 1000, ElementAttack: ea(2000), PhysicsType: Shoot, ElementType: Water},

	// 木属性
	{ID: 118, IndexNumber: "0118", Name: "ヴァクロー", Category: "爪", PhysicsAttack: 2500, ElementAttack: ea(800), PhysicsType: Slash, ElementType: Wood},
	{ID: 119, IndexNumber: "0119", Name: "レディフレイル", Category: "モーニングスター", PhysicsAttack: 2000, ElementAttack: ea(1000), PhysicsType: Blow, ElementType: Wood},
	{ID: 120, IndexNumber: "0120", Name: "エクボウ", Category: "弓", PhysicsAttack: 1600, ElementAttack: ea(1250), PhysicsType: Shoot, ElementType: Wood},

	// 光属性
	{ID: 121, IndexNumber: "0121", Name: "トニツインソード", Category: "双剣", PhysicsAttack: 1600, ElementAttack: ea(1250), PhysicsType: Slash, ElementType: Shine},
	{ID: 122, IndexNumber: "0122", Name: "クラナックル", Category: "拳", PhysicsAttack: 1250, ElementAttack: ea(1600), PhysicsType: Blow, ElementType: Shine},
	{ID: 123, IndexNumber: "0123", Name: "インガン", Category: "銃", PhysicsAttack: 2000, ElementAttack: ea(1000), PhysicsType: Shoot, ElementType: Shine},

	// 闇属性
	{ID: 124, IndexNumber: "0124", Name: "キドラツインソード", Category: "双剣", PhysicsAttack: 1000, ElementAttack: ea(2000), PhysicsType: Slash, ElementType: Dark},
	{ID: 125, IndexNumber: "0125", Name: "シルフレイル", Category: "モーニングスター", PhysicsAttack: 2500, ElementAttack: ea(800), PhysicsType: Blow, ElementType: Dark},
	{ID: 126, IndexNumber: "0126", Name: "デスキャノン", Category: "重銃", PhysicsAttack: 4000, ElementAttack: ea(500), PhysicsType: Shoot, ElementType: Dark},
}
