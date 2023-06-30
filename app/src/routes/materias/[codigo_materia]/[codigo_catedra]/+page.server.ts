import prisma from "$lib/prisma";
import { comparar_cuatrimestre, generar_nombre_catedra, calcular_promedio_docente } from "$lib/utils";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import calificacion_schema from "./schema";
import { error, fail } from "@sveltejs/kit";

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
							calificacion: true,
							comentario: true
						}
					}
				}
			}
		}
	});

	if (!catedra) {
		throw error(404, { message: "Catedra no encontrada" });
	}

	const docentes = catedra.catedradocente
		.map((d) => {
			d.docente.comentario.sort((a, b) => comparar_cuatrimestre(a.cuatrimestre, b.cuatrimestre));
			return { ...d.docente, promedio: calcular_promedio_docente(d.docente.calificacion) };
		})
		.sort((a, b) => b.promedio - a.promedio);

	const cuatrimestres = await prisma.cuatrimestre.findMany();
	const cuatrimestres_ordenados = cuatrimestres.map((c) => c.nombre).sort(comparar_cuatrimestre);

	return {
		catedra: {
			nombre: generar_nombre_catedra(docentes),
			docentes
		},
		cuatrimestres: cuatrimestres_ordenados
	};
}) satisfies PageServerLoad;

export const actions = {
	default: async ({ request }) => {
		const datos_formulario = await request.formData();
		const resultado = calificacion_schema.safeParse(datos_formulario);

		if (!resultado.success) {
			return fail(422, { issues: resultado.error.issues });
		}
	}
} satisfies Actions;
