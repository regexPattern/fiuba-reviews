import db from "$lib/db";
import { materia, plan, planMateria } from "$lib/db/schema";
import type { LayoutServerLoad } from "./$types";
import { eq } from "drizzle-orm";

export const load: LayoutServerLoad = async () => {
  const materias = await db
    .select({
      codigo: materia.codigo,
      nombre: materia.nombre,
    })
    .from(materia)
    .innerJoin(planMateria, eq(materia.codigo, planMateria.codigoMateria))
    .innerJoin(plan, eq(planMateria.codigoPlan, plan.codigo))
    .where(eq(plan.estaVigente, true))
    .groupBy(materia.codigo)
    .orderBy(materia.nombre);

  return { materias };
};
