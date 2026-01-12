import { db } from "$lib/server/db";
import * as schema from "$lib/server/db/schema";
import type { PageServerLoad } from "./$types";
import { and, desc, eq, sql } from "drizzle-orm";
import { redirect } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ params }) => {
  const promedioDocenteExpr = sql<number>`(
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

  const calificacionDocente = db.$with("calificacion_docente").as(
    db
      .select({
        codigoCatedra: schema.catedraDocente.codigoCatedra,
        codigoDocente: schema.docente.codigo,
        promedioDocente: sql<number>`avg(${promedioDocenteExpr})`.as("promedio_docente")
      })
      .from(schema.catedraDocente)
      .innerJoin(schema.docente, eq(schema.docente.codigo, schema.catedraDocente.codigoDocente))
      .innerJoin(
        schema.calificacionDolly,
        eq(schema.calificacionDolly.codigoDocente, schema.docente.codigo)
      )
      .where(eq(schema.docente.codigoMateria, params.codigo_materia))
      .groupBy(schema.catedraDocente.codigoCatedra, schema.docente.codigo)
  );

  const calificacionCatedraExpr = sql<number>`avg(${calificacionDocente.promedioDocente})`;

  const catedrasOrdenadas = await db
    .with(calificacionDocente)
    .select({
      codigo: schema.catedra.codigo,
      codigoMateria: schema.catedra.codigoMateria,
      calificacion: calificacionCatedraExpr.as("calificacion_catedra")
    })
    .from(schema.catedra)
    .innerJoin(calificacionDocente, eq(calificacionDocente.codigoCatedra, schema.catedra.codigo))
    .innerJoin(
      schema.planMateria,
      eq(schema.planMateria.codigoMateria, schema.catedra.codigoMateria)
    )
    .innerJoin(
      schema.plan,
      and(eq(schema.plan.codigo, schema.planMateria.codigoPlan), eq(schema.plan.estaVigente, true))
    )
    .where(eq(schema.catedra.codigoMateria, params.codigo_materia))
    .groupBy(schema.catedra.codigo, schema.catedra.codigoMateria)
    .orderBy(desc(calificacionCatedraExpr));

  console.log(catedrasOrdenadas);

  throw redirect(302, `/${params.codigo_materia}/${catedrasOrdenadas[0].codigo}`);
};
