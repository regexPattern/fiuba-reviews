import "server-only";
import { desc, eq } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";
import { calcularNombreCatedra, PRIORIDAD_SIN_ROL } from "./nombre-catedra";

/**
 * Detalle de una cátedra: docentes ordenados, promedios por criterio y comentarios.
 * Devuelve `null` si la cátedra no existe.
 */
export async function getCatedra(codigo: string) {
  "use cache";
  cacheLife("max");
  cacheTag(`catedra:${codigo}`);

  const cal = schema.resumenCalificacionesDocente;

  const cd = schema.catedraDocente;

  // Ambas queries parten de `catedra_docente`, así que no dependen una de otra y van en paralelo.
  const [rowsDocentes, rowsComentarios] = await Promise.all([
    db
      .select({
        codigo: schema.docente.codigo,
        nombre: schema.docente.nombre,
        nombreSiu: schema.docente.nombreSiu,
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
      .from(cd)
      .innerJoin(schema.docente, eq(schema.docente.codigo, cd.codigoDocente))
      .leftJoin(schema.prioridadRol, eq(schema.prioridadRol.rol, schema.docente.rol))
      .leftJoin(cal, eq(cal.codigoDocente, schema.docente.codigo))
      .where(eq(cd.codigoCatedra, codigo)),
    db
      .select({
        codigoDocente: schema.comentario.codigoDocente,
        codigo: schema.comentario.codigo,
        contenido: schema.comentario.contenido,
        deDolly: schema.comentario.deDolly,
        cuatrimestreNumero: schema.cuatrimestre.numero,
        cuatrimestreAnio: schema.cuatrimestre.anio
      })
      .from(cd)
      .innerJoin(schema.comentario, eq(schema.comentario.codigoDocente, cd.codigoDocente))
      .innerJoin(
        schema.cuatrimestre,
        eq(schema.cuatrimestre.codigo, schema.comentario.codigoCuatrimestre)
      )
      .where(eq(cd.codigoCatedra, codigo))
      .orderBy(desc(schema.comentario.codigoCuatrimestre), desc(schema.comentario.codigo))
  ]);

  if (rowsDocentes.length === 0) {
    return null;
  }

  const comentariosPorDocente = new Map<
    string,
    {
      codigo: number;
      contenido: string;
      cuatrimestre: { numero: number; anio: number };
      deDolly: boolean;
    }[]
  >();

  for (const row of rowsComentarios) {
    const lista = comentariosPorDocente.get(row.codigoDocente) ?? [];
    lista.push({
      codigo: row.codigo,
      contenido: row.contenido,
      cuatrimestre: { numero: row.cuatrimestreNumero, anio: row.cuatrimestreAnio },
      deDolly: row.deDolly
    });
    comentariosPorDocente.set(row.codigoDocente, lista);
  }

  const docentes = rowsDocentes
    .map((row) => ({
      codigo: row.codigo,
      nombre: row.nombre,
      nombreSiu: row.nombreSiu,
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
