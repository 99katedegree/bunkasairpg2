"use client";

import { cn } from "@/utils/cn";
import {
  MonsterDetailResponse,
  MonsterCatalogResponse,
  WeaponResponse,
  ItemResponse,
} from "@/lib/bunkasaiRPGAPI.schemas";
import { MonsterCard } from "@/components/features/battle/monster-card/monster-card";
import { WeaponCard } from "@/components/features/battle/weapon-card/weapon-card";
import { ItemCard } from "@/components/features/battle/item-card/item-card";
import { AssetTypeIcon } from "@/components/shared/asset-type-icon";
import { BgCamera } from "@/components/shared/bg-camera";
import { QuestionIcon } from "@/components/shared/icons/question-icon";
import { Modal } from "@/components/shared/modal";
import { assetBgColor } from "@/utils/asset-bg-color";
import Image from "next/image";
import Link from "next/link";
import { getAuthInit } from "@/utils/auth-init";
import { useEffect, useState } from "react";
import { getMonster, getMonsters } from "@/lib/monster/monster";
import { getMeWeaponIndex } from "@/lib/weapon/weapon";
import { getMeItemIndex } from "@/lib/item/item";

const PAGE_SIZE = 9;

export default function Page() {
  const [category, setCategory] = useState<"monster" | "weapon" | "item">(
    "monster"
  );
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPage, setTotalPage] = useState(1);
  const [monsterIndex, setMonsterIndex] = useState<(MonsterCatalogResponse | null)[]>([]);
  const [weaponIndex, setWeaponIndex] = useState<(WeaponResponse | null)[]>([]);
  const [itemIndex, setItemIndex] = useState<(ItemResponse | null)[]>([]);
  const [monster, setMonster] = useState<MonsterDetailResponse | null>(null);
  const [weapon, setWeapon] = useState<WeaponResponse | null>(null);
  const [item, setItem] = useState<ItemResponse | null>(null);

  useEffect(() => {
    const offset = (currentPage - 1) * PAGE_SIZE;
    const limit = PAGE_SIZE;

    if (category === "monster") {
      getMonsters({ offset, limit }, getAuthInit()).then((res) => {
        if (res.status !== 200) return;
        setMonsterIndex(res.data.monsters);
        setTotalPage(Math.ceil(res.data.total / PAGE_SIZE));
      });
    } else if (category === "weapon") {
      getMeWeaponIndex({ offset, limit }, getAuthInit()).then((res) => {
        if (res.status !== 200) return;
        setWeaponIndex(res.data.weapons);
        setTotalPage(Math.ceil(res.data.total / PAGE_SIZE));
      });
    } else if (category === "item") {
      getMeItemIndex({ offset, limit }, getAuthInit()).then((res) => {
        if (res.status !== 200) return;
        setItemIndex(res.data.items);
        setTotalPage(Math.ceil(res.data.total / PAGE_SIZE));
      });
    }
  }, [category, currentPage]);

  return (
    <>
      <BgCamera />
      <div className="h-[100dvh] w-screen bg-cover bg-center bg-no-repeat flex flex-col items-center font-dotgothic">
        <div className="h-full w-full flex flex-col items-center justify-end">
          <h1 className="text-xl text-white">図鑑</h1>
          <Image
            className="w-[70%] h-[2px]"
            width={100}
            height={100}
            src="/profile-border.png"
            alt="見出しの下線"
          />

          <div className="w-[80%] mt-10 grid grid-cols-3 items-end">
            <div
              className={cn(
                "relative flex justify-center items-center",
                category === "monster" ? "h-14 text-white" : "h-12"
              )}
              onClick={() => {
                setCategory("monster");
                setCurrentPage(1);
              }}
            >
              <Image
                fill
                src={
                  category === "monster"
                    ? "/ranking-selected-btn.png"
                    : "/ranking-btn.png"
                }
                alt="ボタン画像"
              />
              <div className={cn("absolute", category === "monster" && "text-lg")}>
                モンスター
              </div>
            </div>
            <div
              className={cn(
                "relative flex justify-center items-center",
                category === "weapon" ? "h-14 text-white" : "h-12"
              )}
              onClick={() => {
                setCategory("weapon");
                setCurrentPage(1);
              }}
            >
              <Image
                fill
                src={
                  category === "weapon"
                    ? "/ranking-selected-btn.png"
                    : "/ranking-btn.png"
                }
                alt="ボタン画像"
              />
              <div className={cn("absolute", category === "weapon" && "text-lg")}>
                武器
              </div>
            </div>
            <div
              className={cn(
                "relative flex justify-center items-center",
                category === "item" ? "h-14 text-white" : "h-12"
              )}
              onClick={() => {
                setCategory("item");
                setCurrentPage(1);
              }}
            >
              <Image
                fill
                src={
                  category === "item"
                    ? "/ranking-selected-btn.png"
                    : "/ranking-btn.png"
                }
                alt="ボタン画像"
              />
              <div className={cn("absolute", category === "item" && "text-lg")}>
                アイテム
              </div>
            </div>
          </div>

          <div
            className="relative h-[80%] w-full bg-cover bg-center bg-no-repeat flex flex-col items-center [box-shadow:0_-8px_10px_-1px_rgba(0,0,0,0.25)]"
            style={{ backgroundImage: `url(${"/bg-reward.png"})` }}
          >
            <div className="relative h-[93%] w-full p-6">
              <div className="overflow-y-scroll w-full [&::-webkit-scrollbar]:hidden grid grid-cols-3 gap-2 pt-2 h-fit max-h-[85%] pb-12">
                {category === "monster" &&
                  monsterIndex.length > 0 &&
                  monsterIndex.map((m, index) => (
                    <div
                      key={index}
                      className="bg-neutral relative p-1 rounded-2xl shadow-lg h-fit"
                      onClick={() => {
                        if (m === null) return;
                        getMonster(m.id, getAuthInit()).then((res) => {
                          if (res.status === 200) setMonster(res.data.monster);
                        });
                      }}
                    >
                      {m ? (
                        <div className="relative bg-gray-300 rounded-xl aspect-square flex items-center overflow-hidden">
                          <div className="w-full h-full bg-gray-300 rounded-xl" />
                        </div>
                      ) : (
                        <div className="bg-gray-300 aspect-square rounded-xl flex items-center justify-center overflow-hidden">
                          <QuestionIcon className="text-white w-16 h-16" />
                        </div>
                      )}
                    </div>
                  ))}

                {category === "weapon" &&
                  weaponIndex.length > 0 &&
                  weaponIndex.map((w, index) => (
                    <div
                      key={index}
                      className={cn(
                        w ? assetBgColor(w.elementType) : "bg-neutral",
                        "relative p-1 rounded-2xl shadow-lg h-fit"
                      )}
                      onClick={() => {
                        if (w === null) return;
                        setWeapon(w);
                      }}
                    >
                      {w ? (
                        <div className="relative bg-gray-300 rounded-xl aspect-square flex items-center overflow-hidden">
                          <div className="w-full h-full bg-gray-300 rounded-xl" />
                          <div className="absolute w-full px-2 bottom-2 flex gap-2 z-10 justify-end">
                            <AssetTypeIcon type={w.physicsType} size="35%" />
                            <AssetTypeIcon type={w.elementType} size="35%" />
                          </div>
                        </div>
                      ) : (
                        <div className="bg-gray-300 aspect-square rounded-xl flex items-center justify-center overflow-hidden">
                          <QuestionIcon className="text-white w-16 h-16" />
                        </div>
                      )}
                    </div>
                  ))}

                {category === "item" &&
                  itemIndex.length > 0 &&
                  itemIndex.map((i, index) => (
                    <div
                      key={index}
                      className={cn(
                        i ? assetBgColor(i.effectType) : "bg-neutral",
                        "relative p-1 rounded-2xl shadow-lg h-fit"
                      )}
                      onClick={() => {
                        if (i === null) return;
                        setItem(i);
                      }}
                    >
                      {i ? (
                        <div className="relative bg-gray-300 rounded-xl aspect-square flex items-center overflow-hidden">
                          <div className="w-full h-full bg-gray-300 rounded-xl" />
                          <div className="absolute w-full px-2 bottom-2 flex gap-2 z-10 justify-end">
                            <AssetTypeIcon type={i.effectType} size="35%" />
                          </div>
                        </div>
                      ) : (
                        <div className="bg-gray-300 aspect-square rounded-xl flex items-center justify-center overflow-hidden">
                          <QuestionIcon className="text-white w-16 h-16" />
                        </div>
                      )}
                    </div>
                  ))}
              </div>

              {/* Tailwind pagination */}
              <div className="absolute left-1/2 -translate-x-1/2 bottom-[9%] flex items-center gap-1">
                <button
                  className="w-8 h-8 flex items-center justify-center rounded bg-[#E0DBD7] border text-black disabled:opacity-40"
                  onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                  disabled={currentPage === 1}
                >
                  ‹
                </button>
                {Array.from({ length: totalPage }).map((_, idx) => {
                  const page = idx + 1;
                  return (
                    <button
                      key={page}
                      className={cn(
                        "w-8 h-8 flex items-center justify-center rounded border text-sm",
                        page === currentPage
                          ? "bg-linear-to-t to-[#661412] from-[#A72731] text-white"
                          : "bg-[#E0DBD7] text-black"
                      )}
                      onClick={() => setCurrentPage(page)}
                    >
                      {page}
                    </button>
                  );
                })}
                <button
                  className="w-8 h-8 flex items-center justify-center rounded bg-[#E0DBD7] border text-black disabled:opacity-40"
                  onClick={() =>
                    setCurrentPage((p) => Math.min(totalPage, p + 1))
                  }
                  disabled={currentPage === totalPage}
                >
                  ›
                </button>
              </div>

              <Link
                href="/scan"
                className="absolute left-1/2 -translate-x-1/2 -bottom-[3%] w-[30%]"
              >
                <Image
                  className="w-full h-auto"
                  width={100}
                  height={100}
                  src="/back-button.png"
                  alt="戻るボタン"
                />
              </Link>
            </div>
          </div>
        </div>
      </div>

      {monster && (
        <Modal onClose={() => setMonster(null)}>
          <div className="w-screen px-2 text-white">
            <div className="w-full bg-black/70 p-10 rounded-2xl">
              <MonsterCard monster={monster} setMonster={setMonster} />
            </div>
          </div>
        </Modal>
      )}

      {weapon && (
        <Modal onClose={() => setWeapon(null)}>
          <div className="w-screen px-2 text-white">
            <div className="w-full bg-black/70 p-4 rounded-2xl">
              <WeaponCard weapon={weapon} />
            </div>
          </div>
        </Modal>
      )}

      {item && (
        <Modal onClose={() => setItem(null)}>
          <div className="w-screen px-2 text-white">
            <div className="w-full bg-black/70 p-4 rounded-2xl">
              <ItemCard
                item={{
                  ...item,
                  count: 0,
                }}
              />
            </div>
          </div>
        </Modal>
      )}
    </>
  );
}
