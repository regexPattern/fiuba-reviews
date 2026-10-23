import Fuse, { type FuseResultMatch } from "fuse.js";

export type MateriaBuscador = {
  codigo: string;
  nombre: string;
  cantidadCatedras: number;
};

export type ResultadoMateria = MateriaBuscador & {
  coincidencias?: readonly FuseResultMatch[];
};

export function crearBuscadorMaterias(materias: MateriaBuscador[]) {
  const ordenadas = [...materias].sort(
    (a, b) =>
      b.cantidadCatedras - a.cantidadCatedras ||
      a.nombre.localeCompare(b.nombre, "es") ||
      a.codigo.localeCompare(b.codigo)
  );
  const indice = new Fuse(ordenadas, {
    keys: ["nombre", "codigo"],
    ignoreDiacritics: true,
    ignoreLocation: true,
    threshold: 0.35,
    shouldSort: true,
    includeMatches: true
  });

  // Sin búsqueda: cantidad de cátedras y nombre. Con búsqueda: relevancia del
  // fuzzy matching; Fuse desempata por la posición en el listado ordenado.
  return (query: string): ResultadoMateria[] => {
    const texto = query.trim();
    return texto
      ? indice.search(texto).map(({ item, matches }) => ({ ...item, coincidencias: matches }))
      : ordenadas;
  };
}
