import {
  PUBLIC_POSTHOG_API_HOST_URL,
  PUBLIC_POSTHOG_PROJECT_API_KEY,
  PUBLIC_POSTHOG_UI_HOST_URL
} from "$env/static/public";
import type { LayoutLoad } from "./$types";
import { browser } from "$app/environment";
import posthog from "posthog-js";

export const load: LayoutLoad = async ({ data }) => {
  if (browser) {
    posthog.init(PUBLIC_POSTHOG_PROJECT_API_KEY, {
      api_host: PUBLIC_POSTHOG_API_HOST_URL,
      ui_host: PUBLIC_POSTHOG_UI_HOST_URL,
      defaults: "2026-01-30"
    });
  }

  return data;
};
