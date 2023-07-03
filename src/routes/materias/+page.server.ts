import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";

export const load = (async () => {
	const materias = await prisma.materias.findMany();

	return {
		materias: materias.sort((a, b) => a.codigo - b.codigo)
	};
}) satisfies PageServerLoad;
