import { sql } from "drizzle-orm";
import {
  boolean,
  check,
  doublePrecision,
  foreignKey,
  index,
  integer,
  jsonb,
  numeric,
  pgMaterializedView,
  pgTable,
  pgView,
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
    }),
    // El actualizador del SIU resuelve docentes por (materia, nombre_siu).
    index("docente_codigo_materia_nombre_siu_idx").on(table.codigoMateria, table.nombreSiu)
  ]
);

export const cuatrimestre = pgTable(
  "cuatrimestre",
  {
    codigo: serial().primaryKey().notNull(),
    numero: smallint().notNull(),
    anio: smallint().notNull()
  },
  () => [check("cuatrimestre_numero_check", sql`numero = ANY (ARRAY[1, 2])`)]
);

export const materia = pgTable(
  "materia",
  {
    codigo: text().primaryKey().notNull(),
    nombre: text().notNull(),
    codigoCuatrimestreUltimaActualizacion: integer("codigo_cuatrimestre_ultima_actualizacion")
  },
  (table) => [
    foreignKey({
      columns: [table.codigoCuatrimestreUltimaActualizacion],
      foreignColumns: [cuatrimestre.codigo],
      name: "materia_codigo_cuatrimestre_ultima_actualizacion_fkey"
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
    vigente: boolean().notNull()
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
    }).onUpdate("cascade"),
    index("catedra_codigo_materia_activa_idx").on(table.codigoMateria, table.activa)
  ]
);

