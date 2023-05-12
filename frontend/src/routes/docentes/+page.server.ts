import type { Actions, PageServerLoad } from "./$types";
import { fail } from "@sveltejs/kit";
import prisma from "$lib/prisma";

export const actions = {
	default: async ({ request }) => {
		const data = await request.formData();

		const contenido = data.get("contenido");
		const cuatrimestre = data.get("cuatrimestre");
		const codigo_docente = data.get("codigoDocente");

		if (!contenido || !cuatrimestre || !codigo_docente) {
			return fail(400, { contenido, cuatrimestre, codigo_docente, missing: true });
		}

		if (
			typeof contenido != "string" ||
			typeof cuatrimestre != "string" ||
			typeof codigo_docente != "string"
		) {
			return fail(400, { incorrect: true });
		}

		await prisma.comentario.create({
			data: {
				codigo: "testing",
				contenido,
				cuatrimestre,
				codigo_docente,
			},
		});
	},
} satisfies Actions;

export const load = (async () => {
	const docente = await prisma.docente.findMany();
	return { docente };
}) satisfies PageServerLoad;
