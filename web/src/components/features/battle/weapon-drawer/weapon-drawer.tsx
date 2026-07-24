"use client";

import { Modal } from "@/components/features/battle/modal/modal";
import { useRef, useState, useEffect } from "react";
import { WeaponCard } from "@/components/features/battle/weapon-card/weapon-card";
import { useUserStore } from "@/stores/user-store";
import Image from "next/image";
import { WeaponResponse } from "@/lib/bunkasaiRPGAPI.schemas";

export type WeaponDrawerProps = {
  onClose: () => void;
  changeWeapon: (weapon: WeaponResponse) => void;
};

export function WeaponDrawer({ onClose, changeWeapon }: WeaponDrawerProps) {
  const { user, weapons } = useUserStore();
  const [selectedWeapon, setSelectedWeapon] = useState<WeaponResponse | null>(null);
  const [currentIndex, setCurrentIndex] = useState(0);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const itemsPerPage = 2;
  const pageCount = Math.ceil(weapons.length / itemsPerPage);

  const handleDotClick = (index: number) => {
    const container = scrollContainerRef.current;
    if (!container) return;
    const scrollX = container.clientWidth * index;
    container.scrollTo({ left: scrollX, behavior: "smooth" });
  };

  useEffect(() => {
    const container = scrollContainerRef.current;
    if (!container) return;

    const handleScroll = () => {
      const index = Math.round(container.scrollLeft / container.clientWidth);
      setCurrentIndex(index);
    };

    container.addEventListener("scroll", handleScroll);
    return () => container.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <div className="w-full">
      {/* 武器スワイプ領域 */}
      <div
        ref={scrollContainerRef}
        className="
          [&::-webkit-scrollbar]:hidden
          [-ms-overflow-style:none]
          [scrollbar-width:none]
          grid grid-rows-2 grid-flow-col
          auto-cols-[100%]
          overflow-x-auto
          gap-4
          py-2
          w-full
          snap-x snap-mandatory
          scroll-smooth
        "
      >
        {weapons.map((we, index) => (
          <div
            key={index}
            onClick={() => {
              if (user?.weapon?.id !== we.id) setSelectedWeapon(we);
            }}
            className="snap-start"
          >
            <WeaponCard weapon={we} isEquipped={user?.weapon?.id === we.id} />
          </div>
        ))}
      </div>

      {/* ページネーションドット */}
      <div className="flex justify-center gap-2 mt-2">
        {Array.from({ length: pageCount }).map((_, index) => (
          <button
            key={index}
            onClick={() => handleDotClick(index)}
            className={`h-2 w-2 rounded-full transition-all ${
              currentIndex === index
                ? "bg-white scale-125"
                : "bg-gray-500 hover:bg-gray-400"
            }`}
          />
        ))}
      </div>

      <div>
        <button className="flex justify-center w-full my-4">
          <Image
            className="w-[130px] h-auto"
            src={"/back-button.png"}
            alt="戻る"
            width={1000}
            height={1000}
            onClick={onClose}
          />
        </button>
      </div>

      {/* 確認モーダル */}
      {selectedWeapon && (
        <Modal
          onClose={() => {
            setSelectedWeapon(null);
          }}
          onConfirm={() => {
            changeWeapon(selectedWeapon);
            setSelectedWeapon(null);
          }}
          title={`「${selectedWeapon.name}」に\n変更しますか？`}
        />
      )}
    </div>
  );
}
