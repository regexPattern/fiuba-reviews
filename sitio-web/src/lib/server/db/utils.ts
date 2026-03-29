import { and, desc, eq, inArray, sql } from "drizzle-orm";
import { db, schema } from "./index";

const cantidadPlanesExpr = sql<number>`count(distinct ${schema.plan.codigo})::int`;
const cantidadCatedrasExpr = sql<number>`(
  CASE
    WHEN ${schema.materia.cuatrimestreUltimaActualizacion} IS NOT NULL
      THEN (
        SELECT COUNT(*)
        FROM ${schema.catedra}
        WHERE ${schema.catedra.codigoMateria} = ${schema.materia.codigo}
          AND ${schema.catedra.activa} = true
      )
    ELSE (
      SELECT COUNT(*)
      FROM ${schema.catedra}
      WHERE ${schema.catedra.codigoMateria} IN (
        SELECT ${schema.equivalencia.codigoMateriaPlanAnterior}
        FROM ${schema.equivalencia}
        WHERE ${schema.equivalencia.codigoMateriaPlanVigente} = ${schema.materia.codigo}
      )
    )
  END
)::int`;

export const getMateriasDisponibles = async () => {
  const materiasRows = await db
    .select({ codigo: schema.materia.codigo, nombre: schema.materia.nombre })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .where(eq(schema.plan.estaVigente, true))
    .groupBy(schema.materia.codigo, schema.materia.nombre)
    .orderBy(desc(cantidadPlanesExpr), desc(cantidadCatedrasExpr), schema.materia.nombre);

  return materiasRows;
};

type DocenteParaNombreCatedra = { nombre: string; prioridad: number };

export function calcularNombreCatedra(docentes: DocenteParaNombreCatedra[]) {
  const docentesOrdenados = [...docentes].sort((a, b) => {
    if (a.prioridad !== b.prioridad) return a.prioridad - b.prioridad;
    return a.nombre.localeCompare(b.nombre);
  });

  return docentesOrdenados.map((docente) => docente.nombre).join("-");
}

