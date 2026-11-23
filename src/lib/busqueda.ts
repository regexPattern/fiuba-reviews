import Fuse, { type FuseResult } from "fuse.js";

/** Búsqueda fuzzy en las `keys` de cada item. Sin búsqueda devuelve todos, en el orden original. */
export function crearBusqueda<T>(items: T[], keys: string[]) {
  const fuse = new Fuse(items, {
    keys,
    ignoreDiacritics: true,
    ignoreLocation: true,
    threshold: 0.35,
    includeMatches: true,
    // Las palabras separadas por espacios se buscan en cualquier orden y tienen que estar todas.
    useExtendedSearch: true
  });

  return (query: string): FuseResult<T>[] =>
    query.trim() ? fuse.search(query) : items.map((item, refIndex) => ({ item, refIndex }));
}
