import db from "$lib/db";
import { catedra, materia } from "$lib/db/schema";
import { eq, sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const load = (async () => {
	const materias = await db
		.select({
			nombre: materia.nombre,
			codigo: sql<string>`CAST(${materia.codigo} AS TEXT)`,
			codigoEquivalencia: sql<string | null>`CAST(${materia.codigoEquivalencia} AS TEXT)`,
			cantidadCatedras: sql<number>`(
        SELECT COUNT(*)
        FROM ${catedra}
        WHERE ${catedra.codigoMateria} = ${materia.codigo}
      )`,
		})
		.from(materia)
		.innerJoin(catedra, eq(materia.codigo, catedra.codigoMateria))
		.groupBy(materia.codigo)
		.orderBy(materia.codigo);

	return { materias };
}) satisfies PageServerLoad;
