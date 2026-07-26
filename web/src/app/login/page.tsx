"use client";

import { useState } from "react";
import { useAdminLogin } from "@/lib/auth/auth";
import Cookies from "js-cookie";
import { AlertCircle, ArrowRight, Loader2, Shield } from "lucide-react";

export default function AdminLoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const { trigger, isMutating } = useAdminLogin();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      const res = await trigger({ email, password });
      if (res.status === 200) {
        Cookies.set("adminToken", res.data.authToken, { expires: 7, path: "/" });
        window.location.href = "/admin";
      } else {
        setError("メールアドレスまたはパスワードが正しくありません");
      }
    } catch {
      setError("接続に失敗しました");
    }
  };

  return (
    <div className="min-h-screen bg-zinc-950 flex">

      {/* ── Left ── */}
      <div className="hidden lg:flex flex-col flex-1 relative overflow-hidden bg-zinc-900 border-r border-zinc-800 p-14">

        {/* dot grid */}
        <div
          className="absolute inset-0 opacity-[0.15]"
          style={{
            backgroundImage: "radial-gradient(circle, #52525b 1px, transparent 1px)",
            backgroundSize: "28px 28px",
          }}
        />

        {/* glow */}
        <div className="absolute bottom-0 left-0 w-[500px] h-[500px] rounded-full bg-zinc-700/20 blur-[120px] pointer-events-none -translate-x-1/3 translate-y-1/3" />

        {/* content */}
        <div className="relative flex flex-col h-full">
          <div className="flex items-center gap-2.5">
            <Shield size={16} className="text-zinc-400" />
            <span className="text-zinc-400 font-mono text-xs tracking-widest uppercase">Bunkasai RPG</span>
          </div>

          <div className="flex-1 flex flex-col justify-center">
            <p className="text-zinc-600 font-mono text-xs uppercase tracking-[0.2em] mb-4">Admin Console</p>
            <h2 className="text-zinc-100 text-4xl font-bold tracking-tight leading-tight mb-4">
              文化祭RPG<br />管理システム
            </h2>
            <p className="text-zinc-500 text-sm leading-relaxed max-w-xs">
              ゲームセッションの開始・アーカイブ、QRコードの発行など、イベント全体を一元管理します。
            </p>
          </div>

          <p className="text-zinc-700 font-mono text-[10px]">bunkasai-rpg · admin · v1.0</p>
        </div>
      </div>

      {/* ── Right ── */}
      <div className="flex flex-1 lg:max-w-md items-center justify-center p-10 relative">

        {/* subtle top glow */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-64 h-px bg-gradient-to-r from-transparent via-zinc-700 to-transparent" />

        <div className="w-full max-w-sm">
          <div className="mb-10">
            {/* mobile only label */}
            <p className="text-zinc-600 font-mono text-[10px] uppercase tracking-widest mb-3 lg:hidden">
              Bunkasai RPG · Admin
            </p>
            <h1 className="text-zinc-100 text-2xl font-semibold tracking-tight mb-1">ログイン</h1>
            <p className="text-zinc-500 text-xs">管理者アカウントでサインインしてください。</p>
          </div>

          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-1.5">
              <label className="text-zinc-500 text-[10px] font-mono uppercase tracking-widest" htmlFor="email">
                メールアドレス
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                autoComplete="email"
                className="admin-input bg-zinc-900 border border-zinc-800 text-zinc-100 text-sm rounded-lg px-3 py-2.5 outline-none focus:border-zinc-600 focus:bg-zinc-800/60 placeholder:text-zinc-600 transition-all"
                placeholder="admin@example.com"
              />
            </div>

            <div className="flex flex-col gap-1.5">
              <label className="text-zinc-500 text-[10px] font-mono uppercase tracking-widest" htmlFor="password">
                パスワード
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                autoComplete="current-password"
                className="admin-input bg-zinc-900 border border-zinc-800 text-zinc-100 text-sm rounded-lg px-3 py-2.5 outline-none focus:border-zinc-600 focus:bg-zinc-800/60 placeholder:text-zinc-600 transition-all"
                placeholder="••••••••"
              />
            </div>

            {error && (
              <div className="flex items-center gap-2 bg-red-950/40 border border-red-900/40 rounded-lg px-3 py-2.5">
                <AlertCircle size={13} className="text-red-500 shrink-0" />
                <span className="text-red-400 text-xs">{error}</span>
              </div>
            )}

            <button
              type="submit"
              disabled={isMutating}
              className="mt-2 group flex items-center justify-center gap-2 bg-zinc-100 hover:bg-white text-zinc-900 text-sm font-medium rounded-lg py-2.5 transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
            >
              {isMutating
                ? <Loader2 size={14} className="animate-spin" />
                : <ArrowRight size={14} className="transition-transform group-hover:translate-x-0.5" />}
              {isMutating ? "ログイン中..." : "ログイン"}
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}
