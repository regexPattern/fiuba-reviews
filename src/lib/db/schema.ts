import { getTableColumns } from "drizzle-orm";
import {
	doublePrecision,
	foreignKey,
	integer,
	pgTable,
	primaryKey,
	text,
	uuid
} from "drizzle-orm/pg-core";

export const materia = pgTable(
	"materia",
	{
		codigo: integer("codigo").primaryKey().notNull(),
		nombre: text("nombre").notNull(),
		codigoEquivalencia: integer("codigo_equivalencia")
	},
	(table) => ({
		materiaCodigoEquivalenciaFkey: foreignKey({
			columns: [table.codigoEquivalencia],
			foreignColumns: [table.codigo]
		})
	})
);

export const catedra = pgTable("catedra", {
	codigo: uuid("codigo").primaryKey().notNull(),
	codigoMateria: integer("codigo_materia")
		.notNull()
		.references(() => materia.codigo)
});

export const cuatrimestre = pgTable("cuatrimestre", {
	nombre: text("nombre").primaryKey().notNull()
});

export const catedraDocente = pgTable(
	"catedra_docente",
	{
		codigoCatedra: uuid("codigo_catedra")
			.notNull()
			.references(() => catedra.codigo),
		codigoDocente: uuid("codigo_docente")
			.notNull()
			.references(() => docente.codigo)
	},
	(table) => {
		return {
			catedraDocentePkey: primaryKey(table.codigoCatedra, table.codigoDocente)
		};
	}
);

export const docente = pgTable("docente", {
	codigo: uuid("codigo").primaryKey().notNull(),
	nombre: text("nombre").notNull(),
	descripcion: text("descripcion"),
	comentariosUltimaDescripcion: integer("comentarios_ultima_descripcion").notNull()
});

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
	contenido: text("contenido").notNull()
});
