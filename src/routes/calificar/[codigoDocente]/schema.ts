import { z } from "zod";

export default z.object({
	acepta_critica: z.number().min(0).max(5),
	asistencia: z.number().min(0).max(5),
	buen_trato: z.number().min(0).max(5),
	claridad: z.number().min(0).max(5),
	clase_organizada: z.number().min(0).max(5),
	cumple_horario: z.number().min(0).max(5),
	fomenta_participacion: z.number().min(0).max(5),
	panorama_amplio: z.number().min(0).max(5),
	responde_mails: z.number().min(0).max(5),
	comentario: z.string()
});
