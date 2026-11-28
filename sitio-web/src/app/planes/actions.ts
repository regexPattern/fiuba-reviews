"use server";

import { and, eq, sql } from "drizzle-orm";
import { updateTag } from "next/cache";
import { type Cuatrimestre, formatearCuatrimestre } from "@/lib/cuatrimestres";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";
import { logger } from "@/server/logger";
import { parsearMetadata } from "@fiuba-reviews/parser-siu";
import { buscarCarrera } from "./utils";

export type ResultadoOferta =
  | { estado: "inicial" }
  | { estado: "enviada"; carrera: string; cuatrimestre: Cuatrimestre }
  | { estado: "error"; mensaje: string };

/*
 * La página del SIU ronda los 100 000 caracteres en las carreras más grandes (ver los fixtures de
 * `parser-siu`), entonces le damos el doble para dejar margen.
 */
const MAX_CARACTERES_OFERTA = 200_000;

const log = logger.child({ scope: "enviarOferta" });

/**
 * Guarda una oferta pegada desde el SIU para procesarla después. La carrera y el cuatrimestre se
 * vuelven a leer del contenido, sin confiar en lo que detectó el cliente.
 */
export async function enviarOferta(
  _anterior: ResultadoOferta,
  formData: FormData
): Promise<ResultadoOferta> {
  const contenido = formData.get("contenido");

  if (typeof contenido !== "string" || contenido.length > MAX_CARACTERES_OFERTA) {
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

  /*
   * Busca si ya hay una oferta de este cuatrimestre o de uno más nuevo, tanto entre las procesadas
   * como entre las que esperan ser procesadas, así no se reciben duplicadas ni viejas.
   */
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
    log.error(
      { codigoCarrera: carrera.codigo, codigoCuatrimestre: cuatrimestre.codigo, err: error },
      "Error al guardar la oferta."
    );

    return {
      estado: "error",
      mensaje: "No pudimos guardar la oferta. Probá de nuevo más tarde."
    };
  }

  log.info(
    { codigoCarrera: carrera.codigo, codigoCuatrimestre: cuatrimestre.codigo },
    "Oferta guardada."
  );

  /* Invalida la tabla de ofertas, para que muestre esta como "en proceso". */
  updateTag(tags.ofertas);

  return {
    estado: "enviada",
    carrera: carrera.nombre,
    cuatrimestre: metadata.cuatrimestre
  };
}
