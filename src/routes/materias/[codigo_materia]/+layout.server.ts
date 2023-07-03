import db from "$lib/db";
import * as schema from "$lib/db/schema";
import type { LayoutServerLoad } from "./$types";
import { error } from "@sveltejs/kit";
import { eq, sql } from "drizzle-orm";

export const load = (async ({ params }) => {
	const materia = (
		await db
			.select({ codigo: schema.materia.codigo, nombre: schema.materia.nombre })
			.from(schema.materia)
			.where(eq(schema.materia.codigo, Number(params.codigo_materia)))
			.limit(1)
	)[0];

	if (!materia) {
		throw error(404, { message: "Materia no encontrada" });
	}

	const docentesDeCatedrasConPromedio = await db
		.select({
			codigoCatedra: schema.catedra.codigo,
			nombre: schema.docente.nombre,
			promedio: sql<number | null>`
        AVG((${schema.calificacion.aceptaCritica} +
        ${schema.calificacion.asistencia} +
        ${schema.calificacion.buenTrato} +
        ${schema.calificacion.claridad} +
        ${schema.calificacion.claseOrganizada} +
        ${schema.calificacion.cumpleHorarios} +
        ${schema.calificacion.fomentaParticipacion} +
        ${schema.calificacion.panoramaAmplio} +
        ${schema.calificacion.respondeMails}) / 9)
      `
		})
		.from(schema.docente)
		.leftJoin(schema.calificacion, eq(schema.calificacion.codigoDocente, schema.docente.codigo))
		.innerJoin(
			schema.catedraDocente,
			eq(schema.catedraDocente.codigoDocente, schema.docente.codigo)
		)
		.innerJoin(schema.catedra, eq(schema.catedra.codigo, schema.catedraDocente.codigoCatedra))
		.where(eq(schema.catedra.codigoMateria, materia.codigo))
		.groupBy(sql`${schema.catedra.codigo}, ${schema.docente.codigo}`);

	const codigoCatedraADocentes: Map<string, { nombre: string; promedio: number | null }[]> =
		new Map();

	for (const docente of docentesDeCatedrasConPromedio) {
		const docentesDeCatedra = codigoCatedraADocentes.get(docente.codigoCatedra) ?? [];

		docentesDeCatedra.push(docente);

		codigoCatedraADocentes.set(docente.codigoCatedra, docentesDeCatedra);
	}

	const catedras: { codigo: string; nombre: string; promedio: number }[] = [];

	codigoCatedraADocentes.forEach((docentes, codigoCatedra) => {
		const nombreCatedra = docentes.map((d) => d.nombre).join(", ");

		const docentesConCalificaciones = docentes.filter((d) => d.promedio != null).length;
		const promedioCatedra =
			docentes.reduce((acc, d) => acc + (d.promedio || 0), 0) / docentesConCalificaciones || 0;

		catedras.push({
			codigo: codigoCatedra,
			nombre: nombreCatedra,
			promedio: promedioCatedra
		});
	});

	catedras.sort((a, b) => b.promedio - a.promedio);

	return { materia, catedras };
}) satisfies LayoutServerLoad;
