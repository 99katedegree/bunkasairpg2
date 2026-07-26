"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect } from "react";
import Cookies from "js-cookie";
import { getMe } from "@/lib/me/me";
import { getAuthInit } from "@/utils/auth-init";
import { getMeWeaponIndex } from "@/lib/weapon/weapon";
import { getMeItemIndex } from "@/lib/item/item";
import { useUserStore } from "@/stores/user-store";
import { LoadingScreen } from "@/components/shared/loading-screen";
import { Footer } from "@/components/shared/footer";
import type { WeaponResponse, UserItemResponse } from "@/lib/bunkasaiRPGAPI.schemas";

type Props = {
  children: React.ReactNode;
};

export function MainLayout({ children }: Props) {
  const pathname = usePathname();
  const router = useRouter();
  const { user, setUser, setWeapons, setItems } = useUserStore();

  const allowedPaths = ["/admin", "/login"];
  const isPublicPath =
    pathname === "/" ||
    allowedPaths.some(
      (path) => pathname === path || pathname.startsWith(path + "/")
    );

  useEffect(() => {
    if (isPublicPath) return;

    const fetchData = async () => {
      if (!Cookies.get("authToken")) {
        router.push("/?notLoggedIn=1");
        return;
      }

      try {
        const [meRes, weaponsRes, itemsRes] = await Promise.all([
          getMe(getAuthInit()),
          getMeWeaponIndex(undefined, getAuthInit()),
          getMeItemIndex(undefined, getAuthInit()),
        ]);

        if (meRes.status === 200 && "user" in meRes.data) {
          const meUser = meRes.data.user;
          setUser({ ...meUser, maxHitPoint: meUser.hitPoint });
        } else {
          router.push("/?notLoggedIn=1");
          return;
        }

        if (weaponsRes.status === 200 && "weapons" in weaponsRes.data) {
          setWeapons(
            (weaponsRes.data.weapons as (WeaponResponse | null)[]).filter(
              (w): w is WeaponResponse => w !== null
            )
          );
        }

        if (itemsRes.status === 200 && "items" in itemsRes.data) {
          setItems(
            (itemsRes.data.items as (UserItemResponse | null)[]).filter(
              (i): i is UserItemResponse => i !== null
            )
          );
        }
      } catch (err) {
        console.warn(err);
        router.push("/?notLoggedIn=1");
      }
    };

    fetchData();
  }, [isPublicPath]);

  if (isPublicPath) {
    return <>{children}</>;
  }

  if (!user) {
    return <LoadingScreen />;
  }

  return (
    <>
      {children}
      <Footer />
    </>
  );
}
