import type { LayoutLoad } from "./$types";
import posthog from "posthog-js";
import { browser } from "$app/environment";
import { PUBLIC_POSTHOG_PROJECT_API_KEY } from "$env/static/public";

export const load: LayoutLoad = async ({ data }) => {
  if (browser) {
    posthog.init(PUBLIC_POSTHOG_PROJECT_API_KEY, {
      api_host: "https://us.i.posthog.com",
      defaults: "2026-01-30"
    });
  }

  return data;
};
