import prisma from "$lib/prisma";
import { promedioDocente } from "$lib/utils";
import type { LayoutServerLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
	const materia = await prisma.materia.findUnique({
		where: { codigo: Number(params.codigoMateria) }
	});

	if (materia === null) {
		throw error(404, { message: "Materia no encontrada" });
	}

	const catedras = await prisma.catedra.findMany({
		where: {
			codigo_materia: Number(params.codigoMateria)
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

		const nombre = docentes
			.map(({ nombre }) => nombre)
			.sort()
			.join("-");

		docentes = docentes.filter((d) => d.calificacion.length != 0);
		const promedio =
			docentes.reduce((acc, curr) => acc + promedioDocente(curr.calificacion), 0) / docentes.length;

		return { codigo: c.codigo, nombre, promedio };
	});

	return {
		materia,
		catedras: catedrasConPromedio.map((c) => ({
			...c,
			codigo_materia: params.codigoMateria
		}))
	};
}) satisfies LayoutServerLoad;
