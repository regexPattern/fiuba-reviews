export type PatchMateria = {
	codigo: string;
	nombre: string;
	cuatrimestre: {
		numero: number;
		anio: number;
	};
	docentes_pendientes: PatchDocente[];
	docentes_por_catedra: PatchCatedra[];
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
	nombre: string;
	codigo: string | null;
}[];
