import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
  server: {
    DATABASE_URL: z.url(),
    REVALIDATE_SECRET: z.string().min(32)
  },
  runtimeEnv: {
    DATABASE_URL: process.env.DATABASE_URL,
    REVALIDATE_SECRET: process.env.REVALIDATE_SECRET
  },
  emptyStringAsUndefined: true,
  skipValidation: process.env.SKIP_ENV_VALIDATION === "1"
});
