"use client";

import { Search } from "@withicons/react/line";
import { Command } from "cmdk";
import { useRouter } from "next/navigation";
import { type ReactNode, type Ref, useMemo, useState } from "react";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { crearIndiceDeBusqueda } from "@/lib/busqueda";
import { TextoConMatches } from "./texto-con-matches";

export type MateriaBuscador = { codigo: string; nombre: string };

const MAX_RESULTADOS = 30;

export function BuscadorMaterias({
  materias,
  query,
  onQueryChange,
  onSelect,
  inputAction,
  inputRef
}: {
  materias: MateriaBuscador[];
  query: string;
  onQueryChange: (query: string) => void;
  onSelect?: () => void;
  inputAction?: ReactNode;
  inputRef?: Ref<HTMLInputElement>;
}) {
  const router = useRouter();
  const debouncedQuery = useDebouncedValue(query);
  const buscar = useMemo(() => crearIndiceDeBusqueda(materias, ["codigo", "nombre"]), [materias]);
  const todos = useMemo(() => buscar(debouncedQuery), [buscar, debouncedQuery]);
  const resultados = todos.slice(0, MAX_RESULTADOS);
  const restantes = todos.length - resultados.length;

  /*
   * Cuando cambia el texto del input, cmdk resalta el primer resultado, pero
   * por el debounce la lista todavía es la de la búsqueda anterior. Cuando
   * llegan los resultados nuevos, el ítem resaltado puede no estar o no ser el
   * primero.
   *
   * Por eso la selección se guarda junto con la búsqueda en la que se hizo. Si
   * la búsqueda cambió, se resalta el primer resultado nuevo.
   */
  const [seleccion, setSeleccion] = useState({ query: debouncedQuery, codigo: "" });
  const seleccionada =
    (seleccion.query === debouncedQuery && seleccion.codigo) || (resultados[0]?.item.codigo ?? "");

  return (
    <Command
      label="Buscador de materias"
      shouldFilter={false}
      loop
      value={seleccionada}
      onValueChange={(codigo) => setSeleccion({ query: debouncedQuery, codigo })}
      className="flex min-h-0 flex-1 flex-col"
    >
      <div className="flex shrink-0 items-center gap-2 border-b border-stone-300 px-3">
        <Search size={18} aria-hidden="true" className="shrink-0 text-stone-500" />
        <Command.Input
          ref={inputRef}
          autoCapitalize="off"
          aria-label="Buscar materias por código o título"
          placeholder="Buscar materia o código…"
          value={query}
          onValueChange={onQueryChange}
          className="h-12 min-w-0 flex-1 bg-transparent text-stone-900 outline-none placeholder:text-stone-400 md:text-sm"
        />
        {inputAction}
      </div>
      <Command.List
        aria-label={
          debouncedQuery.trim()
            ? "Materias ordenadas por relevancia"
            : "Materias ordenadas por relevancia general"
        }
        /*
         * cmdk mete los ítems dentro de un `div[cmdk-list-sizer]`, así que el `space-y` va en ese
         * div y no en la lista (si no, solo separaría al sizer).
         */
        className="min-h-0 flex-1 scrollbar-gutter-stable overflow-y-auto overscroll-contain p-2 **:[[cmdk-list-sizer]]:space-y-1.5"
        /*
         * En iOS, mientras el input de búsqueda tiene el foco (teclado abierto), deslizar el dedo
         * sobre la lista scrollea la página entera en lugar de la lista. Para evitarlo, apenas el
         * usuario empieza a deslizar sobre la lista le sacamos el foco al input, lo que cierra el
         * teclado y el scroll pasa a la lista.
         */
        onTouchMove={() => {
          const activo = document.activeElement;
          if (activo instanceof HTMLInputElement && activo.hasAttribute("cmdk-input"))
            activo.blur();
        }}
      >
        {resultados.length === 0 && (
          <p className="py-6 text-center text-sm text-stone-400">No se encontraron materias.</p>
        )}
        {resultados.map(({ item: materia, matches }) => (
          <Command.Item
            key={materia.codigo}
            value={materia.codigo}
            onSelect={() => {
              onSelect?.();
              router.push(`/materias/${materia.codigo}`);
            }}
            /*
             * Se usa `scroll-m*` porque al navegar con las flechas, cmdk hace `scrollIntoView`. Con
             * estos márgenes, el ítem respeta el `p-2` de la lista y los 3px del relieve de abajo.
             */
            className="item-lista cursor-pointer scroll-mt-2 scroll-mb-[calc(--spacing(2)+3px)] select-none data-[selected=true]:item-lista-activo"
          >
            <TextoConMatches
              className="min-w-0 flex-1 text-base/6 tracking-tight wrap-break-word"
              texto={materia.nombre}
              indices={matches?.find((match) => match.key === "nombre")?.indices}
            />
            <TextoConMatches
              className="shrink-0 font-mono text-sm/6 text-stone-400"
              texto={materia.codigo}
              indices={matches?.find((match) => match.key === "codigo")?.indices}
            />
          </Command.Item>
        ))}
        {restantes > 0 && (
          <p className="px-2 py-3 text-center text-xs text-stone-500">
            <span className="font-mono tabular-nums">{restantes}</span>{" "}
            {restantes === 1 ? "materia más" : "materias más"} · seguí escribiendo para filtrar
          </p>
        )}
      </Command.List>
    </Command>
  );
}
