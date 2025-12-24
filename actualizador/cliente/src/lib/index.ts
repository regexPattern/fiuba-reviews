export type PatchMateria = {
	codigo: string;
	nombre: string;
	carrera: string;
	cuatrimestre: {
		numero: number;
		anio: number;
	};
	docentes_pendientes: PatchDocente[];
	catedras: PatchCatedra[];
};

export type PatchDocente = {
	nombre: string;
	rol: string;
	matches: MatchDocente[];
};

export type MatchDocente = {
	codigo: string;
	nombre: string;
	score: number;
};

export type PatchCatedra = {
	ya_existente: boolean;
	docentes: {
		nombre: string;
		codigo: string | null;
	}[];
};
