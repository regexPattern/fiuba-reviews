import "server-only";
import { cacheLife, cacheTag } from "next/cache";
import type { SQLWrapper } from "drizzle-orm";
import { desc, eq, inArray, sql } from "drizzle-orm";
import { db, schema } from "@/server/db";
import { promedioDollyExpr } from "@/server/db/sql";
import { calcularNombreCatedra, PRIORIDAD_SIN_ROL } from "./nombre-catedra";

const avg = (col: SQLWrapper) => sql<number>`avg(${col})::double precision`;

/**
 * Detalle de una cátedra: docentes ordenados, promedios por criterio y comentarios.
 * Devuelve `null` si la cátedra no existe.
 */
export async function getCatedra(codigo: string) {
  "use cache";
  cacheLife("max");
  cacheTag(`catedra:${codigo}`);

  const c = schema.calificacionDolly;

  const calificacionesPorDocente = db
    .select({
      codigoDocente: c.codigoDocente,
      cantidad: sql<number>`count(*)::int`.as("cantidad"),
      promedio: avg(promedioDollyExpr).as("promedio"),
      aceptaCritica: avg(c.aceptaCritica).as("acepta_critica"),
      asistencia: avg(c.asistencia).as("asistencia"),
      buenTrato: avg(c.buenTrato).as("buen_trato"),
      claridad: avg(c.claridad).as("claridad"),
      claseOrganizada: avg(c.claseOrganizada).as("clase_organizada"),
      cumpleHorarios: avg(c.cumpleHorarios).as("cumple_horarios"),
      fomentaParticipacion: avg(c.fomentaParticipacion).as("fomenta_participacion"),
      panoramaAmplio: avg(c.panoramaAmplio).as("panorama_amplio"),
      respondeMails: avg(c.respondeMails).as("responde_mails")
    })
    .from(c)
    .groupBy(c.codigoDocente)
    .as("calificaciones_por_docente");

  const cal = calificacionesPorDocente;

  const rowsDocentes = await db
    .select({
      codigo: schema.docente.codigo,
      nombre: schema.docente.nombre,
      rol: schema.docente.rol,
      resumenComentario: schema.docente.resumenComentarios,
      prioridadRol: schema.prioridadRol.prioridad,
      cantidadCalificaciones: cal.cantidad,
      promedio: cal.promedio,
      aceptaCritica: cal.aceptaCritica,
      asistencia: cal.asistencia,
      buenTrato: cal.buenTrato,
      claridad: cal.claridad,
      claseOrganizada: cal.claseOrganizada,
      cumpleHorarios: cal.cumpleHorarios,
      fomentaParticipacion: cal.fomentaParticipacion,
      panoramaAmplio: cal.panoramaAmplio,
      respondeMails: cal.respondeMails
    })
    .from(schema.catedra)
    .innerJoin(
      schema.catedraDocente,
      eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
    )
    .innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
    .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
    .leftJoin(cal, eq(cal.codigoDocente, schema.docente.codigo))
    .where(eq(schema.catedra.codigo, codigo));

  if (rowsDocentes.length === 0) {
    return null;
  }

  const comentariosPorDocente = new Map<
    string,
    {
      codigo: number;
      contenido: string;
      cuatrimestre: { numero: number; anio: number };
      esDeDolly: boolean;
    }[]
  >();

  const rowsComentarios = await db
    .select({
      codigoDocente: schema.comentario.codigoDocente,
      codigo: schema.comentario.codigo,
      contenido: schema.comentario.contenido,
      esDeDolly: schema.comentario.esDeDolly,
      cuatrimestreNumero: schema.cuatrimestre.numero,
      cuatrimestreAnio: schema.cuatrimestre.anio
    })
    .from(schema.comentario)
    .innerJoin(
      schema.cuatrimestre,
      eq(schema.cuatrimestre.codigo, schema.comentario.codigoCuatrimestre)
    )
    .where(
      inArray(
        schema.comentario.codigoDocente,
        rowsDocentes.map((d) => d.codigo)
      )
    )
    .orderBy(desc(schema.comentario.codigoCuatrimestre), desc(schema.comentario.codigo));

  for (const row of rowsComentarios) {
    const lista = comentariosPorDocente.get(row.codigoDocente) ?? [];
    lista.push({
      codigo: row.codigo,
      contenido: row.contenido,
      cuatrimestre: { numero: row.cuatrimestreNumero, anio: row.cuatrimestreAnio },
      esDeDolly: row.esDeDolly
    });
    comentariosPorDocente.set(row.codigoDocente, lista);
  }

  const docentes = rowsDocentes
    .map((row) => ({
      codigo: row.codigo,
      nombre: row.nombre,
      rol: row.rol,
      prioridad: row.prioridadRol ?? PRIORIDAD_SIN_ROL,
      resumenComentario: row.resumenComentario,
      cantidadCalificaciones: row.cantidadCalificaciones ?? 0,
      promedios:
        row.promedio !== null
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
          : null,
      comentarios: comentariosPorDocente.get(row.codigo) ?? []
    }))
    .sort((a, b) => a.prioridad - b.prioridad || a.nombre.localeCompare(b.nombre));

  return { codigo, nombre: calcularNombreCatedra(docentes), docentes };
}

export type Catedra = NonNullable<Awaited<ReturnType<typeof getCatedra>>>;
export type DocenteCatedra = Catedra["docentes"][number];
export type PromediosDocente = NonNullable<DocenteCatedra["promedios"]>;
export type ComentarioDocente = DocenteCatedra["comentarios"][number];
