import { TURNSTILE_SECRET_KEY } from "$env/static/private";
import db from "$lib/db";
import {
	calificacion,
	catedra,
	catedraDocente,
	comentario,
	cuatrimestre,
	docente,
	materia
} from "$lib/db/schema";
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
	default: async ({ params, request }) => {
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
			).length === 1;

		if (form.data.cuatrimestre && form.data.cuatrimestre != "undefined" && !esCuatrimestreValido) {
			return setError(form, "cuatrimestre", `Cuatrimestre '${form.data.cuatrimestre}' no existe`);
		}

		const { success } = await validateToken(
			form.data["cf-turnstile-response"],
			TURNSTILE_SECRET_KEY
		);

		if (!success) {
			return setError(form, "Error al validar CAPTCHA");
		}

		try {
			if (form.data.comentario && form.data.comentario.length > 0 && form.data.cuatrimestre) {
				await db.insert(comentario).values({
					codigoDocente: params.codigoDocente,
					cuatrimestre: form.data.cuatrimestre,
					contenido: form.data.comentario
				});
			}

			await db.insert(calificacion).values({
				codigoDocente: params.codigoDocente,
				aceptaCritica: form.data["acepta-critica"],
				asistencia: form.data["asistencia"],
				buenTrato: form.data["buen-trato"],
				claridad: form.data["claridad"],
				claseOrganizada: form.data["clase-organizada"],
				cumpleHorarios: form.data["cumple-horario"],
				fomentaParticipacion: form.data["fomenta-participacion"],
				panoramaAmplio: form.data["panorama-amplio"],
				respondeMails: form.data["responde-mails"]
			});
		} catch {
			return fail(500);
		}

		return message(form, "Calificación registrada con éxito");
	}
};

interface TokenValidateResponse {
	"error-codes": string[];
	success: boolean;
	action: string;
	cdata: string;
}

async function validateToken(token: string, secret: string) {
	const res = await fetch("https://challenges.cloudflare.com/turnstile/v0/siteverify", {
		method: "POST",
		headers: {
			"content-type": "application/json"
		},
		body: JSON.stringify({
			response: token,
			secret: secret
		})
	});

	const data: TokenValidateResponse = await res.json();

	return {
		success: data.success,
		error: data["error-codes"]?.length ? data["error-codes"][0] : null
	};
}
