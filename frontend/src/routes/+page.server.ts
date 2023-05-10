import type { Materia } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const load = (async ({ fetch }) => {
	const respuesta = await fetch(`${env.BACKEND_URL}/materia`);
	const materias = (await respuesta.json()) as Materia[];
	return { materias };
}) satisfies PageServerLoad;
