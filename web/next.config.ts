import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  env: {
    // API のベースURL。生成された API クライアント（src/lib）が実行時に読む。
    //
    // 本番は Vercel の環境変数で渡す。ここに書いてあるのはローカル開発の
    // 既定値で、未設定のまま動かしたときに undefined/... という URL を
    // 叩きに行くのを防ぐためのもの。
    //
    // NEXT_PUBLIC_ の値はビルド時にコードへ埋め込まれるので、デプロイ先ごとに
    // ビルドし直す必要がある。逆に言えば実行時に差し替えることはできない。
    NEXT_PUBLIC_API_URL:
      process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8085",
  },
};

export default nextConfig;
