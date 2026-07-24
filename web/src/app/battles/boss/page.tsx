"use client";

import {
  BattleLog,
  BattlePage,
} from "@/components/features/battle/battle-page/battle-page";
import { ITEM_TARGET_LABEL_MAP } from "@/constants/item-target-label-map";
import { useUserStore } from "@/stores/user-store";
import { ElementType } from "@/types/element-type";
import { BossBattle } from "@/utils/boss-battle";
import { playSound } from "@/utils/play-sound/play-sound";
import { useEffect, useRef, useState } from "react";
import { startBossBattle, finishBossBattle } from "@/lib/boss/boss";

export default function Page() {
  const [battle, setBattle] = useState<BossBattle | null>(null);
  const tokenRef = useRef<string | null>(null);
  const { user, setUser } = useUserStore();

  const bossImageMap: Record<ElementType, string> = {
    neutral: "/boss-neutral.png",
    flame: "/boss-wood.png",
    water: "/boss-flame.png",
    wood: "/boss-water.png",
    shine: "/boss-dark.png",
    dark: "/boss-shine.png",
  };

  useEffect(() => {
    if (!user) return;

    startBossBattle().then((res) => {
      if (res.status !== 200) return;
      const { token, seed } = res.data;
      tokenRef.current = token;
      setBattle(
        new BossBattle(
          structuredClone({ ...user, maxHitPoint: user.hitPoint }),
          seed
        )
      );
    });
  }, []);

  if (!battle || !user) return null;

  return (
    <BattlePage
      battle={battle}
      onVictory={async (b) => {
        const res = await finishBossBattle({
          token: tokenRef.current!,
          actions: b.getActions(),
        });
        if (res.status !== 200) {
          throw new Error("Boss battle finish failed");
        }
        return {
          level: res.data.level,
          hitPoint: res.data.hitPoint,
          experiencePoint: res.data.experiencePoint,
          drop: null,
          clearTime: res.data.clearTime,
        };
      }}
      monsterAttackLogs={(setBattlePhase, setMonster) => {
        const monster = battle.getMonster();
        const takeDamageData = battle.takeDamage();
        const { physicsType, elementType } = battle.shiftWeakness();
        const logs: BattleLog[] = [
          {
            message: `${monster.name}の姿が変わっていく！`,
            action: () => {
              setMonster({ ...monster, imageUrl: bossImageMap[elementType] });
              playSound("/sounds/debuff.mp3");
            },
          },
          {
            message: `${monster.name}の${ITEM_TARGET_LABEL_MAP[physicsType]}耐性が下がった！`,
            action: () => {
              playSound("/sounds/debuff.mp3");
            },
          },
          {
            message: `${monster.name}の${ITEM_TARGET_LABEL_MAP[elementType]}耐性が下がった！`,
            action: () => {},
          },
          {
            message: `${monster.name}の\n攻撃！`,
            action: () => {
              if (takeDamageData.damage === 0) {
                playSound("/sounds/prevent.mp3");
              } else {
                playSound("/sounds/monster-attack.mp3");
              }
            },
          },
        ];
        if (takeDamageData.damage === 0) {
          logs.push({
            message: `${monster.name}はダメージを与えられなかった！`,
            action: () => {},
          });
        } else {
          logs.push({
            message: `${takeDamageData.damage}のダメージを\n受けた！`,
            action: () => {
              if (takeDamageData.userHitPoint === 0) {
                setBattlePhase({ status: "command", action: null });
                playSound("/sounds/monster-down.mp3");
              }
              setUser({ ...user, hitPoint: takeDamageData.userHitPoint });
            },
          });
        }
        if (takeDamageData.userHitPoint === 0) {
          logs.push({
            message: `${user.name}は死んでしまった！`,
            action: () => {
              location.href = "/scan";
            },
          });
        }
        return logs;
      }}
    />
  );
}
