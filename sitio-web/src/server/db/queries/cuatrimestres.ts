import "server-only";
import { desc } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";

/** Cuántos cuatrimestres se pueden elegir al calificar. */
const CANTIDAD_CALIFICABLES = 2;

/**
 * Cuatrimestres que se pueden elegir al calificar: los últimos registrados en la base, del más
 * nuevo al más viejo.
 */
export async function getCuatrimestresCalificables() {
  "use cache";
  cacheLife("max");
  cacheTag(tags.cuatrimestres);

  const c = schema.cuatrimestre;

  return db
    .select({ codigo: c.codigo, numero: c.numero, anio: c.anio })
    .from(c)
    .orderBy(desc(c.anio), desc(c.numero))
    .limit(CANTIDAD_CALIFICABLES);
}
