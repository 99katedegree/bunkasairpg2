"use client";

import Image from "next/image";
import Link from "next/link";
import { RankingCard } from "@/components/features/ranking/ranking-card/ranking-card";
import { MeRankingCard } from "@/components/features/ranking/me-ranking-card/me-ranking-card";
import { useRankingPage } from "./use-ranking-page";

export function RankingPage() {
  const {
    collectedRankings,
    meCollectedRanking,
    clearRankings,
    meClearRanking,
    mode,
    setMode,
    loading,
  } = useRankingPage();

  if (loading) {
    return (
      <div className="flex items-center justify-center h-full w-full text-white">
        読み込み中...
      </div>
    );
  }

  return (
    <div className="h-full w-full flex flex-col items-center justify-end">
      <h1 className="text-xl text-white">ランキング</h1>
      <Image
        className="w-[70%] h-[2px]"
        width={100}
        height={100}
        src="/profile-border.png"
        alt="見出しの下線"
      />

      {mode === "collected" ? (
        <div className="flex w-[80%] mt-10 items-end">
          <div className="relative w-[50%] h-14 text-white flex justify-center items-center">
            <Image
              fill
              src="/ranking-selected-btn.png"
              alt="ランキングボタン画像"
            />
            <div className="absolute text-lg">コンプリート</div>
          </div>
          <div
            className="relative w-[50%] h-12 flex justify-center items-center"
            onClick={() => setMode("clear")}
          >
            <Image fill src="/ranking-btn.png" alt="ランキングボタン画像" />
            <div className="absolute">ボス討伐</div>
          </div>
        </div>
      ) : (
        <div className="flex w-[80%] mt-10 items-end">
          <div
            className="relative w-[50%] h-12 flex justify-center items-center"
            onClick={() => setMode("collected")}
          >
            <Image fill src="/ranking-btn.png" alt="ランキングボタン画像" />
            <div className="absolute">コンプリート</div>
          </div>
          <div className="relative w-[50%] h-14 text-white flex justify-center items-center">
            <Image
              fill
              src="/ranking-selected-btn.png"
              alt="ランキングボタン画像"
            />
            <div className="absolute text-lg">ボス討伐</div>
          </div>
        </div>
      )}

      <div
        className="relative h-[80%] w-full bg-cover bg-center bg-no-repeat flex flex-col items-center [box-shadow:0_-8px_10px_-1px_rgba(0,0,0,0.25)]"
        style={{ backgroundImage: `url(${"/bg-reward.png"})` }}
      >
        <div className="h-[83%] w-full p-6">
          <div className="overflow-y-scroll h-full w-full [&::-webkit-scrollbar]:hidden">
            <div className="mt-2" />
            {mode === "collected" ? (
              <>
                {collectedRankings.map((collectedRanking) => (
                  <div
                    className="w-full h-18 mb-5"
                    key={collectedRanking.rank}
                  >
                    <RankingCard
                      rank={collectedRanking.rank}
                      imageUrl={collectedRanking.image_url}
                      mode={true}
                      name={collectedRanking.name}
                      value={collectedRanking.collection_rate}
                    />
                  </div>
                ))}
              </>
            ) : (
              <>
                {clearRankings.map((clearRanking) => (
                  <div className="w-full h-18 mb-5" key={clearRanking.rank}>
                    <RankingCard
                      rank={clearRanking.rank}
                      imageUrl={clearRanking.imageUrl}
                      mode={false}
                      name={clearRanking.name}
                      value={clearRanking.clearTime}
                    />
                  </div>
                ))}
              </>
            )}
            <div className="mt-10" />
          </div>
        </div>
        <div className="absolute bottom-[13%] w-full h-20 px-6 z-50">
          {mode === "collected" ? (
            <MeRankingCard
              rank={meCollectedRanking?.rank ?? 0}
              imageUrl={meCollectedRanking?.image_url ?? null}
              mode={true}
              name={meCollectedRanking?.name ?? ""}
              value={meCollectedRanking?.collection_rate ?? 0}
            />
          ) : (
            <MeRankingCard
              rank={meClearRanking?.rank ?? 0}
              imageUrl={meClearRanking?.imageUrl ?? null}
              mode={false}
              name={meClearRanking?.name ?? ""}
              value={meClearRanking?.clearTime ?? ""}
            />
          )}
        </div>

        <Link href="/scan" className="absolute bottom-[2%] w-[30%]">
          <Image
            className="w-full h-auto"
            width={100}
            height={100}
            src="/back-button.png"
            alt="戻るボタン"
          />
        </Link>
      </div>
    </div>
  );
}
