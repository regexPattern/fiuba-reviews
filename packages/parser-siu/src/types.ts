export type Metadata = {
  carrera: string;
  anio: number;
  cuatrimestre: 1 | 2;
};

export type Docente = {
  nombre: string;
  rol: string | null;
};

export type Clase = {
  tipo: string;
  dia: string;
  inicio: string;
  fin: string;
  aula: string | null;
};

export type Comision = {
  nombre: string;
  docentes: Docente[];
  clases: Clase[];
};

export type Materia = {
  codigo: string;
  nombre: string;
  comisiones: Comision[];
};

export type Oferta = Metadata & {
  materias: Materia[];
};
