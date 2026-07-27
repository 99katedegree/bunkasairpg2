export const monsterImages: Record<string, string> = {
  // キーは魔物の図鑑番号（4 桁）。上位 3 桁が種族、末尾 1 桁が歪みの深さ。
  // 内部 UUID ではなく図鑑番号を使うのは、QR や図鑑に印刷されて人が目にする
  // 番号がこちらで、管理画面の突き合わせもこの番号で行うため。
  // "0010": "/monsters/<name>.png",
};
