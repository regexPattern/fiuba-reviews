import { form, getRequestEvent } from "$app/server";
import { db, schema } from "$lib/server/db";
import { validateToken } from "$lib/server/turnstile";
import { error, invalid } from "@sveltejs/kit";
import { eq } from "drizzle-orm";
import * as v from "valibot";

const CALIFICACIONES_VALIDAS = [0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5];
const CODIGO_DOCENTE_REGEX =
  /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

const esquemaFormulario = v.object({
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
  cfTurnstileResponse: v.string()
});

export const calificarDocente = form(
  esquemaFormulario,
  async ({ calificaciones, comentario, cfTurnstileResponse }) => {
    const { success } = await validateToken(cfTurnstileResponse);

    if (!success) {
      invalid("CAPTCHA inválido.");
    }

    const { url } = getRequestEvent();

    const codigoDocente = url.searchParams.get("docente");

    if (codigoDocente === null) {
      error(400, "Parámetro de query 'docente' no encontrado.");
    } else if (!CODIGO_DOCENTE_REGEX.test(codigoDocente)) {
      error(400, "Código de docente inválido.");
    }

    const docente = await db
      .select({ codigo: schema.docente.codigo })
      .from(schema.docente)
      .where(eq(schema.docente.codigo, codigoDocente))
      .limit(1);

    if (docente.length === 0) {
      error(404, "Docente no encontrado.");
    }

    console.log(docente[0], calificaciones, comentario);

    return { success: true };
  }
);
