import { db, schema } from "$lib/server/db";
import type { PageServerLoad } from "./$types";
import { eq } from "drizzle-orm";

export const load: PageServerLoad = async () => {
  const materias = await db
    .selectDistinct({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre
    })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .where(eq(schema.plan.estaVigente, true))
    .orderBy(schema.materia.codigo);

  return { materias };
};
