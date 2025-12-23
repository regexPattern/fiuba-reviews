import { BACKEND_URL } from "$env/static/private";
import type { PatchMateria } from "$lib";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import { error, fail } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ params }) => {
	const res = await fetch(`${BACKEND_URL}/${params.codigoMateria}`);

	if (res.statusText !== "OK") {
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
	default: async ({ params: _, request }) => {
		const formData = await request.formData();

		type Resolucion = {
			nombre_db: string;
			codigo_match: string | null | undefined;
		};

		const resoluciones: Record<string, Resolucion> = {};
		const docentesFaltantes: string[] = [];

		for (const [nombre_siu, resJson] of formData.entries()) {
			const res = JSON.parse(resJson as string) as Resolucion;
			if (res.codigo_match === "__UNDEFINED__") {
				docentesFaltantes.push(nombre_siu);
			} else {
				resoluciones[nombre_siu] = res;
			}
		}

		if (docentesFaltantes.length > 0) {
			docentesFaltantes.sort((a, b) => a.localeCompare(b));
			return fail(400, { docentesFaltantes });
		}

		// await fetch(`${BACKEND_URL}/patches/${params.codigoMateria}`, {
		// 	method: "PATCH",
		// 	headers: {
		// 		"Content-Type": "application/json"
		// 	},
		// 	body: JSON.stringify({
		// 		codigos_ya_resueltos: Array.from(codigosYaResueltos),
		// 		resoluciones_actuales: resolucionesActuales
		// 	})
		// });

		return { success: true };
	}
} satisfies Actions;
