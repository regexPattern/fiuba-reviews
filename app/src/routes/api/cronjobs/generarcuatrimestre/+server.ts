import { CRON_SECRET } from "$env/static/private";
import db from "$lib/db";
import { cuatrimestre } from "$lib/db/schema";
import type { RequestHandler } from "@sveltejs/kit";

interface Cuatrimestre {
  anio: number;
  numero: number;
}

function calcularCuatrimestreAnterior(fechaActual: Date): Cuatrimestre {
  const mesActual = fechaActual.getUTCMonth() + 1;
  // 1er Cuatri Diciembre-Junio, 2do Cuatri Julio-Noviembre.
  const cuatrimestreActual = mesActual <= 6 || mesActual == 12 ? 1 : 2;
  const cuatrimestrePasado = {
    anio: fechaActual.getUTCFullYear(),
    numero: cuatrimestreActual - 1,
  };

  if (cuatrimestrePasado.numero === 0) {
    cuatrimestrePasado.anio--;
    cuatrimestrePasado.numero = 2;
  }

  return cuatrimestrePasado;
}

export const GET: RequestHandler = async ({ request }) => {
  const authHeader = request.headers.get("Authorization");
  if (authHeader !== `Bearer ${CRON_SECRET}`) {
    return new Response("Unauthorized", {
      status: 401,
    });
  }

  const fechaActual = new Date();
  const cuatrimestreAnterior = calcularCuatrimestreAnterior(fechaActual);
  let registro;
  try {
    registro = await db
      .insert(cuatrimestre)
      .values(cuatrimestreAnterior)
      .returning();
  } catch (err: any) {
    const UNIQUE_VIOLATION = "23505";
    if (err?.code == UNIQUE_VIOLATION)
      return new Response("Cuatrimestre duplicado", { status: 400 });
    else throw err;
  }

  return new Response(JSON.stringify(registro), {
    status: 201,
  });
};
