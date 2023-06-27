import type { LayoutServerLoad } from "./$types";
import type { materia } from "@prisma/client";
import type { Catedra } from "./catedras/+server";

export const load = (async ({ fetch, params }) => {
	const res_materia = await fetch(`/materias/${params.codigoMateria}`);
	const res_catedras = await fetch(`/materias/${params.codigoMateria}/catedras`);

	const materia = (await res_materia.json()) as materia;
	const catedras = ((await res_catedras.json()) as Catedra[]).map((c) => ({
		...c,
		codigo_materia: params.codigoMateria
	}));

	return { catedras };
}) satisfies LayoutServerLoad;
