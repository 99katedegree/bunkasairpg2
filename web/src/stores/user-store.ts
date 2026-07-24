import { create } from "zustand";
import type { MeResponse, WeaponResponse, UserItemResponse } from "@/lib/bunkasaiRPGAPI.schemas";

type UserStore = {
  user: (MeResponse & { maxHitPoint: number }) | null;
  weapons: WeaponResponse[];
  items: UserItemResponse[];
  setUser: (user: (MeResponse & { maxHitPoint: number }) | null) => void;
  setWeapons: (weapons: WeaponResponse[]) => void;
  setItems: (items: UserItemResponse[]) => void;
  reset: () => void;
};

export const useUserStore = create<UserStore>((set) => ({
  user: null,
  weapons: [],
  items: [],
  setUser: (user) => set({ user }),
  setWeapons: (weapons) => set({ weapons }),
  setItems: (items) => set({ items }),
  reset: () => set({ user: null, weapons: [], items: [] }),
}));
