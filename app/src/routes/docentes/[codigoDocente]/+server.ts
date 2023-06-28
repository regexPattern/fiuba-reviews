import prisma from "$lib/prisma";
import type { RequestHandler } from "./$types";
import { json, error } from "@sveltejs/kit";

export const GET = (async ({ params }) => {
	const docente = await prisma.docente.findUnique({
		where: { codigo: params.codigoDocente }
	});

	if (docente === null) {
		throw error(404, { message: "Materia no encontrada" });
	}

	return json(docente);
}) satisfies RequestHandler;

export const POST = (async ({ request }) => {
	console.log(await request.text());

	return json({})
}) satisfies RequestHandler;
