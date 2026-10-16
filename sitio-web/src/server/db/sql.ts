import { sql } from "drizzle-orm";
import { calificacionDolly } from "./schema";

/** Promedio de los 9 criterios de una calificación individual. */
export const promedioDollyExpr = sql<number>`(
  ${calificacionDolly.aceptaCritica} +
  ${calificacionDolly.asistencia} +
  ${calificacionDolly.buenTrato} +
  ${calificacionDolly.claridad} +
  ${calificacionDolly.claseOrganizada} +
  ${calificacionDolly.cumpleHorarios} +
  ${calificacionDolly.fomentaParticipacion} +
  ${calificacionDolly.panoramaAmplio} +
  ${calificacionDolly.respondeMails}
) / 9.0`;
