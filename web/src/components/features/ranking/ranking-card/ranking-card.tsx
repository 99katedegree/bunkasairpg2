import Image from "next/image";

interface CommonProps {
  rank: number;
  imageUrl: string | null;
  name: string;
}

export interface ModeTrueProps extends CommonProps {
  mode: true;
  value: number;
}

export interface ModeFalseProps extends CommonProps {
  mode: false;
  value: string;
}

type Props = ModeTrueProps | ModeFalseProps;

export function RankingCard({ rank, name, mode, value, imageUrl }: Props) {
  const rankStyle: { [key: number]: { bgColor: string; imageUrl: string } } = {
    1: {
      bgColor: "linear-gradient(to bottom right, #D4A424, #B37810)",
      imageUrl: "/bg-ranking-1.png",
    },
    2: {
      bgColor: "linear-gradient(to bottom right, #969EAD, #959DA5)",
      imageUrl: "/bg-ranking-2.png",
    },
    3: {
      bgColor: "linear-gradient(to bottom right, #F19E6A, #8B450E)",
      imageUrl: "/bg-ranking-3.png",
    },
  };
  const defaultColor = {
    bgColor: "linear-gradient(to bottom right, #7A7EAB, #393B64)",
    imageUrl: "/bg-ranking-4.png",
  };

  return (
    <div
      className="flex h-full w-full border-2 border-black"
      style={{
        backgroundImage: rankStyle[rank]?.bgColor ?? defaultColor.bgColor,
      }}
    >
      <div className="flex justify-center items-center">
        <div className="relative mr-2">
          <div className="relative w-[87px] h-[90px]">
            <Image
              fill
              src={rankStyle[rank]?.imageUrl ?? defaultColor.imageUrl}
              alt="ランキングカード"
            />
          </div>
          <div className="absolute inset-0 flex items-center justify-center text-4xl text-white z-10 text-shadow-outline">
            {rank}
          </div>
        </div>
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
          <p className="truncate w-36">{name}</p>
        </div>
        <div className="flex justify-between items-center">
          {mode ? (
            <>
              <div className="text-xs">コンプリート率</div>
              <div>{value}%</div>
            </>
          ) : (
            <>
              <div className="text-xs">クリアタイム</div>
              <div>{value}</div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
