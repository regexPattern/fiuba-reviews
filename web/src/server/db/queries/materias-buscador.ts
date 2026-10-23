import "server-only";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";

/** El buscador ordena por cantidad de cátedras o relevancia según la búsqueda. */
export async function getMateriasBuscador() {
  "use cache";
  cacheLife("max");
  cacheTag("materias", "inicio", "buscador-materias");

  const m = schema.materiaVigente;

  return db
    .select({
      codigo: m.codigo,
      nombre: m.nombre,
      cantidadCatedras: m.cantidadCatedras
    })
    .from(m);
}
