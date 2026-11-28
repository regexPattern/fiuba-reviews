export type Cuatrimestre = { numero: number; anio: number };

/** Indica si el cuatrimestre `a` es posterior al `b`. */
export function esPosterior(a: Cuatrimestre, b: Cuatrimestre) {
  return a.anio > b.anio || (a.anio === b.anio && a.numero > b.numero);
}

/** Indica si `a` y `b` son el mismo cuatrimestre. */
export function esMismoCuatrimestre(a: Cuatrimestre, b: Cuatrimestre) {
  return a.numero === b.numero && a.anio === b.anio;
}

/** Formatea el cuatrimestre con el número, una `C` y el año. */
export function formatearCuatrimestre({ numero, anio }: Cuatrimestre) {
  return `${numero}C${anio}`;
}
