import db from "$lib/db";
import { catedra, catedraDocente, comentario, materia } from "$lib/db/schema";
import { desc, eq, sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const prerender = true;

export const load: PageServerLoad = async () => {
	const materiasPopulares = await db
		.select({
			codigo: materia.codigo,
			nombre: materia.nombre,
			cantidadCatedras: sql<number>`(
        SELECT COUNT(*)
        FROM ${catedra}
        WHERE ${catedra.codigoMateria} = ${materia.codigo}
      )`,
			cantidadComentarios: sql<number>`COUNT(${comentario.codigo})`
		})
		.from(materia)
		.innerJoin(catedra, eq(materia.codigo, catedra.codigoMateria))
		.innerJoin(catedraDocente, eq(catedra.codigo, catedraDocente.codigoCatedra))
		.innerJoin(comentario, eq(catedraDocente.codigoDocente, comentario.codigoDocente))
		.groupBy(materia.codigo)
		.orderBy(({ cantidadComentarios }) => desc(cantidadComentarios))
		.limit(20);

	return { materiasPopulares };
};
