# bunkasairpg2 Web

## 技術スタック

| カテゴリ | ライブラリ |
|---|---|
| フレームワーク | [Next.js](https://nextjs.org/) (App Router) |
| バンドラー | Turbopack |
| 言語 | TypeScript |
| スタイリング | [Tailwind CSS v4](https://tailwindcss.com/) |
| Linter / Formatter | [Biome](https://biomejs.dev/) |
| HTTP / API | [orval](https://orval.dev/)（OpenAPI → SWR hooks 生成） + [SWR](https://swr.vercel.app/) |
| 状態管理 | [Zustand](https://zustand-demo.pmnd.rs/) |
| フォーム | [React Hook Form](https://react-hook-form.com/) |
| アニメーション | [Framer Motion](https://www.framer.com/motion/) |
| クッキー | [js-cookie](https://github.com/js-cookie/js-cookie) |
| className | [clsx](https://github.com/lukeed/clsx) + [tailwind-merge](https://github.com/dcastil/tailwind-merge) |

## コマンド

```bash
bun run dev       # 開発サーバー起動
bun run build     # プロダクションビルド
bun run generate  # OpenAPI スキーマから API クライアント生成
bun run lint      # Biome チェック
bun run format    # Biome フォーマット
bun run check     # Biome チェック + 自動修正
```

## API クライアント生成

`../schema/openapi/openapi.yaml` を元に `src/lib/` 配下へ SWR hooks を生成します。
スキーマを変更したら必ず再実行してください。

```bash
bun run generate
```

## アーキテクチャ

```
src
├── app                        # ルーティング（Next.js App Router）
│   ├── page.tsx               # サーバーコンポーネント（データフェッチ・メタデータ）
│   ├── page-client.tsx        # クライアントコンポーネント（インタラクション）
│   ├── page-client.stories.tsx
│   └── page-client.test.tsx
├── components
│   ├── features               # 特定機能に紐づくコンポーネント
│   │   └── {feature}
│   │       └── {component-name}
│   │           ├── {component-name}.tsx
│   │           ├── {component-name}.stories.tsx
│   │           ├── {component-name}.test.tsx
│   │           └── use-{component-name}.ts  # そのコンポーネント専用 hooks
│   └── shared                 # 汎用コンポーネント
│       └── {component-name}
│           ├── {component-name}.tsx
│           ├── {component-name}.stories.tsx
│           ├── {component-name}.test.tsx
│           └── use-{component-name}.ts
├── constants                  # 定数
├── lib                        # orval 生成コード（OpenAPI → SWR hooks）
├── stores                     # Zustand ストア
└── utils                      # 汎用ユーティリティ（cn.ts など）
```

### 規約

- `hooks/` ディレクトリは作らない
- 汎用的な hooks はコンポーネントと同じディレクトリに配置する（`use-{component-name}.ts`）
- `page.tsx` はサーバーコンポーネント、インタラクションは `page-client.tsx` に分離する
