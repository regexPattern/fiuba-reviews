import * as schema from "./schema";
import { sql } from "drizzle-orm";
import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

const client = postgres("postgres://postgres:postgres@localhost:5432");

export default drizzle(client, { logger: false, schema });

export function queryPromedioDocente<T extends number | (number | null)>() {
	return sql<T>`
			AVG((${schema.calificacion.aceptaCritica} +
			${schema.calificacion.asistencia} +
			${schema.calificacion.buenTrato} +
			${schema.calificacion.claridad} +
			${schema.calificacion.claseOrganizada} +
			${schema.calificacion.cumpleHorarios} +
			${schema.calificacion.fomentaParticipacion} +
			${schema.calificacion.panoramaAmplio} +
			${schema.calificacion.respondeMails}) / 9)
	`;
}
