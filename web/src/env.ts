import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
  server: {
    DATABASE_URL: z.string().min(1),
    REVALIDATE_SECRET: z.string().min(1)
  },
  client: {
    NEXT_PUBLIC_BASE_URL: z.url()
  },
  runtimeEnv: {
    DATABASE_URL: process.env.DATABASE_URL,
    REVALIDATE_SECRET: process.env.REVALIDATE_SECRET,
    NEXT_PUBLIC_BASE_URL: process.env.NEXT_PUBLIC_BASE_URL
  },
  skipValidation: process.env.SKIP_ENV_VALIDATION === "1"
});
