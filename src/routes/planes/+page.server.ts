import { contenidoSiu as schema } from "$lib/zod/schema";
import type { Actions, PageServerLoad } from "./$types";
import { message, superValidate } from "sveltekit-superforms/server";

export const load: PageServerLoad = async () => {
  return { form: await superValidate(schema) };
};

export const actions: Actions = {
  default: async ({ request }) => {
    const form = await superValidate(request, schema);

    if (!form.valid) {
      console.log(form.errors);
      return message(form, "Datos inválidos");
    }

    return message(form, "Contenido del SIU registrado con éxito");
  },
};
