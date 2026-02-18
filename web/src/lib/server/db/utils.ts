import { and, eq, inArray, sql } from "drizzle-orm";
import { db, schema } from "./index";

export async function obtenerCatedras(
  codigosMateria: string[],
  opciones?: { soloActivas?: boolean }
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

  const catedrasDocentes = await db
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
      opciones?.soloActivas
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

  for (const row of catedrasDocentes) {
    if (!catedrasDocentesMap.has(row.codigoCatedra)) {
      catedrasDocentesMap.set(row.codigoCatedra, {
        codigoMateria: row.codigoMateria,
        docentes: []
      });
    }

    catedrasDocentesMap
      .get(row.codigoCatedra)!
      .docentes.push({
        nombre: row.nombreDocente,
        promedio: row.promedioDocente ?? 0,
        tieneCalificacion: row.promedioDocente != null,
        prioridad: row.prioridadRol ?? Number.MAX_SAFE_INTEGER
      });
  }

  const catedras = [];

  for (const [codigoCatedra, grupo] of catedrasDocentesMap) {
    const docentesOrdenados = [...grupo.docentes].sort((a, b) => {
      if (a.prioridad !== b.prioridad) return a.prioridad - b.prioridad;
      return a.nombre.localeCompare(b.nombre);
    });

    const nombreCatedra = docentesOrdenados.map((d) => d.nombre).join("-");

    const docentesConCalif = docentesOrdenados.filter((d) => d.tieneCalificacion);
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
