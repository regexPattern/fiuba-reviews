import type { RequestHandler } from "./$types";

import { read } from "$app/server";
import StaticOGImgInicio from "$lib/assets/open-graph/StaticOGImgInicio.svelte";

import { ImageResponse } from "@ethercorps/sveltekit-og";
import { CustomFont, resolveFonts } from "@ethercorps/sveltekit-og/fonts";
import sourceSerif4SemiboldWoff from "@fontsource/source-serif-4/files/source-serif-4-latin-600-normal.woff";

const sourceSerif4SemiboldFontData = () => read(sourceSerif4SemiboldWoff).arrayBuffer();

export const prerender = true;

export const GET: RequestHandler = async () => {
  const fonts = await resolveFonts([
    new CustomFont("Source Serif 4", sourceSerif4SemiboldFontData, { weight: 600 })
  ]);

  return new ImageResponse(StaticOGImgInicio, { width: 1200, height: 630, fonts });
};
