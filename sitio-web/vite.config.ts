import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { sveltekit } from "@sveltejs/kit/vite";

import { sveltekitOG } from "@ethercorps/sveltekit-og/plugin";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig, loadEnv } from "vite";

const sitioWebRoot = dirname(fileURLToPath(import.meta.url));
const fiubaReviewsRoot = resolve(sitioWebRoot, "..");

export default defineConfig(({ mode }) => {
  // Para utilizar el .env.development global que está en la raíz del repositorio, compartido por
  // los dos proyectos y por el compose. No falla si no se encuentra el archivo.

  const fiubaReviewsEnv = loadEnv(mode, fiubaReviewsRoot, "");
  const sitioWebEnv = loadEnv(mode, sitioWebRoot, "");
  const mergedEnv = { ...fiubaReviewsEnv, ...sitioWebEnv };

  Object.assign(process.env, mergedEnv);

  return {
    plugins: [tailwindcss(), sveltekit(), sveltekitOG()]
  };
});
