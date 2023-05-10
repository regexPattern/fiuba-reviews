import type { Catedra } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const load = (async ({ params }) => {
	const respuesta = await fetch(`${env.BACKEND_URL}/materia/${params.codigo}/catedras`);
	const payload = (await respuesta.json()) as {
		nombre_materia: string;
		catedras: Catedra[];
	};

	return {
		nombre_materia: payload.nombre_materia,
		codigo_materia: params.codigo,
		catedras: payload.catedras,
	};
}) satisfies PageServerLoad;
