import db from "$lib/db";
import { calificacion, catedra, catedraDocente, comentario, docente } from "$lib/db/schema";
import { eq, sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const load = (async ({ params }) => {
	const individualDocentesAndComentarios = await db
		.select({
			docente: {
				codigo: docente.codigo,
				nombre: docente.nombre,
				promedio: sql<number>`
(SELECT AVG((
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
  GROUP BY ${docente.codigo})`
			},
			comentario: {
				cuatrimestre: comentario.cuatrimestre,
				contenido: comentario.contenido
			}
		})
		.from(catedra)
		.innerJoin(catedraDocente, eq(catedra.codigo, catedraDocente.codigoCatedra))
		.innerJoin(docente, eq(catedraDocente.codigoDocente, docente.codigo))
		.innerJoin(comentario, eq(docente.codigo, comentario.codigoDocente))
		.where(eq(catedra.codigo, params.codigoCatedra))
		.orderBy(docente.nombre);

	type Docente = (typeof individualDocentesAndComentarios)[number]["docente"];
	type Comentario = (typeof individualDocentesAndComentarios)[number]["comentario"];

	const codigoToDocente: Map<string, Docente> = new Map();
	const codigoDocenteToComentarios: Map<string, Comentario[]> = new Map();

	for (const { docente, comentario } of individualDocentesAndComentarios) {
		if (!codigoToDocente.has(docente.codigo)) {
			codigoToDocente.set(docente.codigo, docente);
		}

		const comentarios = codigoDocenteToComentarios.get(docente.codigo) || [];
		comentarios.push(comentario);
		codigoDocenteToComentarios.set(docente.codigo, comentarios);
	}

	const comentariosByDocente = [];

	for (const [codigo, docente] of codigoToDocente) {
		comentariosByDocente.push({
			...docente,
			comentarios: codigoDocenteToComentarios.get(codigo) || []
		});
	}

	comentariosByDocente.sort((a, b) => a.nombre.localeCompare(b.nombre));

	return { comentariosByDocente };
}) satisfies PageServerLoad;
