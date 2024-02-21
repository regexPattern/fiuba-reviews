import db from "$lib/db";
import { catedra, catedraDocente, cuatrimestre, docente, materia } from "$lib/db/schema";
import { sortCuatrimestres } from "$lib/utils";
import schema from "$lib/zod/schema";
import { error, fail } from "@sveltejs/kit";
import { eq } from "drizzle-orm";
import { message, setError, superValidate } from "sveltekit-superforms/server";

import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";

export const load: PageServerLoad = async ({ params }) => {
	let docentes;

	try {
		docentes = await db
			.select({
				nombreDocente: docente.nombre,
				codigoMateria: materia.codigo,
				codigoCatedra: catedraDocente.codigoCatedra
			})
			.from(docente)
			.innerJoin(catedraDocente, eq(docente.codigo, catedraDocente.codigoDocente))
			.innerJoin(catedra, eq(catedraDocente.codigoCatedra, catedra.codigo))
			.innerJoin(materia, eq(catedra.codigoMateria, materia.codigo))
			.where(eq(docente.codigo, params.codigoDocente))
			.limit(1);
	} catch (e: any) {
		// Si me pasan un código que no es serializable como UUID.
		if (e.code === "22P02") {
			throw error(404, "Docente no encontrado.");
		} else {
			throw e;
		}
	}

	if (docentes.length === 0) {
		throw error(404, "Docente no encontrado");
	}

	const cuatrimestres = await db.select().from(cuatrimestre);
	cuatrimestres.sort((a, b) => -sortCuatrimestres(a.nombre, b.nombre));

	const form = superValidate(schema);

	return {
		...docentes[0],
		cuatrimestres,
		form
	};
};

export const actions: Actions = {
	default: async ({ request }) => {
		const form = await superValidate(request, schema);

		if (!form.valid) {
			console.log(form.errors);
			return message(form, "Datos inválidos");
		}

		const esCuatrimestreValido =
			(
				await db
					.select()
					.from(cuatrimestre)
					.where(eq(cuatrimestre.nombre, form.data.cuatrimestre || ""))
					.limit(1)
			).length === 0;

		if (form.data.cuatrimestre && form.data.cuatrimestre != "undefined" && esCuatrimestreValido) {
			return setError(form, "cuatrimestre", `Cuatrimestre '${form.data.cuatrimestre}' no existe`);
		}

		try {
			// No estoy insertando las calificaciones realmente, pero estoy haciendo
			// el error handling como si lo hiciera. Uso un timeout para fingir la
			// operación de inserción.
			//
			await new Promise((resolve) => setTimeout(resolve, 1000));
			console.log(form.data);
		} catch {
			return fail(500);
		}

		return message(form, "Calificación registrada con éxito");
	}
};
