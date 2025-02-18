import db from "$lib/db";
import * as s from "$lib/db/schema";
import {
  and,
  count,
  countDistinct,
  desc,
  eq,
  exists,
  gt,
  lt,
  sql,
} from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async () => {
  const [
    materiasMasPopulares,
    ultimosComentarios,
    [{ cantidad: cantidadCatedras }],
    [{ cantidad: cantidadComentarios }],
    [{ cantidad: cantidadCalificaciones }],
  ] = await Promise.all([
    db
      .select({
        codigo: s.materia.codigo,
        nombre: s.materia.nombre,
        cantidadPlanesVigentes: countDistinct(s.plan.codigo),
      })
      .from(s.materia)
      .innerJoin(
        s.planMateria,
        eq(s.materia.codigo, s.planMateria.codigoMateria),
      )
      .innerJoin(s.plan, eq(s.planMateria.codigoPlan, s.plan.codigo))
      .where(eq(s.plan.estaVigente, true))
      .groupBy(s.materia.codigo, s.materia.nombre)
      .orderBy(desc(countDistinct(s.plan.codigo)))
      .limit(30),
    db
      .select({
        codigo: s.comentario.codigo,
        contenido: s.comentario.contenido,
        nombreDocente: s.docente.nombre,
      })
      .from(s.comentario)
      .innerJoin(s.docente, eq(s.comentario.codigoDocente, s.docente.codigo))
      .where(
        and(
          eq(
            s.comentario.codigoCuatrimestre,
            db
              .select({ value: sql<number>`max(${s.cuatrimestre.codigo}) - 1` })
              .from(s.cuatrimestre),
          ),
          gt(sql`length(${s.comentario.contenido})`, 100),
          lt(sql`length(${s.comentario.contenido})`, 200),
        ),
      )
      .limit(4),
    db
      .select({ cantidad: countDistinct(s.catedra.codigo) })
      .from(s.catedra)
      .where(
        exists(
          db
            .select({ one: sql`1` })
            .from(s.catedraDocente)
            .innerJoin(
              s.docente,
              eq(s.catedraDocente.codigoDocente, s.docente.codigo),
            )
            .innerJoin(
              s.comentario,
              eq(s.docente.codigo, s.comentario.codigoDocente),
            )
            .where(eq(s.catedraDocente.codigoCatedra, s.catedra.codigo))
            .$dynamic(),
        ),
      ),
    db.select({ cantidad: count() }).from(s.comentario),
    db.select({ cantidad: count() }).from(s.calificacionDolly),
  ]);

  return {
    materiasMasPopulares,
    ultimosComentarios,
    cantidadCatedras,
    cantidadComentarios,
    cantidadCalificaciones,
  };
};
