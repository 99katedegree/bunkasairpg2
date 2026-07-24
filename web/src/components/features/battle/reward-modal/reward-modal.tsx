"use client";

import Image from "next/image";
import { motion } from "framer-motion";
import { calculateExperience } from "@/utils/calculate-experience";
import { useEffect, useState } from "react";
import { usePathname, useRouter } from "next/navigation";

type Props = {
  restLevel: number;
  restExperiencePoint: number;
  level: number;
  experiencePoint: number;
  drop: React.ReactNode;
  clearTime: string;
};

export function RewardModal({
  restLevel,
  restExperiencePoint,
  level,
  experiencePoint,
  drop,
  clearTime,
}: Props) {
  const router = useRouter();
  const pathname = usePathname();
  const expBarAnimation = createExpBarAnimation(
    restLevel,
    level,
    restExperiencePoint,
    experiencePoint
  );
  const levelValue = useCountUp(
    restLevel,
    level,
    expBarAnimation.transition.duration * 100
  );

  return (
    <div className="w-full flex justify-center items-center h-full text-sm">
      <div
        className="relative w-[90%] aspect-[380/605] bg-cover bg-center bg-no-repeat flex justify-center"
        style={{ backgroundImage: `url(${"/bg-reward.png"})` }}
      >
        <div className="w-[90%] flex flex-col items-center">
          <div className="mt-[10%]">
            <Image
              src={"/clear.png"}
              alt="タップ画像"
              width={1000}
              height={1000}
              priority
            />
          </div>
          <div className="bg-[#dac49d] border border-black rounded-xl p-4 w-[90%] flex flex-col items-center mt-[7%]">
            <div className="flex border-b-1 justify-between w-full">
              <div className="flex items-center space-x-3">
                <div>LV</div>
                <div>{levelValue}</div>
              </div>
              <div>次のレベルまで</div>
              <div>{calculateExperience(level + 1) - experiencePoint}</div>
            </div>

            <div className="relative w-[90%] h-3 bg-black border-black border rounded-full my-3 overflow-hidden">
              <motion.div
                className="h-full bg-[#1FA2C6] rounded-full origin-left"
                {...expBarAnimation}
              />
            </div>

            <div className="flex border-b-1 justify-between w-full">
              <div>クリアタイム</div>
              <div>{clearTime}</div>
            </div>
          </div>

          <div className="pt-[5vh]">{drop}</div>
        </div>
      </div>

      <button
        className="fixed bottom-4 left-1/2 -translate-x-1/2"
        onClick={() =>
          pathname === "/battles/boss"
            ? (location.href = "/scan")
            : router.push("/scan")
        }
      >
        <Image
          className="w-[130px] h-auto"
          src={"/back-button.png"}
          alt="戻る"
          width={1000}
          height={1000}
        />
      </button>
    </div>
  );
}

function createExpBarAnimation(
  restLevel: number,
  nowLevel: number,
  restExperiencePoint: number,
  nowExperiencePoint: number
) {
  const restLevelExp = calculateExperience(restLevel);
  const nowLevelExp = calculateExperience(nowLevel);
  const firstWidth =
    (restExperiencePoint - restLevelExp) /
    (calculateExperience(restLevel + 1) - restLevelExp);
  const lastWidth =
    (nowExperiencePoint - nowLevelExp) /
    (calculateExperience(nowLevel + 1) - nowLevelExp);

  const scaleX = [
    firstWidth,
    ...Array.from({ length: nowLevel - restLevel }).flatMap(() => [1, 0]),
    lastWidth,
  ];
  const times = [
    0,
    ...Array.from(
      { length: (nowLevel - restLevel) * 2 },
      (_, i) => 0.1 * (Math.floor(i / 2) + 1)
    ),
    (nowLevel - restLevel + 1) * 0.1,
  ];
  const ease: ("linear" | "easeInOut")[] = [
    "easeInOut",
    ...Array.from({ length: nowLevel - restLevel }).flatMap(
      (): ("linear" | "easeInOut")[] => ["linear", "easeInOut"]
    ),
  ];

  return {
    animate: {
      scaleX,
    },
    transition: {
      duration: scaleX.length / 15,
      times,
      ease,
    },
  };
}

function useCountUp(start: number, end: number, duration: number) {
  const [count, setCount] = useState(start);

  useEffect(() => {
    if (count >= end) return;
    const timer = setTimeout(() => {
      setCount(count + 1);
    }, duration);
    return () => clearTimeout(timer);
  }, [count, end]);

  return count;
}
