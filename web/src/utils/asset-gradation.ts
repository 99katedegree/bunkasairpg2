import type { EffectType } from "@/types/effect-type";
import type { ElementType } from "@/types/element-type";

const assetGradationMap = {
  neutral: "from-neutral to-neutral-dark",
  flame: "from-flame to-flame-dark",
  water: "from-water to-water-dark",
  wood: "from-wood to-wood-dark",
  shine: "from-shine to-shine-dark",
  dark: "from-dark to-dark-dark",
  heal: "from-heal to-heal-dark",
  buff: "from-buff to-buff-dark",
  debuff: "from-debuff to-debuff-dark",
};

export function assetGradation(type: ElementType | EffectType): string {
  return `bg-linear-to-br ${assetGradationMap[type]}`;
}
