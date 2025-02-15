import { DATABASE_URL } from "$env/static/private";
import * as schema from "./schema";
import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

const client = postgres(DATABASE_URL);

export default drizzle(client, { schema });
