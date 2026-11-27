import "server-only";
import { env } from "@/env";
import { createDb } from "@fiuba-reviews/db";

export const db = createDb(env.DATABASE_URL);

export { schema } from "@fiuba-reviews/db";
