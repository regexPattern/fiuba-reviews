import "server-only";
import { and, eq, inArray, sql } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";
import { promedioDollyExpr } from "@/server/db/sql";
import { getMateria } from "./materia";
import { calcularNombreCatedra, PRIORIDAD_SIN_ROL } from "./nombre-catedra";

/**
 * Listado de cátedras de una materia con su promedio, ordenado por promedio desc y nombre asc.
 *
 * Una materia vigente sin cátedras propias muestra las de sus equivalencias, así que se taguea
 * `catedras:<codigo>` por cada materia de la que realmente se leen cátedras.
 */
export async function getCatedrasMateria(codigo: string) {
  "use cache";
  cacheLife("max");

  const materia = await getMateria(codigo);

  if (!materia) {
    return [];
  }

  const codigosMateria = materia.usaEquivalencias
    ? materia.equivalencias.map((e) => e.codigo)
    : [materia.codigo];

  cacheTag(`materia:${codigo}`, ...codigosMateria.map((c) => `catedras:${c}`));

  if (codigosMateria.length === 0) {
    return [];
  }

  const promediosDocente = db
    .select({
      codigoDocente: schema.calificacionDolly.codigoDocente,
      promedio: sql<number>`avg(${promedioDollyExpr})::double precision`.as("promedio")
    })
    .from(schema.calificacionDolly)
    .groupBy(schema.calificacionDolly.codigoDocente)
    .as("promedios_docente");

  const rows = await db
    .select({
      codigoCatedra: schema.catedra.codigo,
      codigoMateria: schema.catedra.codigoMateria,
      nombreDocente: schema.docente.nombre,
      prioridadRol: schema.prioridadRol.prioridad,
      promedioDocente: promediosDocente.promedio
    })
    .from(schema.catedra)
    .innerJoin(
      schema.catedraDocente,
      eq(schema.catedraDocente.codigoCatedra, schema.catedra.codigo)
    )
    .innerJoin(schema.docente, eq(schema.catedraDocente.codigoDocente, schema.docente.codigo))
    .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
    .leftJoin(promediosDocente, eq(promediosDocente.codigoDocente, schema.docente.codigo))
    .where(
      materia.usaEquivalencias
        ? inArray(schema.catedra.codigoMateria, codigosMateria)
        : and(
            inArray(schema.catedra.codigoMateria, codigosMateria),
            eq(schema.catedra.activa, true)
          )
    );

  const grupos = new Map<
    string,
    {
      codigoMateria: string;
      docentes: { nombre: string; prioridad: number; promedio: number | null }[];
    }
  >();

  for (const row of rows) {
    const grupo = grupos.get(row.codigoCatedra) ?? {
      codigoMateria: row.codigoMateria,
      docentes: []
    };
    grupo.docentes.push({
      nombre: row.nombreDocente,
      prioridad: row.prioridadRol ?? PRIORIDAD_SIN_ROL,
      promedio: row.promedioDocente
    });
    grupos.set(row.codigoCatedra, grupo);
  }

  const catedras = [...grupos].map(([codigoCatedra, grupo]) => {
    const calificados = grupo.docentes.filter((d) => d.promedio !== null);
    const calificacion =
      calificados.length === 0
        ? 0
        : calificados.reduce((sum, d) => sum + d.promedio!, 0) / calificados.length;

    return {
      codigo: codigoCatedra,
      codigoMateria: grupo.codigoMateria,
      nombre: calcularNombreCatedra(grupo.docentes),
      calificacion
    };
  });

  return catedras.sort(
    (a, b) => b.calificacion - a.calificacion || a.nombre.localeCompare(b.nombre)
  );
}

export type CatedraListada = Awaited<ReturnType<typeof getCatedrasMateria>>[number];
