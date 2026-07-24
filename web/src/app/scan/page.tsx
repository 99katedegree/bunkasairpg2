"use client";

import { useZxing } from "react-zxing";
import { useRouter } from "next/navigation";

export default function Page() {
  const router = useRouter();
  const { ref } = useZxing({
    onDecodeResult(result) {
      const text = result.rawValue;
      try {
        const url = new URL(text, window.location.origin);
        const currentOrigin = window.location.origin;
        if (
          url.origin === currentOrigin &&
          (url.pathname.startsWith("/battles") ||
            url.pathname.startsWith("/auth/user-login"))
        ) {
          router.push(url.pathname + url.search + url.hash);
        }
      } catch {
        // Invalid URL, ignore
      }
    },
  });

  return (
    <div>
      <div className="fixed inset-0 -z-10">
        <video
          ref={ref}
          className="w-full h-full object-cover"
          autoPlay
          muted
          playsInline
        />
      </div>

      <div className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
        <p className="text-white text-center pb-4">
          QRコードをスキャンしてください。
        </p>
        <div className="w-80 h-80 border-2 border-white" />
      </div>
    </div>
  );
}
