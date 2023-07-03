import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";

export const load = (async () => {
	const materias = await prisma.materias.findMany();

	materias.sort((a, b) => a.codigo - b.codigo);

	return { materias };
}) satisfies PageServerLoad;
