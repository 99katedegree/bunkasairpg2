import { useEffect, useState } from "react";
import { Battle } from "@/utils/battle";
import { useUserStore } from "@/stores/user-store";
import { getMonster } from "@/lib/monster/monster";
import { MonsterDetailResponse } from "@/lib/bunkasaiRPGAPI.schemas";

export function useBattlePage(monsterId: string) {
  const { user } = useUserStore();
  const [battle, setBattle] = useState<Battle | null>(null);
  const [monster, setMonsterState] = useState<MonsterDetailResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) return;

    setLoading(true);
    getMonster(monsterId)
      .then((res) => {
        if (res.status !== 200) {
          setError("モンスターが見つかりませんでした");
          return;
        }
        const monsterData = res.data.monster;
        setMonsterState(monsterData);
        setBattle(
          new Battle(
            structuredClone({ ...user, maxHitPoint: user.hitPoint }),
            {
              ...monsterData,
              maxHitPoint: monsterData.hitPoint,
            }
          )
        );
      })
      .catch(() => setError("データの取得に失敗しました"))
      .finally(() => setLoading(false));
  }, [monsterId, user]);

  return { battle, monster, loading, error };
}
