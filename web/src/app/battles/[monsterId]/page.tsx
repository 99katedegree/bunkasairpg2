"use client";

import {
  BattleLog,
  BattlePage,
} from "@/components/features/battle/battle-page/battle-page";
import { playSound } from "@/utils/play-sound/play-sound";
import { useUserStore } from "@/stores/user-store";
import { Battle } from "@/utils/battle";
import { useParams, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { monsterAttackSound } from "@/utils/play-sound/monster-attack-sound";
import { getMonster } from "@/lib/monster/monster";
import { startBattle, finishBattle } from "@/lib/battle/battle";

export default function Page() {
  const [battle, setBattle] = useState<Battle | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const { user, setUser } = useUserStore();
  const router = useRouter();
  const { monsterId } = useParams<{ monsterId: string }>();

  useEffect(() => {
    if (!user) return;

    Promise.all([
      getMonster(monsterId),
      startBattle({ monsterToken: monsterId }),
    ]).then(([monsterRes, battleRes]) => {
      if (monsterRes.status !== 200 || battleRes.status !== 200) return;

      const monster = monsterRes.data.monster;
      const { token: battleToken } = battleRes.data;
      setToken(battleToken);

      setBattle(
        new Battle(
          structuredClone({ ...user, maxHitPoint: user.hitPoint }),
          {
            ...monster,
            maxHitPoint: monster.hitPoint,
          }
        )
      );
    });
  }, [monsterId]);

  if (!battle || !user) return null;

  return (
    <BattlePage
      battle={battle}
      onVictory={async (b) => {
        const res = await finishBattle({
          token: token!,
          actions: b.getActions(),
        });
        if (res.status !== 200) {
          throw new Error("Battle finish failed");
        }
        const { level, hitPoint, experiencePoint, drop } = res.data;
        return {
          level,
          hitPoint,
          experiencePoint,
          drop: (drop?.type ?? null) as "weapon" | "item" | null,
          clearTime: undefined,
        };
      }}
      monsterAttackLogs={(setBattlePhase) => {
        const monster = battle.getMonster();
        const takeDamageData = battle.takeDamage();
        const logs: BattleLog[] = [
          {
            message: `${monster.name}の\n攻撃！`,
            action: () => {
              monsterAttackSound({ takeDamage: takeDamageData.damage });
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
                setBattlePhase({
                  status: "command",
                  action: null,
                });
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
              setUser({ ...user, hitPoint: user.maxHitPoint });
              router.push("/scan");
            },
          });
        }
        return logs;
      }}
    />
  );
}
