import type { Materia } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const prerender = true;

export const load = (async ({ fetch }) => {
	const respuesta = await fetch(`${env.BACKEND_URL}/materia`);
	const materias = (await respuesta.json()) as Materia[];
	return {
		materias: materias.map((m) => ({ ...m, nombre: m.nombre.toLowerCase() })),
	};
}) satisfies PageServerLoad;
