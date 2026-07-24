import { defineConfig } from "orval";

export default defineConfig({
  bunkasairpg: {
    input: "../schema/openapi/openapi.yaml",
    output: {
      mode: "tags-split",
      target: "./src/lib",
      client: "swr",
      httpClient: "fetch",
      baseUrl: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8085",
    },
  },
});
