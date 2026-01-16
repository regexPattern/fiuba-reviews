import { db, schema } from "$lib/server/db";
import type { PageServerLoad } from "./$types";
import { asc } from "drizzle-orm";

export const prerender = true;

export const load: PageServerLoad = async () => {
  const cuatrimestres = await db
    .select({
      codigo: schema.cuatrimestre.codigo,
      numero: schema.cuatrimestre.numero,
      anio: schema.cuatrimestre.anio
    })
    .from(schema.cuatrimestre)
    .orderBy(asc(schema.cuatrimestre.codigo))
    .limit(4);

  return { cuatrimestres };
}
