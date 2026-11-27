import { z } from "zod";
import { type Cuatrimestre, formatearCuatrimestre, mismoCuatrimestre } from "./cuatrimestres";

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

type Criterio = (typeof CRITERIOS)[number]["nombre"];

export const PUNTAJE_MINIMO = 1;
export const PUNTAJE_MAXIMO = 5;

/** De 1 a 5, de a media estrella. */
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

const NOMBRES = CRITERIOS.map(({ nombre }) => nombre) as [Criterio, ...Criterio[]];

/** Criterios sin puntaje, en el orden de `CRITERIOS`. */
export function criteriosFaltantes(borrador: BorradorCalificacion): Criterio[] {
  return NOMBRES.filter((nombre) => !borrador.puntajes[nombre]);
}

const PuntajeSchema = z.number().min(PUNTAJE_MINIMO).max(PUNTAJE_MAXIMO).multipleOf(0.5);

/**
 * Valida en el servidor lo que manda el panel de calificación: no se confía en el cliente. El
 * cuatrimestre tiene que ser uno de los `permitidos` (los que se ofrecen en el panel), y se
 * devuelve el de `permitidos` tal cual. El comentario se recorta y queda `null` si está vacío.
 *
 * Solo el error de cuatrimestre trae un mensaje para mostrar; el resto usa los de zod.
 */
export function schemaCalificacion<C extends Cuatrimestre>(permitidos: C[]) {
  return z.object({
    codigoDocente: z.uuid(),
    puntajes: z.record(z.enum(NOMBRES), PuntajeSchema),
    cuatrimestre: z.object({ numero: z.number(), anio: z.number() }).transform((c, ctx) => {
      const permitido = permitidos.find((p) => mismoCuatrimestre(p, c));

      if (!permitido) {
        ctx.addIssue({
          code: "custom",
          message: `No se puede calificar el ${formatearCuatrimestre(c)}.`
        });
        return z.NEVER;
      }

      return permitido;
    }),
    comentario: z
      .string()
      .trim()
      .max(LARGO_MAXIMO_COMENTARIO)
      .transform((c) => c || null)
  });
}

export type NivelPuntaje = "muyBueno" | "bueno" | "regular" | "malo";

export function nivelPuntaje(puntaje: number): NivelPuntaje {
  if (puntaje >= 4.5) return "muyBueno";
  if (puntaje >= 3.5) return "bueno";
  if (puntaje >= 2.5) return "regular";
  return "malo";
}
