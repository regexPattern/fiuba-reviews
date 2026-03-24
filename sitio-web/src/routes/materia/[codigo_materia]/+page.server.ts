import type { EntryGenerator, PageServerLoad } from "./$types";
import { obtenerMateriasBuscador } from "$lib/server/db/materias";

export const prerender = true;

export const entries: EntryGenerator = async () => {
  const materias = await obtenerMateriasBuscador();

  return materias.map(({ codigo }) => ({ codigo_materia: codigo }));
};

export const load: PageServerLoad = async ({ parent }) => {
  await parent();
};
