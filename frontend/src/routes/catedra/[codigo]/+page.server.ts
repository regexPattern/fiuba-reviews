import type { Docente } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const prerender = true;

export const load = (async ({ params }) => {
	const respuesta = await fetch(`${env.BACKEND_URL}/catedra/${params.codigo}/docentes`);
	const payload = (await respuesta.json()) as {
		nombre_catedra: string;
		docentes_con_comentarios: Docente[];
	};

	payload.docentes_con_comentarios.sort((a, b) => b.promedio - a.promedio);

	return { nombre_catedra: payload.nombre_catedra, docentes: payload.docentes_con_comentarios };
}) satisfies PageServerLoad;
