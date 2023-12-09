import db from "$lib/db";
import { catedra, catedraDocente, cuatrimestre, docente, materia } from "$lib/db/schema";
import { sortCuatrimestres } from "$lib/utils";
import { error, fail } from "@sveltejs/kit";
import { eq } from "drizzle-orm";
import { message, superValidate } from "sveltekit-superforms/server";

import type { PageServerLoad } from "./$types";
import type { Actions } from "./$types";
import schema from "./schema";

export const prerender = false;

export const load: PageServerLoad = async ({ params }) => {
	const docentes = await db
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

	if (docentes.length === 0) {
		throw error(404, { message: "Docente no encontrado." });
	}

	const cuatrimestres = await db.select().from(cuatrimestre);
	cuatrimestres.sort((a, b) => -sortCuatrimestres(a.nombre, b.nombre));

	const form = superValidate(schema);

	return { ...docentes[0], cuatrimestres, form };
};

export const actions: Actions = {
	default: async (e) => {
		const form = await superValidate(e, schema);

		if (form.valid && form.data.cuatrimestre && form.data.cuatrimestre != "undefined") {
			form.valid =
				(
					await db
						.select()
						.from(cuatrimestre)
						.where(eq(cuatrimestre.nombre, form.data.cuatrimestre))
				).length > 0;
		}

		if (!form.valid) {
			return fail(400, { form });
		}

    console.log(form);

		return message(form, {
			type: "success",
			text: "Calificación registrada con éxito"
		});
	}
};
