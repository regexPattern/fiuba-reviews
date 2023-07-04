import prisma from "$lib/prisma";
import * as utils from "$lib/utils";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import { error, fail } from "@sveltejs/kit";
import { z } from "zod";
import { zfd } from "zod-form-data";

export const load = (async ({ params }) => {
	const catedra = await prisma.catedra.findUnique({
		where: { codigo: params.codigo_catedra },
		select: {
			catedra_docente: {
				select: {
					docente: {
						include: {
							calificacion: true,
							comentario: {
								select: {
									contenido: true,
									cuatrimestre: true
								}
							}
						}
					}
				}
			}
		}
	});

	if (!catedra) {
		throw error(404, { message: "Catedra no encontrada" });
	}

	const docentes = catedra.catedra_docente.map(({ docente: d }) => {
		return {
			...d,
			comentarios: d.comentario.sort((a, b) =>
				utils.cmpCuatrimestre(a.cuatrimestre, b.cuatrimestre)
			),
			promedio: utils.calcPromedioDocente(d)
		};
	});

	const cuatrimestres = await prisma.cuatrimestre.findMany();

	return {
		catedra: {
			nombre: utils.fmtNombreCatedra(docentes),
			docentes: docentes.sort((a, b) => b.promedio - a.promedio)
		},
		cuatrimestres: cuatrimestres.map((c) => c.nombre).sort(utils.cmpCuatrimestre)
	};
}) satisfies PageServerLoad;

export const actions = {
	default: async ({ request }) => {
		const dataFormulario = await request.formData();
		const parse = schema.safeParse(dataFormulario);

		if (!parse.success) {
			return fail(422, { errores: parse.error.issues });
		}
	}
} satisfies Actions;

const calificacionNumerica = zfd.numeric(z.number().min(1).max(5));

const schema = zfd.formData({
	codigo_docente: zfd.text(),
	acepta_critica: calificacionNumerica,
	asistencia: calificacionNumerica,
	buen_trato: calificacionNumerica,
	claridad: calificacionNumerica,
	clase_organizada: calificacionNumerica,
	cumple_horarios: calificacionNumerica,
	fomenta_participacion: calificacionNumerica,
	panorama_amplio: calificacionNumerica,
	responde_mails: calificacionNumerica,
	comentario: zfd.text().optional()
});
