"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { cn } from "@/utils/cn";
import { boss } from "@/constants/boss";
import { assetBgColor } from "@/utils/asset-bg-color";
import { useRouter } from "next/navigation";
import { MonsterDetailResponse } from "@/lib/bunkasaiRPGAPI.schemas";

type Props = {
  monster: MonsterDetailResponse;
  setMonster: (m: MonsterDetailResponse | null) => void;
};

export function MonsterCard({ monster, setMonster }: Props) {
  const router = useRouter();
  const dotBorderClassName =
    "bg-[linear-gradient(to_right,#666_6px,transparent_2px,transparent_5px)] bg-[length:10px_2px] bg-bottom bg-repeat-x h-8 flex items-center";

  return (
    <div>
      <div className="flex flex-col gap-6">
        <div className="p-2 bg-neutral w-fit mx-auto rounded-2xl">
          <div className="relative h-44 w-44 mx-auto bg-gray-300 rounded-xl aspect-square overflow-hidden">
            <div className="w-full h-full bg-gray-300 rounded-xl" />
          </div>
        </div>
        <div className="flex flex-col gap-4">
          <div className="h-full flex flex-col justify-between overflow-hidden">
            <div className="w-full overflow-hidden">
              <motion.div
                className="flex whitespace-nowrap font-bold text-lg"
                style={{ display: "inline-flex" }}
                animate={{ x: ["0%", "-50%"] }}
                transition={{
                  x: {
                    repeat: Infinity,
                    repeatType: "loop",
                    duration: 10,
                    ease: "linear",
                  },
                }}
              >
                <span className="pr-12">{monster.name}</span>
                <span className="pr-12">{monster.name}</span>
                <span className="pr-12">{monster.name}</span>
                <span className="pr-12">{monster.name}</span>
              </motion.div>
            </div>
            {[
              {
                label: "攻撃力",
                value: monster.attack,
                width: monster.attack / boss.attack,
                color: "bg-[red]",
              },
              {
                label: "斬撃属性",
                value: "",
                width: (monster.slash + 1) / 3,
                color: assetBgColor("neutral"),
              },
              {
                label: "打撃属性",
                value: "",
                width: (monster.blow + 1) / 3,
                color: assetBgColor("neutral"),
              },
              {
                label: "射撃属性",
                value: "",
                width: (monster.shoot + 1) / 3,
                color: assetBgColor("neutral"),
              },
              {
                label: "無属性",
                value: "",
                width: (monster.neutral + 1) / 3,
                color: assetBgColor("neutral"),
              },
              {
                label: "火属性",
                value: "",
                width: (monster.flame + 1) / 3,
                color: assetBgColor("flame"),
              },
              {
                label: "水属性",
                value: "",
                width: (monster.water + 1) / 3,
                color: assetBgColor("water"),
              },
              {
                label: "木属性",
                value: "",
                width: (monster.wood + 1) / 3,
                color: assetBgColor("wood"),
              },
              {
                label: "光属性",
                value: "",
                width: (monster.shine + 1) / 3,
                color: assetBgColor("shine"),
              },
              {
                label: "闇属性",
                value: "",
                width: (monster.dark + 1) / 3,
                color: assetBgColor("dark"),
              },
            ].map(({ label, value, width, color }) => (
              <div
                key={label}
                className={cn(dotBorderClassName, "flex items-end px-1")}
              >
                <p className="w-20 text-[#bababa] text-lg pb-0.5">{label}</p>
                <div className="flex-grow h-full flex items-center">
                  <div className="relative bg-gray-300 w-full h-5 rounded-lg overflow-hidden">
                    <div
                      className={`${color} absolute top-0 left-0 h-full rounded-lg z-10`}
                      style={{ width: `${width * 100}%` }}
                    />
                  </div>
                </div>
                <p className="h-full flex items-center justify-end text-base w-10">
                  {value}
                </p>
              </div>
            ))}
          </div>
          <div className="flex gap-5 justify-center">
            <div
              className="relative w-30 h-12 flex justify-center items-center"
              onClick={() => setMonster(null)}
            >
              <Image src={"/profile-btn.png"} fill alt="戻る" />
              <div className="relative text-black text-sm font-bold">戻る</div>
            </div>
            <div
              className="relative w-30 h-12 flex justify-center items-center"
              onClick={() => router.push(`/battles/${monster.id}`)}
            >
              <Image src={"/profile-btn.png"} fill alt="挑戦する" />
              <div className="relative text-black text-sm font-bold">
                挑戦する
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
