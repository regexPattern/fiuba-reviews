import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";

export const load = (async ({ fetch }) => {
	const materias = await prisma.materia.findMany();

	return { materias };
}) satisfies PageServerLoad;
