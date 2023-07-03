import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
	const materia = await prisma.materias.findUnique({
		where: { codigo: Number(params.codigo_materia) || 0 },
		include: {
			other_materias: {
				select: {
					nombre: true,
					codigo: true
				}
			}
		}
	});

	if (!materia) {
		throw error(404, { message: "Materia no encontrada" });
	}

	return {
		materia: {
			codigo: materia.codigo,
			nombre: materia.nombre,
			equivalencias: materia.other_materias
		}
	};
}) satisfies PageServerLoad;
