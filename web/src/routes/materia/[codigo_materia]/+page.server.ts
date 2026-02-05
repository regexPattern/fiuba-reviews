import type { PageServerLoad } from "./$types";
import { redirect } from "@sveltejs/kit";

// export const prerender = true;

export const load: PageServerLoad = async ({ params, parent }) => {
  const { catedras } = await parent();

  if (catedras.length > 0) {
    redirect(307, `/materia/${params.codigo_materia}/${catedras[0].codigo}`);
  }
};
