import { z } from "zod";
import { type Cuatrimestre, esMismoCuatrimestre, formatearCuatrimestre } from "./cuatrimestres";

export const CRITERIOS = [
  { nombre: "aceptaCritica", label: "Acepta crítica" },
  { nombre: "asistencia", label: "Asistencia" },
  { nombre: "buenTrato", label: "Buen trato" },
  { nombre: "claridad", label: "Claridad" },
  { nombre: "claseOrganizada", label: "Clase organizada" },
  { nombre: "cumpleHorarios", label: "Cumple horarios" },
  { nombre: "fomentaParticipacion", label: "Fomenta participación" },
  { nombre: "panoramaAmplio", label: "Panorama amplio" },
  { nombre: "respondeMails", label: "Responde mails" }
] as const;

export type Criterio = (typeof CRITERIOS)[number]["nombre"];

export const MIN_PUNTAJE = 1;
export const MAX_PUNTAJE = 5;

export const MAX_CARACTERES_COMENTARIO = 2000;

const SchemaPuntaje = z.number().min(MIN_PUNTAJE).max(MAX_PUNTAJE).multipleOf(0.5);

/** Crea el schema de validación de una calificación a partir de los `cuatrimestres` permitidos. */
export function crearSchemaCalificacion<C extends Cuatrimestre>(cuatrimestres: C[]) {
  return z.object({
    codigoDocente: z.uuid(),
    puntajes: z.record(z.enum(CRITERIOS.map(({ nombre }) => nombre)), SchemaPuntaje),
    cuatrimestre: z.object({ numero: z.number(), anio: z.number() }).transform((c, ctx) => {
      const permitido = cuatrimestres.find((p) => esMismoCuatrimestre(p, c));

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
      .max(MAX_CARACTERES_COMENTARIO)
      .transform((c) => c || null)
  });
}
