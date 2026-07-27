package main

// ============================================================================
// アイテムマスタ — 全 44 個
//
// buff は  (1 + バフ率)          を攻撃側に掛ける。素直な火力アップ。
// debuff は(1 - 耐性 * (1 - 率)) の内側に入る。つまり「相手の耐性を削る」効果で、
//          率 1.0 なら耐性を完全に無視できる。
//
// この非対称が攻略の肝になっている。
//   ・耐性 0 の相手にデバフを使っても意味がない
//   ・耐性 1.0（無効）の相手でも、デバフ 0.5 を入れれば係数 0.5 まで戻る
//   ・耐性 1.3（吸収）の相手でも、デバフ 0.5 を入れれば係数 0.35 の正ダメージになる
// 「全軸 1.0 でデバフ必須」のモンスターはここを踏まえて配置している。
//
// rate はシミュレータ側で 0.1 刻みに切り捨てられる（math.Floor(rate*10)/10）ので、
// 0.1 刻みで指定すること。
// ============================================================================

var items = []ItemSeed{

	// --- 回復 ---
	{ID: 1, IndexNumber: "I001", Name: "傷薬", EffectType: EffectHeal, Amount: 40},
	{ID: 2, IndexNumber: "I002", Name: "回復薬", EffectType: EffectHeal, Amount: 90},
	{ID: 3, IndexNumber: "I003", Name: "上級回復薬", EffectType: EffectHeal, Amount: 150},
	{ID: 4, IndexNumber: "I004", Name: "特上回復薬", EffectType: EffectHeal, Amount: 230},
	{ID: 5, IndexNumber: "I005", Name: "女神の雫", EffectType: EffectHeal, Amount: 330},
	{ID: 6, IndexNumber: "I006", Name: "星霜の霊薬", EffectType: EffectHeal, Amount: 450},
	{ID: 7, IndexNumber: "I007", Name: "神樹の果実", EffectType: EffectHeal, Amount: 620},
	{ID: 8, IndexNumber: "I008", Name: "蘇生の光", EffectType: EffectHeal, Amount: 810},

	// --- バフ（火力アップ） ---
	{ID: 9, IndexNumber: "I009", Name: "斬撃の砥石", EffectType: EffectBuff, Rate: 0.3, Target: Slash},
	{ID: 10, IndexNumber: "I010", Name: "剛断の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Slash},
	{ID: 11, IndexNumber: "I011", Name: "剛拳の指輪", EffectType: EffectBuff, Rate: 0.3, Target: Blow},
	{ID: 12, IndexNumber: "I012", Name: "破城の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Blow},
	{ID: 13, IndexNumber: "I013", Name: "精密照準器", EffectType: EffectBuff, Rate: 0.3, Target: Shoot},
	{ID: 14, IndexNumber: "I014", Name: "神狙の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Shoot},
	{ID: 15, IndexNumber: "I015", Name: "無銘の霊薬", EffectType: EffectBuff, Rate: 0.3, Target: Neutral},
	{ID: 16, IndexNumber: "I016", Name: "万象の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Neutral},
	{ID: 17, IndexNumber: "I017", Name: "火炎の香油", EffectType: EffectBuff, Rate: 0.3, Target: Flame},
	{ID: 18, IndexNumber: "I018", Name: "業火の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Flame},
	{ID: 19, IndexNumber: "I019", Name: "水霊の護符", EffectType: EffectBuff, Rate: 0.3, Target: Water},
	{ID: 20, IndexNumber: "I020", Name: "深海の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Water},
	{ID: 21, IndexNumber: "I021", Name: "大樹の若芽", EffectType: EffectBuff, Rate: 0.3, Target: Wood},
	{ID: 22, IndexNumber: "I022", Name: "世界樹の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Wood},
	{ID: 23, IndexNumber: "I023", Name: "聖光の香", EffectType: EffectBuff, Rate: 0.3, Target: Shine},
	{ID: 24, IndexNumber: "I024", Name: "天啓の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Shine},
	{ID: 25, IndexNumber: "I025", Name: "深淵の印", EffectType: EffectBuff, Rate: 0.3, Target: Dark},
	{ID: 26, IndexNumber: "I026", Name: "冥府の秘薬", EffectType: EffectBuff, Rate: 0.8, Target: Dark},

	// --- デバフ（耐性削り） ---
	{ID: 27, IndexNumber: "I027", Name: "硬鱗砕き", EffectType: EffectDebuff, Rate: 0.5, Target: Slash},
	{ID: 28, IndexNumber: "I028", Name: "絶断の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Slash},
	{ID: 29, IndexNumber: "I029", Name: "衝撃徹し", EffectType: EffectDebuff, Rate: 0.5, Target: Blow},
	{ID: 30, IndexNumber: "I030", Name: "絶砕の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Blow},
	{ID: 31, IndexNumber: "I031", Name: "貫通弾薬", EffectType: EffectDebuff, Rate: 0.5, Target: Shoot},
	{ID: 32, IndexNumber: "I032", Name: "絶貫の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Shoot},
	{ID: 33, IndexNumber: "I033", Name: "無象の触媒", EffectType: EffectDebuff, Rate: 0.5, Target: Neutral},
	{ID: 34, IndexNumber: "I034", Name: "万象の触媒", EffectType: EffectDebuff, Rate: 1.0, Target: Neutral},
	{ID: 35, IndexNumber: "I035", Name: "滅火の粉", EffectType: EffectDebuff, Rate: 0.5, Target: Flame},
	{ID: 36, IndexNumber: "I036", Name: "絶火の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Flame},
	{ID: 37, IndexNumber: "I037", Name: "断水の塩", EffectType: EffectDebuff, Rate: 0.5, Target: Water},
	{ID: 38, IndexNumber: "I038", Name: "絶水の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Water},
	{ID: 39, IndexNumber: "I039", Name: "枯朽の胞子", EffectType: EffectDebuff, Rate: 0.5, Target: Wood},
	{ID: 40, IndexNumber: "I040", Name: "絶木の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Wood},
	{ID: 41, IndexNumber: "I041", Name: "遮光の灰", EffectType: EffectDebuff, Rate: 0.5, Target: Shine},
	{ID: 42, IndexNumber: "I042", Name: "絶光の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Shine},
	{ID: 43, IndexNumber: "I043", Name: "呪解の香", EffectType: EffectDebuff, Rate: 0.5, Target: Dark},
	{ID: 44, IndexNumber: "I044", Name: "絶闇の劇薬", EffectType: EffectDebuff, Rate: 1.0, Target: Dark},
}
