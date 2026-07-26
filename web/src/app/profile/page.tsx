"use client";

import Image from "next/image";
import { useUserStore } from "@/stores/user-store";
import { getAuthInit } from "@/utils/auth-init";
import { useState } from "react";
import { QrModal } from "@/components/features/profile/qr-modal/qr-modal";
import { WeaponDrawer } from "@/components/features/battle/weapon-drawer/weapon-drawer";
import { ItemDrawer } from "@/components/features/battle/item-drawer/item-drawer";
import { ProfileConsole } from "@/components/features/profile/profile-console/profile-console";
import { WeaponCard } from "@/components/features/battle/weapon-card/weapon-card";
import { playSound } from "@/utils/play-sound/play-sound";
import { updateMe } from "@/lib/me/me";
import { uploadImage } from "@/lib/image/image";
import { WeaponResponse } from "@/lib/bunkasaiRPGAPI.schemas";
import { BgCamera } from "@/components/shared/bg-camera";
import { UserStatus } from "@/components/shared/user-status";
import { UserIcon } from "@/components/shared/icons/user-icon";

export default function Page() {
  const { user, setUser, items, weapons } = useUserStore();
  const [name, setName] = useState(user?.name ?? "");
  const [nameOpen, setNameOpen] = useState(false);
  const [qrModalOpen, setQrModalOpen] = useState(false);
  const [weaponDrawerOpen, setWeaponDrawerOpen] = useState(false);
  const [itemDrawerOpen, setItemDrawerOpen] = useState(false);

  if (!user) return null;

  const handleNameChange = async (newName: string) => {
    if (newName === user.name || !newName) return;
    setName(newName);
    await updateMe({ name: newName }, getAuthInit());
    setUser({ ...user, name: newName });
  };

  const handleImageChange = async (imageFile?: File | null) => {
    if (!imageFile) return;
    const res = await uploadImage({ imageFile, directory: "avatars" }, getAuthInit());
    if (res.status === 200) {
      const imageId = res.data.image.id;
      await updateMe({ avatarImageId: imageId }, getAuthInit());
      const reader = new FileReader();
      reader.onload = () => {
        const imageUrl = reader.result as string;
        setUser({ ...user, avatarImageUrl: imageUrl });
      };
      reader.readAsDataURL(imageFile);
    }
  };

  const handleChangeWeapon = async (weapon: WeaponResponse) => {
    await updateMe({ equippedWeaponId: weapon.id }, getAuthInit());
    setUser({ ...user, weapon });
    playSound("/sounds/weapon-change.mp3");
    setWeaponDrawerOpen(false);
  };

  return (
    <div className="flex justify-center items-center h-[100dvh] w-screen bg-cover bg-center bg-no-repeat text-xl">
      <BgCamera />
      <div className="fixed top-0 w-full p-2 z-30">
        <UserStatus
          name={user.name}
          imageUrl={user.avatarImageUrl ?? ""}
          level={user.level}
          hitPoint={user.hitPoint}
          maxHitPoint={user.maxHitPoint}
        />
      </div>
      <div className="relative mt-2 mx-2 pr-1 w-full h-[490px] bg-cover bg-center bg-no-repeat flex justify-center rounded-2xl overflow-hidden">
        <Image
          className="absolute -z-10 h-full"
          src="/bg-profile.png"
          alt="bg-profile"
          fill
          priority
        />
        <div className="w-full relative flex flex-col items-center px-4">
          <h1 className="mt-12 text-xl">マイページ</h1>
          <Image
            className="w-[70%] h-[2px]"
            width={100}
            height={100}
            src="/profile-border2.png"
            alt="マイページの下線"
          />
          <div className="px-4 w-full">
            <div className="p-2 mt-[2%] w-full rounded-md text-white relative overflow-hidden">
              <div className="absolute inset-0 -z-10">
                <Image
                  className="w-full h-full"
                  src="/bg-weapon.png"
                  alt="bg-weapon"
                  width={500}
                  height={300}
                />
              </div>
              <div className="absolute inset-0 -z-10 bg-black/50" />
              {user.weapon && <WeaponCard weapon={user.weapon} isEquipped />}
            </div>
          </div>

          <div className="mt-[5%] flex w-full px-[5%] justify-between items-end">
            <div
              className="relative flex w-[63%] h-6 rounded-full"
              onClick={() => setNameOpen(!nameOpen)}
            >
              <div className="absolute bottom-[50%] left-[108%] w-5 h-5 aspect-square z-10">
                <Image src="/change-me-name.png" alt="user" fill />
              </div>
              <div className="flex items-center justify-center [clip-path:polygon(17%_0,100%_0,91%_54%,83%_100%,17%_100%,0%_50%)] z-10 w-18 text-sm bg-[#553115] text-white">
                名前
              </div>
              <div className="absolute left-[20%] flex items-center top-0 w-[100%] h-full text-lg bg-white pl-10 rounded-full">
                {!nameOpen ? (
                  <div className="w-full truncate px-2">{user.name}</div>
                ) : (
                  <input
                    type="text"
                    className="w-[80%] truncate px-1 outline-none"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                    onBlur={(e) => handleNameChange(e.target.value)}
                    autoFocus
                  />
                )}
              </div>
            </div>
            <label className="relative rounded-full">
              <div className="rounded-full relative bg-gray-300 aspect-square flex items-center justify-center overflow-hidden text-white">
                {user.avatarImageUrl ? (
                  <Image
                    className="object-cover w-16 h-16"
                    src={user.avatarImageUrl}
                    alt="user"
                    width={150}
                    height={150}
                  />
                ) : (
                  <UserIcon className="w-16 h-16" />
                )}
              </div>
              <div className="absolute top-0 right-0 w-5 h-5 aspect-square z-30">
                <Image src="/change-weapon-btn.png" alt="user" fill />
              </div>
              <input
                type="file"
                className="hidden"
                accept="image/*"
                onChange={(e) =>
                  handleImageChange(e.target.files ? e.target.files[0] : null)
                }
              />
            </label>
          </div>
          <div className="w-full p-4 grid grid-cols-2 gap-2">
            {[
              {
                label: "武器変更",
                onClick: () => setWeaponDrawerOpen(true),
              },
              {
                label: "アイテム一覧",
                onClick: () => setItemDrawerOpen(true),
              },
              {
                label: "QRコード",
                onClick: () => setQrModalOpen(true),
              },
            ].map((i) => (
              <div
                key={i.label}
                className="relative h-12 flex items-center justify-center"
                onClick={i.onClick}
              >
                <Image
                  className="absolute -z-10"
                  src="/profile-btn.png"
                  alt="プロフィールボタン背景画像"
                  fill
                />
                <p className="text-sm font-bold">{i.label}</p>
              </div>
            ))}
          </div>
        </div>
      </div>

      {qrModalOpen && <QrModal setOpen={setQrModalOpen} userId={user.id} />}
      {weaponDrawerOpen && (
        <ProfileConsole setClose={setWeaponDrawerOpen} weapons={weapons}>
          <WeaponDrawer
            onClose={() => setWeaponDrawerOpen(false)}
            changeWeapon={(w) => handleChangeWeapon(w)}
          />
        </ProfileConsole>
      )}
      {itemDrawerOpen && (
        <ProfileConsole setClose={setItemDrawerOpen} items={items}>
          <ItemDrawer
            onClose={() => setItemDrawerOpen(false)}
            useItem={() => {}}
            invalid={true}
          />
        </ProfileConsole>
      )}
    </div>
  );
}
