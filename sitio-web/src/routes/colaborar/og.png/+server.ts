import type { RequestHandler } from "./$types";

import { read } from "$app/server";
import StaticOGImgColaborar from "$lib/assets/open-graph/StaticOGImgColaborar.svelte";

import { ImageResponse } from "@ethercorps/sveltekit-og";
import { CustomFont, resolveFonts } from "@ethercorps/sveltekit-og/fonts";
import sourceSerif4Woff from "@fontsource/source-serif-4/files/source-serif-4-latin-400-normal.woff";
import sourceSerif4SemiboldWoff from "@fontsource/source-serif-4/files/source-serif-4-latin-600-normal.woff";

const sourceSerif4FontData = () => read(sourceSerif4Woff).arrayBuffer();
const sourceSerif4SemiboldFontData = () => read(sourceSerif4SemiboldWoff).arrayBuffer();

export const prerender = true;

export const GET: RequestHandler = async () => {
  const fonts = await resolveFonts([
    new CustomFont("Source Serif 4", sourceSerif4FontData, { weight: 400 }),
    new CustomFont("Source Serif 4", sourceSerif4SemiboldFontData, { weight: 600 })
  ]);

  return new ImageResponse(StaticOGImgColaborar, { width: 1200, height: 630, fonts });
};
