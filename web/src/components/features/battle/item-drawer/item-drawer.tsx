"use client";

import { Modal } from "@/components/features/battle/modal/modal";
import { useState, useRef, useEffect } from "react";
import { ItemCard } from "@/components/features/battle/item-card/item-card";
import Image from "next/image";
import { useUserStore } from "@/stores/user-store";
import { UserItemResponse } from "@/lib/bunkasaiRPGAPI.schemas";

type Props = {
  onClose: () => void;
  useItem: (item: UserItemResponse) => void;
  invalid?: boolean;
};

export function ItemDrawer({ onClose, useItem, invalid = false }: Props) {
  const { items } = useUserStore();

  const [selectedItem, setSelectedItem] = useState<UserItemResponse | null>(null);
  const [currentIndex, setCurrentIndex] = useState(0);
  const scrollContainerRef = useRef<HTMLDivElement>(null);
  const itemsPerPage = 2;
  const pageCount = Math.ceil(items.length / itemsPerPage);

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
      {/* アイテムスワイプ領域 */}
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
        {items.map((item, index) => (
          <div
            key={index}
            onClick={() => {
              if (!invalid) setSelectedItem(item);
            }}
            className="snap-start"
          >
            <ItemCard item={item} />
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

      {/* 確認モーダル */}
      {selectedItem && (
        <Modal
          onClose={() => {
            setSelectedItem(null);
          }}
          onConfirm={() => {
            useItem(selectedItem);
            setSelectedItem(null);
          }}
          title={`本当に「${selectedItem.name}」を\n使用しますか？`}
        />
      )}
    </div>
  );
}
