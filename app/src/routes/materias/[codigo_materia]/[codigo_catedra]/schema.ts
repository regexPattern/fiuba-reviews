import { z } from "zod";
import { zfd } from "zod-form-data";

const calificacion_numerica = zfd.numeric(z.number().min(1).max(5));

export default zfd.formData({
	codigo_docente: zfd.text(),
	acepta_critica: calificacion_numerica,
	asistencia: calificacion_numerica,
	buen_trato: calificacion_numerica,
	claridad: calificacion_numerica,
	clase_organizada: calificacion_numerica,
	cumple_horarios: calificacion_numerica,
	fomenta_participacion: calificacion_numerica,
	panorama_amplio: calificacion_numerica,
	responde_mails: calificacion_numerica,
	comentario: zfd.text().optional()
});
