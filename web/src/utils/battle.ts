import { updateMe } from "@/lib/me/me";
import { useMeItem, changeMeWeapon } from "@/lib/action/action";
import { startBattle, finishBattle } from "@/lib/battle/battle";
import type { ElementType } from "@/types/element-type";
import type { PhysicsType } from "@/types/physics-type";
import { calculateLevel } from "./calculate-level";
import type { MonsterDetailResponse } from "@/lib/bunkasaiRPGAPI.schemas";
import type { BattleAction } from "@/lib/bunkasaiRPGAPI.schemas";

export type User = {
  level: number;
  maxHitPoint: number;
  hitPoint: number;
  experiencePoint: number;

  weapon?: Weapon | null;
};

type Weapon = {
  id: number;
  physicsAttack: number;
  elementAttack: number | null;
  physicsType: PhysicsType;
  elementType: ElementType;
};

type HeelItem = {
  id: number;
  amount?: number;
};

type BuffItem = {
  id: number;
  rate?: number; // 何%増加するか
  target?: PhysicsType | ElementType | string;
};

type DebuffItem = {
  id: number;
  rate?: number; // 何%減少するか
  target?: PhysicsType | ElementType | string;
};

export class Battle {
  protected user: User;
  protected monster: MonsterDetailResponse & { maxHitPoint: number; imageUrl?: string };
  protected buffs = {
    slash: 0.0,
    blow: 0.0,
    shoot: 0.0,
    neutral: 0.0,
    flame: 0.0,
    water: 0.0,
    wood: 0.0,
    shine: 0.0,
    dark: 0.0,
  };
  protected debuffs = {
    slash: 0.0,
    blow: 0.0,
    shoot: 0.0,
    neutral: 0.0,
    flame: 0.0,
    water: 0.0,
    wood: 0.0,
    shine: 0.0,
    dark: 0.0,
  };
  protected isBoss = false;
  protected rng: () => number;
  private actions: BattleAction[] = [];
  private actionCount = 0;

  constructor(
    user: User,
    monster: MonsterDetailResponse & { maxHitPoint: number; imageUrl?: string },
    isBoss: boolean = false,
    rng: () => number = Math.random
  ) {
    this.user = user;
    this.monster = monster;
    this.isBoss = isBoss;
    this.rng = rng;

    if (!this.isBoss) {
      startBattle({ monsterId: this.monster.id });
    }
  }

  public getActions(): BattleAction[] {
    return this.actions;
  }

  public getMonster(): MonsterDetailResponse & { maxHitPoint: number; imageUrl?: string } {
    return this.monster;
  }

  public getIsBoss(): boolean {
    return this.isBoss;
  }

  public attack(): {
    monsterHitPoint: number;
    damage: number;
    monsterResistance: {
      physics: number;
      element: number;
    };
  } {
    if (!this.user.weapon) {
      return {
        monsterHitPoint: this.monster.hitPoint,
        damage: 0,
        monsterResistance: { physics: 1.0, element: 1.0 },
      };
    }
    const physicsType = this.user.weapon.physicsType;
    const elementType = this.user.weapon.elementType;
    const monsterPhysics = this.monster[physicsType] ?? 1.0;
    const monsterElement = this.monster[elementType] ?? 1.0;
    const physics =
      this.user.weapon.physicsAttack *
      (1 + this.buffs[physicsType]) *
      (1 - monsterPhysics * (1 - this.debuffs[physicsType]));
    const element =
      (this.user.weapon.elementAttack || 1) *
      (1 + this.buffs[elementType]) *
      (1 - monsterElement * (1 - this.debuffs[elementType]));
    const baseDamage = physics * element;
    const levelFactor = 0.8 + Math.sqrt(this.user.level) / 5;
    const random = 0.95 + this.rng() * 0.1;

    // √|
    //   (武器物理攻撃力 * (1 + 物理バフ) * (1 - モンスター物理耐性 * (1 - 物理デバフ))
    //    * 武器属性攻撃力 * (1 + 属性バフ) * (1 - モンスター属性耐性 * (1 - 属性デバフ))) |
    //   * レベル補正(0.8 + √(ユーザーLv)/5)
    //   * 乱数補正(0.95〜1.05)
    //   * 符号(元の基本ダメージが負なら -1、正なら 1)
    const damage = Math.floor(
      Math.sqrt(Math.abs(baseDamage)) *
        levelFactor *
        random *
        (baseDamage >= 0 ? 1 : -1)
    );

    this.monster.hitPoint =
      damage < 0
        ? Math.min(this.monster.hitPoint - damage, this.monster.maxHitPoint)
        : Math.max(this.monster.hitPoint - damage, 0);

    if (this.isBoss) {
      this.actions.push({ turn: this.actionCount++, type: "attack" });
    }

    return {
      monsterHitPoint: this.monster.hitPoint,
      damage: damage,
      monsterResistance: {
        physics: monsterPhysics,
        element: monsterElement,
      },
    };
  }

