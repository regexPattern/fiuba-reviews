import "server-only";
import { and, count, desc, eq, exists, gt, lt, sql } from "drizzle-orm";
import { cacheLife, cacheTag } from "next/cache";
import { db, schema } from "@/server/db";

const N_COMENTARIOS = 12;
const MIN_CHARS_COMENTARIO = 100; // inclusivo
const MAX_CHARS_COMENTARIO = 200; // inclusivo
const N_MATERIAS_POPULARES = 10;

/** Comentarios destacados y materias populares de la pantalla de inicio. */
export async function getInicio() {
  "use cache";
  cacheLife("hours");
  cacheTag("inicio");

  const comentariosConRank = db
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
        gt(sql`char_length(${schema.comentario.contenido})`, MIN_CHARS_COMENTARIO - 1),
        lt(sql`char_length(${schema.comentario.contenido})`, MAX_CHARS_COMENTARIO + 1),
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

  const comentarios = await db
    .select({
      codigo: comentariosConRank.codigo,
      contenido: comentariosConRank.contenido,
      nombreDocente: comentariosConRank.nombreDocente,
      codigoMateria: comentariosConRank.codigoMateria,
      nombreMateria: comentariosConRank.nombreMateria
    })
    .from(comentariosConRank)
    .where(eq(comentariosConRank.ordenPorContenido, 1))
    .orderBy(desc(comentariosConRank.fechaCreacion), desc(comentariosConRank.codigo))
    .limit(N_COMENTARIOS);

  const cantidadPlanesVigentes = count(sql`DISTINCT ${schema.plan.codigo}`).as(
    "cantidad_planes_vigentes"
  );
  const cantidadComentarios = count(sql`DISTINCT ${schema.comentario.codigo}`).as(
    "cantidad_comentarios"
  );

  const materiasPopulares = await db
    .select({
      codigo: schema.materia.codigo,
      nombre: schema.materia.nombre,
      cantidadPlanesVigentes,
      cantidadComentarios
    })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(
      schema.plan,
      and(eq(schema.plan.codigo, schema.planMateria.codigoPlan), eq(schema.plan.estaVigente, true))
    )
    .leftJoin(schema.docente, eq(schema.docente.codigoMateria, schema.materia.codigo))
    .leftJoin(schema.comentario, eq(schema.comentario.codigoDocente, schema.docente.codigo))
    .groupBy(schema.materia.codigo, schema.materia.nombre)
    .orderBy(desc(cantidadPlanesVigentes), desc(cantidadComentarios), schema.materia.nombre)
    .limit(N_MATERIAS_POPULARES);

  return { comentarios, materiasPopulares };
}

export type Inicio = Awaited<ReturnType<typeof getInicio>>;
