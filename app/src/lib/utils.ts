import type { calificacion } from "@prisma/client";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function comparar_cuatrimestre(a: string, b: string): number {
	const [cuatri_a, anio_a] = a.split("Q");
	const [cuatri_b, anio_b] = b.split("Q");

	if (anio_a < anio_b) {
		return 1;
	} else if (anio_a > anio_b) {
		return -1;
	} else {
		return cuatri_a <= cuatri_b ? 1 : -1;
	}
}

export function calcular_promedio_docente(calificaciones: calificacion[]): number {
	const total = calificaciones
		.map(
			(c) =>
				(c.acepta_critica +
					c.asistencia +
					c.buen_trato +
					c.claridad +
					c.clase_organizada +
					c.cumple_horarios +
					c.fomenta_participacion +
					c.panorama_amplio +
					c.responde_mails) /
				9
		)
		.reduce((acc, curr) => acc + curr, 0);

	return total / calificaciones.length;
}

export function generar_nombre_catedra(docentes: { nombre: string }[]): string {
	return docentes.sort().join(", ");
}
