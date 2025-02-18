import db from "$lib/db";
import { sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";
import { error } from "@sveltejs/kit";

export const load = (async ({ params }) => {
  let filasInfoDocentesComentarios: {
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
    filasInfoDocentesComentarios = await db.execute(sql`
    SELECT * FROM informacion_comentarios_docentes_catedra(${params.codigoCatedra}::uuid);
  `);
  } catch (e: any) {
    if (e.code === "22P02") {
      throw error(404);
    } else {
      console.error(e);
      throw error(500);
    }
  }

  if (filasInfoDocentesComentarios.length == 0) {
    throw error(404);
  }

  return { docentes: filasInfoDocentesComentarios };
}) satisfies PageServerLoad;
