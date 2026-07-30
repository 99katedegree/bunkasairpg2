"use client";

import { useCachedImage } from "@/utils/use-cached-image";

type Props = {
  src: string | null | undefined;
  alt: string;
  className?: string;
  fallback?: React.ReactNode;
};

// 一度読み込んだ画像 URL は再フェッチしない <img> のラッパー。詳細は
// utils/use-cached-image を参照。src が空、未取得、読み込み失敗のときは
// fallback（省略時は何も出さない）を表示する。
export function CachedImage({ src, alt, className, fallback = null }: Props) {
  const { src: loadedSrc } = useCachedImage(src);

  if (!loadedSrc) {
    return <>{fallback}</>;
  }

  return <img src={loadedSrc} alt={alt} className={className} />;
}
