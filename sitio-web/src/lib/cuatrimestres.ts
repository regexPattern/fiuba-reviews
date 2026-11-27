export type Cuatrimestre = { numero: number; anio: number };

/** Orden cronológico de un cuatrimestre. */
export function ordenCuatrimestre({ numero, anio }: Cuatrimestre) {
  return anio * 2 + numero;
}

export function mismoCuatrimestre(a: Cuatrimestre, b: Cuatrimestre) {
  return a.numero === b.numero && a.anio === b.anio;
}

/** Cuatrimestre en formato corto, como `1C2026`. */
export function formatearCuatrimestre({ numero, anio }: Cuatrimestre) {
  return `${numero}C${anio}`;
}
