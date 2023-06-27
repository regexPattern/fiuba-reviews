import prisma from "$lib/prisma";
import type { RequestHandler } from "./$types";
import { error } from "@sveltejs/kit";

export const GET = (async ({ params }) => {
	const docente = await prisma.docente.findUnique({
		where: { codigo: params.codigoDocente }
	});

	if (docente === null) {
		throw error(404, { message: "Materia no encontrada" });
	}

	return new Response(JSON.stringify(docente));
}) satisfies RequestHandler;

export const POST = (async ({ params }) => {
	const calificacion = await prisma.docente.
}) satisfies RequestHandler;
