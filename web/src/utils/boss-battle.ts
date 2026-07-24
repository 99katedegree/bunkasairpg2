import type { PhysicsType } from "@/types/physics-type";
import { Battle, type User } from "./battle";
import type { ElementType } from "@/types/element-type";
import { boss } from "@/constants/boss";
import { mulberry32 } from "./seeded-rng";

export class BossBattle extends Battle {
  constructor(user: User, seed: number) {
    super(user, { ...boss }, true, mulberry32(seed));
  }

  public shiftWeakness(): { physicsType: PhysicsType; elementType: ElementType } {
    const physicsTypes: PhysicsType[] = ["slash", "blow", "shoot"];
    const elementTypes: ElementType[] = ["neutral", "flame", "water", "wood", "shine", "dark"];
    const physicsType = physicsTypes[Math.floor(this.rng() * 3)];
    const elementType = elementTypes[Math.floor(this.rng() * 6)];

    const resistance = {
      slash: 1.0,
      blow: 1.0,
      shoot: 1.0,
      neutral: 1.0,
      flame: 1.4,
      water: 1.4,
      wood: 1.4,
      shine: 1.4,
      dark: 1.4,
    };
    resistance[physicsType] = 0.7;
    resistance[elementType] = 0.9;
    this.monster = {
      ...this.monster,
      slash: resistance.slash,
      blow: resistance.blow,
      shoot: resistance.shoot,
      neutral: resistance.neutral,
      flame: resistance.flame,
      water: resistance.water,
      wood: resistance.wood,
      shine: resistance.shine,
      dark: resistance.dark,
    };

    return { physicsType, elementType };
  }
}
