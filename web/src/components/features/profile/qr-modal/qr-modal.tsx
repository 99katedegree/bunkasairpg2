"use client";

import { QRCodeSVG } from "qrcode.react";

interface Props {
  setOpen: React.Dispatch<React.SetStateAction<boolean>>;
  userId: string;
}

export function QrModal({ setOpen, userId }: Props) {
  return (
    <div
      className="fixed flex justify-center items-center bg-black/70 w-full h-full z-30"
      onClick={() => setOpen(false)}
    >
      <div
        onClick={(e) => e.stopPropagation()}
        className="w-[80%] aspect-square rounded-3xl bg-white flex justify-center items-center"
      >
        <QRCodeSVG
          value={`${window.location.origin}/auth/user-login/${userId}`}
          className="w-[80%] h-[80%]"
        />
      </div>
    </div>
  );
}
