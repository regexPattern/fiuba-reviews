export type Cuatrimestre = { numero: number; anio: number };

/** Para comparar nombres de carrera sin importar tildes, mayúsculas ni espacios. */
export function normalizarNombre(nombre: string) {
  return nombre
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .replace(/\s+/g, " ")
    .trim()
    .toLowerCase();
}

/** Orden cronológico de un cuatrimestre, para compararlos. */
export function ordenCuatrimestre({ numero, anio }: Cuatrimestre) {
  return anio * 2 + numero;
}

export function formatearCuatrimestre({ numero, anio }: Cuatrimestre) {
  return `${numero}C${anio}`;
}
