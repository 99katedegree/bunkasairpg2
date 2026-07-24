"use client";

import { useEffect } from "react";
import { useRouter, usePathname } from "next/navigation";
import Cookies from "js-cookie";
import Link from "next/link";
import { cn } from "@/utils/cn";

type Props = {
  children: React.ReactNode;
};

const NAV_ITEMS = [
  { href: "/admin/monsters", label: "モンスター" },
  { href: "/admin/users", label: "ユーザー" },
];

export default function AdminLayout({ children }: Props) {
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    const token = Cookies.get("adminToken");
    if (!token) {
      router.replace("/auth/admin-login");
    }
  }, [router]);

  return (
    <div className="min-h-screen bg-gray-50 lg:grid lg:grid-cols-[260px_1fr]">
      {/* Sidebar (desktop only) */}
      <aside className="hidden lg:flex flex-col gap-4 pt-4 pr-0 pl-0">
        <div className="bg-white rounded-r-xl overflow-hidden shadow-md">
          <div className="h-2 bg-[linear-gradient(90deg,hsl(194,74%,56%),hsl(266,74%,56%),hsl(338,74%,56%),hsl(50,74%,56%),hsl(122,74%,56%))]" />
          <div className="p-4">
            <h1 className="text-lg font-bold text-gray-900 mb-4">メニュー</h1>
            <nav className="flex flex-col gap-1">
              {NAV_ITEMS.map(({ href, label }) => (
                <Link
                  key={href}
                  href={href}
                  className={cn(
                    "flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors",
                    pathname.startsWith(href)
                      ? "bg-gray-900 text-white"
                      : "text-gray-700 hover:bg-gray-100"
                  )}
                >
                  {label}
                </Link>
              ))}
            </nav>
          </div>
        </div>
      </aside>

      {/* Main content */}
      <main className="flex-1 p-6">{children}</main>
    </div>
  );
}
