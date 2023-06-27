import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function parseCodigoMateria(codigo: string): number {
	// En caso de que el codigo recibido sea NaN asignamos un codigo inexistente.
	return parseInt(codigo, 10) || -1;
}
