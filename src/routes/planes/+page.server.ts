import type { Actions, PageServerLoad } from "./$types";

import { BACKEND_URL, TURNSTILE_SECRET_KEY } from "$env/static/private";

import { error } from "@sveltejs/kit";
import { formPlanSiu as schema } from "$lib/zod/schema";
import { message, setError, superValidate } from "sveltekit-superforms/server";
import { validarToken } from "$lib/utils";

export const prerender = false;

const carreras = new Set([
  "Ingeniería Civil",
  "Ingeniería Electrónica",
  "Ingeniería Industrial",
  "Ingeniería Mecánica",
  "Ingeniería Naval",
  "Ingeniería Química",
  "Ingeniería de Alimentos",
  "Ingeniería en Agrimensura",
  "Ingeniería en Energía Eléctrica",
  "Ingeniería en Petróleo",
  "Ingeniería en Informática",
]);

function diferencia(a: Set<string>, b: Set<string>) {
  const res = new Set();
  for (const element of a) {
    if (!b.has(element)) {
      res.add(element);
    }
  }
  return res;
}

export const load: PageServerLoad = async () => {
  const planesRegistradosRes = await fetch(`${BACKEND_URL}/planes`);
  const planesRegistrados: {
    carrera: string;
    cuatri: { numero: number; anio: number };
  }[] = await planesRegistradosRes.json();

  const carrerasListas = planesRegistrados.map((p) => p.carrera);
  const carrerasFaltantes = diferencia(carreras, new Set(carrerasListas));

  return {
    planesRegistrados,
    carrerasFaltantes,
    form: await superValidate(schema),
  };
};

export const actions: Actions = {
  default: async ({ request }) => {
    const form = await superValidate(request, schema);

    if (!form.valid) {
      return message(form, "Datos inválidos");
    }

    const { esValido } = await validarToken(
      form.data["cf-turnstile-response"],
      TURNSTILE_SECRET_KEY,
    );

    if (!esValido) {
      return setError(form, "Error al validar CAPTCHA.");
    }

    try {
      await fetch(`${BACKEND_URL}/planes`, {
        method: "POST",
        headers: {
          "Content-Type": " text/plain; charset=UTF-8",
        },
        body: form.data["contenido-siu"],
      });
    } catch (e) {
      console.error(e);
      throw error(500, "Error interno al guardar el contenido del SIU.");
    }

    return message(form, "Contenido del SIU registrado con éxito.");
  },
};
