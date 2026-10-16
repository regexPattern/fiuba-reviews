import posthog from "posthog-js";

posthog.init(process.env.NEXT_PUBLIC_POSTHOG_PROJECT_API_KEY!, {
  api_host: process.env.NEXT_PUBLIC_POSTHOG_API_HOST_URL,
  ui_host: process.env.NEXT_PUBLIC_POSTHOG_UI_HOST_URL,
  defaults: "2026-01-30"
});
