import type { PageServerLoad } from "./$types";

import prisma from "$lib/prisma";

export const load = (async () => {
	const materias = await prisma.materia.findMany();

	return {
		materias: materias.map((m) => {
			return {
				...m,
				search_terms: [
					m.nombre.toLowerCase(),
					m.codigo.toString(),
					m.codigo_equivalencia?.toString()
				]
			};
		})
	};
}) satisfies PageServerLoad;
