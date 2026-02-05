import { db, schema } from "$lib/server/db";
import { UUID_V4_RE } from "$lib/utils";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { desc, eq } from "drizzle-orm";

export const load: PageServerLoad = async ({ url }) => {
  const codigoCatedra = url.searchParams.get("catedra");

  if (codigoCatedra && !UUID_V4_RE.test(codigoCatedra)) {
    error(400, "Código de cátedra inválido.");
  }

  const codigoDocente = url.searchParams.get("docente");

  if (!codigoDocente) {
    error(400, "Parámetro de query 'docente' no encontrado.");
  }

  if (!UUID_V4_RE.test(codigoDocente)) {
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

  const cuatris = await db
    .select({
      codigo: schema.cuatrimestre.codigo,
      numero: schema.cuatrimestre.numero,
      anio: schema.cuatrimestre.anio
    })
    .from(schema.cuatrimestre)
    .orderBy(desc(schema.cuatrimestre.codigo))
    .limit(4);

  if (cuatris.length === 0) {
    error(500, "No se han encontrado cuatrimestres para calificar.");
  }

  return {
    codigoCatedra,
    docente: docente[0],
    cuatris
  };
};
