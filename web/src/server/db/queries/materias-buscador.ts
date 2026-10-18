import "server-only";
import { countDistinct, desc, eq } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";

/** Cuenta los comentarios visibles, respetando la resolución de equivalencias de catedraMostrada. */
export async function getMateriasBuscador() {
  "use cache";
  cacheLife("max");
  cacheTag("materias", "inicio", "buscador-materias");

  const m = schema.materiaVigente;
  const cm = schema.catedraMostrada;
  const cd = schema.catedraDocente;
  // Un docente puede estar en varias cátedras de la misma materia.
  const cantidadComentarios = countDistinct(schema.comentario.codigo);

  return db
    .select({ codigo: m.codigo, nombre: m.nombre, cantidadComentarios })
    .from(m)
    .leftJoin(cm, eq(cm.codigoMateria, m.codigo))
    .leftJoin(cd, eq(cd.codigoCatedra, cm.codigoCatedra))
    .leftJoin(schema.comentario, eq(schema.comentario.codigoDocente, cd.codigoDocente))
    .groupBy(m.codigo, m.nombre)
    .orderBy(desc(cantidadComentarios), m.nombre, m.codigo);
}
