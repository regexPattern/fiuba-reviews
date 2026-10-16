import "server-only";
import { desc, eq, sql } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";
import { getMateria } from "./materia";

/**
 * Listado de cátedras de una materia con su calificación, ordenado por calificación desc y nombre asc.
 *
 * Una materia vigente sin cátedras propias muestra las de sus equivalencias, así que se taguea
 * `catedras:<codigo>` por cada materia de la que realmente se leen cátedras.
 */
export async function getCatedrasMateria(codigo: string) {
  "use cache";
  cacheLife("max");

  const materia = await getMateria(codigo);

  if (!materia) {
    return [];
  }

  cacheTag(`materia:${codigo}`, ...materia.codigosMateriasFuente.map((c) => `catedras:${c}`));

  const cm = schema.catedraMostrada;
  const calificacion = sql<number>`coalesce(${schema.calificacionCatedra.calificacion}, 0)`;

  return db
    .select({
      codigo: cm.codigoCatedra,
      codigoMateria: cm.codigoMateriaCatedra,
      nombre: cm.nombre,
      calificacion: calificacion.as("calificacion")
    })
    .from(cm)
    .leftJoin(
      schema.calificacionCatedra,
      eq(schema.calificacionCatedra.codigoCatedra, cm.codigoCatedra)
    )
    .where(eq(cm.codigoMateria, codigo))
    .orderBy(desc(calificacion), cm.nombre);
}

export type CatedraListada = Awaited<ReturnType<typeof getCatedrasMateria>>[number];
