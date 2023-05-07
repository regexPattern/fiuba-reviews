import type { Materia } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const load = (async ({ fetch }) => {
	const response = await fetch(`${env.BACKEND_URL}/materias`);
	const materias = (await response.json()) as Materia[];

	return {
		materias: materias.map((m) => ({ ...m, nombre: m.nombre.toLowerCase() })),
	};
}) satisfies PageServerLoad;
