import { cn } from "@/utils/cn";
import Image from "next/image";
import { WeaponResponse } from "@/lib/bunkasaiRPGAPI.schemas";
import { assetGradation } from "@/utils/asset-gradation";
import { motion } from "framer-motion";

type Props = {
  weapon: WeaponResponse;
  isEquipped?: boolean;
};

export function WeaponCard({ weapon, isEquipped }: Props) {
  const physicsTypeMap: Record<string, string> = {
    slash: "斬撃",
    blow: "打撃",
    shoot: "射撃",
  };
  const elementsTypeMap: Record<string, string> = {
    neutral: "無",
    flame: "炎",
    water: "水",
    wood: "木",
    shine: "光",
    dark: "闇",
  };

  const dotBorderClassName =
    "bg-[linear-gradient(to_right,#666_6px,transparent_2px,transparent_5px)] bg-[length:10px_2px] bg-bottom bg-repeat-x";

  return (
    <div className="grid grid-cols-[140px_1fr] gap-2">
      <div
        className={cn("p-1 rounded-md", assetGradation(weapon.elementType))}
      >
        <div className="relative aspect-square bg-gray-300">
          {isEquipped && (
            <div className="absolute top-0 left-0 bg-[#666666]/80 z-10">
              装備中
            </div>
          )}
        </div>
      </div>

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
            <span className="pr-12">{weapon.name}</span>
            <span className="pr-12">{weapon.name}</span>
            <span className="pr-12">{weapon.name}</span>
            <span className="pr-12">{weapon.name}</span>
            <span className="pr-12">{weapon.name}</span>
            <span className="pr-12">{weapon.name}</span>
          </motion.div>
        </div>
        {[
          {
            label: "攻撃力",
            value: weapon.physicsAttack,
          },
          {
            label: "属性値",
            value: weapon.elementAttack ?? 0,
          },
          {
            label: "属性",
            value: elementsTypeMap[weapon.elementType],
          },
          {
            label: "物理タイプ",
            value: physicsTypeMap[weapon.physicsType],
          },
        ].map(({ label, value }) => (
          <div
            key={label}
            className={cn(
              dotBorderClassName,
              "flex justify-between items-end px-1"
            )}
          >
            <p className="text-[#bababa] text-sm pb-0.5">{label}</p>
            <p className="text-base">{value}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
