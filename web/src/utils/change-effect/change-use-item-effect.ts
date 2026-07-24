import type { EffectType } from "@/types/effect-type";
import type { EffectMode } from "./change-attack-effect";

interface Props {
  setChangeEffect: React.Dispatch<React.SetStateAction<EffectMode>>;
  effectType: EffectType;
  healedAmount?: number;
}
export function changeUseItemEffect({ setChangeEffect, effectType, healedAmount }: Props) {
  if (healedAmount === 0) return;
  setChangeEffect(effectType);
  setTimeout(() => {
    setChangeEffect("none");
  }, 500);
}