export const comentario = pgTable(
  "comentario",
  {
    codigo: serial().primaryKey().notNull(),
    codigoDocente: uuid("codigo_docente").notNull(),
    codigoCuatrimestre: integer("codigo_cuatrimestre").notNull(),
    contenido: text().notNull(),
    deDolly: boolean("de_dolly").default(false).notNull(),
    fechaCreacion: timestamp("fecha_creacion", { withTimezone: true, mode: "string" }).default(
      sql`(now() AT TIME ZONE 'America/Argentina/Buenos_Aires'::text)`
    ),
    codigoCalificacionDolly: integer("codigo_calificacion_dolly")
  },
  (table) => [
    // Sirve `WHERE codigo_docente IN (...) ORDER BY codigo_cuatrimestre DESC, codigo DESC`.
    index("comentario_docente_cuatrimestre_idx").on(
      table.codigoDocente,
      table.codigoCuatrimestre.desc(),
      table.codigo.desc()
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
      sql`(acepta_critica >= (1)::numeric) AND (acepta_critica <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_asistencia_check",
      sql`(asistencia >= (1)::numeric) AND (asistencia <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_buen_trato_check",
      sql`(buen_trato >= (1)::numeric) AND (buen_trato <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_claridad_check",
      sql`(claridad >= (1)::numeric) AND (claridad <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_clase_organizada_check",
      sql`(clase_organizada >= (1)::numeric) AND (clase_organizada <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_cumple_horarios_check",
      sql`(cumple_horarios >= (1)::numeric) AND (cumple_horarios <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_fomenta_participacion_check",
      sql`(fomenta_participacion >= (1)::numeric) AND (fomenta_participacion <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_panorama_amplio_check",
      sql`(panorama_amplio >= (1)::numeric) AND (panorama_amplio <= (5)::numeric)`
    ),
    check(
      "calificacion_dolly_responde_mails_check",
      sql`(responde_mails >= (1)::numeric) AND (responde_mails <= (5)::numeric)`
    )
  ]
);

/**
 * Resumen de las calificaciones de cada docente (cantidad y promedios por criterio).
 *
 * No se escribe a mano: lo mantienen los triggers de `calificacion_dolly`, así que siempre está
 * al día y evita agrupar las ~30k calificaciones en cada lectura. `promedio` es el promedio de
 * los 9 criterios.
 */
export const resumenCalificacionesDocente = pgTable(
  "resumen_calificaciones_docente",
  {
    codigoDocente: uuid("codigo_docente").primaryKey().notNull(),
    cantidad: integer().notNull(),
    promedio: doublePrecision().notNull(),
    aceptaCritica: doublePrecision("acepta_critica").notNull(),
    asistencia: doublePrecision().notNull(),
    buenTrato: doublePrecision("buen_trato").notNull(),
    claridad: doublePrecision().notNull(),
    claseOrganizada: doublePrecision("clase_organizada").notNull(),
    cumpleHorarios: doublePrecision("cumple_horarios").notNull(),
    fomentaParticipacion: doublePrecision("fomenta_participacion").notNull(),
    panoramaAmplio: doublePrecision("panorama_amplio").notNull(),
    respondeMails: doublePrecision("responde_mails").notNull()
  },
  (table) => [
    foreignKey({
      columns: [table.codigoDocente],
      foreignColumns: [docente.codigo],
      name: "resumen_calificaciones_docente_codigo_docente_fkey"
    }).onDelete("cascade")
  ]
);

export const prioridadRol = pgTable(
  "prioridad_rol",
  { rol: text().primaryKey().notNull(), prioridad: integer().notNull() },
  () => [check("prioridad_rol_prioridad_check", sql`(prioridad >= 1) AND (prioridad <= 10)`)]
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
    index("equivalencia_codigo_materia_plan_anterior_idx").on(table.codigoMateriaPlanAnterior),
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
      columns: [table.codigoCatedra, table.codigoDocente],
      name: "catedra_docente_pkey"
    }),
    index("catedra_docente_codigo_docente_idx").on(table.codigoDocente)
  ]
);

export const planMateria = pgTable(
  "plan_materia",
  {
    codigoPlan: integer("codigo_plan").notNull(),
    codigoMateria: text("codigo_materia").notNull(),
    electiva: boolean().notNull()
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
    primaryKey({ columns: [table.codigoPlan, table.codigoMateria], name: "plan_materia_pkey" }),
    index("plan_materia_codigo_materia_idx").on(table.codigoMateria)
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
      columns: [table.codigoCarrera, table.codigoCuatrimestre],
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
      columns: [table.codigoCarrera, table.codigoCuatrimestre],
      name: "oferta_comisiones_raw_pkey"
    })
  ]
);

// ---------------------------------------------------------------------------------------------
// Vistas. Se crean con SQL en `actualizador-ofertas/migrations/`; acá solo se declaran
// (`existing`) para consultarlas con tipos. SQLx administra las migraciones; Drizzle solo es el ORM.
// ---------------------------------------------------------------------------------------------

/**
 * Materias de planes vigentes (cambia solo al correr el actualizador del SIU).
 * `orden` es la posición de relevancia: planes desc, cátedras desc, nombre asc.
 */
export const materiaVigente = pgMaterializedView("materia_vigente", {
  codigo: text().notNull(),
  nombre: text().notNull(),
  cantidadPlanes: integer("cantidad_planes").notNull(),
  cantidadCatedras: integer("cantidad_catedras").notNull(),
  cuatrimestreNumero: smallint("cuatrimestre_numero"),
  cuatrimestreAnio: smallint("cuatrimestre_anio"),
  /** `true` si no tiene cátedras propias y muestra las de sus equivalencias. */
  usaEquivalencias: boolean("usa_equivalencias").notNull(),
  /** Materias de las que se leen las cátedras (ella misma, o sus equivalencias). */
  codigosMateriasFuente: text("codigos_materias_fuente").array().notNull(),
  equivalencias: jsonb().$type<{ codigo: string; nombre: string }[]>().notNull(),
  orden: integer().notNull()
}).existing();

/** Cátedras que se muestran en cada materia vigente, con su nombre ya armado. */
export const catedraMostrada = pgMaterializedView("catedra_mostrada", {
  codigoMateria: text("codigo_materia").notNull(),
  codigoCatedra: uuid("codigo_catedra").notNull(),
  codigoMateriaCatedra: text("codigo_materia_catedra").notNull(),
  nombre: text().notNull()
}).existing();

/** Calificación de cada cátedra: promedio de los promedios de sus docentes calificados. */
export const calificacionCatedra = pgView("calificacion_catedra", {
  codigoCatedra: uuid("codigo_catedra").notNull(),
  calificacion: doublePrecision().notNull()
}).existing();
