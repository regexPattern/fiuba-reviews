import prisma from "$lib/prisma";
import type { RequestHandler } from "./$types";
import { error, json } from "@sveltejs/kit";

export const GET = (async ({ params }) => {
	const materia = await prisma.materia.findUnique({
		where: { codigo: Number(params.codigoMateria) }
	});

	if (materia === null) {
		throw error(404, { message: "Materia no encontrada" });
	}

	return json(materia);
}) satisfies RequestHandler;
