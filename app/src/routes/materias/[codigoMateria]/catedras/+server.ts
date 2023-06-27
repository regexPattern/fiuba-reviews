import prisma from "$lib/prisma";
import { parseCodigoMateria } from "$lib/utils";
import type { RequestHandler } from "./$types";

export type Catedra = Awaited<ReturnType<typeof getCatedras>>[number];

export const GET = (async ({ params }) => {
	const catedras = await getCatedras(parseCodigoMateria(params.codigoMateria));
	catedras.sort((a, b) => b.promedio - a.promedio);

	return new Response(JSON.stringify(catedras));
}) satisfies RequestHandler;

async function getCatedras(codigoMateria: number) {
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
							promedio: true,
							respuestas: true
						}
					}
				}
			}
		}
	});

	return catedras.map((c) => {
		let docentes = c.catedradocente.map(({ docente }) => ({ ...docente }));

		const nombre = docentes
			.map(({ nombre }) => nombre)
			.sort()
			.join("-");

		docentes = docentes.filter((d) => d.respuestas != 0);
		const promedio = docentes.reduce((acc, curr) => acc + curr.promedio, 0) / docentes.length;

		return { codigo: c.codigo, nombre, promedio };
	});
}
