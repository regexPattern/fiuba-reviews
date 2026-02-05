import { env } from "$env/dynamic/private";
import * as schema from "./schema";
import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

if (!env.DATABASE_URL) {
  throw new Error("variable DATABASE_URL no está definida");
}

const client = postgres(env.DATABASE_URL);

export const db = drizzle(client, { schema });

export { schema };
