import { defineConfig } from "orval";

export default defineConfig({
  bunkasairpg: {
    input: "../schema/openapi/openapi.yaml",
    output: {
      mode: "tags-split",
      target: "./src/lib",
      client: "swr",
      httpClient: "fetch",
      // 生成物に実行時の式をそのまま埋め込む。テンプレートリテラルの中に
      // 入るので、ビルド時に Next.js が NEXT_PUBLIC_API_URL の値へ置き換える。
      // 生成時の値を焼き込むと、デプロイ先ごとに作り直す必要が出てしまう。
      baseUrl: "${process.env.NEXT_PUBLIC_API_URL}",
    },
  },
});
