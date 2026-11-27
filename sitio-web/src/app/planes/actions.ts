"use server";

import { and, eq, sql } from "drizzle-orm";
import { updateTag } from "next/cache";
import { type Cuatrimestre, formatearCuatrimestre } from "@/lib/cuatrimestres";
import { buscarCarrera } from "@/lib/ofertas";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";
import { parsearMetadata } from "@fiuba-reviews/parser-siu";

export type ResultadoOferta =
  | { estado: "inicial" }
  | { estado: "enviada"; carrera: string; cuatrimestre: Cuatrimestre }
  | { estado: "error"; mensaje: string };

/*
 * La página entera del SIU ronda los 100 KB en las carreras más grandes,
 * entonces le damos el doble por seguridad.
 */
const MAX_CARACTERES = 2 * 1_000_000;

/**
 * Guarda una oferta pegada desde el SIU para procesarla después. La carrera y el cuatrimestre
 * se vuelven a sacar del contenido acá (no se confía en lo que detectó el cliente).
 */
export async function enviarOferta(
  _anterior: ResultadoOferta,
  formData: FormData
): Promise<ResultadoOferta> {
  const contenido = formData.get("contenido");

  if (typeof contenido !== "string" || contenido.length > MAX_CARACTERES) {
    return { estado: "error", mensaje: "El contenido enviado no es válido." };
  }

  const [metadata] = parsearMetadata(contenido);

  if (!metadata) {
    return {
      estado: "error",
      mensaje:
        "No reconocemos esto como la oferta del SIU. Copiá la página entera, incluido el encabezado con tu carrera."
    };
  }

  const carrera = buscarCarrera(await db.select().from(schema.carrera), metadata.carrera);

  if (!carrera) {
    return {
      estado: "error",
      mensaje: `No tenemos registrada la carrera «${metadata.carrera}».`
    };
  }

  const [cuatrimestre] = await db
    .select({ codigo: schema.cuatrimestre.codigo })
    .from(schema.cuatrimestre)
    .where(
      and(
        eq(schema.cuatrimestre.numero, metadata.cuatrimestre.numero),
        eq(schema.cuatrimestre.anio, metadata.cuatrimestre.anio)
      )
    )
    .limit(1);

  if (!cuatrimestre) {
    return {
      estado: "error",
      mensaje: `Todavía no habilitamos el ${formatearCuatrimestre(metadata.cuatrimestre)}.`
    };
  }

  // Tanto las procesadas como las que esperan ser procesadas, así no se
  // reciben duplicadas.
  const { numero, anio } = metadata.cuatrimestre;
  const existentes = await Promise.all(
    [schema.ofertaComisiones, schema.ofertaComisionesRaw].map((tabla) =>
      db
        .select({ codigo: tabla.codigoCuatrimestre })
        .from(tabla)
        .innerJoin(schema.cuatrimestre, eq(schema.cuatrimestre.codigo, tabla.codigoCuatrimestre))
        .where(
          and(
            eq(tabla.codigoCarrera, carrera.codigo),
            sql`(${schema.cuatrimestre.anio}, ${schema.cuatrimestre.numero}) >= (${anio}, ${numero})`
          )
        )
        .limit(1)
    )
  );

  const yaExiste = existentes.some((filas) => filas.length > 0);

  if (yaExiste) {
    return {
      estado: "error",
      mensaje: `Ya tenemos la oferta del ${formatearCuatrimestre(metadata.cuatrimestre)} (o una más nueva) para ${carrera.nombre}.`
    };
  }

  try {
    await db.insert(schema.ofertaComisionesRaw).values({
      codigoCarrera: carrera.codigo,
      codigoCuatrimestre: cuatrimestre.codigo,
      contenido
    });
  } catch (error) {
    console.error("[enviarOferta] Error al guardar la oferta.", {
      codigoCarrera: carrera.codigo,
      codigoCuatrimestre: cuatrimestre.codigo,
      error
    });

    return {
      estado: "error",
      mensaje: "No pudimos guardar la oferta. Probá de nuevo más tarde."
    };
  }

  console.info("[enviarOferta] Oferta guardada.", {
    codigoCarrera: carrera.codigo,
    codigoCuatrimestre: cuatrimestre.codigo
  });

  // Invalidamos el cache para que la tabla de la página se regenere y pase a
  // mostrar la oferta como "en proceso".
  updateTag(tags.ofertas);

  return {
    estado: "enviada",
    carrera: carrera.nombre,
    cuatrimestre: metadata.cuatrimestre
  };
}
