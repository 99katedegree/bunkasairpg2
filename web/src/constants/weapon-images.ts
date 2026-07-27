// キーは武器の図鑑番号（4 桁）。DB に登録されている武器だけを載せる。
//
// 素手（図鑑番号 0000）はここに書かない。DB に行がなく、装備なしのときに
// handler/me.go がハードコード値を返す特別扱いの装備なので、このマップに
// 載せると管理画面の突き合わせで「DB に該当する図鑑番号がありません」と
// 毎回警告が出てしまう。素手の画像は /weapons/fist.png に置いて直接参照する。
export const weaponImages: Record<string, string> = {
  // "0001": "/weapons/<name>.png",
};
