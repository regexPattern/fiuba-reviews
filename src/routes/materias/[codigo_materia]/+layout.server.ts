import prisma from "$lib/prisma";
import * as utils from "$lib/utils";
import type { LayoutServerLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
	const codigoMateria = Number(params.codigo_materia) || 0;

	const materia = await prisma.materia.findUnique({
		where: {
			codigo: codigoMateria
		}
	});

	if (!materia) {
		throw error(404, { message: "Materia no encontrada" });
	}

	const catedras = await prisma.catedra.findMany({
		where: {
			codigo_materia: codigoMateria
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

	const catedrasConPromedio = catedras.map((c) => {
		let docentes = c.catedradocente.map(({ docente }) => ({ ...docente }));

		const nombre = utils.fmtNombreCatedra(docentes);
		docentes = docentes.filter((d) => d.calificacion.length != 0);

		const promedio =
			docentes.reduce((acc, curr) => acc + utils.calcPromedioDocente(curr), 0) / docentes.length ||
			0;

		return {
			codigo: c.codigo,
			nombre,
			promedio
		};
	});

	return {
		materia,
		catedras: catedrasConPromedio
			.map((c) => ({ ...c, codigo_materia: params.codigo_materia }))
			.sort((a, b) => b.promedio - a.promedio)
	};
}) satisfies LayoutServerLoad;
