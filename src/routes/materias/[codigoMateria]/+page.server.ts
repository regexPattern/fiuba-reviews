import db from "$lib/db";
import { calificacion, catedra, catedraDocente, docente } from "$lib/db/schema";
import { redirect } from "@sveltejs/kit";
import { desc, eq, sql } from "drizzle-orm";

import type { PageServerLoad } from "./$types";

export const prerender = true;

export const load: PageServerLoad = async ({ params }) => {
  const codigoMateria = parseInt(params.codigoMateria, 10);

  const defaultCatedra = (
    await db
      .select({
        codigo: catedra.codigo,
        promedio: sql<number>`
AVG((
  SELECT AVG((
    ${calificacion.aceptaCritica} 
      + ${calificacion.asistencia} 
      + ${calificacion.buenTrato} 
      + ${calificacion.claridad} 
      + ${calificacion.claseOrganizada} 
      + ${calificacion.cumpleHorarios} 
      + ${calificacion.fomentaParticipacion} 
      + ${calificacion.panoramaAmplio} 
      + ${calificacion.respondeMails}) / 9)
	FROM ${calificacion}
	WHERE ${calificacion.codigoDocente} = ${docente.codigo}
  GROUP BY ${docente.codigo})
)`
      })
      .from(catedra)
      .innerJoin(catedraDocente, eq(catedra.codigo, catedraDocente.codigoCatedra))
      .innerJoin(docente, eq(docente.codigo, catedraDocente.codigoDocente))
      .where(eq(catedra.codigoMateria, codigoMateria))
      .groupBy(catedra.codigo)
      .orderBy(({ promedio }) => desc(promedio))
      .limit(1)
  )[0];

  throw redirect(307, `/materias/${codigoMateria}/${defaultCatedra.codigo}`);
};
