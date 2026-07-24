"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { ElementType } from "@/types/element-type";
import { PhysicsType } from "@/types/physics-type";
import { EffectMode } from "@/utils/change-effect/change-attack-effect";

interface Props {
  elementType: ElementType;
  physicsType: PhysicsType;
  effectMode: EffectMode;
}

export function AttackEffect({ elementType, physicsType, effectMode }: Props) {
  const elementTypeList = {
    neutral: { src: "neutral", alt: "無属性エフェクト" },
    flame: { src: "flame", alt: "火属性エフェクト" },
    water: { src: "water", alt: "水属性エフェクト" },
    wood: { src: "wood", alt: "木属性エフェクト" },
    shine: { src: "shine", alt: "光属性エフェクト" },
    dark: { src: "dark", alt: "闇属性エフェクト" },
  };

  if (effectMode === "monsterHeal") {
    return (
      <motion.div
        className="absolute inset-0 flex items-center justify-center"
        initial={{ opacity: 1, scale: 0.8, rotate: 0 }}
        animate={{
          scale: [0.8, 1, 1, 1.2],
          rotate: [0, -360],
        }}
        transition={{
          duration: 2.5,
          ease: "easeInOut",
        }}
      >
        <Image
          src="/effects/absorption.png"
          alt="モンスター回復エフェクト"
          width={500}
          height={500}
          priority
          className="pointer-events-none"
        />
      </motion.div>
    );
  }

  if (effectMode === "monsterGuard") {
    return (
      <div className="absolute inset-0 flex items-center justify-center">
        <Image
          src="/effects/shield.png"
          alt="モンスター防御エフェクト"
          width={500}
          height={500}
          priority
          className="scale-[1.7]"
        />
      </div>
    );
  }

  if (physicsType === "slash") {
    return (
      <motion.div
        className="absolute inset-0"
        initial={{ opacity: 0, x: 100, y: -100 }}
        animate={{ opacity: 1, x: 0, y: 0 }}
        transition={{
          duration: 0.8,
          ease: "easeOut",
        }}
      >
        <Image
          src={`/effects/slash/${elementTypeList[elementType].src}.png`}
          alt={`斬撃/${elementTypeList[elementType].alt}`}
          fill
          priority
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
        />
      </motion.div>
    );
  }

  if (physicsType === "blow") {
    return (
      <div className="absolute inset-0">
        <Image
          src={`/effects/blow/${elementTypeList[elementType].src}.png`}
          alt={`打撃/${elementTypeList[elementType].alt}`}
          fill
          priority
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
          className="object-contain"
        />
      </div>
    );
  }

  if (physicsType === "shoot") {
    return (
      <motion.div
        className="absolute inset-0 flex items-center justify-center"
        initial={{ opacity: 1, scale: 0.8, rotate: 0 }}
        animate={{
          scale: [0.8, 1, 1, 1.2],
          rotate: [0, -360],
        }}
        transition={{
          duration: 2.5,
          ease: "easeInOut",
        }}
      >
        <Image
          src={`/effects/shoot/${elementTypeList[elementType].src}.png`}
          alt={`射撃/${elementTypeList[elementType].alt}`}
          fill
          priority
          sizes="(max-width: 768px) 100vw, (max-width: 1200px) 50vw, 33vw"
        />
      </motion.div>
    );
  }

  return null;
}
