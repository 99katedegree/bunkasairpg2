"use client";

import { useEffect, useState } from "react";
import { AlertTriangle, CheckCircle2, Loader2, RefreshCw } from "lucide-react";
import { getMonsterSummaries } from "@/lib/monster/monster";
import { getWeaponSummaries } from "@/lib/weapon/weapon";
import { getItemSummaries } from "@/lib/item/item";
import { CachedImage } from "@/components/shared/cached-image";

type AssetType = "monster" | "weapon" | "item";

// バックエンドがマスタの存在と画像URL（/assets/... 配信ルート）の両方を返す。
// 画像をどのファイルに割り当てるかは cmd/seed-master が図鑑番号から機械的に
// 決めるので、フロント側で対応表を持つ必要はない。
type Summary = {
  id: string | number;
  name: string;
  indexNumber: string;
  imageUrl?: string | null;
};

type AssetRow = {
  index: number;
  id: string | number | null;
  name: string | null;
  indexNumber: string;
  imageUrl: string | null;
  notes: string[];
};

function buildRows(summaries: Summary[] | undefined): AssetRow[] {
  const rows = summaries ?? [];
  return rows.map((r, i) => {
    const notes: string[] = [];
    if (!r.imageUrl) notes.push("画像がまだ用意されていません");
    return {
      index: i + 1,
      id: r.id,
      name: r.name,
      indexNumber: r.indexNumber,
      imageUrl: r.imageUrl ?? null,
      notes,
    };
  });
}

const TABS: { key: AssetType; label: string }[] = [
  { key: "monster", label: "モンスター" },
  { key: "weapon", label: "武器" },
  { key: "item", label: "アイテム" },
];

type State = {
  rows: AssetRow[];
  loading: boolean;
  error: string | null;
};

const emptyState = (): State => ({ rows: [], loading: false, error: null });

export function AssetTable({ adminToken }: { adminToken: string }) {
  const [tab, setTab] = useState<AssetType>("monster");
  const [states, setStates] = useState<Record<AssetType, State>>({
    monster: emptyState(),
    weapon: emptyState(),
    item: emptyState(),
  });

  const authInit: RequestInit = { headers: { Authorization: `Bearer ${adminToken}` } };

  const load = async (type: AssetType) => {
    setStates((prev) => ({ ...prev, [type]: { ...prev[type], loading: true, error: null } }));

    try {
      let summaries: Summary[] = [];

      if (type === "monster") {
        const res = await getMonsterSummaries(authInit);
        if (res.status !== 200) throw new Error("取得に失敗しました");
        summaries = res.data.monsters;
      } else if (type === "weapon") {
        const res = await getWeaponSummaries(authInit);
        if (res.status !== 200) throw new Error("取得に失敗しました");
        summaries = res.data.weapons;
      } else {
        const res = await getItemSummaries(authInit);
        if (res.status !== 200) throw new Error("取得に失敗しました");
        summaries = res.data.items;
      }

      setStates((prev) => ({
        ...prev,
        [type]: { rows: buildRows(summaries), loading: false, error: null },
      }));
    } catch {
      setStates((prev) => ({
        ...prev,
        [type]: { rows: [], loading: false, error: "データの取得に失敗しました" },
      }));
    }
  };

  useEffect(() => {
    load(tab);
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [tab]);

  const { rows, loading, error } = states[tab];

  const hasIssue = (row: AssetRow) => row.notes.length > 0;
  const totalIssues = rows.filter(hasIssue).length;

  return (
    <div>
      {/* タブ */}
      <div className="flex gap-1 mb-4">
        {TABS.map(({ key, label }) => (
          <button
            key={key}
            onClick={() => setTab(key)}
            className={`px-3 py-1.5 text-xs rounded-md transition-colors ${
              tab === key
                ? "bg-zinc-700 text-zinc-100"
                : "text-zinc-500 hover:text-zinc-300 hover:bg-zinc-900"
            }`}
          >
            {label}
          </button>
        ))}
        <button
          onClick={() => load(tab)}
          disabled={loading}
          className="ml-auto text-zinc-600 hover:text-zinc-400 transition-colors disabled:opacity-40"
        >
          <RefreshCw size={13} className={loading ? "animate-spin" : ""} />
        </button>
      </div>

      {/* サマリー */}
      {!loading && !error && rows.length > 0 && (
        <div className="flex items-center gap-4 mb-3 text-xs text-zinc-500">
          <span>{rows.length} 件</span>
          {totalIssues > 0 ? (
            <span className="flex items-center gap-1 text-amber-500">
              <AlertTriangle size={11} />
              {totalIssues} 件の問題
            </span>
          ) : (
            <span className="flex items-center gap-1 text-emerald-500">
              <CheckCircle2 size={11} />
              すべて正常
            </span>
          )}
        </div>
      )}

      {/* テーブル */}
      {loading ? (
        <div className="flex items-center gap-2 text-zinc-600 text-xs py-6">
          <Loader2 size={13} className="animate-spin" />
          読み込み中...
        </div>
      ) : error ? (
        <p className="text-red-400 text-xs py-6">{error}</p>
      ) : rows.length === 0 ? (
        <p className="text-zinc-600 text-xs py-6">データがありません</p>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full text-xs">
            <thead>
              <tr className="border-b border-zinc-800">
                <th className="text-left text-zinc-600 font-normal pb-2 pr-4 w-10">#</th>
                <th className="text-left text-zinc-600 font-normal pb-2 pr-4">図鑑番号</th>
                <th className="text-left text-zinc-600 font-normal pb-2 pr-4">名前</th>
                <th className="text-left text-zinc-600 font-normal pb-2 pr-4">ID</th>
                <th className="text-left text-zinc-600 font-normal pb-2 pr-4">画像</th>
                <th className="text-left text-zinc-600 font-normal pb-2">備考</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr
                  key={row.indexNumber}
                  className={`border-b border-zinc-900 ${
                    hasIssue(row) ? "bg-amber-950/20" : ""
                  }`}
                >
                  <td className="py-2 pr-4 text-zinc-600 tabular-nums">{row.index}</td>
                  <td className="py-2 pr-4 text-zinc-300 font-mono">{row.indexNumber}</td>
                  <td className="py-2 pr-4 text-zinc-300">
                    {row.name ?? <span className="text-zinc-700">—</span>}
                  </td>
                  <td className="py-2 pr-4 text-zinc-500 font-mono break-all">
                    {row.id ?? <span className="text-zinc-700">—</span>}
                  </td>
                  <td className="py-2 pr-4">
                    <CachedImage
                      src={row.imageUrl}
                      alt={row.name ?? row.indexNumber}
                      className="h-10 w-10 rounded object-cover bg-zinc-800"
                      fallback={<span className="text-zinc-700">—</span>}
                    />
                  </td>
                  <td className="py-2">
                    {row.notes.length === 0 ? (
                      <CheckCircle2 size={11} className="text-emerald-600" />
                    ) : (
                      <div className="flex flex-col gap-0.5">
                        {row.notes.map((note, i) => (
                          <span key={i} className="flex items-center gap-1 text-amber-400">
                            <AlertTriangle size={10} />
                            {note}
                          </span>
                        ))}
                      </div>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
