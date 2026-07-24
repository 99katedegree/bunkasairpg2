"use client";

import { useState } from "react";
import Cookies from "js-cookie";
import { startGame, archiveGame } from "@/lib/game/game";
import { generateQrPdf } from "@/utils/generate-qr-pdf";
import { cn } from "@/utils/cn";

export default function AdminUsersPage() {
  const [count, setCount] = useState<number>(10);
  const [isStarting, setIsStarting] = useState(false);
  const [isArchiving, setIsArchiving] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const getAuthHeaders = (): RequestInit => {
    const token = Cookies.get("adminToken");
    return {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    };
  };

  const handleStartGame = async () => {
    setSuccessMessage(null);
    setErrorMessage(null);
    setIsStarting(true);

    try {
      const res = await startGame({ count }, getAuthHeaders());
      if (res.status !== 200) {
        setErrorMessage("ゲーム開始に失敗しました");
        return;
      }

      const { userIds } = res.data;
      const baseUrl = window.location.origin;
      const urls = userIds.map((id) => `${baseUrl}/auth/user-login/${id}`);

      await generateQrPdf(urls, "user_login_qrcodes.pdf");
      setSuccessMessage(`${userIds.length}件のユーザーQRコードをダウンロードしました`);
    } catch {
      setErrorMessage("QRコードの生成に失敗しました");
    } finally {
      setIsStarting(false);
    }
  };

  const handleArchive = async () => {
    setSuccessMessage(null);
    setErrorMessage(null);
    setIsArchiving(true);

    try {
      const res = await archiveGame(getAuthHeaders());
      if (res.status === 204) {
        setSuccessMessage("一括アーカイブが完了しました");
      } else {
        setErrorMessage("アーカイブに失敗しました");
      }
    } catch {
      setErrorMessage("アーカイブに失敗しました");
    } finally {
      setIsArchiving(false);
    }
  };

  const isLoading = isStarting || isArchiving;

  return (
    <div className="max-w-2xl">
      <div className="mb-6">
        <div className="h-2 w-24 rounded-full bg-[linear-gradient(90deg,hsl(194,74%,56%),hsl(266,74%,56%),hsl(338,74%,56%),hsl(50,74%,56%),hsl(122,74%,56%))] mb-3" />
        <h1 className="text-2xl font-bold text-gray-900">ユーザー管理</h1>
      </div>

      {successMessage && (
        <div className="mb-4 text-sm text-green-700 bg-green-50 border border-green-200 rounded-lg px-4 py-3">
          {successMessage}
        </div>
      )}
      {errorMessage && (
        <div className="mb-4 text-sm text-red-700 bg-red-50 border border-red-200 rounded-lg px-4 py-3">
          {errorMessage}
        </div>
      )}

      <div className="flex flex-col gap-4">
        {/* 新規ゲームスタート */}
        <section className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">新規ゲームスタート</h2>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-1">
              <label className="text-sm font-medium text-gray-700" htmlFor="user-count">
                ユーザー数
              </label>
              <input
                id="user-count"
                type="number"
                min={1}
                value={count}
                onChange={(e) => setCount(Number(e.target.value))}
                className="border border-gray-300 rounded-lg px-3 py-2 text-sm w-32 focus:outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent"
              />
            </div>
            <p className="text-xs text-gray-500">
              既存ユーザーをアーカイブし、新規ユーザーを作成してログインQRコードPDFをダウンロードします。
            </p>
            <button
              onClick={handleStartGame}
              disabled={isLoading || count < 1}
              className={cn(
                "self-start bg-black text-white rounded-lg px-5 py-2.5 text-sm font-medium transition-opacity",
                (isLoading || count < 1) && "opacity-50 cursor-not-allowed"
              )}
            >
              {isStarting ? "生成中..." : "新規ゲームスタート & QR出力"}
            </button>
          </div>
        </section>

        {/* 一括アーカイブ */}
        <section className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">一括アーカイブ</h2>
          <p className="text-sm text-gray-600 mb-4">
            全ユーザーをアーカイブします。未ログインユーザーのレコードは削除されます。
          </p>
          <button
            onClick={handleArchive}
            disabled={isLoading}
            className={cn(
              "bg-red-600 text-white rounded-lg px-5 py-2.5 text-sm font-medium transition-opacity hover:bg-red-700",
              isLoading && "opacity-50 cursor-not-allowed"
            )}
          >
            {isArchiving ? "アーカイブ中..." : "一括アーカイブ"}
          </button>
        </section>
      </div>
    </div>
  );
}
