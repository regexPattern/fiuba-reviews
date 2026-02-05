import { db, schema } from "$lib/server/db";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { asc, eq } from "drizzle-orm";

const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

export const load: PageServerLoad = async ({ url }) => {
  const codigoCatedra = url.searchParams.get("catedra");

  if (codigoCatedra && !UUID_REGEX.test(codigoCatedra)) {
    error(400, "Código de cátedra inválido.");
  }

  const codigoDocente = url.searchParams.get("docente");

  if (!codigoDocente) {
    error(400, "Parámetro de query 'docente' no encontrado.");
  }

  if (!UUID_REGEX.test(codigoDocente)) {
    error(400, "Código de docente inválido.");
  }

  const docente = await db
    .select({
      codigo: schema.docente.codigo,
      codigoMateria: schema.docente.codigoMateria,
      nombre: schema.docente.nombre,
      nombreSiu: schema.docente.nombreSiu,
      rol: schema.docente.rol
    })
    .from(schema.docente)
    .where(eq(schema.docente.codigo, codigoDocente))
    .limit(1);

  if (docente.length === 0) {
    error(404, "Docente no encontrado.");
  }

  const cuatrimestres = await db
    .select({
      codigo: schema.cuatrimestre.codigo,
      numero: schema.cuatrimestre.numero,
      anio: schema.cuatrimestre.anio
    })
    .from(schema.cuatrimestre)
    .orderBy(asc(schema.cuatrimestre.codigo))
    .limit(4);

  if (cuatrimestres.length === 0) {
    error(500, "No se han encontrado cuatrimestres para calificar.");
  }

  return {
    cuatrimestres,
    docente: docente[0],
    codigoCatedra
  };
};
