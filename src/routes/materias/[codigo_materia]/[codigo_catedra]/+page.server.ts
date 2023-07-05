import db, { queryPromedioDocente } from "$lib/db";
import * as schema from "$lib/db/schema";
import type { PageServerLoad } from "./$types";
import { eq, sql } from "drizzle-orm";

export const load = (async ({ params }) => {
	const docentes = await db
		.select({
			codigo: schema.docente.codigo,
			nombre: schema.docente.nombre,
			descripcion: schema.docente.descripcion,
			promedio: queryPromedioDocente<number>()
		})
		.from(schema.catedra)
		.innerJoin(
			schema.catedraDocente,
			eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
		)
		.innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
		.innerJoin(schema.calificacion, eq(schema.calificacion.codigoDocente, schema.docente.codigo))
		.innerJoin(schema.comentario, eq(schema.comentario.codigoDocente, schema.docente.codigo))
		.where(eq(schema.catedra.codigo, params.codigo_catedra))
		.groupBy(sql`${schema.docente.codigo}, ${schema.docente.nombre}`);

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
		const comentariosDocente = codigoDocenteAComentarios.get(comentario.codigoDocente) || [];
		comentariosDocente.push({
			codigo: comentario.codigo,
			cuatrimestre: comentario.cuatrimestre,
			contenido: comentario.contenido
		});
		codigoDocenteAComentarios.set(comentario.codigoDocente, comentariosDocente);
	}

	const docentesConComentarios = docentes.map((d) => ({
		...d,
		comentarios: codigoDocenteAComentarios.get(d.codigo) || []
	}));

	const cuatrimestres = await db.select().from(schema.cuatrimestre);

	return {
		docentes: docentesConComentarios,
		cuatrimestres: cuatrimestres.sort(ordernarCuatrimestres)
	};
}) satisfies PageServerLoad;

function ordernarCuatrimestres<T extends { nombre: string }>(a: T, b: T) {
	const [cuatriA, anioA] = a.nombre.split("Q");
	const [cuatriB, anioB] = b.nombre.split("Q");

	if (anioA < anioB) {
		return 1;
	} else if (anioA > anioB) {
		return -1;
	} else {
		return cuatriA <= cuatriB ? 1 : -1;
	}
}
