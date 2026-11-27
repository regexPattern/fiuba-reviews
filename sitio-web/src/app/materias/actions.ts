"use server";

import { eq } from "drizzle-orm";
import { updateTag } from "next/cache";
import { validarCalificacion } from "@/lib/calificaciones";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";
import { getCuatrimestresCalificables } from "@/server/db/queries/cuatrimestres";

export type ResultadoCalificacion = { estado: "enviada" } | { estado: "error"; mensaje: string };

/**
 * Guarda la calificación de un docente y, si dejó uno, el comentario. El borrador llega tal cual
 * del cliente y se valida en `validarCalificacion`.
 */
export async function enviarCalificacion(entrada: unknown): Promise<ResultadoCalificacion> {
  const calificables = await getCuatrimestresCalificables();
  const validacion = validarCalificacion(entrada, calificables);

  if (!validacion.ok) {
    return { estado: "error", mensaje: validacion.mensaje };
  }

  const { codigoDocente, puntajes, cuatrimestre, comentario } = validacion.calificacion;

  const [docente] = await db
    .select({ codigoMateria: schema.docente.codigoMateria })
    .from(schema.docente)
    .where(eq(schema.docente.codigo, codigoDocente))
    .limit(1);

  if (!docente) {
    return { estado: "error", mensaje: "No encontramos al docente." };
  }

  // Ya validado: es uno de los calificables.
  const registrado = calificables.find(
    (c) => c.numero === cuatrimestre.numero && c.anio === cuatrimestre.anio
  )!;

  try {
    await db.transaction(async (tx) => {
      const [calificacion] = await tx
        .insert(schema.calificacionDolly)
        .values({ codigoDocente, ...puntajes })
        .returning({ codigo: schema.calificacionDolly.codigo });

      if (comentario) {
        await tx.insert(schema.comentario).values({
          codigoDocente,
          codigoCuatrimestre: registrado.codigo,
          contenido: comentario,
          deDolly: false,
          codigoCalificacionDolly: calificacion.codigo
        });
      }
    });
  } catch (error) {
    console.error("[enviarCalificacion] Error al guardar la calificación.", {
      codigoDocente,
      error
    });

    return {
      estado: "error",
      mensaje: "No pudimos guardar la calificación. Probá de nuevo más tarde."
    };
  }

  console.info("[enviarCalificacion] Calificación guardada.", {
    codigoDocente,
    conComentario: comentario !== null
  });

  // Cambian los promedios y comentarios de cada cátedra del docente, y con ellos el orden de las
  // cátedras de su materia.
  const catedras = await db
    .select({ codigo: schema.catedraDocente.codigoCatedra })
    .from(schema.catedraDocente)
    .where(eq(schema.catedraDocente.codigoDocente, codigoDocente));

  for (const { codigo } of catedras) updateTag(tags.catedra(codigo));
  updateTag(tags.catedras(docente.codigoMateria));

  return { estado: "enviada" };
}
