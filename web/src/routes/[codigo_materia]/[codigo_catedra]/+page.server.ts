import { db, schema } from "$lib/server/db";
import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";
import { and, desc, eq, inArray, sql } from "drizzle-orm";

export const load: PageServerLoad = async ({ params, parent }) => {
  const uuidV4Regex = /^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

  if (!uuidV4Regex.test(params.codigo_catedra)) {
    error(400, "Código de cátedra inválido.");
  }

  const { catedras: catedrasValidas } = await parent();
  const catedra = catedrasValidas.find((c) => c.codigo === params.codigo_catedra);

  if (!catedra) {
    error(404, "Cátedra no encontrada para la materia.");
  }

  const promedioDollyExpr = sql<number>`(
    ${schema.calificacionDolly.aceptaCritica} +
    ${schema.calificacionDolly.asistencia} +
    ${schema.calificacionDolly.buenTrato} +
    ${schema.calificacionDolly.claridad} +
    ${schema.calificacionDolly.claseOrganizada} +
    ${schema.calificacionDolly.cumpleHorarios} +
    ${schema.calificacionDolly.fomentaParticipacion} +
    ${schema.calificacionDolly.panoramaAmplio} +
    ${schema.calificacionDolly.respondeMails}
  ) / 9.0`;

  const calificacionesPorDocente = db
    .select({
      codigoDocente: schema.calificacionDolly.codigoDocente,
      cantidad: sql<number>`count(*)::int`.as("cantidad"),
      promedio: sql<number>`avg(${promedioDollyExpr})::double precision`.as("promedio"),
      aceptaCritica:
        sql<number>`avg(${schema.calificacionDolly.aceptaCritica})::double precision`.as(
          "acepta_critica"
        ),
      asistencia: sql<number>`avg(${schema.calificacionDolly.asistencia})::double precision`.as(
        "asistencia"
      ),
      buenTrato: sql<number>`avg(${schema.calificacionDolly.buenTrato})::double precision`.as(
        "buen_trato"
      ),
      claridad: sql<number>`avg(${schema.calificacionDolly.claridad})::double precision`.as(
        "claridad"
      ),
      claseOrganizada:
        sql<number>`avg(${schema.calificacionDolly.claseOrganizada})::double precision`.as(
          "clase_organizada"
        ),
      cumpleHorarios:
        sql<number>`avg(${schema.calificacionDolly.cumpleHorarios})::double precision`.as(
          "cumple_horarios"
        ),
      fomentaParticipacion:
        sql<number>`avg(${schema.calificacionDolly.fomentaParticipacion})::double precision`.as(
          "fomenta_participacion"
        ),
      panoramaAmplio:
        sql<number>`avg(${schema.calificacionDolly.panoramaAmplio})::double precision`.as(
          "panorama_amplio"
        ),
      respondeMails:
        sql<number>`avg(${schema.calificacionDolly.respondeMails})::double precision`.as(
          "responde_mails"
        )
    })
    .from(schema.calificacionDolly)
    .groupBy(schema.calificacionDolly.codigoDocente)
    .as("calificaciones_por_docente");

  const rowsDocentes = await db
    .select({
      codigo: schema.docente.codigo,
      nombre: schema.docente.nombre,
      rol: schema.docente.rol,
      resumenComentario: schema.docente.resumenComentarios,
      prioridadRol: schema.prioridadRol.prioridad,
      cantidadCalificaciones: calificacionesPorDocente.cantidad,
      promedio: calificacionesPorDocente.promedio,
      aceptaCritica: calificacionesPorDocente.aceptaCritica,
      asistencia: calificacionesPorDocente.asistencia,
      buenTrato: calificacionesPorDocente.buenTrato,
      claridad: calificacionesPorDocente.claridad,
      claseOrganizada: calificacionesPorDocente.claseOrganizada,
      cumpleHorarios: calificacionesPorDocente.cumpleHorarios,
      fomentaParticipacion: calificacionesPorDocente.fomentaParticipacion,
      panoramaAmplio: calificacionesPorDocente.panoramaAmplio,
      respondeMails: calificacionesPorDocente.respondeMails
    })
    .from(schema.catedra)
    .innerJoin(
      schema.catedraDocente,
      eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
    )
    .innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
    .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
    .leftJoin(
      calificacionesPorDocente,
      eq(calificacionesPorDocente.codigoDocente, schema.docente.codigo)
    )
    .where(eq(schema.catedra.codigo, params.codigo_catedra));

  const codigosDocente = rowsDocentes.map((d) => d.codigo);

  const comentariosDocentes = new Map<
    string,
    {
      codigo: number;
      contenido: string;
      cuatrimestre: { numero: number; anio: number };
      esDeDolly: boolean;
    }[]
  >();

  if (codigosDocente.length > 0) {
    const rowsComentarios = await db
      .select({
        codigoDocente: schema.comentario.codigoDocente,
        codigoComentario: schema.comentario.codigo,
        contenido: schema.comentario.contenido,
        esDeDolly: schema.comentario.esDeDolly,
        cuatrimestreNumero: schema.cuatrimestre.numero,
        cuatrimestreAnio: schema.cuatrimestre.anio,
        codigoCuatrimestre: schema.comentario.codigoCuatrimestre
      })
      .from(schema.comentario)
      .innerJoin(
        schema.cuatrimestre,
        eq(schema.cuatrimestre.codigo, schema.comentario.codigoCuatrimestre)
      )
      .where(inArray(schema.comentario.codigoDocente, codigosDocente))
      .orderBy(desc(schema.comentario.codigoCuatrimestre));

    for (const row of rowsComentarios) {
      const comentarios = comentariosDocentes.get(row.codigoDocente) ?? [];
      comentarios.push({
        codigo: row.codigoComentario,
        contenido: row.contenido,
        cuatrimestre: { numero: row.cuatrimestreNumero, anio: row.cuatrimestreAnio },
        esDeDolly: row.esDeDolly
      });
      comentariosDocentes.set(row.codigoDocente, comentarios);
    }
  }

  const docentes = rowsDocentes
    .map((row) => {
      const prioridad = row.prioridadRol ?? Number.MAX_SAFE_INTEGER;

      const promedioCalificaciones =
        row.promedio != null
          ? {
              general: row.promedio,
              aceptaCritica: row.aceptaCritica!,
              asistencia: row.asistencia!,
              buenTrato: row.buenTrato!,
              claridad: row.claridad!,
              claseOrganizada: row.claseOrganizada!,
              cumpleHorarios: row.cumpleHorarios!,
              fomentaParticipacion: row.fomentaParticipacion!,
              panoramaAmplio: row.panoramaAmplio!,
              respondeMails: row.respondeMails!
            }
          : null;

      return {
        nombre: row.nombre,
        codigo: row.codigo,
        rol: row.rol,
        cantidadCalificaciones: row.cantidadCalificaciones ?? 0,
        promedioCalificaciones,
        resumenComentario: row.resumenComentario ?? null,
        comentarios: comentariosDocentes.get(row.codigo) ?? [],
        prioridad
      };
    })
    .sort((a, b) => {
      if (a.prioridad !== b.prioridad) return a.prioridad - b.prioridad;
      return a.nombre.localeCompare(b.nombre);
    });

  return { docentes };
};
