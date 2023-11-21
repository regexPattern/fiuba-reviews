import db from "$lib/db";
import { catedra, catedraDocente, cuatrimestre, docente, materia } from "$lib/db/schema";
import { error, fail } from "@sveltejs/kit";
import { eq } from "drizzle-orm";
import { superValidate } from "sveltekit-superforms/server";

import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import schema from "./schema";

export const prerender = false;

export const load: PageServerLoad = async ({ params }) => {
  const docentes = await db
    .select({
      nombreDocente: docente.nombre,
      codigoMateria: materia.codigo,
      codigoCatedra: catedraDocente.codigoCatedra
    })
    .from(docente)
    .innerJoin(catedraDocente, eq(docente.codigo, catedraDocente.codigoDocente))
    .innerJoin(catedra, eq(catedraDocente.codigoCatedra, catedra.codigo))
    .innerJoin(materia, eq(catedra.codigoMateria, materia.codigo))
    .where(eq(docente.codigo, params.codigoDocente))
    .limit(1);

  if (docentes.length === 0) {
    throw error(404, { message: "Docente no encontrado." });
  }

  const cuatrimestres = await db.select().from(cuatrimestre);
  const form = superValidate(schema);

  return { ...docentes[0], cuatrimestres, form };
};

export const actions: Actions = {
  default: async (e) => {
    const form = await superValidate(e, schema);
    if (!form.valid) {
      return fail(400, { form });
    }
    console.log(form);
    return { form };
  }
};
