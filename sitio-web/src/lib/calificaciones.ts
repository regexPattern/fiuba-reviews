import { z } from "zod";
import { type Cuatrimestre, formatearCuatrimestre } from "./ofertas";

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

export const PUNTAJE_MINIMO = 1;
export const PUNTAJE_MAXIMO = 5;

/** De 1 a 5, de a media estrella. */
export const PUNTAJES = Array.from(
  { length: (PUNTAJE_MAXIMO - PUNTAJE_MINIMO) * 2 + 1 },
  (_, i) => PUNTAJE_MINIMO + i / 2
);
export type Puntaje = number;

const FORMATO_PUNTAJE = new Intl.NumberFormat("es-AR", {
  minimumFractionDigits: 1,
  maximumFractionDigits: 1
});

/** Con coma decimal y siempre un decimal: `4,0`, `3,7`. */
export function formatearPuntaje(puntaje: Puntaje) {
  return FORMATO_PUNTAJE.format(puntaje);
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

/** Calificación lista para guardar: todos los puntajes válidos y el comentario sin espacios de más. */
export type CalificacionValida = {
  codigoDocente: string;
  puntajes: Record<Criterio, Puntaje>;
  cuatrimestre: Cuatrimestre;
  /** `null` si no dejó comentario. */
  comentario: string | null;
};

const puntajeValido = z.number().refine((p) => PUNTAJES.includes(p));

const EnvioCalificacion = z.object({
  codigoDocente: z.uuid(),
  puntajes: z.object(
    Object.fromEntries(CRITERIOS.map(({ nombre }) => [nombre, puntajeValido])) as Record<
      Criterio,
      typeof puntajeValido
    >
  ),
  cuatrimestre: z.object({ numero: z.number(), anio: z.number() }),
  comentario: z.string().max(LARGO_MAXIMO_COMENTARIO)
});

/**
 * Valida en el servidor lo que manda el panel de calificación: no se confía en el cliente. El
 * cuatrimestre tiene que ser uno de los `permitidos` (los que se ofrecen en el panel).
 */
export function validarCalificacion(
  entrada: unknown,
  permitidos: Cuatrimestre[]
): { ok: true; calificacion: CalificacionValida } | { ok: false; mensaje: string } {
  const envio = EnvioCalificacion.safeParse(entrada);

  if (!envio.success) {
    return { ok: false, mensaje: "La calificación enviada no es válida." };
  }

  const { codigoDocente, puntajes, cuatrimestre, comentario } = envio.data;
  const permitido = permitidos.some(
    (c) => c.numero === cuatrimestre.numero && c.anio === cuatrimestre.anio
  );

  if (!permitido) {
    return {
      ok: false,
      mensaje: `No se puede calificar el ${formatearCuatrimestre(cuatrimestre)}.`
    };
  }

  const recortado = comentario.trim();

  return {
    ok: true,
    calificacion: {
      codigoDocente,
      puntajes,
      cuatrimestre: { numero: cuatrimestre.numero, anio: cuatrimestre.anio },
      comentario: recortado === "" ? null : recortado
    }
  };
}

export type NivelPuntaje = "muyBueno" | "bueno" | "regular" | "malo";

export function nivelPuntaje(puntaje: number): NivelPuntaje {
  if (puntaje >= 4.5) return "muyBueno";
  if (puntaje >= 3.5) return "bueno";
  if (puntaje >= 2.5) return "regular";
  return "malo";
}
