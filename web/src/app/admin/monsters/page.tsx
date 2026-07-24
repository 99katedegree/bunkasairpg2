"use client";

import { useState } from "react";
import Cookies from "js-cookie";
import { getMonsterBattleTokens } from "@/lib/monster/monster";
import { generateQrPdf } from "@/utils/generate-qr-pdf";
import { cn } from "@/utils/cn";

export default function AdminMonstersPage() {
  const [isGenerating, setIsGenerating] = useState(false);
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

  const handleGenerateMonsterQRs = async () => {
    setSuccessMessage(null);
    setErrorMessage(null);
    setIsGenerating(true);

    try {
      const res = await getMonsterBattleTokens(getAuthHeaders());
      if (res.status !== 200) {
        setErrorMessage("モンスタートークンの取得に失敗しました");
        return;
      }

      const { tokens } = res.data;
      const baseUrl = window.location.origin;
      const urls = tokens.map((token) => `${baseUrl}/battles/${token}`);

      await generateQrPdf(urls, "monster_qrcodes.pdf");
      setSuccessMessage(`${tokens.length}件のモンスターQRコードをダウンロードしました`);
    } catch {
      setErrorMessage("QRコードの生成に失敗しました");
    } finally {
      setIsGenerating(false);
    }
  };

  return (
    <div className="max-w-2xl">
      <div className="mb-6">
        <div className="h-2 w-24 rounded-full bg-[linear-gradient(90deg,hsl(194,74%,56%),hsl(266,74%,56%),hsl(338,74%,56%),hsl(50,74%,56%),hsl(122,74%,56%))] mb-3" />
        <h1 className="text-2xl font-bold text-gray-900">モンスター管理</h1>
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
        {/* モンスターQR一括出力 */}
        <section className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">モンスターQR一括出力</h2>
          <p className="text-sm text-gray-600 mb-4">
            全モンスターのバトルトークンを取得し、QRコードPDFをダウンロードします。
            各QRコードにはバトル開始URLが埋め込まれています。
          </p>
          <button
            onClick={handleGenerateMonsterQRs}
            disabled={isGenerating}
            className={cn(
              "bg-black text-white rounded-lg px-5 py-2.5 text-sm font-medium transition-opacity",
              isGenerating && "opacity-50 cursor-not-allowed"
            )}
          >
            {isGenerating ? "生成中..." : "モンスターQR一括出力"}
          </button>
        </section>
      </div>
    </div>
  );
}
