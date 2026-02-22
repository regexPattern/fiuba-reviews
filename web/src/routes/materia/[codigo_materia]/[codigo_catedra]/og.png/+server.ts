import type { RequestHandler } from "./$types";
import { read } from "$app/server";
import { ISR_BYPASS_TOKEN } from "$env/static/private";
import DynOGImgCatedra from "$lib/assets/open-graph/DynOGImgCatedra.svelte";
import { ImageResponse } from "@ethercorps/sveltekit-og";
import { CustomFont, resolveFonts } from "@ethercorps/sveltekit-og/fonts";
import sourceSerif4Woff from "@fontsource/source-serif-4/files/source-serif-4-latin-400-normal.woff";
import sourceSerif4SemiboldWoff from "@fontsource/source-serif-4/files/source-serif-4-latin-600-normal.woff";

export const config = { isr: { expiration: false, bypassToken: ISR_BYPASS_TOKEN } };

const sourceSerif4FontData = () => read(sourceSerif4Woff).arrayBuffer();
const sourceSerif4SemiboldFontData = () => read(sourceSerif4SemiboldWoff).arrayBuffer();

export const GET: RequestHandler = async ({ params, url }) => {
  const { codigo_materia } = params;

  const nombreMateria = decodeURIComponent(url.searchParams.get("nombre_materia")!);
  const nombreCatedra = decodeURIComponent(url.searchParams.get("nombre_catedra")!);

  const fonts = await resolveFonts([
    new CustomFont("Source Serif 4", sourceSerif4FontData, { weight: 400 }),
    new CustomFont("Source Serif 4", sourceSerif4SemiboldFontData, { weight: 600 })
  ]);

  return new ImageResponse(
    DynOGImgCatedra,
    { width: 1200, height: 630, fonts },
    {
      materia: { codigo: codigo_materia, nombre: nombreMateria },
      catedra: { nombre: nombreCatedra }
    }
  );
};
