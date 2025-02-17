import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import type { Config } from "@sveltejs/adapter-vercel";

import { TURNSTILE_SECRET_KEY } from "$env/static/private";

import * as dbSchema from "$lib/db/schema";
import db from "$lib/db";
import { desc, eq } from "drizzle-orm";
import { error } from "@sveltejs/kit";
import { formCalificacionDocente as formSchema } from "$lib/zod/schema";
import { message, setError, superValidate } from "sveltekit-superforms/server";
import { validarToken } from "$lib/utils";

export const prerender = false;

export const config: Config = {
  runtime: "nodejs18.x",
};

export const load: PageServerLoad = async ({ params }) => {
  let filasDocentes:
    | {
        nombre: string;
        codigo_materia: string | null;
        codigo_catedra: string;
      }[]
    | undefined;

  try {
    filasDocentes = await db
      .select({
        nombre: dbSchema.docente.nombre,
        codigo_materia: dbSchema.equivalencia.codigoMateriaPlanVigente,
        codigo_catedra: dbSchema.catedraDocente.codigoCatedra,
      })
      .from(dbSchema.docente)
      .innerJoin(
        dbSchema.catedraDocente,
        eq(dbSchema.docente.codigo, dbSchema.catedraDocente.codigoDocente),
      )
      .leftJoin(
        dbSchema.equivalencia,
        eq(
          dbSchema.docente.codigoMateria,
          dbSchema.equivalencia.codigoMateriaPlanAnterior,
        ),
      )
      .where(eq(dbSchema.docente.codigo, params.codigoDocente))
      .limit(1);
  } catch (e: any) {
    // Si pasan un código que no es serializable como UUID.
    if (e.code === "22P02") {
      throw error(404, "Docente no encontrado.");
    } else {
      throw e;
    }
  }

  if (filasDocentes.length === 0) {
    throw error(404, "Docente no encontrado.");
  }

  const filasCuatrimestres = await db
    .select()
    .from(dbSchema.cuatrimestre)
    .orderBy(
      desc(dbSchema.cuatrimestre.anio),
      desc(dbSchema.cuatrimestre.numero),
    )
    .limit(4);

  const form = await superValidate(formSchema);

  return {
    docente: filasDocentes[0],
    cuatrimestres: filasCuatrimestres,
    form,
  };
};

export const actions: Actions = {
  default: async ({ params, request }) => {
    const form = await superValidate(request, formSchema);

    if (!form.valid) {
      return message(form, "Datos inválidos.");
    }

    const esCuatrimestreValido =
      (
        await db
          .select()
          .from(dbSchema.cuatrimestre)
          .where(eq(dbSchema.cuatrimestre, form.data.cuatrimestre))
          .limit(1)
      ).length === 1;

    if (form.data.cuatrimestre && !esCuatrimestreValido) {
      return setError(
        form,
        "cuatrimestre",
        `Cuatrimestre '${form.data.cuatrimestre}' no existe.`,
      );
    }

    const { esValido } = await validarToken(
      form.data["cf-turnstile-response"],
      TURNSTILE_SECRET_KEY,
    );

    if (!esValido) {
      return setError(form, "Error al validar CAPTCHA.");
    }

    if (
      form.data.comentario &&
      form.data.comentario.length > 0 &&
      form.data.cuatrimestre
    ) {
      try {
        await db.insert(dbSchema.comentario).values({
          codigoDocente: params.codigoDocente,
          codigoCuatrimestre: form.data.cuatrimestre,
          contenido: form.data.comentario,
        });
      } catch (e) {
        console.error(e);
        throw error(500, "Error interno al guardar el comentario.");
      }
    }

    try {
      await db.insert(dbSchema.calificacionDolly).values({
        codigoDocente: params.codigoDocente,
        aceptaCritica: form.data["acepta-critica"],
        asistencia: form.data["asistencia"],
        buenTrato: form.data["buen-trato"],
        claridad: form.data["claridad"],
        claseOrganizada: form.data["clase-organizada"],
        cumpleHorarios: form.data["cumple-horario"],
        fomentaParticipacion: form.data["fomenta-participacion"],
        panoramaAmplio: form.data["panorama-amplio"],
        respondeMails: form.data["responde-mails"],
      });
    } catch (e) {
      console.error(e);
      throw error(500, "Error interno al guardar la calificación.");
    }

    return message(form, "Calificación registrada con éxito.");
  },
};
