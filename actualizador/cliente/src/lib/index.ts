export type PatchMateria = {
	codigo: string;
	nombre: string;
	docentes_sin_resolver: PatchDocente[];
	catedras: PatchCatedra[];
	cuatrimestre: {
		numero: number;
		anio: number;
	};
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
	codigo: number;
	docentes: {
		nombre: string;
		codigo_ya_resuelto: string | null;
	}[];
	resuelta: boolean;
};
