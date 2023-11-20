import { superValidate } from "sveltekit-superforms/server";

import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import schema from "./schema";

export const prerender = false;

export const load: PageServerLoad = async () => {
	const form = superValidate(schema);

	return { form };
};

export const actions: Actions = {
	default: async (e) => {
    const form = await superValidate(e, schema);
    console.log(form);
    return { form };
  }
};
