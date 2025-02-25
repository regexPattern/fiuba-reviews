import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import type { Config } from "@sveltejs/adapter-vercel";

import { TURNSTILE_SECRET_KEY } from "$env/static/private";

import * as s from "$lib/db/schema";
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
  let docentes:
    | {
        nombre: string;
        codigo_materia: string | null;
        codigo_catedra: string;
      }[]
    | undefined;

  try {
    docentes = await db
      .select({
        nombre: s.docente.nombre,
        codigo_materia: s.equivalencia.codigoMateriaPlanVigente,
        codigo_catedra: s.catedraDocente.codigoCatedra,
      })
      .from(s.docente)
      .innerJoin(
        s.catedraDocente,
        eq(s.docente.codigo, s.catedraDocente.codigoDocente),
      )
      .leftJoin(
        s.equivalencia,
        eq(s.docente.codigoMateria, s.equivalencia.codigoMateriaPlanAnterior),
      )
      .where(eq(s.docente.codigo, params.codigoDocente))
      .limit(1);
  } catch (e: any) {
    if (e.code === "22P02") {
      throw error(404);
    } else {
      console.error(e);
      throw error(500);
    }
  }

  if (docentes.length === 0) {
    throw error(404);
  }

  const cuatrimestres = await db
    .select()
    .from(s.cuatrimestre)
    .orderBy(desc(s.cuatrimestre.anio), desc(s.cuatrimestre.numero))
    .limit(4);

  const form = await superValidate(formSchema);

  return {
    docente: docentes[0],
    cuatrimestres,
    form,
  };
};

export const actions: Actions = {
  default: async ({ params, request }) => {
    const form = await superValidate(request, formSchema);

    if (!form.valid) {
      return message(form, "Datos inválidos.");
    }

    if (form.data.comentario.length > 0 && form.data.cuatrimestre != 0) {
      const esCuatrimestreValido =
        (
          await db
            .select()
            .from(s.cuatrimestre)
            .where(eq(s.cuatrimestre.codigo, Number(form.data.cuatrimestre)))
            .limit(1)
        ).length === 1;

      if (!esCuatrimestreValido) {
        return setError(
          form,
          "cuatrimestre",
          `Cuatrimestre '${form.data.cuatrimestre}' no existe.`,
        );
      }
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
        await db.insert(s.comentario).values({
          codigoDocente: params.codigoDocente,
          codigoCuatrimestre: form.data.cuatrimestre,
          contenido: form.data.comentario,
        });
      } catch (e) {
        console.error("Error registrando el comentario.");
        console.error(e);
        throw error(500, "Error interno al guardar el comentario.");
      }
    }

    try {
      await db.insert(s.calificacionDolly).values({
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
      console.error("Error registrando la calificación.");
      console.error(e);
      throw error(500, "Error interno al guardar la calificación.");
    }

    return message(form, "Calificación registrada con éxito.");
  },
};
