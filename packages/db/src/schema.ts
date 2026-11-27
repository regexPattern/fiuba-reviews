import {
  boolean,
  doublePrecision,
  integer,
  jsonb,
  pgMaterializedView,
  pgTable,
  pgView,
  serial,
  smallint,
  text,
  timestamp,
  uuid
} from "drizzle-orm/pg-core";

export const cuatrimestre = pgTable("cuatrimestre", {
  codigo: serial().primaryKey(),
  numero: smallint().notNull(),
  anio: smallint().notNull()
});

export const carrera = pgTable("carrera", {
  codigo: serial().primaryKey(),
  nombre: text().notNull()
});

export const plan = pgTable("plan", {
  codigo: serial().primaryKey(),
  codigoCarrera: integer().notNull(),
  anio: smallint().notNull(),
  vigente: boolean().notNull()
});

/** Oferta de comisiones de una carrera en un cuatrimestre, ya procesada. */
export const ofertaComisiones = pgTable("oferta_comisiones", {
  codigoCarrera: integer().notNull(),
  codigoCuatrimestre: integer().notNull(),
  contenido: jsonb().notNull()
});

/** Oferta tal cual la pegaron desde el SIU, pendiente de procesar. */
export const ofertaComisionesRaw = pgTable("oferta_comisiones_raw", {
  codigoCarrera: integer().notNull(),
  codigoCuatrimestre: integer().notNull(),
  contenido: text().notNull()
});

export const materia = pgTable("materia", {
  codigo: text().primaryKey(),
  nombre: text().notNull(),
  codigoCuatrimestreUltimaActualizacion: integer()
});

export const catedra = pgTable("catedra", {
  codigo: uuid().primaryKey(),
  codigoMateria: text().notNull(),
  activa: boolean().notNull()
});

export const docente = pgTable("docente", {
  codigo: uuid().primaryKey(),
  nombre: text().notNull(),
  codigoMateria: text().notNull(),
  nombreSiu: text(),
  rol: text(),
  resumenComentarios: text()
});

export const catedraDocente = pgTable("catedra_docente", {
  codigoCatedra: uuid().notNull(),
  codigoDocente: uuid().notNull()
});

export const prioridadRol = pgTable("prioridad_rol", {
  rol: text().primaryKey(),
  prioridad: integer().notNull()
});

export const comentario = pgTable("comentario", {
  codigo: serial().primaryKey(),
  codigoDocente: uuid().notNull(),
  codigoCuatrimestre: integer().notNull(),
  contenido: text().notNull(),
  deDolly: boolean().notNull(),
  fechaCreacion: timestamp({ withTimezone: true, mode: "string" })
});

export const resumenCalificacionesDocente = pgTable("resumen_calificaciones_docente", {
  codigoDocente: uuid().primaryKey(),
  cantidad: integer().notNull(),
  promedio: doublePrecision().notNull(),
  aceptaCritica: doublePrecision().notNull(),
  asistencia: doublePrecision().notNull(),
  buenTrato: doublePrecision().notNull(),
  claridad: doublePrecision().notNull(),
  claseOrganizada: doublePrecision().notNull(),
  cumpleHorarios: doublePrecision().notNull(),
  fomentaParticipacion: doublePrecision().notNull(),
  panoramaAmplio: doublePrecision().notNull(),
  respondeMails: doublePrecision().notNull()
});

export const materiaVigente = pgMaterializedView("materia_vigente", {
  codigo: text().notNull(),
  nombre: text().notNull(),
  cantidadPlanes: integer().notNull(),
  cantidadCatedras: integer().notNull(),
  cuatrimestreNumero: smallint(),
  cuatrimestreAnio: smallint(),
  usaEquivalencias: boolean().notNull(),
  codigosMateriasFuente: text().array().notNull(),
  equivalencias: jsonb().$type<{ codigo: string; nombre: string }[]>().notNull(),
  orden: integer().notNull()
}).existing();

export const catedraMostrada = pgMaterializedView("catedra_mostrada", {
  codigoMateria: text().notNull(),
  codigoCatedra: uuid().notNull(),
  codigoMateriaCatedra: text().notNull(),
  nombre: text().notNull()
}).existing();

export const calificacionCatedra = pgView("calificacion_catedra", {
  codigoCatedra: uuid().notNull(),
  calificacion: doublePrecision().notNull()
}).existing();
