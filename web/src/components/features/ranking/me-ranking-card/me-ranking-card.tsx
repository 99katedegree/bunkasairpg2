import Image from "next/image";
import { ModeTrueProps, ModeFalseProps } from "@/components/features/ranking/ranking-card/ranking-card";

type Props = ModeTrueProps | ModeFalseProps;

export function MeRankingCard({ rank, name, mode, value, imageUrl }: Props) {
  const meStyle = {
    bgColor: "linear-gradient(to bottom right, white, #A6A6A6)",
  };

  return (
    <div
      className="flex h-full w-full border-2 border-black [box-shadow:0_-8px_4px_-1px_rgba(0,0,0,0.6)]"
      style={{ backgroundImage: meStyle.bgColor }}
    >
      <div className="flex justify-center items-center text-2xl text-white text-shadow-outline w-[20%]">
        {rank === 0 ? "なし" : rank}
      </div>
      <div className="relative aspect-square w-auto h-full bg-gray-300 border border-white">
        {imageUrl ? (
          <Image
            className="object-cover w-full h-full"
            src={imageUrl}
            alt="user"
            width={150}
            height={150}
          />
        ) : (
          <div className="w-full h-full bg-gray-300" />
        )}
      </div>
      <div className="flex-grow p-2 flex flex-col justify-between">
        <div className="p-1 text-center bg-[#1F1D00] border border-[#CCC4AF] text-[#B29F8E] rounded-lg text-sm">
          <p className="truncate w-40">{name}</p>
        </div>
        <div className="flex justify-between items-center">
          {mode ? (
            <>
              <div className="text-xs">コンプリート率</div>
              <div>{value}%</div>
            </>
          ) : (
            <>
              {rank === 0 ? (
                <>
                  <div className="text-xs">クリアタイム</div>
                  <div>--:--:--</div>
                </>
              ) : (
                <>
                  <div className="text-xs">クリアタイム</div>
                  <div>{value}</div>
                </>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  );
}
