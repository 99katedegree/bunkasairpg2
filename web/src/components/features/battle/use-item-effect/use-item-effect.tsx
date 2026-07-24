"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { EffectMode } from "@/utils/change-effect/change-attack-effect";

interface Props {
  effectMode: EffectMode;
}

export function UseItemEffect({ effectMode }: Props) {
  const motionList: { [key in EffectMode]?: number } = {
    heal: -100,
    buff: -100,
    debuff: 100,
  };

  return (
    <motion.div
      className="absolute inset-0 flex items-center justify-center"
      initial={{ opacity: 1, y: 0 }}
      animate={{ opacity: 0, y: motionList[effectMode] }}
      transition={{
        duration: 1.5,
        ease: "easeOut",
      }}
    >
      <Image
        src={`/effects/${effectMode}.png`}
        alt="アイテム使用エフェクト"
        width={500}
        height={500}
        priority
        className="scale-[1.7]"
      />
    </motion.div>
  );
}
