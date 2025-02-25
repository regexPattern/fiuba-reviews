import type { LayoutServerLoad } from "./$types";

import db from "$lib/db";
import { equivalencia, materia } from "$lib/db/schema";
import { error } from "@sveltejs/kit";
import { aliasedTable, eq, sql } from "drizzle-orm";

export const load = (async ({ params }) => {
  const materias = await db
    .select({
      codigo: materia.codigo,
      nombre: materia.nombre,
    })
    .from(materia)
    .where(eq(materia.codigo, params.codigoMateria));

  const equivalencias = await db
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

  if (materias.length === 0) {
    throw error(404);
  }

  return {
    materia: materias[0],
    equivalencias,
    catedras: streamCatedrasMateria(params.codigoMateria),
  };
}) satisfies LayoutServerLoad;

async function streamCatedrasMateria(codigoMateria: string) {
  const catedras = await db.execute<{
    codigo: string;
    nombre: string;
    promedio: number | null;
  }>(
    sql`
SELECT *
FROM catedras_por_equivalencia(${codigoMateria})
ORDER BY promedio DESC NULLS LAST
`,
  );

  return catedras;
}
