import "server-only";
import { desc, eq, sql } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { tags } from "@/server/cache";
import { db, schema } from "@/server/db";
import { getMateria } from "./materias";

/** Cátedras de una materia, ordenadas por calificación desc y nombre asc. */
export async function getCatedrasMateria(codigo: string) {
  "use cache";
  cacheLife("max");

  const materia = await getMateria(codigo);
  if (!materia) return [];

  // Una materia que usa equivalencias muestra las cátedras de esas otras materias.
  cacheTag(tags.materia(codigo), ...materia.codigosMateriasFuente.map(tags.catedras));

  const cm = schema.catedraMostrada;
  const cc = schema.calificacionCatedra;
  const calificacion = sql<number>`coalesce(${cc.calificacion}, 0)`;

  return db
    .select({
      codigo: cm.codigoCatedra,
      nombre: cm.nombre,
      calificacion: calificacion.as("calificacion")
    })
    .from(cm)
    .leftJoin(cc, eq(cc.codigoCatedra, cm.codigoCatedra))
    .where(eq(cm.codigoMateria, codigo))
    .orderBy(desc(calificacion), cm.nombre);
}

export type CatedraListada = Awaited<ReturnType<typeof getCatedrasMateria>>[number];

/** Docentes de una cátedra con sus promedios y comentarios, o `null` si no existe. */
export async function getCatedra(codigo: string) {
  "use cache";
  cacheLife("max");
  cacheTag(tags.catedra(codigo));

  const { catedraDocente: cd, docente: d, comentario: c, cuatrimestre: cu } = schema;
  const r = schema.resumenCalificacionesDocente;
  const pr = schema.prioridadRol;

  const [docentes, comentarios] = await Promise.all([
    db
      .select({
        codigo: d.codigo,
        nombre: d.nombre,
        nombreSiu: d.nombreSiu,
        rol: d.rol,
        resumenComentarios: d.resumenComentarios,
        cantidadCalificaciones: sql<number>`coalesce(${r.cantidad}, 0)`,
        // Drizzle devuelve `null` en vez del objeto si el docente no tiene calificaciones.
        promedios: {
          general: r.promedio,
          aceptaCritica: r.aceptaCritica,
          asistencia: r.asistencia,
          buenTrato: r.buenTrato,
          claridad: r.claridad,
          claseOrganizada: r.claseOrganizada,
          cumpleHorarios: r.cumpleHorarios,
          fomentaParticipacion: r.fomentaParticipacion,
          panoramaAmplio: r.panoramaAmplio,
          respondeMails: r.respondeMails
        }
      })
      .from(cd)
      .innerJoin(d, eq(d.codigo, cd.codigoDocente))
      .leftJoin(pr, eq(pr.rol, d.rol))
      .leftJoin(r, eq(r.codigoDocente, d.codigo))
      .where(eq(cd.codigoCatedra, codigo))
      .orderBy(sql`${pr.prioridad} asc nulls last`, d.nombre),
    db
      .select({
        codigo: c.codigo,
        codigoDocente: c.codigoDocente,
        contenido: c.contenido,
        deDolly: c.deDolly,
        cuatrimestre: { numero: cu.numero, anio: cu.anio }
      })
      .from(cd)
      .innerJoin(c, eq(c.codigoDocente, cd.codigoDocente))
      .innerJoin(cu, eq(cu.codigo, c.codigoCuatrimestre))
      .where(eq(cd.codigoCatedra, codigo))
      .orderBy(desc(c.codigoCuatrimestre), desc(c.codigo))
  ]);

  if (docentes.length === 0) return null;

  const comentariosPorDocente = Map.groupBy(comentarios, (comentario) => comentario.codigoDocente);

  return {
    codigo,
    docentes: docentes.map((docente) => ({
      ...docente,
      comentarios: comentariosPorDocente.get(docente.codigo) ?? []
    }))
  };
}

export type Catedra = NonNullable<Awaited<ReturnType<typeof getCatedra>>>;
export type DocenteCatedra = Catedra["docentes"][number];
export type PromediosDocente = NonNullable<DocenteCatedra["promedios"]>;
export type ComentarioDocente = DocenteCatedra["comentarios"][number];
