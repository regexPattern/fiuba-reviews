import type { calificacion } from "@prisma/client";

// Comparador de cuatrimestres que los ordena de mas a menos reciente.
export function cmpCuatrimestre(a: string, b: string) {
	const [cuatriA, anioA] = a.split("Q");
	const [cuatriB, anioB] = b.split("Q");

	if (anioA < anioB) {
		return 1;
	} else if (anioA > anioB) {
		return -1;
	} else {
		return cuatriA <= cuatriB ? 1 : -1;
	}
}

// Calculador del promedio de un docente en base a su collecion de
// calificaciones.
export function calcPromedioDocente(docente: { calificacion: calificacion[] }) {
	const total = docente.calificacion
		.map((c) => {
			const params = [
				c.acepta_critica,
				c.asistencia,
				c.buen_trato,
				c.claridad,
				c.clase_organizada,
				c.cumple_horarios,
				c.fomenta_participacion,
				c.panorama_amplio,
				c.responde_mails
			];

			return params.reduce((acc, curr) => acc + curr) / params.length || 0;
		})
		.reduce((acc, curr) => acc + curr, 0);

	return total / docente.calificacion.length || 0;
}

// Estariza el formato en el que se muestra el nombre de una catedra, generado
// a partir de los nombres de los docentes que la componen.
export function fmtNombreCatedra(nombresDocentes: { nombre: string }[]): string {
	return nombresDocentes.sort().join(", ");
}
