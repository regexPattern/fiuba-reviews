import "@valibot/i18n/es";
import type { LayoutServerLoad } from "./$types";
import { desc, eq, sql } from "drizzle-orm";
import posthog from "posthog-js";
import * as v from "valibot";
import { browser } from "$app/environment";
import { PUBLIC_POSTHOG_PROJECT_API_KEY } from "$env/static/public";
import { db, schema } from "$lib/server/db";

v.setGlobalConfig({ lang: "es" });

export const load: LayoutServerLoad = async () => {
  if (browser) {
    posthog.init(PUBLIC_POSTHOG_PROJECT_API_KEY, {
      api_host: "https://us.i.posthog.com",
      defaults: "2026-01-30"
    });
  }

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

  const materias = await db
    .select({ codigo: schema.materia.codigo, nombre: schema.materia.nombre })
    .from(schema.materia)
    .innerJoin(schema.planMateria, eq(schema.planMateria.codigoMateria, schema.materia.codigo))
    .innerJoin(schema.plan, eq(schema.plan.codigo, schema.planMateria.codigoPlan))
    .where(eq(schema.plan.estaVigente, true))
    .groupBy(schema.materia.codigo, schema.materia.nombre)
    .orderBy(desc(cantidadPlanesExpr), desc(cantidadCatedrasExpr), schema.materia.nombre);

  return { materias };
};
