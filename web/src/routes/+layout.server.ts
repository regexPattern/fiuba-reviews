import { db, schema } from "$lib/server/db";
import { desc, eq, sql } from "drizzle-orm";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async () => {
  const cantidadPlanesExpr = sql<number>`count(distinct ${schema.plan.codigo})::int`;

  const materias = await db
    .select({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre
    })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .where(eq(schema.plan.estaVigente, true))
    .groupBy(schema.materia.codigo, schema.materia.nombre)
    .orderBy(desc(cantidadPlanesExpr), schema.materia.codigo);

  return { materias };
};
