import { error, invalid } from "@sveltejs/kit";
import { form, getRequestEvent } from "$app/server";
import { db, schema as dbSchema } from "$lib/server/db";
import { validateToken } from "$lib/server/turnstile";
import { UUID_V4_RE } from "$lib/utils";

import { invalidateByTag } from "@vercel/functions";
import { eq } from "drizzle-orm";
import * as v from "valibot";

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

export const formAction = form(formSchema, async (fields) => {
  const { url } = getRequestEvent();

  const { success, error: captchaError } = await validateToken(fields.cfTurnstileResponse);

  if (!success) {
    console.warn(`Usuario falló el CAPTCHA.`, {
      error: captchaError
    });
    invalid(`CAPTCHA inválido: ${captchaError}`);
  }

  const codigoDocente = url.searchParams.get("docente");

  if (codigoDocente === null) {
    error(400, "Parámetro de query 'docente' no encontrado.");
  } else if (!UUID_V4_RE.test(codigoDocente)) {
    error(400, "Código de docente inválido.");
  }

  const docentesRows = await db
    .select({ codigoMateria: dbSchema.docente.codigoMateria })
    .from(dbSchema.docente)
    .where(eq(dbSchema.docente.codigo, codigoDocente))
    .limit(1);

  if (docentesRows.length === 0) {
    error(404, "Docente no encontrado.");
  }

  const cuatrimestresRows = await db
    .select({ codigo: dbSchema.cuatrimestre.codigo })
    .from(dbSchema.cuatrimestre)
    .where(eq(dbSchema.cuatrimestre.codigo, fields.cuatrimestre))
    .limit(1);

  if (cuatrimestresRows.length === 0) {
    error(404, "Cuatrimestre no encontrado.");
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
      console.error("error al insertar calificación", {
        codigoDocente,
        codigoCuatrimestre: fields.cuatrimestre,
        error: e
      });

      error(500, "Error interno al guardar tu calificación.");
    }

    console.info("calificación insertada correctamente", {
      codigoDocente,
      codigoCuatrimestre: fields.cuatrimestre,
      codigoCalificacionDolly: calificacionInsertada.codigo
    });

    const contenidoComentario = fields.comentario.trim();

    if (contenidoComentario.length > 0) {
      try {
        const [comentarioInsertado] = await tx
          .insert(dbSchema.comentario)
          .values({
            codigoDocente,
            codigoCuatrimestre: fields.cuatrimestre,
            contenido: contenidoComentario,
            codigoCalificacionDolly: calificacionInsertada.codigo
          })
          .returning({ codigo: dbSchema.comentario.codigo });

        console.info("comentario insertado correctamente", {
          codigoDocente,
          codigoCuatrimestre: fields.cuatrimestre,
          codigoCalificacionDolly: calificacionInsertada.codigo,
          codigoComentario: comentarioInsertado.codigo
        });
      } catch (e) {
        console.error("error al insertar comentario", {
          codigoDocente,
          codigoCuatrimestre: fields.cuatrimestre,
          codigoCalificacionDolly: calificacionInsertada.codigo,
          error: e
        });

        error(500, "Error interno al guardar tu comentario.");
      }
    }
  });

  // Cuando se agrega una calificación a un docente de una cátedra se invalida el cache de todas
  // las cátedras de la materia. Esto tiene que ser así porque se tiene que volver a compilar el
  // layout que enmarca a todas las cátedras (el layout de la materia), porque acá se muestra un
  // listado de cátedras con sus promedios, por lo que, al agregar una calificación, se tiene que
  // recomputar el promedio de esta cátedra y tiene que mostrarse igual para todas las cátedras que
  // muestren el listado.

  const codigoMateria = docentesRows[0].codigoMateria;
  const tagInvalidacionCache = `materia-${codigoMateria}`;

  try {
    await invalidateByTag(tagInvalidacionCache);
  } catch (e) {
    console.warn(`error invalidando cache de materia ${codigoMateria}`, {
      codigoDocente,
      tagInvalidacionCache,
      error: e
    });
  }

  return { success: true };
});
