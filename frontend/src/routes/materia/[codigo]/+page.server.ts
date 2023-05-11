import type { Catedra } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const prerender = true;

export const load = (async ({ params }) => {
	const respuesta = await fetch(`${env.BACKEND_URL}/materia/${params.codigo}/catedras`);
	const payload = (await respuesta.json()) as {
		nombre_materia: string;
		codigos_equivalencias: number[],
		catedras: Catedra[];
	};

	payload.catedras.sort((a, b) => b.promedio - a.promedio);

	return {
		nombre_materia: payload.nombre_materia,
		codigos_equivalencias: payload.codigos_equivalencias,
		codigo_materia: params.codigo,
		catedras: payload.catedras,
	};
}) satisfies PageServerLoad;
