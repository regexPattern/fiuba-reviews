import db from "$lib/db";
import { equivalencia, materia } from "$lib/db/schema";
import type { LayoutServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { aliasedTable, eq, sql, inArray } from "drizzle-orm";

export const load = (async ({ params }) => {
  const filasMaterias = await db
    .select({
      codigo: materia.codigo,
      nombre: materia.nombre,
    })
    .from(materia)
    .where(eq(materia.codigo, params.codigoMateria));

  if (filasMaterias.length === 0) {
    error(404, "Materia no encontrada.");
  }

  const filasEquivalencias = await db
    .select({
      codigo: sql<string>`m2.codigo`,
      nombre: sql<string>`m2.nombre`,
    })
    .from(equivalencia)
    .innerJoin(
      aliasedTable(materia, "m2"),
      eq(equivalencia.codigoMateriaPlanAnterior, sql`m2.codigo`)
    )
    .where(eq(equivalencia.codigoMateriaPlanVigente, params.codigoMateria));

  const filasCatedras = await db.execute(
    sql`
SELECT *
FROM catedras_por_equivalencia(${params.codigoMateria})
ORDER BY promedio DESC
`
  );

  console.log(filasCatedras);

  return {
    materia: filasMaterias[0],
    equivalencias: filasEquivalencias,
    catedras: filasCatedras,
  };
}) satisfies LayoutServerLoad;
