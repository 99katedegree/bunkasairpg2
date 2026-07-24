export function hpBgColor(hitPoint: number, maxHitPoint: number) {
  const ratio = hitPoint / maxHitPoint;
  if (ratio > 0.5) return "oklch(89.7% 0.196 126.665)"; // 緑
  if (ratio > 0.25) return "oklch(90.5% 0.182 98.111)"; // 黄
  return "oklch(70.4% 0.191 22.216)"; // 赤
}
