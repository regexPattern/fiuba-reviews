import { sql } from "drizzle-orm";
import {
  boolean,
  check,
  foreignKey,
  index,
  integer,
  jsonb,
  numeric,
  pgTable,
  primaryKey,
  serial,
  smallint,
  text,
  timestamp,
  uuid
} from "drizzle-orm/pg-core";

export const docente = pgTable(
  "docente",
  {
    codigo: uuid().defaultRandom().primaryKey().notNull(),
    nombre: text().notNull(),
    codigoMateria: text("codigo_materia").notNull(),
    resumenComentarios: text("resumen_comentarios"),
    comentariosUltimoResumen: integer("comentarios_ultimo_resumen").default(0).notNull(),
    nombreSiu: text("nombre_siu"),
    rol: text()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoMateria],
      foreignColumns: [materia.codigo],
      name: "docente_codigo_materia_fkey"
    })
  ]
);

export const cuatrimestre = pgTable(
  "cuatrimestre",
  {
    codigo: serial().primaryKey().notNull(),
    numero: smallint().notNull(),
    anio: smallint().notNull()
  },
  (table) => [check("cuatrimestre_numero_check", sql`numero = ANY (ARRAY[1, 2])`)]
);

export const materia = pgTable(
  "materia",
  {
    codigo: text().primaryKey().notNull(),
    nombre: text().notNull(),
    cuatrimestreUltimaActualizacion: integer("cuatrimestre_ultima_actualizacion")
  },
  (table) => [
    foreignKey({
      columns: [table.cuatrimestreUltimaActualizacion],
      foreignColumns: [cuatrimestre.codigo],
      name: "materia_cuatrimestre_ultima_actualizacion_fkey"
    })
  ]
);

export const carrera = pgTable("carrera", {
  codigo: serial().primaryKey().notNull(),
  nombre: text().notNull()
});

export const plan = pgTable(
  "plan",
  {
    codigo: serial().primaryKey().notNull(),
    codigoCarrera: integer("codigo_carrera").notNull(),
    anio: smallint().notNull(),
    estaVigente: boolean("esta_vigente").notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoCarrera],
      foreignColumns: [carrera.codigo],
      name: "plan_codigo_carrera_fkey"
    })
  ]
);

export const catedra = pgTable(
  "catedra",
  {
    codigo: uuid().defaultRandom().primaryKey().notNull(),
    codigoMateria: text("codigo_materia").notNull(),
    activa: boolean().default(false).notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoMateria],
      foreignColumns: [materia.codigo],
      name: "catedra_codigo_materia_fkey"
    }).onUpdate("cascade")
  ]
);

export const comentario = pgTable(
  "comentario",
  {
    codigo: serial().primaryKey().notNull(),
    codigoDocente: uuid("codigo_docente").notNull(),
    codigoCuatrimestre: integer("codigo_cuatrimestre").notNull(),
    contenido: text().notNull(),
    esDeDolly: boolean("es_de_dolly").default(false).notNull(),
    fechaCreacion: timestamp("fecha_creacion", { withTimezone: true, mode: "string" }).default(
      sql`(now() AT TIME ZONE 'America/Argentina/Buenos_Aires'::text)`
    ),
    codigoCalificacionDolly: integer("codigo_calificacion_dolly")
  },
  (table) => [
    index("comentario_codigo_docente_idx").using(
      "btree",
      table.codigoDocente.asc().nullsLast().op("uuid_ops")
    ),
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "comentario_codigo_docente_fkey"
    }),
    foreignKey({
      columns: [table.codigoCuatrimestre],
      foreignColumns: [cuatrimestre.codigo],
      name: "comentario_codigo_cuatrimestre_fkey"
    }),
    foreignKey({
      columns: [table.codigoCalificacionDolly],
      foreignColumns: [calificacionDolly.codigo],
      name: "comentario_codigo_calificacion_dolly_fkey"
    }).onDelete("set null")
  ]
);

