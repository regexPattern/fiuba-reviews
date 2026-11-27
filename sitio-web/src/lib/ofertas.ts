export type Cuatrimestre = { numero: number; anio: number };

/** Normaliza un nombre ignorando tildes, mayúsculas y espacios de más. */
export function normalizarNombre(nombre: string) {
  return nombre
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .replace(/\s+/g, " ")
    .trim()
    .toLowerCase();
}

/** Orden cronológico de un cuatrimestre. */
export function ordenCuatrimestre({ numero, anio }: Cuatrimestre) {
  return anio * 2 + numero;
}

/** Cuatrimestre en formato corto, como `1C2026`. */
export function formatearCuatrimestre({ numero, anio }: Cuatrimestre) {
  return `${numero}C${anio}`;
}
