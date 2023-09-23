import db from "$lib/db";
import { materia } from "$lib/db/schema";
import { sql } from "drizzle-orm";

import type { LayoutServerLoad } from "./$types";

export const prerender = true;

export const load: LayoutServerLoad = async () => {
  const materias = await db
    .select({
      nombre: materia.nombre,
      codigo: sql<string>`CAST(${materia.codigo} AS TEXT)`,
      codigoEquivalencia: sql<string | null>`CAST(${materia.codigoEquivalencia} AS TEXT)`
    })
    .from(materia)
    .orderBy(materia.codigo);

  return { materias };
};
