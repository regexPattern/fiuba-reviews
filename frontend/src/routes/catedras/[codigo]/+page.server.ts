import type { Catedra, Comentario, Docente } from "$lib/types";
import type { PageServerLoad } from "./$types";

import { env } from "$env/dynamic/private";

export const load = (async ({ params }) => {
	let response = await fetch(`${env.BACKEND_URL}/catedras/${params.codigo}`);
	const catedra = (await response.json()) as Catedra;

	response = await fetch(`${env.BACKEND_URL}/catedras/${params.codigo}/docentes`);
	const codigos_docentes = (await response.json()) as string[];

	const calificaciones = await Promise.all(codigos_docentes.map(async (codigo) => {
		const response = await fetch(`${env.BACKEND_URL}/docentes/${codigo}`);
		return (await response.json()) as Docente;
	}));

	const docentes = await Promise.all(calificaciones.map(async (docente) => {
		const response = await fetch(`${env.BACKEND_URL}/comentarios/${docente.codigo}`);
		const comentarios = (await response.json()) as Comentario[];
		return { ...docente, comentarios };
	}));

	return { catedra, docentes };
}) satisfies PageServerLoad;
