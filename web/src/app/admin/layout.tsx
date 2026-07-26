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
      <div className="fixed top-0 left-0 right-0 z-10">
        <div className="flex justify-end px-6 py-4">
          <button
            onClick={handleLogout}
            className="flex items-center gap-1.5 text-zinc-600 hover:text-zinc-400 text-xs transition-colors"
          >
            <LogOut size={13} />
            ログアウト
          </button>
        </div>
        <div className="w-screen border-b border-zinc-800" />
      </div>
      <main className="flex-1 overflow-auto pt-14">
        {children}
      </main>
    </div>
  );
}
