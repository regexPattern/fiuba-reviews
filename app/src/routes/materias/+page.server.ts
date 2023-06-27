import type { PageServerLoad } from "./$types";
import type { materia } from "@prisma/client";

export const load = (async ({ fetch }) => {
	const res = await fetch("/materias");
	const materias = (await res.json()) as materia[];

	return { materias };
}) satisfies PageServerLoad;
