import db from "$lib/db";
import { sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
  return {
    docentes: streamDocentes(params.codigoCatedra),
  };
}) satisfies PageServerLoad;

async function streamDocentes(codigoCatedra: string) {
  let docentes: {
    codigo: string;
    nombre: string;
    calificaciones: {
      promedio_general: number;
      acepta_critica: number;
      asistencia: number;
      buen_trato: number;
      claridad: number;
      clase_organizada: number;
      cumple_horarios: number;
      fomenta_participacion: number;
      panorama_amplio: number;
      responde_mails: number;
    } | null;
    cantidad_calificaciones: number;
    resumen_comentarios: string | null;
    comentarios: {
      codigo: number;
      contenido: string;
      cuatrimestre: string;
    }[];
  }[];

  try {
    docentes = await db.execute(sql`
    SELECT * FROM informacion_comentarios_docentes_catedra(${codigoCatedra}::uuid);
  `);
  } catch (e: any) {
    if (e.code === "22P02") {
      throw error(404);
    } else {
      console.error(e);
      throw error(500);
    }
  }

  if (docentes.length == 0) {
    throw error(404);
  }

  return { docentes: docentes };
}
