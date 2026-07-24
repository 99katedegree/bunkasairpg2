"use client";

import { cn } from "@/utils/cn";
import { assetGradation } from "@/utils/asset-gradation";
import { playSound } from "@/utils/play-sound/play-sound";
import Image from "next/image";
import { useState, useEffect } from "react";
import { EffectType } from "@/types/effect-type";
import { ElementType } from "@/types/element-type";
import { PhysicsType } from "@/types/physics-type";

export type Drop =
  | {
      type: "weapon";
      physicsType: PhysicsType;
      elementType: ElementType;
    }
  | {
      type: "item";
      effectType: EffectType;
    };

type Props = {
  drop: Drop;
};

export function TreasureBoxButton({ drop }: Props) {
  const [showClosedBox, setShowClosedBox] = useState(true);
  const [showOpenedBox, setShowOpenedBox] = useState(false);
  const [showDrop, setShowDrop] = useState(false);

  useEffect(() => {
    const openTimer = setTimeout(() => {
      setShowClosedBox(false);
      setShowOpenedBox(true);
      playSound("/sounds/reward-open.mp3");
    }, 1000);

    const dropTimer = setTimeout(() => {
      setShowOpenedBox(false);
      setShowDrop(true);
    }, 1500);

    return () => {
      clearTimeout(openTimer);
      clearTimeout(dropTimer);
    };
  }, []);

  return (
    <div className="flex flex-col items-center">
      {showClosedBox && (
        <Image
          src="/treasure-box-close.png"
          width={180}
          height={180}
          alt="treasure-box-close"
          priority
        />
      )}

      {showOpenedBox && (
        <Image
          src="/treasure-box-open.png"
          width={180}
          height={180}
          alt="treasure-box-open"
          priority
        />
      )}

      <div
        className="pt-4 transition-opacity duration-700 opacity-100"
        style={{ opacity: showDrop ? 1 : 0 }}
      >
        {drop.type === "weapon" && (
          <div
            className={cn(
              "p-1 rounded-md w-28 transition-opacity duration-700",
              assetGradation(drop.elementType)
            )}
          >
            <div className="relative aspect-square bg-gray-300" />
          </div>
        )}

        {drop.type === "item" && (
          <div
            className={cn(
              "p-1 rounded-md w-28 transition-opacity duration-700",
              assetGradation(drop.effectType)
            )}
          >
            <div className="relative aspect-square bg-gray-300" />
          </div>
        )}
      </div>
    </div>
  );
}
