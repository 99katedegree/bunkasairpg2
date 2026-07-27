"use client";

import { useState } from "react";
import Cookies from "js-cookie";
import {
  Users,
  Swords,
  Play,
  Archive,
  QrCode,
  Hash,
  CheckCircle2,
  XCircle,
  Loader2,
  ChevronRight,
  X,
  TriangleAlert,
  ImageOff,
} from "lucide-react";
import { startGame, archiveGame } from "@/lib/game/game";
import { getMonsterBattleTokens } from "@/lib/monster/monster";
import { generateQrPdf } from "@/utils/generate-qr-pdf";
import { AssetTable } from "@/components/admin/asset-table";

type LogEntry = { type: "ok" | "err" | "info"; msg: string };

const LogIcon = ({ type }: { type: LogEntry["type"] }) => {
  if (type === "ok") return <CheckCircle2 size={12} className="text-emerald-500 shrink-0 mt-0.5" />;
  if (type === "err") return <XCircle size={12} className="text-red-500 shrink-0 mt-0.5" />;
  return <ChevronRight size={12} className="text-zinc-600 shrink-0 mt-0.5" />;
};

export default function AdminPage() {
  const [userCount, setUserCount] = useState(10);
  const [isStarting, setIsStarting] = useState(false);
  const [isArchiving, setIsArchiving] = useState(false);
  const [isGenerating, setIsGenerating] = useState(false);
  const [log, setLog] = useState<LogEntry[]>([]);
  const [showStartModal, setShowStartModal] = useState(false);
  const [showArchiveModal, setShowArchiveModal] = useState(false);

  const authHeaders = (): RequestInit => ({
    headers: { Authorization: `Bearer ${Cookies.get("adminToken")}` },
  });

  const push = (type: LogEntry["type"], msg: string) =>
    setLog((prev) => [...prev, { type, msg }]);

  const apiErrors = (data: unknown): string[] =>
    (data as { errors?: string[] })?.errors ?? [];

  const handleStartGame = async () => {
    setShowStartModal(false);
    setIsStarting(true);
    push("info", `${userCount}人のユーザーで新しいゲームを開始しています...`);
    try {
      const res = await startGame({ count: userCount }, authHeaders());
      if (res.status !== 200) {
        const codes = apiErrors(res.data);
        if (codes.includes("INVALID_COUNT")) push("err", "ユーザー数は1以上で指定してください");
        else push("err", "ゲームの開始に失敗しました");
        return;
      }
      const { userIds } = res.data;
      const baseUrl = window.location.origin;
      await generateQrPdf(userIds.map((id) => `${baseUrl}/login/${id}`), "user_login_qrcodes.pdf");
      push("ok", `${userIds.length}件のユーザーQRコードを出力しました`);
    } catch {
      push("err", "QRコードの生成に失敗しました");
    } finally {
      setIsStarting(false);
    }
  };

  const handleArchive = async () => {
    setShowArchiveModal(false);
    setIsArchiving(true);
    push("info", "全ユーザーをアーカイブしています...");
    try {
      await archiveGame(authHeaders());
      push("ok", "全ユーザーをアーカイブしました");
    } catch {
      push("err", "アーカイブに失敗しました");
    } finally {
      setIsArchiving(false);
    }
  };

  const handleMonsterQR = async () => {
    setIsGenerating(true);
    push("info", "モンスターのバトルトークンを取得しています...");
    try {
      const res = await getMonsterBattleTokens(authHeaders());
      if (res.status !== 200) {
        const codes = apiErrors(res.data);
        if (codes.includes("NO_MONSTERS")) push("err", "モンスターが1件も登録されていません");
        else push("err", "トークンの取得に失敗しました");
        return;
      }
      const { tokens } = res.data;
      const baseUrl = window.location.origin;
      await generateQrPdf(tokens.map((t) => `${baseUrl}/battles/${t}`), "monster_qrcodes.pdf");
      push("ok", `${tokens.length}件のモンスターQRコードを出力しました`);
    } catch {
      push("err", "QRコードの生成に失敗しました");
    } finally {
      setIsGenerating(false);
    }
  };

  const isAnyLoading = isStarting || isArchiving || isGenerating;

  return (
    <>
      <div className="p-8 flex gap-6 h-full">
        {/* ── Left: main content ── */}
        <div className="flex-1 min-w-0">
          <div className="mb-10">
            <h1 className="text-zinc-100 text-2xl font-semibold tracking-tight">ダッシュボード</h1>
            <p className="text-zinc-500 text-sm mt-1">ゲームセッションとアセットを管理します。</p>
          </div>

          <div className="grid grid-cols-1 gap-px bg-zinc-800 rounded-xl overflow-hidden border border-zinc-800">

            {/* ── Users ── */}
            <section className="bg-zinc-950 p-6">
              <div className="flex items-center gap-2 mb-5">
                <Users size={15} className="text-zinc-400" />
                <h2 className="text-zinc-300 text-sm font-medium">ユーザー</h2>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div className="bg-zinc-900 border border-zinc-800 rounded-lg p-4">
                  <p className="text-zinc-200 text-xs font-medium mb-0.5">新規ゲーム開始</p>
                  <p className="text-zinc-500 text-xs mb-4 leading-relaxed">
                    既存ユーザーをアーカイブし、ログイン用QRコードを生成します。
                  </p>
                  <button
                    onClick={() => setShowStartModal(true)}
                    disabled={isAnyLoading}
                    className="flex items-center gap-1.5 bg-zinc-100 hover:bg-white text-zinc-900 text-xs font-medium rounded-md px-3 py-1.5 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                  >
                    {isStarting ? <Loader2 size={12} className="animate-spin" /> : <Play size={12} />}
                    {isStarting ? "生成中..." : "開始 & QR出力"}
                  </button>
                </div>

                <div className="bg-zinc-900 border border-zinc-800 rounded-lg p-4">
                  <p className="text-zinc-200 text-xs font-medium mb-0.5">一括アーカイブ</p>
                  <p className="text-zinc-500 text-xs mb-4 leading-relaxed">
                    全ユーザーをアーカイブします。未ログインのユーザーは削除されます。
                  </p>
                  <button
                    onClick={() => setShowArchiveModal(true)}
                    disabled={isAnyLoading}
                    className="flex items-center gap-1.5 bg-red-950/60 hover:bg-red-950 border border-red-900/50 text-red-400 text-xs font-medium rounded-md px-3 py-1.5 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                  >
                    {isArchiving ? <Loader2 size={12} className="animate-spin" /> : <Archive size={12} />}
                    {isArchiving ? "アーカイブ中..." : "一括アーカイブ"}
                  </button>
                </div>
              </div>
            </section>

            <div className="h-px bg-zinc-800" />

            {/* ── Monsters ── */}
            <section className="bg-zinc-950 p-6">
              <div className="flex items-center gap-2 mb-5">
                <Swords size={15} className="text-zinc-400" />
                <h2 className="text-zinc-300 text-sm font-medium">モンスター</h2>
              </div>
              <div className="bg-zinc-900 border border-zinc-800 rounded-lg p-4 max-w-sm">
                <p className="text-zinc-200 text-xs font-medium mb-0.5">QRコード一括出力</p>
                <p className="text-zinc-500 text-xs mb-4 leading-relaxed">
                  全モンスターのバトルQRコードをPDFでダウンロードします。
                </p>
                <button
                  onClick={handleMonsterQR}
                  disabled={isAnyLoading}
                  className="flex items-center gap-1.5 bg-zinc-100 hover:bg-white text-zinc-900 text-xs font-medium rounded-md px-3 py-1.5 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
                >
                  {isGenerating ? <Loader2 size={12} className="animate-spin" /> : <QrCode size={12} />}
                  {isGenerating ? "生成中..." : "QR PDF出力"}
                </button>
              </div>
            </section>

            <div className="h-px bg-zinc-800" />

            {/* ── Assets ── */}
            <section className="bg-zinc-950 p-6">
              <div className="flex items-center gap-2 mb-5">
                <ImageOff size={15} className="text-zinc-400" />
                <h2 className="text-zinc-300 text-sm font-medium">アセット管理</h2>
              </div>
              <AssetTable adminToken={Cookies.get("adminToken") ?? ""} />
            </section>

          </div>
        </div>

        {/* ── Right: log ── */}
        <div className="w-72 shrink-0 bg-zinc-900 border border-zinc-800 rounded-xl p-4 font-mono flex flex-col">
          <p className="text-zinc-600 text-[10px] uppercase tracking-widest mb-3">ログ</p>
          <div className="flex-1 flex flex-col gap-1.5 overflow-y-auto">
            {log.length === 0 ? (
              <p className="text-zinc-700 text-xs">操作ログがここに表示されます。</p>
            ) : (
              log.map((entry, i) => (
                <div key={i} className="flex items-start gap-2 text-xs">
                  <LogIcon type={entry.type} />
                  <span className={
                    entry.type === "ok" ? "text-emerald-400" :
                    entry.type === "err" ? "text-red-400" :
                    "text-zinc-400"
                  }>
                    {entry.msg}
                  </span>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* ── 新規ゲーム開始モーダル ── */}
      {showStartModal && (
        <div
          className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
          onClick={() => setShowStartModal(false)}
        >
          <div
            className="bg-zinc-900 border border-zinc-700 rounded-xl p-6 w-80 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between mb-5">
              <div className="flex items-center gap-2">
                <Play size={14} className="text-zinc-400" />
                <h2 className="text-zinc-100 text-sm font-medium">新規ゲーム開始</h2>
              </div>
              <button onClick={() => setShowStartModal(false)} className="text-zinc-600 hover:text-zinc-400 transition-colors">
                <X size={15} />
              </button>
            </div>

            <p className="text-zinc-500 text-xs leading-relaxed mb-5">
              既存ユーザーをアーカイブし、指定した人数のログイン用QRコードを生成します。
            </p>

            <div className="mb-5">
              <label className="text-zinc-500 text-xs font-mono uppercase tracking-widest block mb-2">
                ユーザー数
              </label>
              <div className="flex items-center gap-2 bg-zinc-800 border border-zinc-700 rounded-md px-3 focus-within:border-zinc-500 transition-colors">
                <Hash size={12} className="text-zinc-600 shrink-0" />
                <input
                  type="number"
                  min={1}
                  value={userCount}
                  onChange={(e) => setUserCount(Number(e.target.value))}
                  onWheel={(e) => e.currentTarget.blur()}
                  className="text-zinc-100 text-sm font-mono bg-transparent py-2 w-full outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
                  autoFocus
                />
                <span className="text-zinc-600 text-xs shrink-0">人</span>
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className="shrink-0 text-zinc-500 hover:text-zinc-300 cursor-pointer transition-colors"
                  onClick={(e) => {
                    const { top, height } = e.currentTarget.getBoundingClientRect();
                    const clickY = e.clientY - top;
                    setUserCount((v) => Math.max(1, clickY < height / 2 ? v + 1 : v - 1));
                  }}
                >
                  <path d="m7 15 5 5 5-5" />
                  <path d="m7 9 5-5 5 5" />
                </svg>
              </div>
            </div>

            <div className="flex gap-2">
              <button
                onClick={() => setShowStartModal(false)}
                className="flex-1 text-xs text-zinc-400 hover:text-zinc-200 border border-zinc-700 rounded-md py-2 transition-colors"
              >
                キャンセル
              </button>
              <button
                onClick={handleStartGame}
                disabled={userCount < 1}
                className="flex-1 flex items-center justify-center gap-1.5 bg-zinc-100 hover:bg-white text-zinc-900 text-xs font-medium rounded-md py-2 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
              >
                <Play size={12} />
                開始 & QR出力
              </button>
            </div>
          </div>
        </div>
      )}

      {/* ── アーカイブ確認モーダル ── */}
      {showArchiveModal && (
        <div
          className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
          onClick={() => setShowArchiveModal(false)}
        >
          <div
            className="bg-zinc-900 border border-zinc-700 rounded-xl p-6 w-80 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="flex items-center justify-between mb-5">
              <div className="flex items-center gap-2">
                <TriangleAlert size={14} className="text-red-400" />
                <h2 className="text-zinc-100 text-sm font-medium">一括アーカイブの確認</h2>
              </div>
              <button onClick={() => setShowArchiveModal(false)} className="text-zinc-600 hover:text-zinc-400 transition-colors">
                <X size={15} />
              </button>
            </div>

            <p className="text-zinc-500 text-xs leading-relaxed mb-5">
              全ユーザーをアーカイブします。<br />
              未ログインのユーザーは<span className="text-red-400 font-medium">削除されます</span>。この操作は取り消せません。
            </p>

            <div className="flex gap-2">
              <button
                onClick={() => setShowArchiveModal(false)}
                className="flex-1 text-xs text-zinc-400 hover:text-zinc-200 border border-zinc-700 rounded-md py-2 transition-colors"
              >
                キャンセル
              </button>
              <button
                onClick={handleArchive}
                className="flex-1 flex items-center justify-center gap-1.5 bg-red-950/60 hover:bg-red-950 border border-red-900/50 text-red-400 text-xs font-medium rounded-md py-2 transition-colors"
              >
                <Archive size={12} />
                アーカイブする
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
