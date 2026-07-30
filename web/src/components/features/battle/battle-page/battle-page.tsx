"use client";

import { useState, useEffect } from "react";
import { Battle } from "@/utils/battle";
import { usePathname, useRouter } from "next/navigation";
import { useUserStore } from "@/stores/user-store";
import {
  WeaponDrawer,
} from "@/components/features/battle/weapon-drawer/weapon-drawer";
import { ItemDrawer } from "@/components/features/battle/item-drawer/item-drawer";
import { Modal } from "@/components/features/battle/modal/modal";
import Image from "next/image";
import { cn } from "@/utils/cn";
import { WeaponCard } from "@/components/features/battle/weapon-card/weapon-card";
import { RewardModal } from "@/components/features/battle/reward-modal/reward-modal";
import { motion } from "framer-motion";
import { BattleConsole } from "@/components/features/battle/battle-console/battle-console";
import { ITEM_TARGET_LABEL_MAP } from "@/constants/item-target-label-map";
import {
  Drop,
  TreasureBoxButton,
} from "@/components/features/battle/treasure-box-button/treasure-box-button";
import { hpBgColor } from "@/utils/hp-bg-color";
import { AttackEffect } from "@/components/features/battle/attack-effect/attack-effect";
import { UseItemEffect } from "@/components/features/battle/use-item-effect/use-item-effect";
import { playAttackSound } from "@/utils/play-sound/play-attack-sound";
import { playSound } from "@/utils/play-sound/play-sound";
import { useItemSound } from "@/utils/play-sound/use-item-sound";
import { EffectMode, changeAttackEffect } from "@/utils/change-effect/change-attack-effect";
import { changeUseItemEffect } from "@/utils/change-effect/change-use-item-effect";
import {
  MonsterDetailResponse,
  UserItemResponse,
  WeaponResponse,
} from "@/lib/bunkasaiRPGAPI.schemas";

export type BattlePhase =
  | { status: "first"; action: null | "weapon" | "item" }
  | { status: "command"; action: null | "weapon" | "item" };

export type BattleLog = {
  message: string;
  action: () => void;
};

type VictoryResult = {
  level: number;
  hitPoint: number;
  experiencePoint: number;
  drop: "weapon" | "item" | null;
  clearTime?: number;
};

type Props = {
  battle: Battle;
  monsterAttackLogs: (
    setBattlePhase: (bp: BattlePhase) => void,
    setMonster: (m: MonsterDetailResponse & { maxHitPoint: number; imageUrl?: string | null }) => void
  ) => BattleLog[];
  // 報酬はサーバーが決めるため必須。クライアント側で計算する経路は持たない。
  onVictory: (battle: Battle) => Promise<VictoryResult>;
};

