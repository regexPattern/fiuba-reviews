import { browser, dev } from "$app/environment";
import { env } from "$env/dynamic/public";
import { injectSpeedInsights } from "@vercel/speed-insights/sveltekit";
import posthog from "posthog-js";

export const load = async () => {
	if (!dev && browser) {
		if (env.PUBLIC_POSTHOG_PROJECT_API_KEY) {
			posthog.init(env.PUBLIC_POSTHOG_PROJECT_API_KEY, {
				api_host: "https://us.i.posthog.com",
				person_profiles: "never"
			});
		}
		injectSpeedInsights();
	}
};
