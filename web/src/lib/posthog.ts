import { PUBLIC_POSTHOG_PROJECT_API_KEY } from "$env/static/public";
import posthog from "posthog-js";

export const posthogInit = () => {
  posthog.init(PUBLIC_POSTHOG_PROJECT_API_KEY, {
    api_host: "https://proxy.fiuba-reviews.com",
    ui_host: "https://us.i.posthog.com",
    defaults: "2026-01-30"
  });
};
