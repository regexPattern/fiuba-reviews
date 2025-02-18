import { z } from "zod";

const campoNumerico = z
  .number()
  .min(1, { message: "Valor mínimo 1" })
  .max(5, { message: "Valor máximo 5" });

export const formCalificacionDocente = z
  .object({
    ["acepta-critica"]: campoNumerico,
    ["asistencia"]: campoNumerico,
    ["buen-trato"]: campoNumerico,
    ["claridad"]: campoNumerico,
    ["clase-organizada"]: campoNumerico,
    ["cumple-horario"]: campoNumerico,
    ["fomenta-participacion"]: campoNumerico,
    ["panorama-amplio"]: campoNumerico,
    ["responde-mails"]: campoNumerico,
    ["comentario"]: z.string(),
    ["cuatrimestre"]: z.number().optional(),
    ["cf-turnstile-response"]: z.string(),
  })
  .refine((data) => (data.comentario.length > 0 ? data.cuatrimestre : true), {
    message: "Cuatrimestre requerido.",
    path: ["cuatrimestre"],
  });

export const formPlanSiu = z.object({
  ["carrera"]: z.string(),
  ["contenido-siu"]: z
    .string()
    .min(1, { message: "Contenido del SIU no puede estar vacío." }),
  ["cf-turnstile-response"]: z.string(),
});
