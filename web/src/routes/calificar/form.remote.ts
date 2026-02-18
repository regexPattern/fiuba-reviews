import { error, invalid } from "@sveltejs/kit";
import { eq } from "drizzle-orm";
import * as v from "valibot";
import { form, getRequestEvent } from "$app/server";
import { db, schema as dbSchema } from "$lib/server/db";
import { validateToken } from "$lib/server/turnstile";
import { UUID_V4_RE } from "$lib/utils";

const CALIFICACIONES_VALIDAS = [0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5];

const formSchema = v.object({
  calificaciones: v.object({
    aceptaCritica: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    asistencia: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    buenTrato: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    claridad: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    claseOrganizada: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    cumpleHorarios: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    fomentaParticipacion: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    panoramaAmplio: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS)),
    respondeMails: v.pipe(v.number(), v.values(CALIFICACIONES_VALIDAS))
  }),
  comentario: v.pipe(
    v.string(),
    v.trim(),
    v.check(
      (comentario) => comentario.length === 0 || comentario.length >= 20,
      "El comentario debe tener al menos 20 caracteres."
    )
  ),
  cuatrimestre: v.pipe(v.number(), v.integer()),
  cfTurnstileResponse: v.string()
});

export const submitForm = form(formSchema, async (fields) => {
  const { success } = await validateToken(fields.cfTurnstileResponse);

  if (!success) {
    console.log("Usuario falló el CAPTCHA.");
    invalid("CAPTCHA inválido.");
  }

  const { url } = getRequestEvent();

  const codigoDocente = url.searchParams.get("docente");

  if (codigoDocente === null) {
    error(400, "Parámetro de query 'docente' no encontrado.");
  } else if (!UUID_V4_RE.test(codigoDocente)) {
    error(400, "Código de docente inválido.");
  }

  const docente = await db
    .select({ codigo: dbSchema.docente.codigo })
    .from(dbSchema.docente)
    .where(eq(dbSchema.docente.codigo, codigoDocente))
    .limit(1);

  if (docente.length === 0) {
    error(404, "Docente no encontrado.");
  }

  const cuatrimestre = await db
    .select({ codigo: dbSchema.cuatrimestre.codigo })
    .from(dbSchema.cuatrimestre)
    .where(eq(dbSchema.cuatrimestre.codigo, fields.cuatrimestre))
    .limit(1);

  if (cuatrimestre.length === 0) {
    error(400, "Cuatrimestre no encontrado.");
  }

  await db.transaction(async (tx) => {
    let calificacionInsertada: { codigo: number };

    try {
      [calificacionInsertada] = await tx
        .insert(dbSchema.calificacionDolly)
        .values({
          codigoDocente,
          aceptaCritica: fields.calificaciones.aceptaCritica.toString(),
          asistencia: fields.calificaciones.asistencia.toString(),
          buenTrato: fields.calificaciones.buenTrato.toString(),
          claridad: fields.calificaciones.claridad.toString(),
          claseOrganizada: fields.calificaciones.claseOrganizada.toString(),
          cumpleHorarios: fields.calificaciones.cumpleHorarios.toString(),
          fomentaParticipacion: fields.calificaciones.fomentaParticipacion.toString(),
          panoramaAmplio: fields.calificaciones.panoramaAmplio.toString(),
          respondeMails: fields.calificaciones.respondeMails.toString()
        })
        .returning({ codigo: dbSchema.calificacionDolly.codigo });
    } catch (e) {
      console.error("[calificarDocente] Error al insertar calificación", {
        codigoDocente,
        codigoCuatrimestre: fields.cuatrimestre,
        error: e
      });

      error(500, "No pudimos guardar tu calificación. Intentá nuevamente en unos minutos.");
    }

    console.info("[calificarDocente] Calificación insertada correctamente", {
      codigoDocente,
      codigoCuatrimestre: fields.cuatrimestre,
      codigoCalificacionDolly: calificacionInsertada.codigo
    });

    if (fields.comentario.length > 0) {
      try {
        const [comentarioInsertado] = await tx
          .insert(dbSchema.comentario)
          .values({
            codigoDocente: codigoDocente + "asldfafs",
            codigoCuatrimestre: fields.cuatrimestre,
            contenido: fields.comentario,
            codigoCalificacionDolly: calificacionInsertada.codigo
          })
          .returning({ codigo: dbSchema.comentario.codigo });

        console.info("[calificarDocente] Comentario insertado correctamente", {
          codigoDocente,
          codigoCuatrimestre: fields.cuatrimestre,
          codigoCalificacionDolly: calificacionInsertada.codigo,
          codigoComentario: comentarioInsertado.codigo
        });
      } catch (e) {
        console.error("[calificarDocente] Error al insertar comentario.", {
          codigoDocente,
          codigoCuatrimestre: fields.cuatrimestre,
          codigoCalificacionDolly: calificacionInsertada.codigo,
          error: e
        });

        error(500, "No pudimos guardar tu comentario. Intentá nuevamente en unos minutos.");
      }
    }
  });

  return { success: true };
});
