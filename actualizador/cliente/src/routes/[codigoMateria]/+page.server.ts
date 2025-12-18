import { BACKEND_URL } from "$env/static/private";
import type { PatchMateria } from "$lib";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
	const res = await fetch(`${BACKEND_URL}/patches/${params.codigoMateria}`);
	const patch = (await res.json()) as PatchMateria;

	patch.docentes.sort((a, b) => (a.nombre === b.nombre ? 0 : a.nombre > b.nombre ? 1 : -1));

	const docentesNuevos = new Set<string>();

	for (const doc of patch.docentes) {
		if (doc.matches.length === 0) {
			docentesNuevos.add(doc.nombre);
		}
	}

	return { patch };
};

export const actions = {
	default: async ({ params, request }) => {
		const formData = await request.formData();

		const resoluciones: Record<string, { nombre_db: string; codigo_match: string }> = {};
		for (const [key, value] of formData.entries()) {
			const parsed = JSON.parse(value as string);
			resoluciones[key] = parsed;
		}

		await fetch(`${BACKEND_URL}/patches/${params.codigoMateria}`, {
			method: "PATCH",
			headers: {
				"Content-Type": "application/json"
			},
			body: JSON.stringify(resoluciones)
		});
	}
} satisfies Actions;
