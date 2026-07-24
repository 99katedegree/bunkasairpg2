"use client";

import { motion } from "framer-motion";
import Image from "next/image";
import { UserItemResponse, WeaponResponse } from "@/lib/bunkasaiRPGAPI.schemas";

export type ProfileConsoleProps = {
  children: React.ReactNode;
  setClose: (close: boolean) => void;
  items?: UserItemResponse[];
  weapons?: WeaponResponse[];
};

export function ProfileConsole({
  children,
  setClose,
  items,
  weapons,
}: ProfileConsoleProps) {
  if (weapons?.length === 0) {
    return (
      <div
        className="fixed w-full h-full z-30"
        onClick={() => {
          setClose(false);
        }}
      >
        <div
          className="fixed bottom-0 w-full bg-black/70 text-white p-2 z-20"
          onClick={(e) => e.stopPropagation()}
        >
          <div
            className="flex justify-center items-center h-72"
            onClick={() => setClose(false)}
          >
            武器がありません。
          </div>
          <motion.div
            className="absolute right-5 bottom-5"
            animate={{
              opacity: [1, 0, 1],
            }}
            transition={{
              duration: 1,
              repeat: Infinity,
              ease: "easeInOut",
            }}
          >
            <Image
              src={"/triangle.svg"}
              alt="タップ画像"
              width={30}
              height={30}
              priority
            />
          </motion.div>
        </div>
      </div>
    );
  }

  if (items?.length === 0) {
    return (
      <div
        className="fixed w-full h-full z-30"
        onClick={() => {
          setClose(false);
        }}
      >
        <div
          className="fixed bottom-0 w-full bg-black/70 text-white p-2 z-20"
          onClick={(e) => e.stopPropagation()}
        >
          <div
            className="flex justify-center items-center h-72"
            onClick={() => setClose(false)}
          >
            アイテムがありません。
          </div>
          <motion.div
            className="absolute right-5 bottom-5"
            animate={{
              opacity: [1, 0, 1],
            }}
            transition={{
              duration: 1,
              repeat: Infinity,
              ease: "easeInOut",
            }}
          >
            <Image
              src={"/triangle.svg"}
              alt="タップ画像"
              width={30}
              height={30}
              priority
            />
          </motion.div>
        </div>
      </div>
    );
  }

  return (
    <div
      className="fixed w-full h-full z-30"
      onClick={() => {
        setClose(false);
      }}
    >
      <div
        className="fixed bottom-0 w-full bg-black/70 text-white p-2 z-20"
        onClick={(e) => e.stopPropagation()}
      >
        {children}
      </div>
    </div>
  );
}
