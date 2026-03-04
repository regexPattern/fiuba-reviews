import type { RequestHandler } from "./$types";
import { read } from "$app/server";
import { ISR_BYPASS_TOKEN } from "$env/static/private";
import DynOGImgCatedra from "$lib/assets/open-graph/DynOGImgCatedra.svelte";
import { db, schema } from "$lib/server/db";
import { calcularNombreCatedra } from "$lib/server/db/utils";
import { ImageResponse } from "@ethercorps/sveltekit-og";
import { CustomFont, resolveFonts } from "@ethercorps/sveltekit-og/fonts";
import { addCacheTag } from "@vercel/functions";
import { and, eq } from "drizzle-orm";
import sourceSerif4Woff from "@fontsource/source-serif-4/files/source-serif-4-latin-400-normal.woff";
import sourceSerif4SemiboldWoff from "@fontsource/source-serif-4/files/source-serif-4-latin-600-normal.woff";

export const config = { isr: { expiration: false, bypassToken: ISR_BYPASS_TOKEN } };

const sourceSerif4FontData = () => read(sourceSerif4Woff).arrayBuffer();
const sourceSerif4SemiboldFontData = () => read(sourceSerif4SemiboldWoff).arrayBuffer();

export const GET: RequestHandler = async ({ params }) => {
  const { codigo_materia, codigo_catedra } = params;

  await addCacheTag(["og-images", `og-materia-${codigo_materia}`, `og-catedra-${codigo_catedra}`]);

  const materias = await db
    .select({ nombre: schema.materia.nombre })
    .from(schema.materia)
    .where(eq(schema.materia.codigo, codigo_materia));

  const nombreMateria = materias[0].nombre;

  const catedraDocenteRows = await db
    .select({ nombre: schema.docente.nombre, prioridadRol: schema.prioridadRol.prioridad })
    .from(schema.catedra)
    .innerJoin(
      schema.catedraDocente,
      eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
    )
    .innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
    .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
    .where(
      and(
        eq(schema.catedra.codigo, codigo_catedra),
        eq(schema.catedra.codigoMateria, codigo_materia)
      )
    );

  const nombreCatedra = calcularNombreCatedra(
    catedraDocenteRows.map((row) => ({
      nombre: row.nombre,
      prioridad: row.prioridadRol ?? Number.MAX_SAFE_INTEGER
    }))
  );

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
