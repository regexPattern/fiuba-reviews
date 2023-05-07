import type { Materia, Catedra } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const load = (async ({ params }) => {
	let response = await fetch(`${env.BACKEND_URL}/materias/${params.codigo}`);
	const materia = (await response.json()) as Materia;

	response = await fetch(`${env.BACKEND_URL}/materias/${params.codigo}/catedras`);
	const catedras = (await response.json()) as Catedra[];

	return { materia, catedras };
}) satisfies PageServerLoad;
