import { browser, dev } from "$app/environment";
import { env } from "$env/dynamic/public";
import posthog from "posthog-js";
import type { LayoutLoad } from "./$types";

export const load: LayoutLoad = async ({ data }) => {
	if (!dev && browser) {
		if (env.PUBLIC_POSTHOG_PROJECT_API_KEY) {
			posthog.init(env.PUBLIC_POSTHOG_PROJECT_API_KEY, {
				api_host: "https://us.i.posthog.com",
				person_profiles: "never"
			});
		}
	}

	return data;
};
