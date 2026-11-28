import { normalizar } from "@/lib/utils";

/**
 * Busca entre `carreras` la que corresponde al nombre leído del SIU, ignorando tildes, mayúsculas
 * y espacios de más. Devuelve `undefined` si no está registrada.
 */
export function buscarCarrera<T extends { nombre: string }>(carreras: T[], nombre: string) {
  const buscado = normalizar(nombre);
  return carreras.find((c) => normalizar(c.nombre) === buscado);
}
