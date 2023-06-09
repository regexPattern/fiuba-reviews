import type { PageServerLoad } from "./$types";

import materias from "$lib/materias";

export const load = (async () => {
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
