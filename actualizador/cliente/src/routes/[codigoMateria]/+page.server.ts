import { BACKEND_URL } from "$env/static/private";
import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";

type PatchMateria = {
	codigo: string;
	nombre: string;
	docentes: {
		nombre: string;
		rol: string;
		matches: {
			codigo: string;
			nombre: string;
			similitud: number;
		}[];
	}[];
	catedras: {
		codigo: number;
		docentes: {
			nombre: string;
			rol: string;
		}[];
	}[];
	cuatrimestre: {
		numero: number;
		anio: number;
	};
};

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
	default: async ({ request }) => {
		console.log(await request.formData());
	}
} satisfies Actions;
