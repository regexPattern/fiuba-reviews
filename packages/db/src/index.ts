import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";
import * as schema from "./schema";

/** Crea un cliente de Drizzle sobre `postgres` con el esquema de FIUBA Reviews. */
export function createDb(url: string) {
  return drizzle(postgres(url), { schema, casing: "snake_case" });
}

export type Db = ReturnType<typeof createDb>;

export { schema };
