import "server-only";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";

/** Materias de planes vigentes, ordenadas por relevancia. Para el inicio, sitemap y static params. */
export async function getMateriasVigentes() {
  "use cache";
  cacheLife("max");
  cacheTag("materias");

  return db
    .select({ codigo: schema.materiaVigente.codigo, nombre: schema.materiaVigente.nombre })
    .from(schema.materiaVigente)
    .orderBy(schema.materiaVigente.orden);
}
