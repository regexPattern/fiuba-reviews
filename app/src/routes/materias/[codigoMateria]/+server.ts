import prisma from "$lib/prisma";
import { parseCodigoMateria } from "$lib/utils";
import type { RequestHandler } from "./$types";
import { error } from "@sveltejs/kit";

export const GET = (async ({ params }) => {
	const materia = await prisma.materia.findUnique({
		where: { codigo: parseCodigoMateria(params.codigoMateria) }
	});

	if (materia === null) {
		throw error(404, { message: "Materia no encontrada" });
	}

	return new Response(JSON.stringify(materia));
}) satisfies RequestHandler;
