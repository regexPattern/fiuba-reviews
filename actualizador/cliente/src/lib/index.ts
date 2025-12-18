export type PatchMateria = {
	codigo: string;
	nombre: string;
	docentes: PatchDocente[];
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
	similitud: number;
};

export type PatchCatedra = {
	codigo: number;
	docentes: {
		nombre: string;
		rol: string;
	}[];
};
