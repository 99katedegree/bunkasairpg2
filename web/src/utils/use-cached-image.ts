"use client";

import useSWRImmutable from "swr/immutable";

// SWR のグローバルキャッシュは URL（キー）ごとに fetcher を 1 回しか呼ばない。
// useSWRImmutable は revalidateOnFocus / revalidateOnReconnect / revalidateIfStale を
// すべて無効化した設定なので、一度取得できた画像はタブを切り替えても
// 再マウントしても再フェッチされない。同じ URL を要求する箇所が複数あっても
// 実際にブラウザへ読み込みに行くのは最初の 1 回だけになる。
//
// 実データの取得は Image() によるプリロードに任せる。バックエンドの
// /assets ルートが Cache-Control: immutable を返し、URL には内容のハッシュが
// 含まれる（cmd/seed-master 参照）ので、同じ URL は中身も変わらないという
// 前提が成り立っている。
function preload(url: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const img = new window.Image();
    img.onload = () => resolve(url);
    img.onerror = () => reject(new Error(`画像の読み込みに失敗しました: ${url}`));
    img.src = url;
  });
}

// url が null/undefined のときは SWR のキーも null にして、そもそも何も取得しない。
export function useCachedImage(url: string | null | undefined) {
  const { data, error, isLoading } = useSWRImmutable(url ?? null, preload);
  return { src: data, isLoading, error };
}
