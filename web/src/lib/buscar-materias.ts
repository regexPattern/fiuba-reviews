import Fuse, { type FuseResultMatch } from "fuse.js";

export type MateriaBuscador = {
  codigo: string;
  nombre: string;
  cantidadCatedras: number;
  usaEquivalencias: boolean;
};

export type ResultadoMateria = MateriaBuscador & {
  coincidencias?: readonly FuseResultMatch[];
};

export function crearBuscadorMaterias(materias: MateriaBuscador[]) {
  const ordenadas = [...materias].sort(
    (a, b) =>
      Number(!b.usaEquivalencias && b.cantidadCatedras > 0) -
        Number(!a.usaEquivalencias && a.cantidadCatedras > 0) ||
      b.cantidadCatedras - a.cantidadCatedras ||
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

  // El fuzzy matching decide qué coincide; las cátedras propias tienen prioridad,
  // luego ordenamos por cantidad de cátedras y nombre.
  return (query: string): ResultadoMateria[] => {
    const texto = query.trim();
    return texto
      ? indice.search(texto).map(({ item, matches }) => ({ ...item, coincidencias: matches }))
      : ordenadas;
  };
}
