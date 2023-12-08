import { z } from "zod";

const parametroCalificacion = z
	.number()
	.min(1, { message: "Valor mínimo 1" })
	.max(5, { message: "Valor máximo 5" });

export default z
	.object({
		// acepta_critica: parametroCalificacion,
		// asistencia: parametroCalificacion,
		// buen_trato: parametroCalificacion,
		// claridad: parametroCalificacion,
		// clase_organizada: parametroCalificacion,
		// cumple_horario: parametroCalificacion,
		// fomenta_participacion: parametroCalificacion,
		// panorama_amplio: parametroCalificacion,
		responde_mails: parametroCalificacion,
		comentario: z.string(),
		cuatrimestre: z.string().optional()
	})
	.refine((data) => (data.comentario.length > 0 ? data.cuatrimestre : true), {
		message: "Seleccione el cuatrimestre al que corresponde el comentario",
		path: ["cuatrimestre"]
	});
