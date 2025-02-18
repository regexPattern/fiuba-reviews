import {
  pgTable,
  pgPolicy,
  check,
  serial,
  smallint,
  foreignKey,
  integer,
  boolean,
  index,
  uuid,
  text,
  timestamp,
  unique,
  doublePrecision,
  primaryKey,
} from "drizzle-orm/pg-core";
import { sql } from "drizzle-orm";

export const cuatrimestre = pgTable(
  "cuatrimestre",
  {
    codigo: serial().primaryKey().notNull(),
    numero: smallint().notNull(),
    anio: smallint().notNull(),
  },
  (table) => [
    pgPolicy("Enable read access for all users", {
      as: "permissive",
      for: "select",
      to: ["public"],
      using: sql`true`,
    }),
    check("cuatrimestre_numero_check", sql`numero = ANY (ARRAY[1, 2])`),
  ],
);

export const plan = pgTable(
  "plan",
  {
    codigo: serial().primaryKey().notNull(),
    codigoCarrera: integer("codigo_carrera").notNull(),
    anio: smallint().notNull(),
    estaVigente: boolean("esta_vigente").notNull(),
  },
  (table) => [
    foreignKey({
      columns: [table.codigoCarrera],
      foreignColumns: [carrera.codigo],
      name: "plan_codigo_carrera_fkey",
    }),
  ],
);

export const comentario = pgTable(
  "comentario",
  {
    codigo: serial().primaryKey().notNull(),
    codigoDocente: uuid("codigo_docente").notNull(),
    codigoCuatrimestre: integer("codigo_cuatrimestre").notNull(),
    contenido: text().notNull(),
    esDeDolly: boolean("es_de_dolly").default(false).notNull(),
    fechaCreacion: timestamp("fecha_creacion", {
      withTimezone: true,
      mode: "string",
    }).default(
      sql`(now() AT TIME ZONE 'America/Argentina/Buenos_Aires'::text)`,
    ),
  },
  (table) => [
    index("idx_comentario_codigo_docente").using(
      "btree",
      table.codigoDocente.asc().nullsLast().op("uuid_ops"),
    ),
    foreignKey({
      columns: [table.codigoCuatrimestre],
      foreignColumns: [cuatrimestre.codigo],
      name: "comentario_codigo_cuatrimestre_fkey",
    }),
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "comentario_codigo_docente_fkey",
    }),
  ],
);

export const materia = pgTable("materia", {
  codigo: text().primaryKey().notNull(),
  nombre: text().notNull(),
});

export const catedra = pgTable(
  "catedra",
  {
    codigo: uuid().defaultRandom().primaryKey().notNull(),
    codigoMateria: text("codigo_materia").notNull(),
  },
  (table) => [
    index("idx_catedra_codigo_materia").using(
      "btree",
      table.codigoMateria.asc().nullsLast().op("text_ops"),
    ),
    foreignKey({
      columns: [table.codigoMateria],
      foreignColumns: [materia.codigo],
      name: "catedra_codigo_materia_fkey",
    }),
  ],
);

export const carrera = pgTable("carrera", {
  codigo: serial().primaryKey().notNull(),
  nombre: text().notNull(),
});

export const docente = pgTable(
  "docente",
  {
    codigo: uuid().defaultRandom().primaryKey().notNull(),
    nombre: text().notNull(),
    codigoMateria: text("codigo_materia").notNull(),
    resumenComentarios: text("resumen_comentarios"),
    comentariosUltimoResumen: integer("comentarios_ultimo_resumen")
      .default(0)
      .notNull(),
  },
  (table) => [
    foreignKey({
      columns: [table.codigoMateria],
      foreignColumns: [materia.codigo],
      name: "docente_codigo_materia_fkey_cascade",
    }),
    unique("docente_nombre_codigo_materia_key").on(
      table.nombre,
      table.codigoMateria,
    ),
  ],
);

