import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
  server: {
    DATABASE_URL: z.string().min(1),
    TURNSTILE_SECRET_KEY: z.string().min(1),
    REVALIDATE_SECRET: z.string().min(1)
  },
  client: {
    NEXT_PUBLIC_BASE_URL: z.url(),
    NEXT_PUBLIC_TURNSTILE_SITE_KEY: z.string().min(1),
    NEXT_PUBLIC_POSTHOG_PROJECT_API_KEY: z.string().min(1),
    NEXT_PUBLIC_POSTHOG_API_HOST_URL: z.url(),
    NEXT_PUBLIC_POSTHOG_UI_HOST_URL: z.url()
  },
  runtimeEnv: {
    DATABASE_URL: process.env.DATABASE_URL,
    TURNSTILE_SECRET_KEY: process.env.TURNSTILE_SECRET_KEY,
    REVALIDATE_SECRET: process.env.REVALIDATE_SECRET,
    NEXT_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_BASE_URL,
    NEXT_PUBLIC_TURNSTILE_SITE_KEY: process.env.NEXT_PUBLIC_TURNSTILE_SITE_KEY,
    NEXT_PUBLIC_POSTHOG_PROJECT_API_KEY: process.env.NEXT_PUBLIC_POSTHOG_PROJECT_API_KEY,
    NEXT_PUBLIC_POSTHOG_API_HOST_URL: process.env.NEXT_PUBLIC_POSTHOG_API_HOST_URL,
    NEXT_PUBLIC_POSTHOG_UI_HOST_URL: process.env.NEXT_PUBLIC_POSTHOG_UI_HOST_URL
  },
  skipValidation: process.env.SKIP_ENV_VALIDATION === "1"
});