export async function obtenerCatedrasMaterias(
  codigosMateria: string[],
  soloActivas: boolean = false
) {
  if (codigosMateria.length === 0) {
    return [];
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

  const promediosDocenteSubquery = db
    .select({
      codigoDocente: schema.calificacionDolly.codigoDocente,
      promedio: sql<number>`avg(${promedioDollyExpr})::double precision`.as("promedio")
    })
    .from(schema.calificacionDolly)
    .groupBy(schema.calificacionDolly.codigoDocente)
    .as("promedios_docente");

  const catedraDocenteRows = await db
    .select({
      codigoCatedra: schema.catedra.codigo,
      codigoMateria: schema.catedra.codigoMateria,
      codigoDocente: schema.docente.codigo,
      nombreDocente: schema.docente.nombre,
      rolDocente: schema.docente.rol,
      prioridadRol: schema.prioridadRol.prioridad,
      promedioDocente: promediosDocenteSubquery.promedio
    })
    .from(schema.catedra)
    .innerJoin(
      schema.catedraDocente,
      eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
    )
    .innerJoin(schema.docente, eq(schema.catedraDocente.codigoDocente, schema.docente.codigo))
    .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
    .leftJoin(
      promediosDocenteSubquery,
      eq(promediosDocenteSubquery.codigoDocente, schema.docente.codigo)
    )
    .where(
      soloActivas
        ? and(
            inArray(schema.catedra.codigoMateria, codigosMateria),
            eq(schema.catedra.activa, true)
          )
        : inArray(schema.catedra.codigoMateria, codigosMateria)
    );

  const catedrasDocentesMap = new Map<
    string,
    {
      codigoMateria: string;
      docentes: {
        nombre: string;
        promedio: number;
        tieneCalificacion: boolean;
        prioridad: number;
      }[];
    }
  >();

  for (const row of catedraDocenteRows) {
    if (!catedrasDocentesMap.has(row.codigoCatedra)) {
      catedrasDocentesMap.set(row.codigoCatedra, {
        codigoMateria: row.codigoMateria,
        docentes: []
      });
    }

    catedrasDocentesMap.get(row.codigoCatedra)!.docentes.push({
      nombre: row.nombreDocente,
      promedio: row.promedioDocente ?? 0,
      tieneCalificacion: row.promedioDocente != null,
      prioridad: row.prioridadRol ?? Number.MAX_SAFE_INTEGER
    });
  }

  const catedras = [];

  for (const [codigoCatedra, grupo] of catedrasDocentesMap) {
    const nombreCatedra = calcularNombreCatedra(grupo.docentes);
    const docentesConCalif = grupo.docentes.filter((d) => d.tieneCalificacion);
    const calificacion =
      docentesConCalif.length === 0
        ? 0
        : docentesConCalif.reduce((sum, d) => sum + d.promedio, 0) / docentesConCalif.length;

    catedras.push({
      codigo: codigoCatedra,
      codigoMateria: grupo.codigoMateria,
      nombre: nombreCatedra,
      calificacion
    });
  }

  catedras.sort((a, b) => {
    if (a.calificacion !== b.calificacion) {
      return b.calificacion - a.calificacion;
    } else {
      return a.nombre.localeCompare(b.nombre);
    }
  });

  return catedras;
}

type PromedioCalificacionesDocente = {
  general: number;
  aceptaCritica: number;
  asistencia: number;
  buenTrato: number;
  claridad: number;
  claseOrganizada: number;
  cumpleHorarios: number;
  fomentaParticipacion: number;
  panoramaAmplio: number;
  respondeMails: number;
};

type ComentarioDocente = {
  codigo: number;
  contenido: string;
  cuatrimestre: { numero: number; anio: number };
  esDeDolly: boolean;
};

export type DocenteDetalleCatedra = {
  nombre: string;
  codigo: string;
  rol: string | null;
  cantidadCalificaciones: number;
  promedioCalificaciones: PromedioCalificacionesDocente | null;
  resumenComentario: string | null;
  comentarios: ComentarioDocente[];
};

export async function obtenerDetallesCatedras(codigosCatedra: string[]) {
  if (codigosCatedra.length === 0) {
    return {} satisfies Record<string, DocenteDetalleCatedra[]>;
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

  const calificacionesDocentesRows = db
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

  const docentesRows = await db
    .select({
      codigoCatedra: schema.catedra.codigo,
      codigoDocente: schema.docente.codigo,
      nombre: schema.docente.nombre,
      rol: schema.docente.rol,
      resumenComentario: schema.docente.resumenComentarios,
      prioridadRol: schema.prioridadRol.prioridad,
      cantidadCalificaciones: calificacionesDocentesRows.cantidad,
      promedio: calificacionesDocentesRows.promedio,
      aceptaCritica: calificacionesDocentesRows.aceptaCritica,
      asistencia: calificacionesDocentesRows.asistencia,
      buenTrato: calificacionesDocentesRows.buenTrato,
      claridad: calificacionesDocentesRows.claridad,
      claseOrganizada: calificacionesDocentesRows.claseOrganizada,
      cumpleHorarios: calificacionesDocentesRows.cumpleHorarios,
      fomentaParticipacion: calificacionesDocentesRows.fomentaParticipacion,
      panoramaAmplio: calificacionesDocentesRows.panoramaAmplio,
      respondeMails: calificacionesDocentesRows.respondeMails
    })
    .from(schema.catedra)
    .innerJoin(
      schema.catedraDocente,
      eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
    )
    .innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
    .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
    .leftJoin(
      calificacionesDocentesRows,
      eq(calificacionesDocentesRows.codigoDocente, schema.docente.codigo)
    )
    .where(inArray(schema.catedra.codigo, codigosCatedra));

  const codigosDocente = [...new Set(docentesRows.map((d) => d.codigoDocente))];

  const comentariosDocentes = new Map<string, ComentarioDocente[]>();

  if (codigosDocente.length > 0) {
    const comentariosRows = await db
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

    for (const row of comentariosRows) {
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

  const docentesPorCatedra: Record<string, (DocenteDetalleCatedra & { prioridad: number })[]> = {};

  for (const row of docentesRows) {
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

    const docente = {
      nombre: row.nombre,
      codigo: row.codigoDocente,
      rol: row.rol,
      cantidadCalificaciones: row.cantidadCalificaciones ?? 0,
      promedioCalificaciones,
      resumenComentario: row.resumenComentario ?? null,
      comentarios: comentariosDocentes.get(row.codigoDocente) ?? [],
      prioridad: row.prioridadRol ?? Number.MAX_SAFE_INTEGER
    };

    if (!docentesPorCatedra[row.codigoCatedra]) {
      docentesPorCatedra[row.codigoCatedra] = [];
    }
    docentesPorCatedra[row.codigoCatedra].push(docente);
  }

  const resultado: Record<string, DocenteDetalleCatedra[]> = {};

  for (const [codigoCatedra, docentes] of Object.entries(docentesPorCatedra)) {
    resultado[codigoCatedra] = docentes
      .sort((a, b) => {
        if (a.prioridad !== b.prioridad) return a.prioridad - b.prioridad;
        return a.nombre.localeCompare(b.nombre);
      })
      .map(({ prioridad, ...docente }) => docente);
  }

  return resultado;
}
