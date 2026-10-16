import Fuse from "fuse.js";

export type MateriaBuscador = {
  codigo: string;
  nombre: string;
  cantidadComentarios: number;
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
    shouldSort: false
  });

  // El fuzzy matching decide qué coincide; la cantidad de comentarios decide el orden.
  return (query: string): MateriaBuscador[] => {
    const texto = query.trim();
    return texto ? indice.search(texto).map(({ item }) => item) : ordenadas;
  };
}
