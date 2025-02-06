import {
	boolean,
	doublePrecision,
	foreignKey,
	integer,
	pgTable,
	primaryKey,
	smallint,
	text,
	unique,
	uuid
} from "drizzle-orm/pg-core";

export const materia = pgTable(
	"materia",
	{
		codigo: smallint("codigo").primaryKey().notNull(),
		nombre: text("nombre").notNull(),
		codigoEquivalencia: smallint("codigo_equivalencia")
	},
	(table) => {
		return {
			materiaCodigoEquivalenciaFkey: foreignKey({
				columns: [table.codigoEquivalencia],
				foreignColumns: [table.codigo],
				name: "materia_codigo_equivalencia_fkey"
			})
		};
	}
);

export const catedra = pgTable("catedra", {
	codigo: uuid("codigo").primaryKey().notNull(),
	codigoMateria: smallint("codigo_materia")
		.notNull()
		.references(() => materia.codigo)
});

export const docente = pgTable(
	"docente",
	{
		codigo: uuid("codigo").primaryKey().notNull(),
		nombre: text("nombre").notNull(),
		codigoMateria: smallint("codigo_materia")
			.notNull()
			.references(() => materia.codigo),
		resumenComentarios: text("resumen_comentarios"),
		comentariosUltimoResumen: integer("comentarios_ultimo_resumen").default(0).notNull()
	},
	(table) => {
		return {
			docenteNombreCodigoMateriaKey: unique("docente_nombre_codigo_materia_key").on(
				table.nombre,
				table.codigoMateria
			)
		};
	}
);

export const calificacion = pgTable("calificacion", {
	codigo: uuid("codigo").defaultRandom().primaryKey().notNull(),
	codigoDocente: uuid("codigo_docente")
		.notNull()
		.references(() => docente.codigo),
	aceptaCritica: doublePrecision("acepta_critica").notNull(),
	asistencia: doublePrecision("asistencia").notNull(),
	buenTrato: doublePrecision("buen_trato").notNull(),
	claridad: doublePrecision("claridad").notNull(),
	claseOrganizada: doublePrecision("clase_organizada").notNull(),
	cumpleHorarios: doublePrecision("cumple_horarios").notNull(),
	fomentaParticipacion: doublePrecision("fomenta_participacion").notNull(),
	panoramaAmplio: doublePrecision("panorama_amplio").notNull(),
	respondeMails: doublePrecision("responde_mails").notNull()
});

export const comentario = pgTable("comentario", {
	codigo: uuid("codigo").defaultRandom().primaryKey().notNull(),
	codigoDocente: uuid("codigo_docente")
		.notNull()
		.references(() => docente.codigo),
	cuatrimestre: text("cuatrimestre")
		.notNull()
		.references(() => cuatrimestre.nombre),
	contenido: text("contenido").notNull(),
	esDeDolly: boolean("es_de_dolly").default(false).notNull()
});

export const cuatrimestre = pgTable("cuatrimestre", {
	nombre: text("nombre").primaryKey().notNull()
});

export const catedraDocente = pgTable(
	"catedra_docente",
	{
		codigoCatedra: uuid("codigo_catedra")
			.notNull()
			.references(() => catedra.codigo, { onDelete: "cascade" }),
		codigoDocente: uuid("codigo_docente")
			.notNull()
			.references(() => docente.codigo)
	},
	(table) => {
		return {
			catedraDocentePkey: primaryKey({
				columns: [table.codigoCatedra, table.codigoDocente],
				name: "catedra_docente_pkey"
			})
		};
	}
);
