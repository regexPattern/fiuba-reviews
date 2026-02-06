import { db, schema } from "$lib/server/db";
import type { PageServerLoad } from "./$types";
import { and, desc, eq, exists, gt, lt, sql } from "drizzle-orm";

const N_COMENTARIOS = 20;
const MIN_CHARS_COMENTARIO = 100; // inclusivo
const MAX_CHARS_COMENTARIO = 200; // inclusivo

export const load: PageServerLoad = async () => {
  const comentariosUnificados = db
    .select({
      codigo: schema.comentario.codigo,
      contenido: schema.comentario.contenido,
      fechaCreacion: schema.comentario.fechaCreacion,
      ordenPorContenido: sql<number>`
        row_number() over (
          partition by ${schema.comentario.contenido}
          order by ${schema.comentario.fechaCreacion} desc nulls last, ${schema.comentario.codigo} desc
        )
      `.as("orden_por_contenido")
    })
    .from(schema.comentario)
    .innerJoin(schema.docente, eq(schema.docente.codigo, schema.comentario.codigoDocente))
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
      codigo: comentariosUnificados.codigo,
      contenido: comentariosUnificados.contenido
    })
    .from(comentariosUnificados)
    .where(eq(comentariosUnificados.ordenPorContenido, 1))
    .orderBy(desc(comentariosUnificados.fechaCreacion), desc(comentariosUnificados.codigo))
    .limit(N_COMENTARIOS);

  return { comentarios };
};