export const calificacionDolly = pgTable(
  "calificacion_dolly",
  {
    codigo: serial().primaryKey().notNull(),
    codigoDocente: uuid("codigo_docente").notNull(),
    aceptaCritica: doublePrecision("acepta_critica").notNull(),
    asistencia: doublePrecision().notNull(),
    buenTrato: doublePrecision("buen_trato").notNull(),
    claridad: doublePrecision().notNull(),
    claseOrganizada: doublePrecision("clase_organizada").notNull(),
    cumpleHorarios: doublePrecision("cumple_horarios").notNull(),
    fomentaParticipacion: doublePrecision("fomenta_participacion").notNull(),
    panoramaAmplio: doublePrecision("panorama_amplio").notNull(),
    respondeMails: doublePrecision("responde_mails").notNull(),
  },
  (table) => [
    index("idx_calificacion_dolly_codigo_docente").using(
      "btree",
      table.codigoDocente.asc().nullsLast().op("uuid_ops"),
    ),
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "calificacion_dolly_codigo_docente_fkey",
    }),
    check(
      "check_rango_calificacion",
      sql`((acepta_critica >= (1)::double precision) AND (acepta_critica <= (5)::double precision)) AND ((asistencia >= (1)::double precision) AND (asistencia <= (5)::double precision)) AND ((buen_trato >= (1)::double precision) AND (buen_trato <= (5)::double precision)) AND ((claridad >= (1)::double precision) AND (claridad <= (5)::double precision)) AND ((clase_organizada >= (1)::double precision) AND (clase_organizada <= (5)::double precision)) AND ((cumple_horarios >= (1)::double precision) AND (cumple_horarios <= (5)::double precision)) AND ((fomenta_participacion >= (1)::double precision) AND (fomenta_participacion <= (5)::double precision)) AND ((panorama_amplio >= (1)::double precision) AND (panorama_amplio <= (5)::double precision)) AND ((responde_mails >= (1)::double precision) AND (responde_mails <= (5)::double precision))`,
    ),
  ],
);

export const equivalencia = pgTable(
  "equivalencia",
  {
    codigoMateriaPlanVigente: text("codigo_materia_plan_vigente").notNull(),
    codigoMateriaPlanAnterior: text("codigo_materia_plan_anterior").notNull(),
  },
  (table) => [
    foreignKey({
      columns: [table.codigoMateriaPlanAnterior],
      foreignColumns: [materia.codigo],
      name: "equivalencia_codigo_materia_plan_anterior_fkey",
    }),
    foreignKey({
      columns: [table.codigoMateriaPlanVigente],
      foreignColumns: [materia.codigo],
      name: "equivalencia_codigo_materia_plan_vigente_fkey",
    }),
    primaryKey({
      columns: [
        table.codigoMateriaPlanVigente,
        table.codigoMateriaPlanAnterior,
      ],
      name: "equivalencia_pkey",
    }),
    check(
      "equivalencia_check",
      sql`codigo_materia_plan_vigente <> codigo_materia_plan_anterior`,
    ),
  ],
);

export const catedraDocente = pgTable(
  "catedra_docente",
  {
    codigoCatedra: uuid("codigo_catedra").notNull(),
    codigoDocente: uuid("codigo_docente").notNull(),
  },
  (table) => [
    index("idx_catedra_docente_codigo_catedra").using(
      "btree",
      table.codigoCatedra.asc().nullsLast().op("uuid_ops"),
    ),
    index("idx_catedra_docente_codigo_docente").using(
      "btree",
      table.codigoDocente.asc().nullsLast().op("uuid_ops"),
    ),
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "catedra_docente_codigo_docente_fkey",
    }),
    foreignKey({
      columns: [table.codigoCatedra],
      foreignColumns: [catedra.codigo],
      name: "catedra_docente_codigo_catedra_fkey",
    }).onDelete("cascade"),
    primaryKey({
      columns: [table.codigoCatedra, table.codigoDocente],
      name: "catedra_docente_pkey",
    }),
  ],
);

export const planMateria = pgTable(
  "plan_materia",
  {
    codigoPlan: integer("codigo_plan").notNull(),
    codigoMateria: text("codigo_materia").notNull(),
    esElectiva: boolean("es_electiva").notNull(),
  },
  (table) => [
    foreignKey({
      columns: [table.codigoMateria],
      foreignColumns: [materia.codigo],
      name: "plan_materia_codigo_materia_fkey",
    }),
    foreignKey({
      columns: [table.codigoPlan],
      foreignColumns: [plan.codigo],
      name: "plan_materia_codigo_plan_fkey",
    }),
    primaryKey({
      columns: [table.codigoPlan, table.codigoMateria],
      name: "plan_materia_pkey",
    }),
  ],
);
