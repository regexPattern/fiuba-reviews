import prisma from "$lib/prisma";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
	const catedra = await prisma.catedra.findUnique({
		where: {
			codigo: params.codigoCatedra
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
		throw error(404, { message: "Catedra no encontrada" });
	}

	const docentes = catedra.catedradocente
		.map((d) => d.docente)
		.sort((a, b) => b.promedio - a.promedio);

	const nombre = docentes
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
		catedra: { nombre, docentes }
	};
}) satisfies PageServerLoad;

export const actions = {
	default: async ({ request }) => {
		const formData = await request.formData();
		const codigoDocente = formData.get("codigo");
		const comentario = formData.get("comentario");
	}
} satisfies Actions;
