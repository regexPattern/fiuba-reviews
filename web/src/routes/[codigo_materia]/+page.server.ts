import { redirect } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ params, parent }) => {
  const { catedras } = await parent();

  if (catedras.length > 0) {
    redirect(307, `/${params.codigo_materia}/${catedras[0].codigo}`);
  }
};
