"use server";

import { and, eq } from "drizzle-orm";
import { updateTag } from "next/cache";
import {
  type Cuatrimestre,
  extraerMetadataOferta,
  formatearCuatrimestre,
  normalizarNombre,
  ordenCuatrimestre
} from "@/lib/ofertas";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";

export type ResultadoEnvio =
  | { estado: "inicial" }
  | { estado: "enviada"; carrera: string; cuatrimestre: Cuatrimestre }
  | { estado: "error"; mensaje: string };

/** La página entera del SIU ronda los 100 KB en las carreras más grandes. */
const MAX_CARACTERES = 1_000_000;

/**
 * Guarda una oferta pegada desde el SIU para procesarla después. La carrera y el cuatrimestre
 * se vuelven a sacar del contenido acá: no se confía en lo que detectó el cliente.
 *
 * TODO: validar un CAPTCHA (Turnstile) antes de escribir en la base.
 */
export async function enviarOferta(
  _anterior: ResultadoEnvio,
  formData: FormData
): Promise<ResultadoEnvio> {
  const contenido = formData.get("contenido");

  if (typeof contenido !== "string" || contenido.length > MAX_CARACTERES) {
    return { estado: "error", mensaje: "El contenido enviado no es válido." };
  }

  const metadata = extraerMetadataOferta(contenido);

  if (!metadata) {
    return {
      estado: "error",
      mensaje: "No reconocimos la oferta. Asegurate de copiar la página entera del SIU."
    };
  }

  const carreras = await db.select().from(schema.carrera);
  const carrera = carreras.find(
    (c) => normalizarNombre(c.nombre) === normalizarNombre(metadata.carrera)
  );

  if (!carrera) {
    return { estado: "error", mensaje: `No tenemos registrada la carrera «${metadata.carrera}».` };
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

  // Tanto las procesadas como las que esperan ser procesadas: así no se reciben duplicadas.
  const existentes = await Promise.all(
    [schema.ofertaComisiones, schema.ofertaComisionesRaw].map((tabla) =>
      db
        .select({ numero: schema.cuatrimestre.numero, anio: schema.cuatrimestre.anio })
        .from(tabla)
        .innerJoin(schema.cuatrimestre, eq(schema.cuatrimestre.codigo, tabla.codigoCuatrimestre))
        .where(eq(tabla.codigoCarrera, carrera.codigo))
    )
  );

  const yaExiste = existentes
    .flat()
    .some((c) => ordenCuatrimestre(c) >= ordenCuatrimestre(metadata.cuatrimestre));

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

    return { estado: "error", mensaje: "No pudimos guardar la oferta. Probá de nuevo más tarde." };
  }

  console.info("[enviarOferta] Oferta guardada.", {
    codigoCarrera: carrera.codigo,
    codigoCuatrimestre: cuatrimestre.codigo
  });

  // La tabla de la página pasa a mostrar la oferta como "en proceso".
  updateTag(tags.ofertas);

  return { estado: "enviada", carrera: carrera.nombre, cuatrimestre: metadata.cuatrimestre };
}
