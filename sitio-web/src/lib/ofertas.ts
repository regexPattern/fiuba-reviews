import { normalizarTexto } from "./texto";

/** La carrera con ese nombre, sin importar tildes, mayúsculas ni espacios. */
export function buscarCarrera<T extends { nombre: string }>(carreras: T[], nombre: string) {
  const buscado = normalizarTexto(nombre);
  return carreras.find((c) => normalizarTexto(c.nombre) === buscado);
}
