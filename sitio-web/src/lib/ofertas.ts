import { parseMetadata } from "@fiuba-reviews/parser-siu";

export type Cuatrimestre = { numero: number; anio: number };

/**
 * Carrera y cuatrimestre de una oferta pegada desde el SIU, o `null` si no se reconoce. Si la
 * página trae varios períodos lectivos (filtro en `Todos`), se toma el primero.
 */
export function extraerMetadataOferta(
  contenido: string
): { carrera: string; cuatrimestre: Cuatrimestre } | null {
  try {
    const [{ carrera, anio, cuatrimestre }] = parseMetadata(contenido);
    return { carrera, cuatrimestre: { numero: cuatrimestre, anio } };
  } catch {
    return null;
  }
}

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
