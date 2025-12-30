import { BACKEND_URL } from "$env/static/private";
import type { PatchMateria } from "$lib";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import { error, redirect } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ params }) => {
	const res = await fetch(`${BACKEND_URL}/${params.codigoMateria}`);

	if (res.status >= 400) {
		const errMsg = await res.text();
		error(res.status, { message: errMsg });
	}

	const patch = (await res.json()) as PatchMateria;

	patch.docentes_pendientes.sort((a, b) => {
		const matchDiff = b.matches.length - a.matches.length;
		if (matchDiff !== 0) {
			return matchDiff;
		}
		return a.nombre.localeCompare(b.nombre);
	});

	return { patch };
};

export const actions = {
	default: async ({ params, request }) => {
		const formData = await request.formData();

		type Resolucion = {
			nombre_db: string;
			codigo_match: string | null;
		};

		const resoluciones = new Map<string, Resolucion>();

		for (const [nombre_siu, resJson] of formData.entries()) {
			const res = JSON.parse(resJson as string) as Resolucion;
			switch (res.codigo_match) {
				case "":
					continue;
				case "__CREATE__":
					res.codigo_match = null;
			}
			resoluciones.set(nombre_siu, res);
		}

		const body = JSON.stringify(
			Array.from(resoluciones, ([nombreSiu, res]) => ({
				nombre_siu: nombreSiu,
				...res
			}))
		);

		const res = await fetch(`${BACKEND_URL}/${params.codigoMateria}`, {
			method: "PATCH",
			headers: { "Content-Type": "application/json" },
			body
		});

		if (res.status >= 400) {
			const errMsg = await res.text();
			error(res.status, { message: errMsg });
		}

		redirect(303, "/success");
	}
} satisfies Actions;
