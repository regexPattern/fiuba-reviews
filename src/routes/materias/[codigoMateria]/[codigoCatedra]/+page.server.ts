import db from "$lib/db";
import type { PageServerLoad } from "./$types";
import { sql } from "drizzle-orm";

export const load = (async ({ params }) => {
  const filasInfoDocentesComentarios = await db.execute<{
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
  }>(sql`
    SELECT * FROM informacion_comentarios_docentes_catedra(${params.codigoCatedra}::uuid);
  `);

  return { docentes: filasInfoDocentesComentarios };
}) satisfies PageServerLoad;