  public changeWeapon(weapon: Weapon): void {
    this.user.weapon = weapon;
    changeMeWeapon({ weaponId: weapon.id });
    if (this.isBoss) {
      this.actions.push({ turn: this.actionCount++, type: "change-weapon", weaponId: weapon.id });
    }
  }

  public useHealItem(item: HeelItem): void {
    if (!this.isBoss) {
      useMeItem({ itemId: item.id });
    }
    if (this.isBoss) {
      this.actions.push({ turn: this.actionCount++, type: "use-item", itemId: item.id });
    }
    this.user.hitPoint = Math.min(
      this.user.hitPoint + (item.amount ?? 0),
      this.user.maxHitPoint
    );
  }

  public useBuffItem(item: BuffItem): void {
    if (!this.isBoss) {
      useMeItem({ itemId: item.id });
    }
    if (this.isBoss) {
      this.actions.push({ turn: this.actionCount++, type: "use-item", itemId: item.id });
    }
    if (item.target && item.rate !== undefined) {
      const target = item.target as keyof typeof this.buffs;
      if (target in this.buffs) {
        this.buffs[target] += Math.floor(item.rate * 10) / 10;
      }
    }
  }

  public useDebuffItem(item: DebuffItem): void {
    if (!this.isBoss) {
      useMeItem({ itemId: item.id });
    }
    if (this.isBoss) {
      this.actions.push({ turn: this.actionCount++, type: "use-item", itemId: item.id });
    }
    if (item.target && item.rate !== undefined) {
      const target = item.target as keyof typeof this.debuffs;
      if (target in this.debuffs) {
        this.debuffs[target] += Math.floor(item.rate * 10) / 10;
      }
    }
  }

  public takeDamage(): {
    userHitPoint: number;
    damage: number;
  } {
    const random = 0.95 + this.rng() * 0.1;
    const levelFactor = 1 + Math.sqrt(this.user.level) / 1.7;
    // モンスター攻撃力 * 乱数(0.95〜1.05) / (1+√(ユーザーレベル)/1.7)
    const damage = Math.floor((this.monster.attack * random) / levelFactor);

    this.user.hitPoint = Math.max(this.user.hitPoint - damage, 0);
    return {
      userHitPoint: this.user.hitPoint,
      damage: damage,
    };
  }

  public reward(weaponIds: number[]): {
    level: number;
    hitPoint: number;
    experiencePoint: number;
    drop: "weapon" | "item" | null;
  } {
    let drop: "weapon" | "item" | null = null;
    if (
      this.monster.weapon?.id &&
      !weaponIds.includes(this.monster.weapon.id)
    ) {
      drop = "weapon";
    } else if (this.monster.item?.id) {
      drop = "item";
    }

    const level = calculateLevel(
      this.user.experiencePoint + this.monster.experiencePoint
    );
    if (this.user.level < level) {
      const increasedHitPoint = Array.from({
        length: level - this.user.level,
      }).reduce<number>(
        (sum) => sum + [6, 7, 8, 9, 10][Math.floor(Math.random() * 5)],
        0
      );
      this.user.level = level;
      this.user.maxHitPoint += increasedHitPoint;
    }
    this.user.experiencePoint += this.monster.experiencePoint;
    if (process.env.NEXT_PUBLIC_ALLOW_LEVEL_UP === "true") {
      updateMe({
        level: this.user.level,
        hitPoint: this.user.maxHitPoint,
        experiencePoint: this.user.experiencePoint,
      });
    }

    return {
      level: this.user.level,
      hitPoint: this.user.maxHitPoint,
      experiencePoint: this.user.experiencePoint,
      drop: drop,
    };
  }
}
