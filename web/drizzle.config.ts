import { defineConfig } from "drizzle-kit";

if (!process.env.DATABASE_URL) throw new Error("variable DATABASE_URL no está definida");

export default defineConfig({
  schema: "./src/lib/server/db/schema.ts",
  dialect: "postgresql",
  dbCredentials: { url: process.env.DATABASE_URL },
  verbose: true,
  strict: true
});
