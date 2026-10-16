import "server-only";
import { cacheLife, cacheTag } from "next/cache";
import { and, eq } from "drizzle-orm";
import { db, schema } from "@/server/db";

/**
 * Info de una materia vigente: nombre, cuatrimestre de última actualización y equivalencias.
 * Devuelve `null` si la materia no está en un plan vigente.
 */
export async function getMateria(codigo: string) {
  "use cache";
  cacheLife("max");
  cacheTag(`materia:${codigo}`, "materias");

  const [materia] = await db
    .select({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre,
      cuatrimestreAnio: schema.cuatrimestre.anio,
      cuatrimestreNumero: schema.cuatrimestre.numero
    })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .leftJoin(
      schema.cuatrimestre,
      eq(schema.cuatrimestre.codigo, schema.materia.cuatrimestreUltimaActualizacion)
    )
    .where(and(eq(schema.plan.estaVigente, true), eq(schema.materia.codigo, codigo)))
    .limit(1);

  if (!materia) {
    return null;
  }

  const equivalencias = await db
    .select({ codigo: schema.materia.codigo, nombre: schema.materia.nombre })
    .from(schema.equivalencia)
    .innerJoin(
      schema.materia,
      eq(schema.materia.codigo, schema.equivalencia.codigoMateriaPlanAnterior)
    )
    .where(eq(schema.equivalencia.codigoMateriaPlanVigente, materia.codigo))
    .orderBy(schema.materia.codigo);

  const cuatrimestre =
    materia.cuatrimestreAnio !== null && materia.cuatrimestreNumero !== null
      ? { numero: materia.cuatrimestreNumero, anio: materia.cuatrimestreAnio }
      : null;

  return {
    codigo: materia.codigo,
    nombre: materia.nombre,
    cuatrimestre,
    equivalencias,
    /** `true` si la materia no tiene cátedras propias y muestra las de sus equivalencias. */
    usaEquivalencias: cuatrimestre === null
  };
}

export type Materia = NonNullable<Awaited<ReturnType<typeof getMateria>>>;
