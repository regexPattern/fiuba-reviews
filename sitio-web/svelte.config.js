import adapter from "@sveltejs/adapter-vercel";
import { vitePreprocess } from "@sveltejs/vite-plugin-svelte";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    alias: { ["$ui"]: "src/lib/ui" },
    adapter: adapter({ runtime: "nodejs24.x" }),
    experimental: { remoteFunctions: true }
  },
  compilerOptions: { experimental: { async: true } }
};

export default config;
