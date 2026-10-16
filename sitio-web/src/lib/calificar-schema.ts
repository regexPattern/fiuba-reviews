import { z } from "zod";
import { CALIFICACIONES_VALIDAS, CRITERIOS } from "./criterios";
import { UUID_RE } from "./utils";

const calificacion = z
  .number()
  .refine(
    (n) => (CALIFICACIONES_VALIDAS as readonly number[]).includes(n),
    "Calificación inválida."
  );

export const MIN_CHARS_COMENTARIO = 20;

export const calificarSchema = z.object({
  codigoDocente: z.string().regex(UUID_RE, "Código de docente inválido."),
  codigoCatedra: z.string().regex(UUID_RE, "Código de cátedra inválido.").nullable(),
  /** Materia desde la que se llegó (puede ser la vigente que muestra la cátedra por equivalencia). */
  codigoMateria: z
    .string()
    .regex(/^[0-9A-Za-z.-]+$/, "Código de materia inválido.")
    .nullable(),
  calificaciones: z.object(
    Object.fromEntries(CRITERIOS.map(({ name }) => [name, calificacion])) as Record<
      (typeof CRITERIOS)[number]["name"],
      typeof calificacion
    >
  ),
  comentario: z
    .string()
    .trim()
    .refine(
      (c) => c.length === 0 || c.length >= MIN_CHARS_COMENTARIO,
      `El comentario debe tener al menos ${MIN_CHARS_COMENTARIO} caracteres.`
    ),
  cuatrimestre: z.number({ error: "Seleccioná un cuatrimestre." }).int(),
  cfTurnstileResponse: z.string().min(1, "Completá el CAPTCHA.")
});

export type CalificarInput = z.infer<typeof calificarSchema>;
