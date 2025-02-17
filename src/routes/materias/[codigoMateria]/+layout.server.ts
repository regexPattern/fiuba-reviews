import type { LayoutServerLoad } from "./$types";

import db from "$lib/db";
import { equivalencia, materia } from "$lib/db/schema";
import { error } from "@sveltejs/kit";
import { aliasedTable, eq, sql } from "drizzle-orm";

export const load = (async ({ params }) => {
  const filasMaterias = await db
    .select({
      codigo: materia.codigo,
      nombre: materia.nombre,
    })
    .from(materia)
    .where(eq(materia.codigo, params.codigoMateria));

  if (filasMaterias.length === 0) {
    throw error(404, "Materia no encontrada.");
  }

  const filasEquivalencias = await db
    .select({
      codigo: sql<string>`m2.codigo`,
      nombre: sql<string>`m2.nombre`,
    })
    .from(equivalencia)
    .innerJoin(
      aliasedTable(materia, "m2"),
      eq(equivalencia.codigoMateriaPlanAnterior, sql`m2.codigo`),
    )
    .where(eq(equivalencia.codigoMateriaPlanVigente, params.codigoMateria));

  const filasCatedras = await db.execute<{
    codigo: string;
    nombre: string;
    promedio: number | null;
  }>(
    sql`
SELECT *
FROM catedras_por_equivalencia(${params.codigoMateria})
ORDER BY promedio DESC NULLS LAST
`,
  );

  return {
    materia: filasMaterias[0],
    equivalencias: filasEquivalencias,
    catedras: filasCatedras,
  };
}) satisfies LayoutServerLoad;
