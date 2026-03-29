import "@valibot/i18n/es";
import type { LayoutServerLoad } from "./$types";
import { getMateriasDisponibles } from "$lib/server/db/utils";
import * as v from "valibot";

v.setGlobalConfig({ lang: "es" });

export const load: LayoutServerLoad = async ({ url }) => {
  return {
    materias: await getMateriasDisponibles(),
    showBuscador: url.pathname !== "/"
  };
};
