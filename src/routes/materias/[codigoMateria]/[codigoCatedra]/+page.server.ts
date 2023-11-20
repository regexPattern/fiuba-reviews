import db from "$lib/db";
import { calificacion, catedraDocente, comentario, docente } from "$lib/db/schema";
import { eq, sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const prerender = false;

export const load = (async ({ params }) => {
	const docentes = await db
		.select({
			codigo: docente.codigo,
			nombre: docente.nombre,
			descripcion: docente.descripcion,
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
)`,
			cantidadCalificaciones: sql<number>`COUNT(${calificacion.codigo})`,
			promedios: {
				aceptaCritica: sql<number>`AVG(${calificacion.aceptaCritica})`,
				asistencia: sql<number>`AVG(${calificacion.asistencia})`,
				buenTrato: sql<number>`AVG(${calificacion.buenTrato})`,
				claridad: sql<number>`AVG(${calificacion.claridad})`,
				claseOrganizada: sql<number>`AVG(${calificacion.claseOrganizada})`,
				cumpleHorarios: sql<number>`AVG(${calificacion.cumpleHorarios})`,
				fomentaParticipacion: sql<number>`AVG(${calificacion.fomentaParticipacion})`,
				panoramaAmplio: sql<number>`AVG(${calificacion.panoramaAmplio})`,
				respondeMails: sql<number>`AVG(${calificacion.respondeMails})`
			}
		})
		.from(docente)
		.innerJoin(catedraDocente, eq(docente.codigo, catedraDocente.codigoDocente))
		.innerJoin(calificacion, eq(docente.codigo, calificacion.codigoDocente))
		.where(eq(catedraDocente.codigoCatedra, params.codigoCatedra))
		.groupBy(docente.codigo, docente.nombre)
		.orderBy(docente.nombre);

	const comentarios = await db
		.select({
			codigoDocente: docente.codigo,
			codigo: comentario.codigo,
			cuatrimestre: comentario.cuatrimestre,
			contenido: comentario.contenido
		})
		.from(comentario)
		.innerJoin(docente, eq(comentario.codigoDocente, docente.codigo))
		.innerJoin(catedraDocente, eq(docente.codigo, catedraDocente.codigoDocente))
		.where(eq(catedraDocente.codigoCatedra, params.codigoCatedra));

	const codigoDocenteToComentario: Map<string, typeof comentarios> = new Map();

	for (const com of comentarios) {
		const comentarios = codigoDocenteToComentario.get(com.codigoDocente) || [];
		comentarios.push(com);
		codigoDocenteToComentario.set(com.codigoDocente, comentarios);
	}

	return {
		docentes: docentes.map((doc) => {
			const comentarios = codigoDocenteToComentario.get(doc.codigo) || [];
			comentarios.sort(sortComentariosByCuatrimestre);

			return { ...doc, comentarios };
		})
	};
}) satisfies PageServerLoad;

function sortComentariosByCuatrimestre<T extends { cuatrimestre: string }>(a: T, b: T) {
	const [cuatriA, anioA] = a.cuatrimestre.split("Q");
	const [cuatriB, anioB] = b.cuatrimestre.split("Q");

	if (anioA === anioB) {
		if (cuatriA === cuatriB) {
			return 0;
		} else {
			return cuatriA > cuatriB ? -1 : 1;
		}
	} else {
		return anioA > anioB ? -1 : 1;
	}
}
