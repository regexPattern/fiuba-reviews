import db from "$lib/db";
import { materia, plan, planMateria } from "$lib/db/schema";
import { countDistinct, desc, eq } from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const prerender = true;

export const load: PageServerLoad = async () => {
  const filasMateriasMasPopulares = await db
    .select({
      codigo: materia.codigo,
      nombre: materia.nombre,
      cantidadPlanesVigentes: countDistinct(plan.codigo),
    })
    .from(materia)
    .innerJoin(planMateria, eq(materia.codigo, planMateria.codigoMateria))
    .innerJoin(plan, eq(planMateria.codigoPlan, plan.codigo))
    .where(eq(plan.estaVigente, true))
    .groupBy(materia.codigo, materia.nombre)
    .orderBy(desc(countDistinct(plan.codigo)))
    .limit(20);

  return { materiasMasPopulares: filasMateriasMasPopulares };
};
