import type { PageServerLoad } from "./$types";
import prisma from "$lib/prisma";

export const prerender = true;

export const load = (async () => {
	const materias = await prisma.materia.findMany();
	return { materias };
}) satisfies PageServerLoad;
