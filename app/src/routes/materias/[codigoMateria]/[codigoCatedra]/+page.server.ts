import prisma from "$lib/prisma";
import { promedioDocente } from "$lib/utils";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import { error } from "@sveltejs/kit";

function ordenarCuatrimestre(a: string, b: nombre): number {
	const [cuatriA, anioA] = a.split("Q");
	const [cuatriB, anioB] = b.split("Q");

	if (anioA < anioB) {
		return 1;
	} else if (anioA > anioB) {
		return -1;
	} else {
		return cuatriA <= cuatriB ? 1 : -1;
	}
}

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

	const cuatrimestres = await prisma.cuatrimestre.findMany();

	const docentes = catedra.catedradocente
		.map((d) => ({ ...d.docente, promedio: promedioDocente(d.docente.calificacion) }))
		.sort((a, b) => b.promedio - a.promedio);

	const nombre = docentes
		.map((d) => d.nombre)
		.sort()
		.join("-");

	for (const d of docentes) {
		d.comentario.sort((a, b) => ordenarCuatrimestre(a.cuatrimestre, b.cuatrimestre));
	}

	return {
		catedra: { nombre, docentes },
		cuatrimestres: cuatrimestres.map((c) => c.nombre).sort(ordenarCuatrimestre)
	};
}) satisfies PageServerLoad;

export const actions = {
	default: async ({ request }) => {
		const formData = await request.formData();

		const codigoDocente = formData.get("codigo");
		const cuatrimestre = formData.get("cuatrimestre");

		if (codigoDocente === null) {
			throw error(422, { message: "Código de docente requerido" });
		}

		if (cuatrimestre === null) {
			throw error(422, { message: "Cuatrimestre requerido" });
		}

		await prisma.calificacion.create({
			data: {
				codigo_docente: codigoDocente.toString(),
				acepta_critica: Number(formData.get("acepta_critica")) || 0,
				asistencia: Number(formData.get("asistencia")) || 0,
				buen_trato: Number(formData.get("buen_trato")) || 0,
				claridad: Number(formData.get("claridad")) || 0,
				clase_organizada: Number(formData.get("clase_organizada")) || 0,
				cumple_horarios: Number(formData.get("cumple_horarios")) || 0,
				fomenta_participacion: Number(formData.get("fomenta_participacion")) || 0,
				panorama_amplio: Number(formData.get("panorama_amplio")) || 0,
				responde_mails: Number(formData.get("responde_mails")) || 0
			}
		});

		const comentario = formData.get("comentario");

		if (comentario != null) {
			await prisma.comentario.create({
				data: {
					codigo_docente: codigoDocente.toString(),
					contenido: comentario.toString(),
					cuatrimestre: cuatrimestre.toString()
				}
			});
		}
	}
} satisfies Actions;
