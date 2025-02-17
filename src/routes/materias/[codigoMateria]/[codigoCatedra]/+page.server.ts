import db from "$lib/db";
import {
  calificacionDolly,
  catedraDocente,
  comentario,
  cuatrimestre,
  docente,
} from "$lib/db/schema";
import type { PageServerLoad } from "./$types";
import { eq, inArray, sql } from "drizzle-orm";

export const load = (async ({ params }) => {
  const filasDocentes = await db.execute<{
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
  }>(sql`
      SELECT
        d.codigo,
        d.nombre,
        CASE
          WHEN pd.promedio_general IS NOT NULL THEN
            JSON_BUILD_OBJECT(
              'promedio_general', pd.promedio_general,
              'resumen_comentarios', pd.resumen_comentarios,
              'acepta_critica', pd.acepta_critica,
              'asistencia', pd.asistencia,
              'buen_trato', pd.buen_trato,
              'claridad', pd.claridad,
              'clase_organizada', pd.clase_organizada,
              'cumple_horarios', pd.cumple_horarios,
              'fomenta_participacion', pd.fomenta_participacion,
              'panorama_amplio', pd.panorama_amplio,
              'responde_mails', pd.responde_mails
            )
          ELSE
            NULL
        END AS calificaciones,
        d.resumen_comentarios,
        COUNT(cdolly.codigo) AS cantidad_calificaciones
      FROM
        ${docente} d
      INNER JOIN
        ${catedraDocente} cd ON d.codigo = cd.codigo_docente
      INNER JOIN LATERAL
        promedio_docente_cal_dolly(d.codigo) pd ON TRUE
      LEFT JOIN
        ${calificacionDolly} cdolly ON d.codigo = cdolly.codigo_docente
      WHERE
        cd.codigo_catedra = ${params.codigoCatedra}
      GROUP BY
        d.codigo, d.nombre, pd.promedio_general, pd.resumen_comentarios,
        pd.acepta_critica, pd.asistencia, pd.buen_trato, pd.claridad,
        pd.clase_organizada, pd.cumple_horarios, pd.fomenta_participacion,
        pd.panorama_amplio, pd.responde_mails
      ORDER BY
        d.nombre ASC;
    `);

  const filasComentarios = await db
    .select({
      codigo: comentario.codigo,
      codigo_docente: comentario.codigoDocente,
      contenido: comentario.contenido,
      cuatrimestre: sql<string>`${cuatrimestre.numero} + ${cuatrimestre.anio}`,
    })
    .from(comentario)
    .innerJoin(
      cuatrimestre,
      eq(comentario.codigoCuatrimestre, cuatrimestre.codigo)
    )
    .where(
      inArray(
        comentario.codigoDocente,
        filasDocentes.map((d) => d.codigo)
      )
    );

  const docentesAComentarios = new Map<string, typeof filasComentarios>();

  for (const com of filasComentarios) {
    const comentarios = docentesAComentarios.get(com.codigo_docente) || [];
    comentarios.push(com);
    docentesAComentarios.set(com.codigo_docente, comentarios);
  }

  return {
    docentes: filasDocentes,
    comentarios: docentesAComentarios,
  };
}) satisfies PageServerLoad;
