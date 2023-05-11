import type { PageServerLoad } from "./$types";
import prisma from "$lib/prisma";

export const prerender = true;

export const load = (async ({ params }) => {
	const materia = await prisma.materia.findUniqueOrThrow({
		where: {
			codigo: Number(params.codigo),
		},
		include: {
			other_materia: {
				select: {
					codigo: true,
				},
			},
			catedra: {
				select: {
					codigo: true,
					nombre: true,
					promedio: true,
				},
			},
		},
	});

	return { materia };
}) satisfies PageServerLoad;
