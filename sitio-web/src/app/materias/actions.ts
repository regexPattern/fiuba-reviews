"use server";

import { eq } from "drizzle-orm";
import { updateTag } from "next/cache";
import { crearSchemaCalificacion } from "@/lib/calificaciones";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";
import { getCuatrimestresCalificables } from "@/server/db/queries/cuatrimestres";
import { logger } from "@/server/logger";

const log = logger.child({ scope: "enviarCalificacion" });

type ResultadoCalificacion = { estado: "enviada" } | { estado: "error"; mensaje: string };

/** Guarda la calificación de un docente y, si tiene, su comentario. */
export async function enviarCalificacion(entrada: unknown): Promise<ResultadoCalificacion> {
  const calificables = await getCuatrimestresCalificables();
  const SchemaCalificacion = crearSchemaCalificacion(calificables);
  const validacion = SchemaCalificacion.safeParse(entrada, {
    error: () => "La calificación enviada no es válida."
  });

  if (!validacion.success) {
    return { estado: "error", mensaje: validacion.error.issues[0].message };
  }

  const { codigoDocente, puntajes, cuatrimestre, comentario } = validacion.data;

  const [docente] = await db
    .select({ codigoMateria: schema.docente.codigoMateria })
    .from(schema.docente)
    .where(eq(schema.docente.codigo, codigoDocente))
    .limit(1);

  if (!docente) {
    return { estado: "error", mensaje: "No encontramos al docente." };
  }

  try {
    await db.transaction(async (tx) => {
      const [calificacion] = await tx
        .insert(schema.calificacionDolly)
        .values({ codigoDocente, ...puntajes })
        .returning({ codigo: schema.calificacionDolly.codigo });

      if (comentario) {
        await tx.insert(schema.comentario).values({
          codigoDocente,
          codigoCuatrimestre: cuatrimestre.codigo,
          contenido: comentario,
          deDolly: false,
          codigoCalificacionDolly: calificacion.codigo
        });
      }
    });
  } catch (error) {
    log.error({ codigoDocente, err: error }, "Error al guardar la calificación.");

    return {
      estado: "error",
      mensaje: "No pudimos guardar la calificación. Probá de nuevo más tarde."
    };
  }

  log.info({ codigoDocente, conComentario: comentario !== null }, "Calificación guardada.");

  /*
   * Invalida cada cátedra del docente, porque cambian sus promedios y comentarios, y el listado de
   * cátedras de la materia, porque puede cambiar el orden.
   */
  const catedras = await db
    .select({ codigo: schema.catedraDocente.codigoCatedra })
    .from(schema.catedraDocente)
    .where(eq(schema.catedraDocente.codigoDocente, codigoDocente));

  for (const { codigo } of catedras) updateTag(tags.catedra(codigo));
  updateTag(tags.catedras(docente.codigoMateria));

  return { estado: "enviada" };
}
