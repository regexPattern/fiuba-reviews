import { redirect } from "@sveltejs/kit";

import type { PageServerLoad } from "./$types";

export const prerender = true;

export const load: PageServerLoad = async ({ parent }) => {
	const layoutData = await parent();

	const codigoMateria = layoutData.materia.codigo;
	const defaultCatedra = layoutData.catedras[0];

	throw redirect(307, `/materias/${codigoMateria}/${defaultCatedra.codigo}`);
};
