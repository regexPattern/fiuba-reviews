import type { PageServerLoad } from "./$types";
import prisma from "$lib/prisma";

export const prerender = true;

export const load = (async ({ params }) => {
	const catedra = await prisma.catedra.findUniqueOrThrow({
		where: {
			codigo: params.codigo,
		},
		select: {
			nombre: true,
			codigo_materia: true,
			catedradocente: {
				select: {
					docentes: {
						include: {
							comentario: true,
						},
					},
				},
			},
		},
	});

	const docentes = catedra.catedradocente.map((rel) => {
		rel.docentes.comentario.sort((a, b) => {
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

		return rel.docentes;
	});

	return { catedra: { nombre: catedra.nombre, codigo_materia: catedra.codigo_materia }, docentes };
}) satisfies PageServerLoad;