export const calificacionDolly = pgTable(
  "calificacion_dolly",
  {
    codigo: serial().primaryKey().notNull(),
    codigoDocente: uuid("codigo_docente").notNull(),
    aceptaCritica: numeric("acepta_critica", { precision: 2, scale: 1 }).notNull(),
    asistencia: numeric({ precision: 2, scale: 1 }).notNull(),
    buenTrato: numeric("buen_trato", { precision: 2, scale: 1 }).notNull(),
    claridad: numeric({ precision: 2, scale: 1 }).notNull(),
    claseOrganizada: numeric("clase_organizada", { precision: 2, scale: 1 }).notNull(),
    cumpleHorarios: numeric("cumple_horarios", { precision: 2, scale: 1 }).notNull(),
    fomentaParticipacion: numeric("fomenta_participacion", { precision: 2, scale: 1 }).notNull(),
    panoramaAmplio: numeric("panorama_amplio", { precision: 2, scale: 1 }).notNull(),
    respondeMails: numeric("responde_mails", { precision: 2, scale: 1 }).notNull(),
    fechaCreacion: timestamp("fecha_creacion", { withTimezone: true, mode: "string" }).default(
      sql`(now() AT TIME ZONE 'America/Argentina/Buenos_Aires'::text)`
    )
  },
  (table) => [
    index("calificacion_dolly_codigo_docente_idx").using(
      "btree",
      table.codigoDocente.asc().nullsLast().op("uuid_ops")
    ),
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "calificacion_dolly_codigo_docente_fkey"
    }),
    check(
      "calificacion_dolly_acepta_critica_check",
      sql`(acepta_critica >= (0)::numeric) AND (acepta_critica <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_asistencia_check",
      sql`(asistencia >= (0)::numeric) AND (asistencia <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_buen_trato_check",
      sql`(buen_trato >= (0)::numeric) AND (buen_trato <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_claridad_check",
      sql`(claridad >= (0)::numeric) AND (claridad <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_clase_organizada_check",
      sql`(clase_organizada >= (0)::numeric) AND (clase_organizada <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_cumple_horarios_check",
      sql`(cumple_horarios >= (0)::numeric) AND (cumple_horarios <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_fomenta_participacion_check",
      sql`(fomenta_participacion >= (0)::numeric) AND (fomenta_participacion <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_panorama_amplio_check",
      sql`(panorama_amplio >= (0)::numeric) AND (panorama_amplio <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_responde_mails_check",
      sql`(responde_mails >= (0)::numeric) AND (responde_mails <= (5)::numeric)`
    )
  ]
);

export const prioridadRol = pgTable(
  "prioridad_rol",
  { rol: text().primaryKey().notNull(), prioridad: integer().notNull() },
  (table) => [check("prioridad_rol_prioridad_check", sql`(prioridad >= 1) AND (prioridad <= 10)`)]
);

export const equivalencia = pgTable(
  "equivalencia",
  {
    codigoMateriaPlanVigente: text("codigo_materia_plan_vigente").notNull(),
    codigoMateriaPlanAnterior: text("codigo_materia_plan_anterior").notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoMateriaPlanVigente],
      foreignColumns: [materia.codigo],
      name: "equivalencia_codigo_materia_plan_vigente_fkey"
    }).onUpdate("cascade"),
    foreignKey({
      columns: [table.codigoMateriaPlanAnterior],
      foreignColumns: [materia.codigo],
      name: "equivalencia_codigo_materia_plan_anterior_fkey"
    }),
    primaryKey({
      columns: [table.codigoMateriaPlanVigente, table.codigoMateriaPlanAnterior],
      name: "equivalencia_pkey"
    }),
    check("equivalencia_check", sql`codigo_materia_plan_vigente <> codigo_materia_plan_anterior`)
  ]
);

export const catedraDocente = pgTable(
  "catedra_docente",
  {
    codigoCatedra: uuid("codigo_catedra").notNull(),
    codigoDocente: uuid("codigo_docente").notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoCatedra],
      foreignColumns: [catedra.codigo],
      name: "catedra_docente_codigo_catedra_fkey"
    }).onDelete("cascade"),
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "catedra_docente_codigo_docente_fkey"
    }),
    primaryKey({
      columns: [table.codigoDocente, table.codigoCatedra],
      name: "catedra_docente_pkey"
    })
  ]
);

export const planMateria = pgTable(
  "plan_materia",
  {
    codigoPlan: integer("codigo_plan").notNull(),
    codigoMateria: text("codigo_materia").notNull(),
    esElectiva: boolean("es_electiva").notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoPlan],
      foreignColumns: [plan.codigo],
      name: "plan_materia_codigo_plan_fkey"
    }),
    foreignKey({
      columns: [table.codigoMateria],
      foreignColumns: [materia.codigo],
      name: "plan_materia_codigo_materia_fkey"
    }).onUpdate("cascade"),
    primaryKey({ columns: [table.codigoPlan, table.codigoMateria], name: "plan_materia_pkey" })
  ]
);

export const ofertaComisiones = pgTable(
  "oferta_comisiones",
  {
    codigoCarrera: integer("codigo_carrera").notNull(),
    codigoCuatrimestre: integer("codigo_cuatrimestre").notNull(),
    contenido: jsonb().notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoCarrera],
      foreignColumns: [carrera.codigo],
      name: "oferta_comisiones_codigo_carrera_fkey"
    }),
    foreignKey({
      columns: [table.codigoCuatrimestre],
      foreignColumns: [cuatrimestre.codigo],
      name: "oferta_comisiones_codigo_cuatrimestre_fkey"
    }),
    primaryKey({
      columns: [table.codigoCuatrimestre, table.codigoCarrera],
      name: "oferta_comisiones_pkey"
    })
  ]
);

export const ofertaComisionesRaw = pgTable(
  "oferta_comisiones_raw",
  {
    codigoCarrera: integer("codigo_carrera").notNull(),
    codigoCuatrimestre: integer("codigo_cuatrimestre").notNull(),
    contenido: text().notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoCarrera],
      foreignColumns: [carrera.codigo],
      name: "oferta_comisiones_raw_codigo_carrera_fkey"
    }),
    foreignKey({
      columns: [table.codigoCuatrimestre],
      foreignColumns: [cuatrimestre.codigo],
      name: "oferta_comisiones_raw_codigo_cuatrimestre_fkey"
    }),
    primaryKey({
      columns: [table.codigoCuatrimestre, table.codigoCarrera],
      name: "oferta_comisiones_raw_pkey"
    })
  ]
);
