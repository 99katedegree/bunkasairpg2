import { useEffect, useState } from "react";
import Cookies from "js-cookie";

export type ClearRankingItem = {
  rank: number;
  name: string;
  imageUrl: string | null;
  clearTime: string;
};

export type CollectedRankingItem = {
  rank: number;
  name: string;
  image_url: string | null;
  collection_rate: number;
};

type UserClearRankingResponse = {
  userRanking: ClearRankingItem;
  rankings: ClearRankingItem[];
};

type UserCollectedRankingResponse = {
  userRanking: CollectedRankingItem;
  rankings: CollectedRankingItem[];
};

async function fetchWithAuth<T>(url: string): Promise<T> {
  const authToken = Cookies.get("authToken");
  const res = await fetch(url, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });
  if (!res.ok) {
    throw new Error(`Failed to fetch: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export function useRankingPage() {
  const [collectedRankings, setCollectedRankings] = useState<
    CollectedRankingItem[]
  >([]);
  const [meCollectedRanking, setMeCollectedRanking] =
    useState<CollectedRankingItem | undefined>();
  const [clearRankings, setClearRankings] = useState<ClearRankingItem[]>([]);
  const [meClearRanking, setMeClearRanking] = useState<
    ClearRankingItem | undefined
  >();
  const [mode, setMode] = useState<"collected" | "clear">("collected");
  const [loadingClearRanking, setLoadingClearRanking] = useState(true);
  const [loadingCollectedRanking, setLoadingCollectedRanking] = useState(true);

  useEffect(() => {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL;

    fetchWithAuth<UserClearRankingResponse>(`${apiUrl}/user/clear-ranking`)
      .then((res) => {
        setMeClearRanking(res.userRanking);
        setClearRankings(res.rankings);
      })
      .finally(() => setLoadingClearRanking(false));

    fetchWithAuth<UserCollectedRankingResponse>(
      `${apiUrl}/user/collected-ranking`
    )
      .then((res) => {
        setMeCollectedRanking(res.userRanking);
        setCollectedRankings(res.rankings);
      })
      .finally(() => setLoadingCollectedRanking(false));
  }, []);

  const loading = loadingClearRanking || loadingCollectedRanking;

  return {
    collectedRankings,
    meCollectedRanking,
    clearRankings,
    meClearRanking,
    mode,
    setMode,
    loading,
  };
}
