import db, { queryPromedioDocente } from "$lib/db";
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

	const catedrasConDocentesYPromedio = await db
		.select({
			codigo: schema.catedra.codigo,
			nombreDocente: schema.docente.nombre,
			promedioDocente: queryPromedioDocente<number | null>()
		})
		.from(schema.catedra)
		.innerJoin(
			schema.catedraDocente,
			eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
		)
		.innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
		.leftJoin(schema.calificacion, eq(schema.calificacion.codigoDocente, schema.docente.codigo))
		.where(eq(schema.catedra.codigoMateria, materia.codigo))
		.groupBy(sql`${schema.catedra.codigo}, ${schema.docente.nombre}`);

	const codigoCatedraADocentes: Map<string, { nombre: string; promedio: number | null }[]> =
		new Map();

	for (const catedra of catedrasConDocentesYPromedio) {
		const docentes = codigoCatedraADocentes.get(catedra.codigo) ?? [];
		docentes.push({ nombre: catedra.nombreDocente, promedio: catedra.promedioDocente });
		codigoCatedraADocentes.set(catedra.codigo, docentes);
	}

	const catedras: { codigo: string; nombre: string; promedio: number }[] = [];

	codigoCatedraADocentes.forEach((docentes, codigoCatedra) => {
		const nombreCatedra = docentes.map((d) => d.nombre).join(", ");

		const docentesConCalificaciones = docentes.filter((d) => d.promedio != null).length;
		const promedio =
			docentes.reduce((acc, d) => acc + (d.promedio || 0), 0) / docentesConCalificaciones || 0;

		catedras.push({
			codigo: codigoCatedra,
			nombre: nombreCatedra,
			promedio
		});
	});

	catedras.sort((a, b) => b.promedio - a.promedio);

	return { materia, catedras };
}) satisfies LayoutServerLoad;
