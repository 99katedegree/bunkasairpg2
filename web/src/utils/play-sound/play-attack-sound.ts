import type { PhysicsType } from "../../types/physics-type";
import { playSound } from "./play-sound";

interface Props {
  physicsType: PhysicsType;
  attackDamage: number;
}
export function playAttackSound({
  physicsType,
  attackDamage,
}: Props) {
  if (attackDamage === 0) {
    playSound("/sounds/prevent.mp3");
  } else if (attackDamage < 0) {
    playSound("/sounds/heal.mp3");
  } else {
    switch (physicsType) {
      case "slash":
        playSound("/sounds/slash.mp3");
        break;
      case "blow":
        playSound("/sounds/blow.mp3");
        break;
      case "shoot":
        playSound("/sounds/shoot.mp3");
        break;
    }
  }
}
