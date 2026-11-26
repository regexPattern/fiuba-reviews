import Fuse from "fuse.js";

/** Rango `[inicio, fin]` (inclusive) del texto original que coincide con la búsqueda. */
export type Rango = readonly [number, number];
export type Match = { key?: string; indices: readonly Rango[] };
export type Resultado<T> = { item: T; refIndex: number; matches?: readonly Match[] };

/**
 * Pasa a minúsculas y saca las tildes caracter por caracter, guardando para cada caracter del
 * texto normalizado su posición en el original (para subrayar el texto tal cual se muestra).
 */
function normalizar(texto: string) {
  let normalizado = "";
  const posiciones: number[] = [];
  for (let i = 0; i < texto.length; i++) {
    const caracter = texto[i].normalize("NFD").replace(/\p{M}/gu, "").toLowerCase();
    normalizado += caracter;
    for (let j = 0; j < caracter.length; j++) posiciones.push(i);
  }
  return { normalizado, posiciones };
}

const esInicioDePalabra = (texto: string, i: number) =>
  i === 0 || !/[\p{L}\p{N}]/u.test(texto[i - 1]);

/**
 * Búsqueda por palabras en las `keys` de cada item, en dos niveles:
 *
 * 1. Exactas: cada palabra de la búsqueda aparece tal cual (en cualquier orden, sin importar tildes
 *    ni mayúsculas) en alguna key. Primero las que coinciden con el comienzo de una palabra
 *    ("Análisis Matemático" antes que "Criptoanálisis") y, a igualdad, en el orden original.
 * 2. Aproximadas: después, lo que encuentra Fuse tolerando errores de tipeo ("gonsalez" →
 *    "González"), ordenado por parecido.
 *
 * Sin búsqueda devuelve todos, en el orden original.
 */
export function crearBusqueda<T>(items: T[], keys: (keyof T & string)[]) {
  const indexados = items.map((item, refIndex) => ({
    item,
    refIndex,
    campos: keys.map((key) => ({ key, ...normalizar(String(item[key] ?? "")) }))
  }));

  const fuse = new Fuse(items, {
    keys,
    ignoreDiacritics: true,
    ignoreLocation: true,
    threshold: 0.35,
    includeMatches: true,
    minMatchCharLength: 2,
    // Las palabras separadas por espacios se buscan en cualquier orden y tienen que estar todas.
    useExtendedSearch: true
  });

  function buscarExactas(palabras: string[]) {
    const resultados: (Resultado<T> & { puntaje: number })[] = [];
    for (const { item, refIndex, campos } of indexados) {
      const rangos = new Map<string, Rango[]>();
      let puntaje = 0;
      let coincidenTodas = true;

      for (const palabra of palabras) {
        // Por cada palabra, la mejor coincidencia: al comienzo de una palabra si la hay.
        let mejor: { key: string; inicio: number; fin: number; alInicio: boolean } | undefined;
        for (const { key, normalizado, posiciones } of campos) {
          for (
            let i = normalizado.indexOf(palabra);
            i !== -1;
            i = normalizado.indexOf(palabra, i + 1)
          ) {
            const alInicio = esInicioDePalabra(normalizado, i);
            if (!mejor || (alInicio && !mejor.alInicio)) {
              const fin = posiciones[i + palabra.length - 1];
              mejor = { key, inicio: posiciones[i], fin, alInicio };
            }
            if (alInicio) break;
          }
          if (mejor?.alInicio) break;
        }
        if (!mejor) {
          coincidenTodas = false;
          break;
        }
        if (!mejor.alInicio) puntaje++;
        rangos.set(mejor.key, [...(rangos.get(mejor.key) ?? []), [mejor.inicio, mejor.fin]]);
      }

      if (!coincidenTodas) continue;
      const matches = [...rangos].map(([key, indices]) => ({ key, indices: unirRangos(indices) }));
      resultados.push({ item, refIndex, matches, puntaje });
    }

    return resultados
      .sort((a, b) => a.puntaje - b.puntaje || a.refIndex - b.refIndex)
      .map(({ item, refIndex, matches }) => ({ item, refIndex, matches }));
  }

  return (query: string): Resultado<T>[] => {
    const palabras = normalizar(query).normalizado.split(/\s+/).filter(Boolean);
    if (palabras.length === 0) return items.map((item, refIndex) => ({ item, refIndex }));

    const exactas = buscarExactas(palabras);
    const encontrados = new Set(exactas.map(({ refIndex }) => refIndex));
    const aproximadas = fuse.search(query).filter(({ refIndex }) => !encontrados.has(refIndex));
    return [...exactas, ...aproximadas];
  };
}

/** Ordena y une rangos superpuestos, para poder subrayarlos de corrido. */
function unirRangos(rangos: Rango[]): Rango[] {
  const ordenados = [...rangos].sort((a, b) => a[0] - b[0]);
  const unidos: [number, number][] = [];
  for (const [inicio, fin] of ordenados) {
    const ultimo = unidos.at(-1);
    if (ultimo && inicio <= ultimo[1] + 1) ultimo[1] = Math.max(ultimo[1], fin);
    else unidos.push([inicio, fin]);
  }
  return unidos;
}
