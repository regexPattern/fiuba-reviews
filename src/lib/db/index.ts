import { DATABASE_URL } from "$env/static/private";
import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

import * as schema from "./schema";

const client = postgres(DATABASE_URL);

export default drizzle(client, { schema });
