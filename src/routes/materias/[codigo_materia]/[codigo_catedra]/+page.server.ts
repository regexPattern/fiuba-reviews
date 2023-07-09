import db from "$lib/db";
import * as schema from "$lib/db/schema";
import type { PageServerLoad } from "./$types";
import { eq, sql } from "drizzle-orm";

export const load = (async ({ params }) => {
	const docentes = await db
		.select({
			codigo: schema.docente.codigo,
			nombre: schema.docente.nombre,
			descripcion: schema.docente.descripcion,
			respuestas: sql<number>`COUNT(${schema.calificacion})`,
			calificaciones: {
				aceptaCritica: sql<number>`AVG(${schema.calificacion.aceptaCritica})`,
				asistencia: sql<number>`AVG(${schema.calificacion.asistencia})`,
				buenTrato: sql<number>`AVG(${schema.calificacion.buenTrato})`,
				claridad: sql<number>`AVG(${schema.calificacion.claridad})`,
				claseOrganizada: sql<number>`AVG(${schema.calificacion.claseOrganizada})`,
				cumpleHorarios: sql<number>`AVG(${schema.calificacion.cumpleHorarios})`,
				fomentaParticipacion: sql<number>`AVG(${schema.calificacion.fomentaParticipacion})`,
				panoramaAmplio: sql<number>`AVG(${schema.calificacion.panoramaAmplio})`,
				respondeMails: sql<number>`AVG(${schema.calificacion.respondeMails})`
			}
		})
		.from(schema.docente)
		.innerJoin(schema.calificacion, eq(schema.calificacion.codigoDocente, schema.docente.codigo))
		.innerJoin(
			schema.catedraDocente,
			eq(schema.catedraDocente.codigoDocente, schema.docente.codigo)
		)
		.innerJoin(schema.catedra, eq(schema.catedra.codigo, schema.catedraDocente.codigoCatedra))
		.where(eq(schema.catedra.codigo, params.codigo_catedra))
		.groupBy(schema.docente.codigo);

	const comentarios = await db
		.select({
			codigo: schema.comentario.codigo,
			codigoDocente: schema.comentario.codigoDocente,
			cuatrimestre: schema.comentario.cuatrimestre,
			contenido: schema.comentario.contenido
		})
		.from(schema.comentario)
		.innerJoin(schema.docente, eq(schema.docente.codigo, schema.comentario.codigoDocente))
		.innerJoin(
			schema.catedraDocente,
			eq(schema.catedraDocente.codigoDocente, schema.docente.codigo)
		)
		.where(eq(schema.catedraDocente.codigoCatedra, params.codigo_catedra));

	const codigoDocenteAComentarios: Map<
		string,
		{ codigo: string; cuatrimestre: string; contenido: string }[]
	> = new Map();

	for (const comentario of comentarios) {
		const comentariosDeDocente = codigoDocenteAComentarios.get(comentario.codigoDocente) || [];

		comentariosDeDocente.push({
			codigo: comentario.codigo,
			cuatrimestre: comentario.cuatrimestre,
			contenido: comentario.contenido
		});

		codigoDocenteAComentarios.set(comentario.codigoDocente, comentariosDeDocente);
	}

	const docentesConPromedioYComentarios = docentes.map((d) => {
		const calificaciones = Object.values(d.calificaciones);
		const promedioDocente =
			calificaciones.reduce((acc, curr) => acc + curr, 0) / calificaciones.length;

		return {
			...d,
			promedio: promedioDocente || 0,
			comentarios: codigoDocenteAComentarios.get(d.codigo) || []
		};
	});

	const cuatrimestres = await db.select().from(schema.cuatrimestre);

	docentesConPromedioYComentarios.sort((a, b) => a.nombre.localeCompare(b.nombre));
	cuatrimestres.sort((a, b) => ordernarCuatrimestres(a.nombre, b.nombre));

	for (const d of docentesConPromedioYComentarios) {
		d.comentarios.sort((a, b) => ordernarCuatrimestres(a.cuatrimestre, b.cuatrimestre));
	}

	return {
		catedras: null,
		docentes: docentesConPromedioYComentarios,
		cuatrimestres
	};
}) satisfies PageServerLoad;

function ordernarCuatrimestres(a: string, b: string) {
	const [cuatriA, anioA] = a.split("Q");
	const [cuatriB, anioB] = b.split("Q");

	if (anioA < anioB) {
		return 1;
	} else if (anioA > anioB) {
		return -1;
	} else {
		return cuatriA <= cuatriB ? 1 : -1;
	}
}
