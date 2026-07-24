import { playSound } from "./play-sound";

interface Props {
  takeDamage: number;
}
export function monsterAttackSound({ takeDamage }: Props) {
  if (takeDamage === 0) {
    playSound("/sounds/prevent.mp3");
  } else {
    playSound("/sounds/monster-attack.mp3");
  }
}
