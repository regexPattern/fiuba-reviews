import type { RequestHandler } from "./$types";
import { read } from "$app/server";
import { ISR_BYPASS_TOKEN } from "$env/static/private";
import DynOGImgMateria from "$lib/assets/open-graph/DynOGImgMateria.svelte";
import { db, schema } from "$lib/server/db";
import { ImageResponse } from "@ethercorps/sveltekit-og";
import { CustomFont, resolveFonts } from "@ethercorps/sveltekit-og/fonts";
import { addCacheTag } from "@vercel/functions";
import { eq } from "drizzle-orm";
import sourceSerif4Woff from "@fontsource/source-serif-4/files/source-serif-4-latin-400-normal.woff";
import sourceSerif4SemiboldWoff from "@fontsource/source-serif-4/files/source-serif-4-latin-600-normal.woff";

export const config = { isr: { expiration: false, bypassToken: ISR_BYPASS_TOKEN } };

const sourceSerif4FontData = () => read(sourceSerif4Woff).arrayBuffer();
const sourceSerif4SemiboldFontData = () => read(sourceSerif4SemiboldWoff).arrayBuffer();

export const GET: RequestHandler = async ({ params }) => {
  const { codigo_materia } = params;

  await addCacheTag(["og-images", `og-materia-${codigo_materia}`]);

  const materias = await db
    .select({ nombre: schema.materia.nombre })
    .from(schema.materia)
    .where(eq(schema.materia.codigo, codigo_materia));

  const nombreMateria = materias[0].nombre;

  const fonts = await resolveFonts([
    new CustomFont("Source Serif 4", sourceSerif4FontData, { weight: 400 }),
    new CustomFont("Source Serif 4", sourceSerif4SemiboldFontData, { weight: 600 })
  ]);

  return new ImageResponse(
    DynOGImgMateria,
    { width: 1200, height: 630, fonts },
    { materia: { codigo: codigo_materia, nombre: nombreMateria } }
  );
};
