import type { calificacion } from "@prisma/client";
import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function promedioDocente(calificaciones: calificacion[]): number {
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

export function compararCuatrimestre(a: string, b: string): number {
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
