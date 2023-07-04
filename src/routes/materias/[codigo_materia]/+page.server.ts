import db from "$lib/db";
import * as schema from "$lib/db/schema";
import type { PageServerLoad } from "./$types";
import { eq } from "drizzle-orm";

export const load = (async ({ params }) => {
	const equivalencias = await db
		.select({ nombre: schema.materia.nombre, codigo: schema.materia.codigo })
		.from(schema.materia)
		.where(eq(schema.materia.codigoEquivalencia, Number(params.codigo_materia)));

  return { equivalencias };
}) satisfies PageServerLoad;
