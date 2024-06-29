import { browser, dev } from "$app/environment";
import { env } from "$env/dynamic/private";
import { injectSpeedInsights } from "@vercel/speed-insights/sveltekit";
import posthog from "posthog-js";

export const load = async () => {
	if (!dev && browser) {
		posthog.init(env.POSTHOG_PROJECT_API_KEY, {
			api_host: "https://us.i.posthog.com",
			person_profiles: "never"
		});
		injectSpeedInsights();
	}
};
