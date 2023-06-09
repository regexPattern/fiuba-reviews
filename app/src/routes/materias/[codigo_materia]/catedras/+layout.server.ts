import type { LayoutServerLoad } from "./$types";

import prisma from "$lib/prisma";

export const load = (async ({ params }) => {
	const catedras_docentes = await prisma.catedra.findMany({
		where: {
			codigo_materia: parseInt(params.codigo_materia, 10)
		},
		include: {
			catedradocente: {
				include: {
					docente: true
				}
			}
		}
	});

	const catedras = catedras_docentes.map((c) => {
		let docentes = c.catedradocente.map((cd) => cd.docente);
		const nombre_docentes = docentes.map((d) => d.nombre);
		nombre_docentes.sort();

		docentes = docentes.filter((d) => d.respuestas != 0);
		const promedio = docentes.reduce((curr, p) => curr + p.promedio, 0) / docentes.length;

		return {
			codigo: c.codigo,
			codigo_materia: c.codigo_materia,
			nombre: nombre_docentes.join("-"),
			promedio
		};
	});

	catedras.sort((a, b) => b.promedio - a.promedio);

	return { catedras };
}) satisfies LayoutServerLoad;
