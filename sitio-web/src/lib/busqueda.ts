import Fuse, { type FuseResult } from "fuse.js";
import { normalizar } from "./utils";

/** Índices `[inicio, fin]` (inclusivo) del texto que matchea la query de
 * búsqueda. Sirven para resaltar el match en la UI. */
export type Rango = readonly [number, number];

/** Crea el índice para búsquedas fuzzy sobre los `items` según las `keys`, y
 * devuelve la función de búsqueda. */
export function crearIndiceDeBusqueda<T>(items: T[], keys: (keyof T & string)[]) {
  const fuse = new Fuse(items, {
    keys,
    ignoreDiacritics: true,
    ignoreLocation: true,
    threshold: 0.35,
    includeMatches: true,
    /* Las palabras separadas por espacios se buscan en cualquier orden y tienen que estar todas. */
    useExtendedSearch: true
  });

  return (query: string): FuseResult<T>[] => {
    const palabras = normalizar(query).split(/\s+/).filter(Boolean);
    if (palabras.length === 0) return items.map((item, refIndex) => ({ item, refIndex }));

    /* TODO: Revisar la performance de siempre normalizar el texto. */
    const esExacta = (item: T) =>
      palabras.every((p) => keys.some((k) => normalizar(String(item[k] ?? "")).includes(p)));

    const resultados = fuse.search(query);
    const exactas = resultados
      .filter((r) => esExacta(r.item))
      .sort((a, b) => a.refIndex - b.refIndex);
    const aproximadas = resultados.filter((r) => !esExacta(r.item));
    return [...exactas, ...aproximadas];
  };
}
