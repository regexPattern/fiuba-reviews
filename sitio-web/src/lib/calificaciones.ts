import type { Cuatrimestre } from "@/lib/ofertas";

export const CRITERIOS = [
  { nombre: "aceptaCritica", etiqueta: "Acepta crítica" },
  { nombre: "asistencia", etiqueta: "Asistencia" },
  { nombre: "buenTrato", etiqueta: "Buen trato" },
  { nombre: "claridad", etiqueta: "Claridad" },
  { nombre: "claseOrganizada", etiqueta: "Clase organizada" },
  { nombre: "cumpleHorarios", etiqueta: "Cumple horarios" },
  { nombre: "fomentaParticipacion", etiqueta: "Fomenta participación" },
  { nombre: "panoramaAmplio", etiqueta: "Panorama amplio" },
  { nombre: "respondeMails", etiqueta: "Responde mails" }
] as const;

export type Criterio = (typeof CRITERIOS)[number]["nombre"];

/** De 0,5 a 5, de a media estrella. */
export const PUNTAJES = Array.from({ length: 10 }, (_, i) => (i + 1) / 2);
export type Puntaje = number;

export const PUNTAJE_MAXIMO = 5;

/** Con coma decimal y sin `,0`: `4`, `4,5`. */
export function formatearPuntaje(puntaje: Puntaje) {
  return puntaje.toLocaleString("es-AR", { maximumFractionDigits: 1 });
}

export const LARGO_MAXIMO_COMENTARIO = 2000;

/** Lo que se va completando en el panel de calificación, antes de enviarlo. */
export type BorradorCalificacion = {
  puntajes: Partial<Record<Criterio, Puntaje>>;
  cuatrimestre: Cuatrimestre | null;
  comentario: string;
};

export const BORRADOR_VACIO: BorradorCalificacion = {
  puntajes: {},
  cuatrimestre: null,
  comentario: ""
};

/** Criterios sin puntaje, en el orden de `CRITERIOS`. */
export function criteriosFaltantes(borrador: BorradorCalificacion): Criterio[] {
  return CRITERIOS.map(({ nombre }) => nombre).filter((nombre) => !borrador.puntajes[nombre]);
}

export function borradorCompleto(borrador: BorradorCalificacion) {
  return (
    criteriosFaltantes(borrador).length === 0 &&
    borrador.cuatrimestre !== null &&
    borrador.comentario.length <= LARGO_MAXIMO_COMENTARIO
  );
}

/** Si el borrador tiene algo cargado, para no tirarlo sin avisar. */
export function borradorTocado(borrador: BorradorCalificacion) {
  return (
    Object.keys(borrador.puntajes).length > 0 ||
    borrador.cuatrimestre !== null ||
    borrador.comentario.trim() !== ""
  );
}

/**
 * Cuatrimestres que se pueden elegir al calificar, del más nuevo al más viejo. El primero es el que
 * está en curso: el 1C de marzo a julio y el 2C de agosto a diciembre. En el verano (enero y febrero)
 * el último es el 2C del año anterior.
 */
export function cuatrimestresRecientes(hoy: Date, cantidad = 2): Cuatrimestre[] {
  const mes = hoy.getMonth();
  let actual: Cuatrimestre =
    mes >= 7
      ? { numero: 2, anio: hoy.getFullYear() }
      : mes >= 2
        ? { numero: 1, anio: hoy.getFullYear() }
        : { numero: 2, anio: hoy.getFullYear() - 1 };

  const cuatrimestres: Cuatrimestre[] = [];
  for (let i = 0; i < cantidad; i++) {
    cuatrimestres.push(actual);
    actual =
      actual.numero === 2 ? { numero: 1, anio: actual.anio } : { numero: 2, anio: actual.anio - 1 };
  }
  return cuatrimestres;
}

export type NivelPuntaje = "muyBueno" | "bueno" | "regular" | "malo";

export function nivelPuntaje(puntaje: number): NivelPuntaje {
  if (puntaje >= 4.5) return "muyBueno";
  if (puntaje >= 3.5) return "bueno";
  if (puntaje >= 2.5) return "regular";
  return "malo";
}
