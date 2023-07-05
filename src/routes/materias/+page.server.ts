import db from "$lib/db";
import { catedra, materia } from "$lib/db/schema";
import type { PageServerLoad } from "./$types";
import { asc, eq, exists, isNotNull, or } from "drizzle-orm";

export const prerender = true;

export const load = (async () => {
	const materias = await db
		.select()
		.from(materia)
		.where(
			or(
				exists(db.select().from(catedra).where(eq(catedra.codigoMateria, materia.codigo))),
				isNotNull(materia.codigoEquivalencia)
			)
		)
		.orderBy(asc(materia.codigo));

	return { materias };
}) satisfies PageServerLoad;
