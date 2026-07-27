package main

// ============================================================================
// 武器マスタ — 全 120 本
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
// ============================================================================

// bareHands は DB に保存しない初期装備。
// api/internal/adapter/handler/me.go が装備なしのときに返すハードコード値と同じ内容で、
// バランス検証の出発点として使う。ここを変えたら me.go も合わせること。
var bareHands = WeaponSeed{
	ID: 0, IndexNumber: "W000", Name: "素手", Category: "拳",
	PhysicsAttack: 10, ElementAttack: nil, PhysicsType: Blow, ElementType: Neutral,
	Note: "DB には入れない。装備なしのとき me.go が返す値のミラー",
}

var weapons = []WeaponSeed{

	// --- T1 始まりの草原 (Lv1-4) 積 3,600 ---
	{ID: 1, IndexNumber: "W001", Name: "木刀", Category: "剣", PhysicsAttack: 105, ElementAttack: ea(35), PhysicsType: Slash, ElementType: Neutral},
	{ID: 2, IndexNumber: "W002", Name: "石のハンマー", Category: "ハンマー", PhysicsAttack: 170, ElementAttack: ea(21), PhysicsType: Blow, ElementType: Neutral},
	{ID: 3, IndexNumber: "W003", Name: "狩人の弓", Category: "弓", PhysicsAttack: 66, ElementAttack: ea(55), PhysicsType: Shoot, ElementType: Wood},
	{ID: 4, IndexNumber: "W004", Name: "焚き火の短剣", Category: "短剣", PhysicsAttack: 66, ElementAttack: ea(55), PhysicsType: Slash, ElementType: Flame},
	{ID: 5, IndexNumber: "W005", Name: "蔦巻きの拳", Category: "拳", PhysicsAttack: 105, ElementAttack: ea(35), PhysicsType: Blow, ElementType: Wood},
	{ID: 6, IndexNumber: "W006", Name: "錆びた銃", Category: "銃", PhysicsAttack: 105, ElementAttack: ea(35), PhysicsType: Shoot, ElementType: Neutral},

	// --- T2 霧の森 (Lv5-9) 積 9,500 ---
	{ID: 7, IndexNumber: "W007", Name: "鉄の剣", Category: "剣", PhysicsAttack: 170, ElementAttack: ea(56), PhysicsType: Slash, ElementType: Neutral},
	{ID: 8, IndexNumber: "W008", Name: "潮鳴の双剣", Category: "双剣", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Slash, ElementType: Water},
	{ID: 9, IndexNumber: "W009", Name: "業火の斧", Category: "斧", PhysicsAttack: 275, ElementAttack: ea(34), PhysicsType: Slash, ElementType: Flame},
	{ID: 10, IndexNumber: "W010", Name: "猛牛メイス", Category: "メイス", PhysicsAttack: 170, ElementAttack: ea(56), PhysicsType: Blow, ElementType: Neutral},
	{ID: 11, IndexNumber: "W011", Name: "雷鳴の拳", Category: "拳", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Blow, ElementType: Shine},
	{ID: 12, IndexNumber: "W012", Name: "連装銃", Category: "銃", PhysicsAttack: 170, ElementAttack: ea(56), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 13, IndexNumber: "W013", Name: "森人の弓", Category: "弓", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Shoot, ElementType: Wood},
	{ID: 14, IndexNumber: "W014", Name: "影爪", Category: "爪", PhysicsAttack: 105, ElementAttack: ea(89), PhysicsType: Slash, ElementType: Dark},

	// --- T3 廃鉱山 (Lv10-14) 積 25,000 ---
	{ID: 15, IndexNumber: "W015", Name: "竜牙剣", Category: "剣", PhysicsAttack: 275, ElementAttack: ea(91), PhysicsType: Slash, ElementType: Flame},
	{ID: 16, IndexNumber: "W016", Name: "氷結弓", Category: "弓", PhysicsAttack: 175, ElementAttack: ea(145), PhysicsType: Shoot, ElementType: Water},
	{ID: 17, IndexNumber: "W017", Name: "大樹の斧", Category: "斧", PhysicsAttack: 445, ElementAttack: ea(56), PhysicsType: Slash, ElementType: Wood},
	{ID: 18, IndexNumber: "W018", Name: "聖光メイス", Category: "メイス", PhysicsAttack: 275, ElementAttack: ea(91), PhysicsType: Blow, ElementType: Shine},
	{ID: 19, IndexNumber: "W019", Name: "重撃ハンマー", Category: "ハンマー", PhysicsAttack: 445, ElementAttack: ea(56), PhysicsType: Blow, ElementType: Neutral},
	{ID: 20, IndexNumber: "W020", Name: "貫通重銃", Category: "重銃", PhysicsAttack: 445, ElementAttack: ea(56), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 21, IndexNumber: "W021", Name: "灼熱両手銃", Category: "両手銃", PhysicsAttack: 175, ElementAttack: ea(145), PhysicsType: Shoot, ElementType: Flame},
	{ID: 22, IndexNumber: "W022", Name: "霜月の短剣", Category: "短剣", PhysicsAttack: 175, ElementAttack: ea(145), PhysicsType: Slash, ElementType: Water},
	{ID: 23, IndexNumber: "W023", Name: "星砕きのモーニングスター", Category: "モーニングスター", PhysicsAttack: 275, ElementAttack: ea(91), PhysicsType: Blow, ElementType: Dark},

	// --- T4 忘却の湖畔 (Lv15-19) 積 50,000 ---
	{ID: 24, IndexNumber: "W024", Name: "白刃の剣", Category: "剣", PhysicsAttack: 385, ElementAttack: ea(130), PhysicsType: Slash, ElementType: Neutral},
	{ID: 25, IndexNumber: "W025", Name: "双牙・蒼波", Category: "双剣", PhysicsAttack: 245, ElementAttack: ea(205), PhysicsType: Slash, ElementType: Water},
	{ID: 26, IndexNumber: "W026", Name: "磁鉄のハンマー", Category: "ハンマー", PhysicsAttack: 630, ElementAttack: ea(79), PhysicsType: Blow, ElementType: Neutral},
	{ID: 27, IndexNumber: "W027", Name: "樹霊メイス", Category: "メイス", PhysicsAttack: 385, ElementAttack: ea(130), PhysicsType: Blow, ElementType: Wood},
	{ID: 28, IndexNumber: "W028", Name: "硝煙の銃", Category: "銃", PhysicsAttack: 385, ElementAttack: ea(130), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 29, IndexNumber: "W029", Name: "毒棘の爪", Category: "爪", PhysicsAttack: 245, ElementAttack: ea(205), PhysicsType: Slash, ElementType: Dark},
	{ID: 30, IndexNumber: "W030", Name: "陽光の弓", Category: "弓", PhysicsAttack: 245, ElementAttack: ea(205), PhysicsType: Shoot, ElementType: Shine},
	{ID: 31, IndexNumber: "W031", Name: "溶岩斧", Category: "斧", PhysicsAttack: 630, ElementAttack: ea(79), PhysicsType: Slash, ElementType: Flame},
	{ID: 32, IndexNumber: "W032", Name: "鉛玉の重銃", Category: "重銃", PhysicsAttack: 630, ElementAttack: ea(79), PhysicsType: Shoot, ElementType: Neutral},

	// --- T5 崩れた城塞 (Lv20-26) 積 95,000 ---
	{ID: 33, IndexNumber: "W033", Name: "覇王剣", Category: "剣", PhysicsAttack: 535, ElementAttack: ea(180), PhysicsType: Slash, ElementType: Neutral},
	{ID: 34, IndexNumber: "W034", Name: "双牙・紅蓮", Category: "双剣", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Slash, ElementType: Flame},
	{ID: 35, IndexNumber: "W035", Name: "深淵の爪", Category: "爪", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Slash, ElementType: Dark},
	{ID: 36, IndexNumber: "W036", Name: "巨神ハンマー", Category: "ハンマー", PhysicsAttack: 870, ElementAttack: ea(110), PhysicsType: Blow, ElementType: Neutral},
	{ID: 37, IndexNumber: "W037", Name: "雷光モーニングスター", Category: "モーニングスター", PhysicsAttack: 535, ElementAttack: ea(180), PhysicsType: Blow, ElementType: Shine},
	{ID: 38, IndexNumber: "W038", Name: "蔦縛りメイス", Category: "メイス", PhysicsAttack: 535, ElementAttack: ea(180), PhysicsType: Blow, ElementType: Wood},
	{ID: 39, IndexNumber: "W039", Name: "蒼海の重銃", Category: "重銃", PhysicsAttack: 870, ElementAttack: ea(110), PhysicsType: Shoot, ElementType: Water},
	{ID: 40, IndexNumber: "W040", Name: "極光両手銃", Category: "両手銃", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Shoot, ElementType: Shine},
	{ID: 41, IndexNumber: "W041", Name: "星屑の弓", Category: "弓", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 42, IndexNumber: "W042", Name: "氷牙の短剣", Category: "短剣", PhysicsAttack: 340, ElementAttack: ea(280), PhysicsType: Slash, ElementType: Water},

	// --- T6 灼熱の火山道 (Lv27-33) 積 160,000 ---
	{ID: 43, IndexNumber: "W043", Name: "焔喰いの大剣", Category: "剣", PhysicsAttack: 695, ElementAttack: ea(230), PhysicsType: Slash, ElementType: Flame},
	{ID: 44, IndexNumber: "W044", Name: "溶鉄のハンマー", Category: "ハンマー", PhysicsAttack: 1130, ElementAttack: ea(140), PhysicsType: Blow, ElementType: Neutral},
	{ID: 45, IndexNumber: "W045", Name: "業火の拳", Category: "拳", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Blow, ElementType: Flame},
	{ID: 46, IndexNumber: "W046", Name: "火砕の斧", Category: "斧", PhysicsAttack: 1130, ElementAttack: ea(140), PhysicsType: Slash, ElementType: Flame},
	{ID: 47, IndexNumber: "W047", Name: "熾火の双剣", Category: "双剣", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Slash, ElementType: Flame},
	{ID: 48, IndexNumber: "W048", Name: "灰塵の銃", Category: "銃", PhysicsAttack: 695, ElementAttack: ea(230), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 49, IndexNumber: "W049", Name: "硫煙の両手銃", Category: "両手銃", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Shoot, ElementType: Flame},
	{ID: 50, IndexNumber: "W050", Name: "陽炎の爪", Category: "爪", PhysicsAttack: 440, ElementAttack: ea(365), PhysicsType: Slash, ElementType: Shine},
	{ID: 51, IndexNumber: "W051", Name: "黒煙のモーニングスター", Category: "モーニングスター", PhysicsAttack: 695, ElementAttack: ea(230), PhysicsType: Blow, ElementType: Dark},
	{ID: 52, IndexNumber: "W052", Name: "火脈の重銃", Category: "重銃", PhysicsAttack: 1130, ElementAttack: ea(140), PhysicsType: Shoot, ElementType: Flame},

	// --- T7 凍てつく尖塔 (Lv34-41) 積 260,000 ---
	{ID: 53, IndexNumber: "W053", Name: "氷結の大剣", Category: "剣", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Slash, ElementType: Water},
	{ID: 54, IndexNumber: "W054", Name: "霜牙の双剣", Category: "双剣", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Water},
	{ID: 55, IndexNumber: "W055", Name: "氷壁のハンマー", Category: "ハンマー", PhysicsAttack: 1440, ElementAttack: ea(180), PhysicsType: Blow, ElementType: Water},
	{ID: 56, IndexNumber: "W056", Name: "極寒のメイス", Category: "メイス", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Blow, ElementType: Water},
	{ID: 57, IndexNumber: "W057", Name: "白銀の弓", Category: "弓", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Shoot, ElementType: Water},
	{ID: 58, IndexNumber: "W058", Name: "樹氷の斧", Category: "斧", PhysicsAttack: 1440, ElementAttack: ea(180), PhysicsType: Slash, ElementType: Wood},
	{ID: 59, IndexNumber: "W059", Name: "極光の重銃", Category: "重銃", PhysicsAttack: 1440, ElementAttack: ea(180), PhysicsType: Shoot, ElementType: Shine},
	{ID: 60, IndexNumber: "W060", Name: "凍蝶の短剣", Category: "短剣", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Water},
	{ID: 61, IndexNumber: "W061", Name: "氷晶の爪", Category: "爪", PhysicsAttack: 560, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Water},
	{ID: 62, IndexNumber: "W062", Name: "雪嶺の拳", Category: "拳", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Blow, ElementType: Neutral},
	{ID: 63, IndexNumber: "W063", Name: "極夜のモーニングスター", Category: "モーニングスター", PhysicsAttack: 885, ElementAttack: ea(295), PhysicsType: Blow, ElementType: Dark},

	// --- T8 黄昏の大樹海 (Lv42-49) 積 420,000 ---
	{ID: 64, IndexNumber: "W064", Name: "世界樹の弓", Category: "弓", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Shoot, ElementType: Wood},
	{ID: 65, IndexNumber: "W065", Name: "大樹守の斧", Category: "斧", PhysicsAttack: 1830, ElementAttack: ea(230), PhysicsType: Slash, ElementType: Wood},
	{ID: 66, IndexNumber: "W066", Name: "樹海のメイス", Category: "メイス", PhysicsAttack: 1120, ElementAttack: ea(375), PhysicsType: Blow, ElementType: Wood},
	{ID: 67, IndexNumber: "W067", Name: "黄昏の大剣", Category: "剣", PhysicsAttack: 1120, ElementAttack: ea(375), PhysicsType: Slash, ElementType: Dark},
	{ID: 68, IndexNumber: "W068", Name: "木漏れ日の両手銃", Category: "両手銃", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Shoot, ElementType: Shine},
	{ID: 69, IndexNumber: "W069", Name: "蔦絡みの重銃", Category: "重銃", PhysicsAttack: 1830, ElementAttack: ea(230), PhysicsType: Shoot, ElementType: Wood},
	{ID: 70, IndexNumber: "W070", Name: "腐葉の爪", Category: "爪", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Slash, ElementType: Dark},
	{ID: 71, IndexNumber: "W071", Name: "根喰いのハンマー", Category: "ハンマー", PhysicsAttack: 1830, ElementAttack: ea(230), PhysicsType: Blow, ElementType: Wood},
	{ID: 72, IndexNumber: "W072", Name: "蜜蝋の銃", Category: "銃", PhysicsAttack: 1120, ElementAttack: ea(375), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 73, IndexNumber: "W073", Name: "深緑の双剣", Category: "双剣", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Slash, ElementType: Wood},
	{ID: 74, IndexNumber: "W074", Name: "落陽の拳", Category: "拳", PhysicsAttack: 710, ElementAttack: ea(590), PhysicsType: Blow, ElementType: Shine},

	// --- T9 虚無の回廊 (Lv50-58) 積 650,000 ---
	{ID: 75, IndexNumber: "W075", Name: "虚無の大剣", Category: "剣", PhysicsAttack: 1400, ElementAttack: ea(465), PhysicsType: Slash, ElementType: Neutral},
	{ID: 76, IndexNumber: "W076", Name: "空虚のハンマー", Category: "ハンマー", PhysicsAttack: 2280, ElementAttack: ea(285), PhysicsType: Blow, ElementType: Neutral},
	{ID: 77, IndexNumber: "W077", Name: "無銘の双剣", Category: "双剣", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Slash, ElementType: Neutral},
	{ID: 78, IndexNumber: "W078", Name: "静寂の弓", Category: "弓", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 79, IndexNumber: "W079", Name: "虚ろの重銃", Category: "重銃", PhysicsAttack: 2280, ElementAttack: ea(285), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 80, IndexNumber: "W080", Name: "忘却の爪", Category: "爪", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Slash, ElementType: Dark},
	{ID: 81, IndexNumber: "W081", Name: "沈黙のメイス", Category: "メイス", PhysicsAttack: 1400, ElementAttack: ea(465), PhysicsType: Blow, ElementType: Neutral},
	{ID: 82, IndexNumber: "W082", Name: "白紙の短剣", Category: "短剣", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Slash, ElementType: Neutral},
	{ID: 83, IndexNumber: "W083", Name: "消失の両手銃", Category: "両手銃", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Shoot, ElementType: Dark},
	{ID: 84, IndexNumber: "W084", Name: "無音の拳", Category: "拳", PhysicsAttack: 885, ElementAttack: ea(735), PhysicsType: Blow, ElementType: Dark},
	{ID: 85, IndexNumber: "W085", Name: "虚数のモーニングスター", Category: "モーニングスター", PhysicsAttack: 1400, ElementAttack: ea(465), PhysicsType: Blow, ElementType: Neutral},

	// --- T10 星霜の墓標 (Lv59-68) 積 950,000 ---
	{ID: 86, IndexNumber: "W086", Name: "星霜の大剣", Category: "剣", PhysicsAttack: 1690, ElementAttack: ea(565), PhysicsType: Slash, ElementType: Shine},
	{ID: 87, IndexNumber: "W087", Name: "墓標のハンマー", Category: "ハンマー", PhysicsAttack: 2760, ElementAttack: ea(345), PhysicsType: Blow, ElementType: Dark},
	{ID: 88, IndexNumber: "W088", Name: "星屑の双剣", Category: "双剣", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Slash, ElementType: Shine},
	{ID: 89, IndexNumber: "W089", Name: "冥星の重銃", Category: "重銃", PhysicsAttack: 2760, ElementAttack: ea(345), PhysicsType: Shoot, ElementType: Dark},
	{ID: 90, IndexNumber: "W090", Name: "彗尾の弓", Category: "弓", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Shoot, ElementType: Shine},
	{ID: 91, IndexNumber: "W091", Name: "銀河の斧", Category: "斧", PhysicsAttack: 2760, ElementAttack: ea(345), PhysicsType: Slash, ElementType: Shine},
	{ID: 92, IndexNumber: "W092", Name: "深宙のメイス", Category: "メイス", PhysicsAttack: 1690, ElementAttack: ea(565), PhysicsType: Blow, ElementType: Dark},
	{ID: 93, IndexNumber: "W093", Name: "白夜の短剣", Category: "短剣", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Slash, ElementType: Shine},
	{ID: 94, IndexNumber: "W094", Name: "暁光の両手銃", Category: "両手銃", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Shoot, ElementType: Shine},
	{ID: 95, IndexNumber: "W095", Name: "黒点の拳", Category: "拳", PhysicsAttack: 1070, ElementAttack: ea(890), PhysicsType: Blow, ElementType: Dark},
	{ID: 96, IndexNumber: "W096", Name: "流星のモーニングスター", Category: "モーニングスター", PhysicsAttack: 1690, ElementAttack: ea(565), PhysicsType: Blow, ElementType: Shine},

	// --- T11 天空の廃神殿 (Lv69-79) 積 1,400,000 ---
	{ID: 97, IndexNumber: "W097", Name: "神鳴の拳", Category: "拳", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Blow, ElementType: Shine},
	{ID: 98, IndexNumber: "W098", Name: "天蓋の大剣", Category: "剣", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Slash, ElementType: Shine},
	{ID: 99, IndexNumber: "W099", Name: "雷神のハンマー", Category: "ハンマー", PhysicsAttack: 3350, ElementAttack: ea(420), PhysicsType: Blow, ElementType: Shine},
	{ID: 100, IndexNumber: "W100", Name: "聖樹の弓", Category: "弓", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Shoot, ElementType: Wood},
	{ID: 101, IndexNumber: "W101", Name: "神威の斧", Category: "斧", PhysicsAttack: 3350, ElementAttack: ea(420), PhysicsType: Slash, ElementType: Flame},
	{ID: 102, IndexNumber: "W102", Name: "天羽の双剣", Category: "双剣", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Slash, ElementType: Water},
	{ID: 103, IndexNumber: "W103", Name: "御雷の重銃", Category: "重銃", PhysicsAttack: 3350, ElementAttack: ea(420), PhysicsType: Shoot, ElementType: Shine},
	{ID: 104, IndexNumber: "W104", Name: "神楽鈴のメイス", Category: "メイス", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Blow, ElementType: Neutral},
	{ID: 105, IndexNumber: "W105", Name: "天啓の両手銃", Category: "両手銃", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Shoot, ElementType: Neutral},
	{ID: 106, IndexNumber: "W106", Name: "神代の爪", Category: "爪", PhysicsAttack: 1300, ElementAttack: ea(1080), PhysicsType: Slash, ElementType: Wood},
	{ID: 107, IndexNumber: "W107", Name: "祝詞の銃", Category: "銃", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Shoot, ElementType: Shine},
	{ID: 108, IndexNumber: "W108", Name: "天壌のモーニングスター", Category: "モーニングスター", PhysicsAttack: 2050, ElementAttack: ea(685), PhysicsType: Blow, ElementType: Wood},

	// --- T12 終末の空 (Lv80-90) 積 2,000,000 ---
	{ID: 109, IndexNumber: "W109", Name: "終焉の大剣", Category: "剣", PhysicsAttack: 2450, ElementAttack: ea(815), PhysicsType: Slash, ElementType: Dark},
	{ID: 110, IndexNumber: "W110", Name: "天焦がす斧", Category: "斧", PhysicsAttack: 4000, ElementAttack: ea(500), PhysicsType: Slash, ElementType: Flame},
	{ID: 111, IndexNumber: "W111", Name: "万象のモーニングスター", Category: "モーニングスター", PhysicsAttack: 2450, ElementAttack: ea(815), PhysicsType: Blow, ElementType: Neutral},
	{ID: 112, IndexNumber: "W112", Name: "世界喰らいのハンマー", Category: "ハンマー", PhysicsAttack: 4000, ElementAttack: ea(500), PhysicsType: Blow, ElementType: Neutral},
	{ID: 113, IndexNumber: "W113", Name: "極点の両手銃", Category: "両手銃", PhysicsAttack: 1550, ElementAttack: ea(1290), PhysicsType: Shoot, ElementType: Water},
	{ID: 114, IndexNumber: "W114", Name: "虚空の重銃", Category: "重銃", PhysicsAttack: 4000, ElementAttack: ea(500), PhysicsType: Shoot, ElementType: Dark},
	{ID: 115, IndexNumber: "W115", Name: "終末の弓", Category: "弓", PhysicsAttack: 1550, ElementAttack: ea(1290), PhysicsType: Shoot, ElementType: Wood},
	{ID: 116, IndexNumber: "W116", Name: "深淵の双剣", Category: "双剣", PhysicsAttack: 1550, ElementAttack: ea(1290), PhysicsType: Slash, ElementType: Dark},
	{ID: 117, IndexNumber: "W117", Name: "滅びの拳", Category: "拳", PhysicsAttack: 1550, ElementAttack: ea(1290), PhysicsType: Blow, ElementType: Flame},
	{ID: 118, IndexNumber: "W118", Name: "創世のメイス", Category: "メイス", PhysicsAttack: 2450, ElementAttack: ea(815), PhysicsType: Blow, ElementType: Shine},
	{ID: 119, IndexNumber: "W119", Name: "断罪の短剣", Category: "短剣", PhysicsAttack: 1550, ElementAttack: ea(1290), PhysicsType: Slash, ElementType: Shine},
	{ID: 120, IndexNumber: "W120", Name: "因果の銃", Category: "銃", PhysicsAttack: 2450, ElementAttack: ea(815), PhysicsType: Shoot, ElementType: Neutral},
}
