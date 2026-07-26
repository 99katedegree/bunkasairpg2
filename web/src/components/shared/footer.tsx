"use client";

import { useEffect, useState } from "react";
import Image from "next/image";
import Link from "next/link";

export function Footer() {
  const [isKeyboardVisible, setIsKeyboardVisible] = useState(false);

  useEffect(() => {
    const initialInnerHeight = window.innerHeight;
    const initialViewportHeight =
      window.visualViewport?.height || initialInnerHeight;

    const checkKeyboard = () => {
      const viewportHeight =
        window.visualViewport?.height || window.innerHeight;

      const diffInner = initialInnerHeight - window.innerHeight;
      const diffViewport = initialViewportHeight - viewportHeight;

      // どちらかが150px以上小さくなったらキーボード表示中とみなす
      const isSmall = diffInner > 150 || diffViewport > 150;
      setIsKeyboardVisible(isSmall);
    };

    window.addEventListener("resize", checkKeyboard);
    window.visualViewport?.addEventListener("resize", checkKeyboard);
    window.visualViewport?.addEventListener("scroll", checkKeyboard);
    window.addEventListener("focus", checkKeyboard);
    window.addEventListener("blur", checkKeyboard);

    return () => {
      window.removeEventListener("resize", checkKeyboard);
      window.visualViewport?.removeEventListener("resize", checkKeyboard);
      window.visualViewport?.removeEventListener("scroll", checkKeyboard);
      window.removeEventListener("focus", checkKeyboard);
      window.removeEventListener("blur", checkKeyboard);
    };
  }, []);

  if (isKeyboardVisible) return null;

  return (
    <div className="fixed bottom-0 w-full">
      <div className="relative h-16">
        <div className="absolute z-10 w-full h-24 flex justify-center overflow-hidden">
          <Image
            className="w-[110%] h-full"
            width={100}
            height={100}
            src="/tab-bar-backgrand.png"
            alt="tab-bar"
          />
        </div>
        <div className="absolute bottom-2 z-20 w-full flex justify-center gap-3">
          <Link className="w-18 h-auto" href="/scan">
            <Image
              className="w-full h-full"
              width={100}
              height={100}
              src="/btn-home.png"
              alt="home"
            />
          </Link>
          <Link className="w-18 h-auto" href="/encyclopedia">
            <Image
              className="w-full h-full"
              width={100}
              height={100}
              src="/btn-encyclopedia.png"
              alt="encyclopedia"
            />
          </Link>
          <Link className="w-18 h-auto" href="/rankings">
            <Image
              className="w-full h-full"
              width={100}
              height={100}
              src="/btn-ranking.png"
              alt="ranking"
            />
          </Link>
          <Link className="w-18 h-auto" href="/profile">
            <Image
              className="w-full h-full"
              width={100}
              height={100}
              src="/btn-profile.png"
              alt="profile"
            />
          </Link>
        </div>
      </div>
    </div>
  );
}
