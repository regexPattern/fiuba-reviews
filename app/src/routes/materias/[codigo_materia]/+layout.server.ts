import prisma from "$lib/prisma";
import { calcular_promedio_docente } from "$lib/utils";
import type { LayoutServerLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
	const codigo_materia = Number(params.codigo_materia) || 0;

	const materia = await prisma.materia.findUnique({
		where: {
			codigo: codigo_materia
		}
	});

	if (materia === null) {
		throw error(404, { message: "Materia no encontrada" });
	}

	const catedras = await prisma.catedra.findMany({
		where: {
			codigo_materia
		},
		include: {
			catedradocente: {
				include: {
					docente: {
						select: {
							nombre: true,
							calificacion: true
						}
					}
				}
			}
		}
	});

	const catedras_con_promedio = catedras.map((c) => {
		let docentes = c.catedradocente.map(({ docente }) => ({ ...docente }));

		const nombre = docentes
			.map(({ nombre }) => nombre)
			.sort()
			.join(", ");

		// Al momento de calcular el promedio de la catedra, no se toman en cuenta
		// los docentes que no tienen calificaciones.
		docentes = docentes.filter((d) => d.calificacion.length != 0);

		const promedio =
			docentes.reduce((acc, curr) => acc + calcular_promedio_docente(curr.calificacion), 0) / docentes.length;

		return {
			codigo: c.codigo,
			nombre,
			promedio
		};
	});

	const catedras_ordenadas_por_promedio = catedras_con_promedio
		.map((c) => ({ ...c, codigo_materia: params.codigo_materia }))
		.sort((a, b) => b.promedio - a.promedio);

	return {
		materia,
		catedras: catedras_ordenadas_por_promedio
	};
}) satisfies LayoutServerLoad;
