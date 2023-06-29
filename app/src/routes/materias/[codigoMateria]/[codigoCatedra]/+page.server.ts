import prisma, { cuatrimestres } from "$lib/prisma";
import { compararCuatrimestre, promedioDocente } from "$lib/utils";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import { fail, error } from "@sveltejs/kit";
import { z } from "zod";
import { zfd } from "zod-form-data";

export const load = (async ({ params }) => {
	const catedra = await prisma.catedra.findUnique({
		where: { codigo: params.codigoCatedra },
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

	if (catedra === null) {
		throw error(404, { message: "Catedra no encontrada" });
	}

	const docentes = catedra.catedradocente
		.map((d) => {
			d.docente.comentario.sort((a, b) => compararCuatrimestre(a.cuatrimestre, b.cuatrimestre));
			return { ...d.docente, promedio: promedioDocente(d.docente.calificacion) };
		})
		.sort((a, b) => b.promedio - a.promedio);

	const nombre = docentes
		.map((d) => d.nombre)
		.sort()
		.join("-");

	return {
		catedra: { nombre, docentes },
		cuatrimestres: cuatrimestres.map((c) => c.nombre).sort(compararCuatrimestre)
	};
}) satisfies PageServerLoad;

const puntaje1Al5Schema = zfd.numeric(z.number().min(1).max(5));
const calificacionSchema = zfd.formData({
	codigo_docente: zfd.text(),
	acepta_critica: puntaje1Al5Schema,
  asistencia: puntaje1Al5Schema,
  buen_trato: puntaje1Al5Schema,
  claridad: puntaje1Al5Schema,
  clase_organizada: puntaje1Al5Schema,
  cumple_horarios: puntaje1Al5Schema,
  fomenta_participacion: puntaje1Al5Schema,
  panorama_amplio: puntaje1Al5Schema,
  responde_mails: puntaje1Al5Schema,
  comentario: zfd.text().optional(),
});

export const actions = {
	default: async ({ request }) => {
		const formData = await request.formData();
    const result = calificacionSchema.safeParse(formData);

    if (!result.success) {
      return fail(422, { issues: result.error.issues });
    }
	}
} satisfies Actions;
