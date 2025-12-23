import type { PageLoad } from "./$types";
import { redirect } from "@sveltejs/kit";

export const load: PageLoad = async ({ parent }) => {
	const data = await parent();
	redirect(307, `/${data.patches[0].codigo}`);
};
