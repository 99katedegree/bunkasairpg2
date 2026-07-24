export type EffectMode = "none" | "attack" | "monsterHeal" | "monsterGuard" | "heal" | "buff" | "debuff";
interface Props {
  setChangeEffect: React.Dispatch<React.SetStateAction<EffectMode>>;
  attackDamage: number;
}
export function changeAttackEffect({ setChangeEffect, attackDamage }: Props) {
  if (attackDamage === 0) {
    setChangeEffect("monsterGuard");
  } else if (attackDamage < 0) {
    setChangeEffect("monsterHeal");
  } else {
    setChangeEffect("attack");
  }
  setTimeout(() => {
    setChangeEffect("none");
  }, 500);
}
