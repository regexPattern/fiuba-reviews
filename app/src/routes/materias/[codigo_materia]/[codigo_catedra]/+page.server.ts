import type { PageServerLoad } from "./$types";

import prisma from "$lib/prisma";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
	const catedra = await prisma.catedra.findUnique({
		where: {
			codigo: params.codigo_catedra
		},
		select: {
			catedradocente: {
				select: {
					docente: {
						include: {
							comentario: true
						}
					}
				}
			}
		}
	});

	if (catedra === null) {
		throw error(404, { message: "Not found" });
	}

	const docentes = catedra.catedradocente
		.map((d) => d.docente)
		.sort((a, b) => b.promedio - a.promedio);

	const nombre_catedra = docentes
		.map((d) => d.nombre)
		.sort()
		.join("-");

	for (const d of docentes) {
		d.comentario.sort((a, b) => {
			const [cuatriA, anioA] = a.cuatrimestre.split("Q");
			const [cuatriB, anioB] = b.cuatrimestre.split("Q");

			if (anioA < anioB) {
				return 1;
			} else if (anioA > anioB) {
				return -1;
			} else {
				return cuatriA <= cuatriB ? 1 : -1;
			}
		});
	}

	return {
		codigo_materia: params.codigo_materia,
		codigo_catedra: params.codigo_catedra,
		nombre_catedra,
		docentes
	};
}) satisfies PageServerLoad;
