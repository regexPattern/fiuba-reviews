import type { PageServerLoad } from "./$types";
import { db, schema } from "$lib/server/db";
import { and, count, desc, eq, exists, gt, lt, sql } from "drizzle-orm";

export const prerender = true;

// Parámetros para elegir comentarios suficientemente relevantes para mostrar en la landing page.
const nComentarios = 12;
const minCharsComentario = 100; // inclusivo
const maxCharsComentario = 200; // inclusivo

// Idem pero para el listado de materias (se toman las materias de planes vigentes con la mayor cantidad de comentarios).
const nMateriasPopulares = 10;

const cantidadPlanesVigentesExpr = count(sql`DISTINCT ${schema.plan.codigo}`).as(
  "cantidad_planes_vigentes"
);
const cantidadComentariosExpr = count(sql`DISTINCT ${schema.comentario.codigo}`).as(
  "cantidad_comentarios"
);

export const load: PageServerLoad = async () => {
  const comentariosUnificados = db
    .select({
      codigo: schema.comentario.codigo,
      contenido: schema.comentario.contenido,
      fechaCreacion: schema.comentario.fechaCreacion,
      nombreDocente: sql<string>`${schema.docente.nombre}`.as("nombre_docente"),
      codigoMateria: schema.docente.codigoMateria,
      nombreMateria: sql<string>`${schema.materia.nombre}`.as("nombre_materia"),
      ordenPorContenido: sql<number>`
        row_number() over (
          partition by ${schema.comentario.contenido}
          order by ${schema.comentario.fechaCreacion} desc nulls last, ${schema.comentario.codigo} desc
        )
      `.as("orden_por_contenido")
    })
    .from(schema.comentario)
    .innerJoin(schema.docente, eq(schema.docente.codigo, schema.comentario.codigoDocente))
    .innerJoin(schema.materia, eq(schema.materia.codigo, schema.docente.codigoMateria))
    .where(
      and(
        gt(sql`char_length(${schema.comentario.contenido})`, minCharsComentario - 1),
        lt(sql`char_length(${schema.comentario.contenido})`, maxCharsComentario + 1),
        exists(
          db
            .select({ uno: sql`1` })
            .from(schema.planMateria)
            .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
            .where(
              and(
                eq(schema.planMateria.codigoMateria, schema.docente.codigoMateria),
                eq(schema.plan.estaVigente, true)
              )
            )
        )
      )
    )
    .as("comentarios_con_rank");

  const comentariosRows = await db
    .select({
      codigo: comentariosUnificados.codigo,
      contenido: comentariosUnificados.contenido,
      nombreDocente: comentariosUnificados.nombreDocente,
      codigoMateria: comentariosUnificados.codigoMateria,
      nombreMateria: comentariosUnificados.nombreMateria
    })
    .from(comentariosUnificados)
    .where(eq(comentariosUnificados.ordenPorContenido, 1))
    .orderBy(desc(comentariosUnificados.fechaCreacion), desc(comentariosUnificados.codigo))
    .limit(nComentarios);

  const materiasPopularesRows = await db
    .select({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre,
      cantidadCatedras: count(sql`DISTINCT ${schema.catedra.codigo}`).as("cantidad_catedras"),
      cantidadDocentes: count(sql`DISTINCT ${schema.docente.codigo}`).as("cantidad_docentes"),
      cantidadPlanesVigentes: cantidadPlanesVigentesExpr,
      cantidadComentarios: cantidadComentariosExpr
    })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(
      schema.plan,
      and(eq(schema.plan.codigo, schema.planMateria.codigoPlan), eq(schema.plan.estaVigente, true))
    )
    .leftJoin(schema.catedra, eq(schema.catedra.codigoMateria, schema.materia.codigo))
    .leftJoin(schema.docente, eq(schema.docente.codigoMateria, schema.materia.codigo))
    .leftJoin(schema.comentario, eq(schema.comentario.codigoDocente, schema.docente.codigo))
    .groupBy(schema.materia.codigo, schema.materia.nombre)
    .orderBy(desc(cantidadPlanesVigentesExpr), desc(cantidadComentariosExpr), schema.materia.nombre)
    .limit(nMateriasPopulares);

  return {
    comentarios: comentariosRows,
    materiasPopulares: materiasPopularesRows
  };
};
