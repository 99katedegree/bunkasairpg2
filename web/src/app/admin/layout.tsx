"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";
import { LogOut } from "lucide-react";

type Props = { children: React.ReactNode };

export default function AdminLayout({ children }: Props) {
  const router = useRouter();

  useEffect(() => {
    if (!Cookies.get("adminToken")) {
      router.replace("/login");
    }
  }, [router]);

  const handleLogout = () => {
    Cookies.remove("adminToken");
    router.replace("/login");
  };

  return (
    <div className="min-h-screen bg-zinc-950 flex flex-col">
      {/* 本文の上に固定で乗るので背景が要る。透明のままだとスクロールした
          本文が透けて、下の境界線を突き抜けて見える。 */}
      <div className="fixed top-0 left-0 right-0 z-10 bg-zinc-950">
        <div className="flex justify-end px-6 py-4">
          <button
            onClick={handleLogout}
            className="flex items-center gap-1.5 text-zinc-600 hover:text-zinc-400 text-xs transition-colors"
          >
            <LogOut size={13} />
            ログアウト
          </button>
        </div>
        {/* w-screen だとスクロールバーの分だけ親からはみ出して横スクロールが
            生まれるので、親の幅に合わせる。 */}
        <div className="w-full border-b border-zinc-800" />
      </div>
      <main className="flex-1 overflow-auto pt-14">
        {children}
      </main>
    </div>
  );
}
