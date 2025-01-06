import db from "$lib/db";
import { calificacion, catedra, catedraDocente, docente, materia } from "$lib/db/schema";
import { error } from "@sveltejs/kit";
import { and, eq, isNull, sql } from "drizzle-orm";

import type { LayoutServerLoad } from "./$types";

export const load = (async ({ params }) => {
	const codigoMateria = parseInt(params.codigoMateria, 10);

	const materias = await db
		.select({
			nombre: materia.nombre,
			codigo: materia.codigo
		})
		.from(materia)
		.where(and(eq(materia.codigo, codigoMateria), isNull(materia.codigoEquivalencia)))
		.innerJoin(catedra, eq(materia.codigo, catedra.codigoMateria))
		.limit(1);

	if (materias.length === 0) {
		throw error(404, "Materia no encontrada.");
	}

	return {
		materia: materias[0],
		streamed: {
			catedras: fetchCatedrasMateria(codigoMateria)
		}
	};
}) satisfies LayoutServerLoad;

async function fetchCatedrasMateria(codigoMateria: number) {
	return await db
		.select({
			codigo: catedra.codigo,
			nombre: sql<string>`STRING_AGG(${docente.nombre}, '-' ORDER BY ${docente.nombre} ASC)`,
			promedio: sql<number | null>`
AVG((
  SELECT AVG((
    ${calificacion.aceptaCritica} 
      + ${calificacion.asistencia} 
      + ${calificacion.buenTrato} 
      + ${calificacion.claridad} 
      + ${calificacion.claseOrganizada} 
      + ${calificacion.cumpleHorarios} 
      + ${calificacion.fomentaParticipacion} 
      + ${calificacion.panoramaAmplio} 
      + ${calificacion.respondeMails}) / 9)
	FROM ${calificacion}
	WHERE ${calificacion.codigoDocente} = ${docente.codigo}
  GROUP BY ${docente.codigo})
)`
		})
		.from(catedra)
		.innerJoin(catedraDocente, eq(catedra.codigo, catedraDocente.codigoCatedra))
		.innerJoin(docente, eq(docente.codigo, catedraDocente.codigoDocente))
		.where(eq(catedra.codigoMateria, codigoMateria))
		.groupBy(catedra.codigo)
		.orderBy(({ promedio }) => sql`${promedio} DESC NULLS LAST`);
}
