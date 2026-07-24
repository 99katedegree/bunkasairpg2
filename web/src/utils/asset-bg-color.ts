import type { EffectType } from "@/types/effect-type";
import type { ElementType } from "@/types/element-type";
import type { PhysicsType } from "@/types/physics-type";

const assetBgColorMap = {
  slash: "bg-slash",
  blow: "bg-blow",
  shoot: "bg-shoot",
  neutral: "bg-neutral",
  flame: "bg-flame",
  water: "bg-water",
  wood: "bg-wood",
  shine: "bg-shine",
  dark: "bg-dark",
  heal: "bg-heal",
  buff: "bg-buff",
  debuff: "bg-debuff",
};

export function assetBgColor(type: PhysicsType | ElementType | EffectType): string {
  return assetBgColorMap[type];
}
