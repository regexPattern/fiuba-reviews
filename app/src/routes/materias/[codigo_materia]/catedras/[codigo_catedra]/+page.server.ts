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

	let docentes = catedra.catedradocente.map((d) => d.docente);

	docentes.sort((a, b) => a.nombre.localeCompare(b.nombre));
	docentes.sort((a, b) => {
		if (a.respuestas === 0) {
			return 1;
		} else if (b.respuestas === 0) {
			return -1;
		} else {
			return 0;
		}
	});

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

	return { docentes };
}) satisfies PageServerLoad;
