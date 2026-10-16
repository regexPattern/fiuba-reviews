import "server-only";
import { eq } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";

/**
 * Info de una materia vigente: nombre, cuatrimestre de última actualización y equivalencias.
 * Devuelve `null` si la materia no está en un plan vigente.
 */
export async function getMateria(codigo: string) {
  "use cache";
  cacheLife("max");
  cacheTag(`materia:${codigo}`, "materias");

  const m = schema.materiaVigente;

  const [materia] = await db
    .select({
      codigo: m.codigo,
      nombre: m.nombre,
      cuatrimestreNumero: m.cuatrimestreNumero,
      cuatrimestreAnio: m.cuatrimestreAnio,
      usaEquivalencias: m.usaEquivalencias,
      codigosMateriasFuente: m.codigosMateriasFuente,
      equivalencias: m.equivalencias
    })
    .from(m)
    .where(eq(m.codigo, codigo))
    .limit(1);

  if (!materia) {
    return null;
  }

  return {
    codigo: materia.codigo,
    nombre: materia.nombre,
    cuatrimestre:
      materia.cuatrimestreAnio !== null && materia.cuatrimestreNumero !== null
        ? { numero: materia.cuatrimestreNumero, anio: materia.cuatrimestreAnio }
        : null,
    equivalencias: materia.equivalencias,
    /** `true` si la materia no tiene cátedras propias y muestra las de sus equivalencias. */
    usaEquivalencias: materia.usaEquivalencias,
    /** Materias de las que se leen las cátedras (ella misma, o sus equivalencias). */
    codigosMateriasFuente: materia.codigosMateriasFuente
  };
}

export type Materia = NonNullable<Awaited<ReturnType<typeof getMateria>>>;