export function BattlePage({ battle, monsterAttackLogs, onVictory }: Props) {
  const { user, setUser, weapons, items, setItems } =
    useUserStore();
  const router = useRouter();
  const pathname = usePathname();
  const [monster, setMonster] = useState<
    MonsterDetailResponse & { maxHitPoint: number; imageUrl?: string | null }
  >(
    structuredClone({
      ...battle.getMonster(),
      maxHitPoint: battle.getMonster().hitPoint,
    })
  );
  const [battlePhase, setBattlePhase] = useState<BattlePhase>({
    status: "first",
    action: null,
  });
  const [battleQueue, setBattleQueue] = useState<BattleLog[]>([]);
  const [isStandBy, setIsStandBy] = useState(false);
  const [isItemUsed, setIsItemUsed] = useState(false);
  const [effectMode, setEffectMode] = useState<EffectMode>("none");
  const [showAwayModal, setShowAwayModal] = useState(false);
  const [reward, setReward] = useState<{
    restLevel: number;
    restExperiencePoint: number;
    level: number;
    experiencePoint: number;
    drop: Drop | null;
    clearTime: string;
  } | null>(null);
  const [startDate, setStartDate] = useState<Date>();

  useEffect(() => {
    setIsStandBy(!!battleQueue.length);
  }, [battleQueue]);

  if (!monster || !battle || !user) return null;

  // 攻撃メソッド
  const handleAttack = () => {
    const attackData = battle.attack();
    const logs: BattleLog[] = [
      {
        message: `${monster.name}を\n攻撃した！`,
        action: () => {
          if (
            !(
              attackData.monsterResistance.physics < 0 ||
              attackData.monsterResistance.element < 0
            )
          ) {
            playAttackSound({
              physicsType: user.weapon!.physicsType,
              attackDamage: attackData.damage,
            });
            changeAttackEffect({
              setChangeEffect: setEffectMode,
              attackDamage: attackData.damage,
            });
          }
        },
      },
    ];

    if (attackData.damage === 0) {
      logs.push({
        message: "モンスターの防御に\n阻まれた！",
        action: () => {},
      });
    } else if (attackData.damage < 0) {
      logs.push({
        message: `${-attackData.damage}ダメージが\n吸収された！`,
        action: () =>
          setMonster({
            ...monster,
            hitPoint: attackData.monsterHitPoint,
          }),
      });
    } else {
      if (
        attackData.monsterResistance.physics < 0 ||
        attackData.monsterResistance.element < 0
      ) {
        logs.push({
          message: "弱点を突いた！",
          action: () => {
            playSound("/sounds/weakness.mp3");
            changeAttackEffect({
              setChangeEffect: setEffectMode,
              attackDamage: attackData.damage,
            });
          },
        });
      }
      logs.push({
        message: `${attackData.damage}のダメージを\n与えた！`,
        action: () => {
          setMonster({
            ...monster,
            hitPoint: attackData.monsterHitPoint,
          });
          if (attackData.monsterHitPoint === 0) {
            playSound("/sounds/monster-down.mp3");
          }
        },
      });
    }

    if (attackData.monsterHitPoint !== 0) {
      setIsItemUsed(false);
      setBattleQueue([
        ...logs,
        ...monsterAttackLogs(setBattlePhase, setMonster),
      ]);
      return;
    }

    {
      const rewardPromise = onVictory(battle);
      logs.push({
        message: `${monster.name}を\n倒した！`,
        action: () => {
          rewardPromise.then((rewardData) => {
            setReward({
              restLevel: user.level,
              restExperiencePoint: user.experiencePoint,
              level: rewardData.level,
              experiencePoint: rewardData.experiencePoint,
              drop: null,
              clearTime:
                rewardData.clearTime !== undefined
                  ? formatMilliseconds(rewardData.clearTime)
                  : "??:??:??",
            });
            setUser({
              ...user,
              level: rewardData.level,
              experiencePoint: rewardData.experiencePoint,
              hitPoint: rewardData.hitPoint,
              maxHitPoint: rewardData.hitPoint,
            });
          });
        },
      });
    }
    setBattleQueue(logs);
  };

  // 武器変更ログ
  const changeWeaponLogs = (weapon: WeaponResponse) => {
    battle.changeWeapon(weapon);
    const logs: BattleLog[] = [
      {
        message: `${weapon.name}を\n装備した！`,
        action: () => {
          setUser({ ...user, weapon: weapon });
          playSound("/sounds/weapon-change.mp3");
        },
      },
    ];
    return logs;
  };

  // アイテム使用ログ
  const useItemLogs = (item: UserItemResponse) => {
    if (item.count === 1) {
      setItems(items.filter((prev) => prev.id !== item.id));
    } else {
      setItems(
        items.map((prev) =>
          prev.id === item.id ? { ...prev, count: prev.count - 1 } : prev
        )
      );
    }

    const logs: BattleLog[] = [
      {
        message: `${item.name}を\n使った！`,
        action: () => {
          if (item.effectType === "heal") {
            useItemSound({
              effectType: item.effectType,
              healedAmount: Math.min(
                item.amount ?? 0,
                user.maxHitPoint - user.hitPoint
              ),
            });
            changeUseItemEffect({
              setChangeEffect: setEffectMode,
              effectType: item.effectType,
              healedAmount: Math.min(
                item.amount ?? 0,
                user.maxHitPoint - user.hitPoint
              ),
            });
          } else {
            useItemSound({ effectType: item.effectType });
            changeUseItemEffect({
              setChangeEffect: setEffectMode,
              effectType: item.effectType,
            });
          }
        },
      },
    ];

    if (item.effectType === "heal") {
      battle.useHealItem(item);
      const healedAmount = Math.min(
        item.amount ?? 0,
        user.maxHitPoint - user.hitPoint
      );
      logs.push({
        message:
          healedAmount === 0
            ? "何も起こらなかった！"
            : `HPが${healedAmount}回復した！`,
        action: () => {
          setUser({ ...user, hitPoint: user.hitPoint + healedAmount });
        },
      });
    }
    if (item.effectType === "buff") {
      battle.useBuffItem(item);
      logs.push({
        message: `${user.name}の${
          ITEM_TARGET_LABEL_MAP[item.target ?? ""]
        }火力が${Math.floor((item.rate ?? 0) * 100)}%上昇した!`,
        action: () => {},
      });
    }
    if (item.effectType === "debuff") {
      battle.useDebuffItem(item);
      logs.push({
        message: `${monster.name}の${
          ITEM_TARGET_LABEL_MAP[item.target ?? ""]
        }耐性が${Math.floor((item.rate ?? 0) * 100)}%低下した！`,
        action: () => {},
      });
    }
    return logs;
  };

  const dotBorderClassName =
    "bg-[linear-gradient(to_right,#666_6px,transparent_2px,transparent_5px)] bg-[length:10px_2px] bg-bottom bg-repeat-x";
  const buttonGradationClassName =
    "bg-[linear-gradient(to_right,rgba(102,102,102,0)_0%,rgba(102,_102,_102,_0.8)_20%,rgba(102,102,102,0.8)_80%,rgba(102,102,102,0)_100%)]";

  return (
    <div>
      {/* モンスター画面 */}
      <div
        className={`h-[calc(100vh_-_320px)] pt-6 flex flex-col items-center justify-center transition-opacity duration-[2000ms]`}
        style={{ opacity: monster.hitPoint > 0 ? 1 : 0 }}
      >
        <div className="relative w-[24vh] h-auto aspect-square">
          <div className="w-full h-full bg-gray-300" />
          {["attack", "monsterHeal", "monsterGuard"].includes(effectMode) && user.weapon && (
            <AttackEffect
              elementType={user.weapon.elementType}
              physicsType={user.weapon.physicsType}
              effectMode={effectMode}
            />
          )}
        </div>
        <div className="w-[70%] bg-white/60 p-3">
          <div className="relative w-full bg-white border border-black h-3 flex items-center">
            <div
              className="h-full transition-all duration-300"
              style={{
                width: `${(monster.hitPoint / monster.maxHitPoint) * 100}%`,
                backgroundColor: hpBgColor(
                  monster.hitPoint,
                  monster.maxHitPoint
                ),
              }}
            />
            <div className="absolute inset-0 top-[100%] flex items-center text-base justify-end text-black">
              {monster.hitPoint}/{monster.maxHitPoint}
            </div>
          </div>
        </div>
      </div>

      {battlePhase.status === "first" && (
        <>
          <BattleConsole>
            <div className="flex flex-col gap-2">
              {user.weapon && <WeaponCard weapon={user.weapon} isEquipped />}
              <div className="flex flex-col">
                <div className={cn("h-1", dotBorderClassName)} />
                {[
                  {
                    label: "アイテム一覧",
                    onClick: () => {
                      if (items.length === 0) {
                        setBattleQueue([
                          {
                            message: "アイテムがありません。",
                            action: () => {},
                          },
                        ]);
                        return;
                      }
                      setBattlePhase({ status: "first", action: "item" });
                    },
                  },
                  {
                    label: "装備変更",
                    onClick: () => {
                      if (weapons.length === 0) {
                        setBattleQueue([
                          {
                            message: "武器がありません。",
                            action: () => {},
                          },
                        ]);
                        return;
                      }
                      setBattlePhase({ status: "first", action: "weapon" });
                    },
                  },
                ].map(({ label, onClick }) => (
                  <div key={label}>
                    <button
                      className={cn(
                        buttonGradationClassName,
                        "text-base font-bold w-full h-10"
                      )}
                      onClick={onClick}
                    >
                      {label}
                    </button>
                    <div className={cn(dotBorderClassName, "h-[2px]")} />
                  </div>
                ))}
              </div>

              <div className="flex grow justify-between px-4">
                <button onClick={() => router.push("/scan")}>
                  <Image
                    className="w-[130px] h-auto"
                    src={"/back-button.png"}
                    alt="戻る"
                    width={1000}
                    height={1000}
                  />
                </button>
                <button
                  onClick={() => {
                    setBattleQueue([
                      {
                        message: `${monster.name}が\n現れた！`,
                        action: () =>
                          setBattlePhase({
                            status: "command",
                            action: null,
                          }),
                      },
                    ]);
                    setStartDate(new Date());
                  }}
                >
                  <Image
                    className="w-[130px] h-auto"
                    src={"/start-button.png"}
                    alt="開始ボタン"
                    width={1000}
                    height={1000}
                  />
                </button>
              </div>
            </div>
          </BattleConsole>
          {battlePhase.action === "weapon" && (
            <BattleConsole>
              <WeaponDrawer
                onClose={() =>
                  setBattlePhase({ status: "first", action: null })
                }
                changeWeapon={(w) => {
                  setBattlePhase({ status: "first", action: null });
                  const logs = changeWeaponLogs(w);
                  setBattleQueue(logs);
                }}
              />
            </BattleConsole>
          )}
          {battlePhase.action === "item" && (
            <BattleConsole>
              <ItemDrawer
                onClose={() =>
                  setBattlePhase({ status: "first", action: null })
                }
                useItem={() => {
                  setBattlePhase({ status: "first", action: null });
                  setBattleQueue([
                    {
                      message: "戦闘前にアイテムは使えない！",
                      action: () => {},
                    },
                  ]);
                }}
              />
            </BattleConsole>
          )}
        </>
      )}

      {battlePhase.status === "command" && (
        <>
          <BattleConsole>
            <div className="flex flex-col gap-2">
              {user.weapon && <WeaponCard weapon={user.weapon} isEquipped />}
              <div>
                <div className={cn(dotBorderClassName, "h-[2px]")} />
                {[
                  { label: "攻撃", onClick: handleAttack },
                  {
                    label: "アイテム",
                    onClick: () => {
                      if (items.length === 0) {
                        setBattleQueue([
                          {
                            message: "アイテムがありません。",
                            action: () => {},
                          },
                        ]);
                        return;
                      }
                      setBattlePhase({ status: "command", action: "item" });
                    },
                  },
                  {
                    label: "装備の変更",
                    onClick: () => {
                      if (weapons.length === 0) {
                        setBattleQueue([
                          {
                            message: "武器がありません。",
                            action: () => {},
                          },
                        ]);
                        return;
                      }
                      setBattlePhase({ status: "command", action: "weapon" });
                    },
                  },
                  { label: "逃げる", onClick: () => setShowAwayModal(true) },
                ].map(({ label, onClick }) => (
                  <div key={label}>
                    <button
                      className={cn(
                        buttonGradationClassName,
                        "text-base font-bold w-full h-10"
                      )}
                      onClick={onClick}
                    >
                      {label}
                    </button>
                    <div className={cn(dotBorderClassName, "h-[2px]")} />
                  </div>
                ))}
              </div>
            </div>
          </BattleConsole>

          {battlePhase.action === "weapon" && (
            <BattleConsole>
              <WeaponDrawer
                onClose={() =>
                  setBattlePhase({ status: "command", action: null })
                }
                changeWeapon={(w) => {
                  setBattlePhase({ status: "command", action: null });
                  const logs = changeWeaponLogs(w);
                  setBattleQueue(logs);
                }}
              />
            </BattleConsole>
          )}
          {battlePhase.action === "item" && (
            <BattleConsole>
              <ItemDrawer
                onClose={() =>
                  setBattlePhase({ status: "command", action: null })
                }
                useItem={(i) => {
                  if (isItemUsed) {
                    setBattlePhase({ status: "command", action: null });
                    setBattleQueue([
                      {
                        message: "アイテムは1ターンに1回しか使えない！",
                        action: () => {},
                      },
                    ]);
                    return;
                  }
                  setIsItemUsed(true);
                  setBattlePhase({ status: "command", action: null });
                  const logs = useItemLogs(i);
                  setBattleQueue(logs);
                }}
              />
            </BattleConsole>
          )}
        </>
      )}

      {isStandBy && (
        <BattleConsole>
          <div
            className="relative h-[340px] flex items-center justify-center px-10"
            onClick={() => {
              if (battleQueue.length === 0) return;
              const [current, ...rest] = battleQueue;
              current.action();
              setBattleQueue(rest);
            }}
          >
            <p className="text-center text-lg">{battleQueue[0]?.message}</p>
            <motion.div
              className="absolute right-5 bottom-5"
              animate={{
                opacity: [1, 0, 1],
              }}
              transition={{
                duration: 1,
                repeat: Infinity,
                ease: "easeInOut",
              }}
            >
              <Image
                src={"/triangle.svg"}
                alt="タップ画像"
                width={30}
                height={30}
                priority
              />
            </motion.div>
          </div>
        </BattleConsole>
      )}

      {showAwayModal && (
        <Modal
          onClose={() => setShowAwayModal(false)}
          onConfirm={() => {
            setShowAwayModal(false);
            setBattleQueue([
              {
                message: `${user.name}は逃げ出した！`,
                action: () => {
                  if (pathname === "/battles/boss") {
                    location.href = "/scan";
                    return;
                  }
                  setUser({ ...user, hitPoint: user.maxHitPoint });
                  router.push("/scan");
                },
              },
            ]);
          }}
          title={`本当に逃げますか？`}
        />
      )}

      {reward && startDate && (
        <div className="fixed top-0 w-screen h-screen z-10 bg-black/70">
          <RewardModal
            restLevel={reward.restLevel}
            restExperiencePoint={reward.restExperiencePoint}
            level={reward.level}
            experiencePoint={reward.experiencePoint}
            drop={
              <div>
                {reward.drop && <TreasureBoxButton drop={reward.drop} />}
              </div>
            }
            clearTime={reward.clearTime}
          />
        </div>
      )}

      {["heal", "buff", "debuff"].includes(effectMode) && (
        <div className="fixed top-0 w-screen h-screen -z-10">
          <UseItemEffect effectMode={effectMode} />
        </div>
      )}
    </div>
  );
}

function formatMilliseconds(ms: number): string {
  const minutes = Math.floor(ms / 60000);
  const seconds = Math.floor((ms % 60000) / 1000);
  const centiseconds = Math.floor((ms % 1000) / 10);
  return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(
    2,
    "0"
  )}:${String(centiseconds).padStart(2, "0")}`;
}

