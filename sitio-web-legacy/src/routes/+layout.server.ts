import "@valibot/i18n/es";
import type { LayoutServerLoad } from "./$types";
import { obtenerMateriasBuscador } from "$lib/server/db/materias";
import * as v from "valibot";

v.setGlobalConfig({ lang: "es" });

export const load: LayoutServerLoad = async ({ url }) => {
  const materiaRows = await obtenerMateriasBuscador();

  return {
    materias: materiaRows,
    showBuscador: url.pathname !== "/"
  };
};
