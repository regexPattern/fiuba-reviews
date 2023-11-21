import { z } from "zod";

const parametroCalificacion = z
  .number()
  .min(1, { message: "Valor mínimo 1" })
  .max(5, { message: "Valor máximo 5" });

export default z.object({
  acepta_critica: parametroCalificacion,
  asistencia: parametroCalificacion,
  buen_trato: parametroCalificacion,
  claridad: parametroCalificacion,
  clase_organizada: parametroCalificacion,
  cumple_horario: parametroCalificacion,
  fomenta_participacion: parametroCalificacion,
  panorama_amplio: parametroCalificacion,
  responde_mails: parametroCalificacion,
  comentario: z.string()
});
