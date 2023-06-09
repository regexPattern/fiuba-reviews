import type { LayoutServerLoad } from "./$types";

import materias from "$lib/materias";
import prisma from "$lib/prisma";
import { error } from "@sveltejs/kit";

const codigo_materias_validos = materias.map((m) => m.codigo.toString());

export const load = (async ({ params }) => {
	if (!codigo_materias_validos.includes(params.codigo_materia)) {
		throw error(404, { message: "Materia no encontrada" });
	}

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
		const nombre_catedra = docentes.map((d) => d.nombre).sort().join("-");

		docentes = docentes.filter((d) => d.respuestas != 0);
		const promedio = docentes.reduce((curr, p) => curr + p.promedio, 0) / docentes.length;

		return {
			codigo: c.codigo,
			codigo_materia: c.codigo_materia,
			nombre: nombre_catedra,
			promedio
		};
	});

	catedras.sort((a, b) => b.promedio - a.promedio);

	return { catedras };
}) satisfies LayoutServerLoad;
