export type Materia = {
	codigo: number;
	nombre: string;
	codigo_equivalencia?: number;
};

export type Catedra = {
	codigo: string;
	nombre: string;
	promedio: number;
};

export type Docente = {
	codigo: string;
	nombre: string;
	respuestas: number;
	promedio: number;
	acepta_critica?: number;
	asistencia?: number;
	buen_trato?: number;
	claridad?: number;
	clase_organizada?: number;
	cumple_horarios?: number;
	fomenta_participacion?: number;
	panorama_amplio?: number;
	responde_mails?: number;
	comentarios: Comentario[];
};

export type Comentario = {
	codigo: string;
	codigo_docente: string;
	cuatrimestre: string;
	contenido: string;
};
