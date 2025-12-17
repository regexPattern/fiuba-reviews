import { BACKEND_URL } from "$env/static/private";
import type { LayoutServerLoad } from "./$types";

type PatchMateria = {
	codigo: string;
	nombre: string;
	docentes: number;
};

export const load: LayoutServerLoad = async () => {
	const res = await fetch(`${BACKEND_URL}/patches`);
	const patches = (await res.json()) as PatchMateria[];
	return { patches };
};
