import { sql } from "drizzle-orm";

export function exprPromedioDollyPorFila(cd: {
  aceptaCritica: any;
  asistencia: any;
  buenTrato: any;
  claridad: any;
  claseOrganizada: any;
  cumpleHorarios: any;
  fomentaParticipacion: any;
  panoramaAmplio: any;
  respondeMails: any;
}) {
  return sql<number>`(
    ${cd.aceptaCritica} +
    ${cd.asistencia} +
    ${cd.buenTrato} +
    ${cd.claridad} +
    ${cd.claseOrganizada} +
    ${cd.cumpleHorarios} +
    ${cd.fomentaParticipacion} +
    ${cd.panoramaAmplio} +
    ${cd.respondeMails}
  ) / 9.0`;
}
