import prisma from "$lib/prisma";
import type { RequestHandler } from "./$types";

export const GET = (async () => {
	const materias = await prisma.materia.findMany();

	return new Response(JSON.stringify(materias));
}) satisfies RequestHandler;
