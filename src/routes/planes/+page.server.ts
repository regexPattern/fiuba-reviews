import { BACKEND_URL, TURNSTILE_SECRET_KEY } from "$env/static/private";
import { contenidoSiu as schema } from "$lib/zod/schema";
import { fail } from "@sveltejs/kit";
import type { Actions, PageServerLoad } from "./$types";
import { message, setError, superValidate } from "sveltekit-superforms/server";

type PlanRegistrado = {
  carrera: string;
  cuatri: { numero: number; anio: number };
};

const CARRERAS = new Set([
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

export const load: PageServerLoad = async () => {
  const planesRegistradosRes = await fetch(`${BACKEND_URL}/planes`);
  const planesRegistrados: PlanRegistrado[] = await planesRegistradosRes.json();

  const carrerasListas = planesRegistrados.map((p) => p.carrera);
  const carrerasFaltantes = CARRERAS.difference(new Set(carrerasListas));

  return {
    form: await superValidate(schema),
    planesRegistrados,
    carrerasFaltantes,
  };
};

export const actions: Actions = {
  default: async ({ request }) => {
    const form = await superValidate(request, schema);

    if (!form.valid) {
      console.log(form.errors);
      return message(form, "Datos inválidos");
    }

    const { success } = await validateToken(
      form.data["cf-turnstile-response"],
      TURNSTILE_SECRET_KEY
    );

    if (!success) {
      return setError(form, "Error al validar CAPTCHA");
    }

    try {
      await fetch(`${BACKEND_URL}/planes`, {
        method: "POST",
        headers: {
          "Content-Type": " text/plain; charset=UTF-8",
        },
        body: form.data["contenido-siu"],
      });
    } catch {
      return fail(500);
    }

    return message(form, "Contenido del SIU registrado con éxito");
  },
};

interface TokenValidateResponse {
  "error-codes": string[];
  success: boolean;
  action: string;
  cdata: string;
}

async function validateToken(token: string, secret: string) {
  const res = await fetch(
    "https://challenges.cloudflare.com/turnstile/v0/siteverify",
    {
      method: "POST",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({
        response: token,
        secret: secret,
      }),
    }
  );

  const data: TokenValidateResponse = await res.json();

  return {
    success: data.success,
    error: data["error-codes"]?.length ? data["error-codes"][0] : null,
  };
}
