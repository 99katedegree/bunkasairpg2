import type { EffectType } from "@/types/effect-type";
import { playSound } from "./play-sound";

interface Props {
  effectType: EffectType;
  healedAmount?: number;
}
export function useItemSound({ effectType, healedAmount }: Props) {
  switch (effectType) {
    case "heal":
      if (healedAmount === 0) {
        playSound("/sounds/not-available.mp3");
      } else {
        playSound("/sounds/heal.mp3");
      }
      break;
    case "buff":
      playSound("/sounds/buff.mp3");
      break;
    case "debuff":
      playSound("/sounds/debuff.mp3");
      break;
  }
}
