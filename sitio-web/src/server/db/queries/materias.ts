import "server-only";
import { eq } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";

const m = schema.materiaVigente;

/** Materias de planes vigentes ordenadas por relevancia. */
export async function getMateriasVigentes() {
  "use cache";
  cacheLife("max");
  cacheTag(tags.materias);

  return db.select({ codigo: m.codigo, nombre: m.nombre }).from(m).orderBy(m.orden);
}

/** Info de una materia vigente, o `null` si no está en un plan vigente. */
export async function getMateria(codigo: string) {
  "use cache";
  cacheLife("max");
  cacheTag(tags.materias, tags.materia(codigo));

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

  if (!materia) return null;

  const { cuatrimestreNumero: numero, cuatrimestreAnio: anio, ...resto } = materia;

  return { ...resto, cuatrimestre: numero !== null && anio !== null ? { numero, anio } : null };
}

export type Materia = NonNullable<Awaited<ReturnType<typeof getMateria>>>;
