"use client";

import { useUserLogin } from "@/lib/auth/auth";
import { useParams } from "next/navigation";
import { useEffect } from "react";
import Cookies from "js-cookie";

export default function Page() {
  const { userId } = useParams<{ userId: string }>();
  const { trigger } = useUserLogin();

  useEffect(() => {
    if (!userId) return;

    trigger({ id: userId })
      .then((res) => {
        if (res.status === 200) {
          Cookies.set("authToken", res.data.authToken, {
            expires: 7,
            path: "/",
          });
          window.location.href = "/scan";
        } else {
          window.location.href = "/?notLoggedIn=1";
        }
      })
      .catch(() => {
        window.location.href = "/?notLoggedIn=1";
      });
  }, [userId]);

  return (
    <div className="flex items-center justify-center h-screen w-screen bg-black text-white">
      <p>ログイン中...</p>
    </div>
  );
}
