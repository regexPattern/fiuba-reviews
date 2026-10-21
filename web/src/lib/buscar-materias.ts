import Fuse, { type FuseResultMatch } from "fuse.js";

export type MateriaBuscador = {
  codigo: string;
  nombre: string;
  cantidadComentarios: number;
};

export type ResultadoMateria = MateriaBuscador & {
  coincidencias?: readonly FuseResultMatch[];
};

export function crearBuscadorMaterias(materias: MateriaBuscador[]) {
  const ordenadas = [...materias].sort(
    (a, b) =>
      b.cantidadComentarios - a.cantidadComentarios ||
      a.nombre.localeCompare(b.nombre, "es") ||
      a.codigo.localeCompare(b.codigo)
  );
  const indice = new Fuse(ordenadas, {
    keys: ["nombre", "codigo"],
    ignoreDiacritics: true,
    ignoreLocation: true,
    threshold: 0.35,
    shouldSort: false,
    includeMatches: true
  });

  // El fuzzy matching decide qué coincide; la cantidad de comentarios decide el orden.
  return (query: string): ResultadoMateria[] => {
    const texto = query.trim();
    return texto
      ? indice.search(texto).map(({ item, matches }) => ({ ...item, coincidencias: matches }))
      : ordenadas;
  };
}
