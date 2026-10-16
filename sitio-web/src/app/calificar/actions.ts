"use server";

import type { CalificarInput } from "@/lib/calificar-schema";
import { eq } from "drizzle-orm";
import { updateTag } from "next/cache";
import { calificarSchema } from "@/lib/calificar-schema";
import { db, schema } from "@/server/db";
import { validarTokenTurnstile } from "@/server/turnstile";

type Resultado = { ok: true; redirectTo: string } | { ok: false; error: string };

export async function calificarDocente(input: CalificarInput): Promise<Resultado> {
  const parsed = calificarSchema.safeParse(input);

  if (!parsed.success) {
    return { ok: false, error: "Datos inválidos." };
  }

  const { codigoDocente, codigoCatedra, codigoMateria, calificaciones, comentario, cuatrimestre } =
    parsed.data;

  const captcha = await validarTokenTurnstile(parsed.data.cfTurnstileResponse);

  if (!captcha.success) {
    console.warn("Usuario falló el CAPTCHA.", { error: captcha.error });
    return { ok: false, error: `CAPTCHA inválido: ${captcha.error ?? "desconocido"}` };
  }

  const [docente] = await db
    .select({ codigoMateria: schema.docente.codigoMateria })
    .from(schema.docente)
    .where(eq(schema.docente.codigo, codigoDocente))
    .limit(1);

  if (!docente) {
    return { ok: false, error: "Docente no encontrado." };
  }

  const [cuatri] = await db
    .select({ codigo: schema.cuatrimestre.codigo })
    .from(schema.cuatrimestre)
    .where(eq(schema.cuatrimestre.codigo, cuatrimestre))
    .limit(1);

  if (!cuatri) {
    return { ok: false, error: "Cuatrimestre no encontrado." };
  }

  try {
    await db.transaction(async (tx) => {
      const [calificacion] = await tx
        .insert(schema.calificacionDolly)
        .values({
          codigoDocente,
          aceptaCritica: calificaciones.aceptaCritica.toString(),
          asistencia: calificaciones.asistencia.toString(),
          buenTrato: calificaciones.buenTrato.toString(),
          claridad: calificaciones.claridad.toString(),
          claseOrganizada: calificaciones.claseOrganizada.toString(),
          cumpleHorarios: calificaciones.cumpleHorarios.toString(),
          fomentaParticipacion: calificaciones.fomentaParticipacion.toString(),
          panoramaAmplio: calificaciones.panoramaAmplio.toString(),
          respondeMails: calificaciones.respondeMails.toString()
        })
        .returning({ codigo: schema.calificacionDolly.codigo });

      if (comentario.length > 0) {
        await tx.insert(schema.comentario).values({
          codigoDocente,
          codigoCuatrimestre: cuatrimestre,
          contenido: comentario,
          codigoCalificacionDolly: calificacion.codigo
        });
      }
    });
  } catch (e) {
    console.error("error al guardar la calificación", { codigoDocente, cuatrimestre, error: e });
    return { ok: false, error: "Error interno al guardar tu calificación." };
  }

  // Un docente puede estar en varias cátedras: hay que invalidar todas.
  const catedrasDocente = await db
    .select({ codigo: schema.catedraDocente.codigoCatedra })
    .from(schema.catedraDocente)
    .where(eq(schema.catedraDocente.codigoDocente, codigoDocente));

  for (const { codigo } of catedrasDocente) {
    updateTag(`catedra:${codigo}`);
  }

  // Listado de cátedras (con promedios) de la materia del docente y de las materias vigentes que
  // la muestran por equivalencia.
  updateTag(`catedras:${docente.codigoMateria}`);

  if (comentario.length > 0) {
    updateTag("inicio");
  }

  // Si vino una cátedra, tiene que ser una de las del docente; si no, se manda a la materia.
  const catedraValida =
    codigoCatedra && catedrasDocente.some((c) => c.codigo === codigoCatedra) ? codigoCatedra : null;

  const materiaDestino = codigoMateria ?? docente.codigoMateria;

  return {
    ok: true,
    redirectTo: catedraValida
      ? `/materia/${materiaDestino}?catedra=${catedraValida}`
      : `/materia/${materiaDestino}`
  };
}
