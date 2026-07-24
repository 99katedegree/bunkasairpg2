import Image from "next/image";
import { UserIcon } from "./icons/user-icon";
import { cn } from "@/utils/cn";
import { hpBgColor } from "@/utils/hp-bg-color";

type UserStatusProps = {
  name: string | null;
  imageUrl: string;
  level: number;
  hitPoint: number;
  maxHitPoint: number;
};

export function UserStatus({
  name,
  imageUrl,
  level = 0,
  hitPoint = 0,
  maxHitPoint = 0,
}: UserStatusProps) {
  return (
    <div className="flex justify-center pt-2 px-8">
      <div className="w-full relative h-18 flex items-center p-2">
        <Image
          className="absolute -z-10"
          src="/user-status-backbrand.png"
          alt="backbrand"
          fill
        />

        <div className="w-full grid grid-cols-[56px_1fr] gap-2 items-center">
          <div className="rounded-full relative bg-gray-300 aspect-square flex items-center justify-center overflow-hidden text-white">
            {imageUrl ? (
              <Image
                className="object-cover w-full h-auto"
                src={imageUrl}
                alt="user"
                width={150}
                height={150}
              />
            ) : (
              <UserIcon className="w-full h-full" />
            )}
          </div>
          <div className="flex flex-col gap-2">
            <div className="flex justify-between px-1">
              <p
                className={cn(
                  "whitespace-nowrap text-base font-bold truncate w-40",
                  hitPoint === 0 && "text-red-600"
                )}
              >
                {name}
              </p>
              <p>
                <span className="text-sm">{"Lv. "}</span>
                {level}
              </p>
            </div>
            <div className="relative w-full h-4 bg-gray-300 flex justify-end items-center">
              <div
                className="absolute top-0 left-0 h-full transition-all duration-300"
                style={{
                  width: `calc(${hitPoint} / ${maxHitPoint} * 100%)`,
                  backgroundColor: hpBgColor(hitPoint, maxHitPoint),
                }}
              />
              <p className="pr-2 z-10 text-sm">
                {hitPoint}/{maxHitPoint}
              </p>
            </div>
          </div>
        </div>

        <div className="absolute -z-10 -bottom-0.5 right-[99%] w-8 h-[84px]">
          <Image src="/user-status-left.png" alt="left" fill />
        </div>

        <div className="absolute -z-10 -bottom-0.5 left-[99%] w-8 h-[84px]">
          <Image src="/user-status-right.png" alt="right" fill />
        </div>
      </div>
    </div>
  );
}
