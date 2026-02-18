// La funcionalidad de búsqueda se puede acceder desde varios componentes diferentes, que pueden o
// no estar montados en el DOM al mismo tiempo. Por ejemplo: la pantalla de inicio tiene un botón
// insertado en el contenido principal de la pantalla que activa esta funcionalidad; la navbar en
// mobile y en desktop son dos componentes diferentes debido a las dificultades para posicionar los
// botones si fueran un único elemento y se usaran únicamente reglas de estilo. Si cada uno de
// estos elementos renderizara el mismo dialog, con la misma funcionalidad de búsqueda (mount de
// fuze, keydown event listener, etc.) estaríamos cargando toda esta funcionalidad varias veces en
// memoria, aunque no se use.
//
// Para solucionar esto lo que se hace es que los diferentes puntos de acceso a la funcionalidad de
// búsqueda sirven únicamente como triggers. Todos estos triggers desembocan en la funcionalidad
// implementada en este módulo, con el fin de optimizar el bundle size que se le envía al cliente.
//
import Fuse from "fuse.js";

const FUZZY_SEARCH_THRESHOLD = 0.25;
const FUZZY_SEARCH_DEBOUNCE_TIMEOUT_MS = 300;

export type MateriaBuscador = { codigo: string; nombre: string };

let open = $state(false);
let materias = $state<MateriaBuscador[]>([]);
let queryValue = $state("");
let queryDebounced = $state("");
let debounceTimeoutHandler: ReturnType<typeof setTimeout> | null = null;

let fuse = $derived(
  new Fuse(materias, {
    ignoreDiacritics: true,
    ignoreFieldNorm: true,
    shouldSort: true,
    includeScore: true,
    threshold: FUZZY_SEARCH_THRESHOLD,
    keys: ["codigo", "nombre"]
  })
);

let materiasFiltradas = $derived.by(() => {
  if (queryDebounced.trim() === "") {
    return materias;
  }

  return fuse
    .search(queryDebounced)
    .sort((a, b) => (a.score ?? 0) - (b.score ?? 0) || a.refIndex - b.refIndex)
    .map((resultado) => resultado.item);
});

function setMaterias(nuevasMaterias: MateriaBuscador[]) {
  materias = nuevasMaterias;
}

function limpiarQuery() {
  queryValue = "";
  queryDebounced = "";

  if (debounceTimeoutHandler) {
    clearTimeout(debounceTimeoutHandler);
    debounceTimeoutHandler = null;
  }
}

function setQueryValue(nuevoValor: string) {
  queryValue = nuevoValor;

  if (debounceTimeoutHandler) {
    clearTimeout(debounceTimeoutHandler);
    debounceTimeoutHandler = null;
  }

  if (nuevoValor.trim() === "") {
    queryDebounced = "";
    return;
  }

  debounceTimeoutHandler = setTimeout(() => {
    queryDebounced = nuevoValor;
  }, FUZZY_SEARCH_DEBOUNCE_TIMEOUT_MS);
}

function abrirBuscador() {
  open = true;
}

function cerrarBuscador() {
  open = false;
}

export const buscadorMaterias = {
  get open() {
    return open;
  },
  set open(abierto: boolean) {
    open = abierto;
  },
  get queryValue() {
    return queryValue;
  },
  get materiasFiltradas() {
    return materiasFiltradas;
  },
  setMaterias,
  setQueryValue,
  limpiarQuery,
  abrirBuscador,
  cerrarBuscador
};
