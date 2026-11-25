"use client";

import { Search } from "@withicons/react/plush";
import { Command } from "cmdk";
import { useRouter } from "next/navigation";
import { type ReactNode, useMemo, useState } from "react";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { crearBusqueda } from "@/lib/busqueda";
import { TextoConMatches } from "./texto-con-matches";

export type MateriaBuscador = { codigo: string; nombre: string };

const MAX_RESULTADOS = 30;

export function BuscadorMaterias({
  materias,
  query,
  onQueryChange,
  onSelect,
  inputAction
}: {
  materias: MateriaBuscador[];
  query: string;
  onQueryChange: (query: string) => void;
  onSelect?: () => void;
  inputAction?: ReactNode;
}) {
  const router = useRouter();
  const debouncedQuery = useDebouncedValue(query);
  const buscar = useMemo(() => crearBusqueda(materias, ["codigo", "nombre"]), [materias]);
  const todos = useMemo(() => buscar(debouncedQuery), [buscar, debouncedQuery]);
  const resultados = todos.slice(0, MAX_RESULTADOS);
  const restantes = todos.length - resultados.length;

  // La selección se controla para que cada búsqueda nueva arranque en el primer resultado: cmdk
  // intenta elegirlo al cambiar el texto, pero los resultados llegan después, con el debounce.
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
        {/* 16px en móvil evita el zoom de iOS al enfocar. */}
        <Command.Input
          autoCapitalize="off"
          aria-label="Buscar materias por código o título"
          placeholder="Buscar materia o código…"
          value={query}
          onValueChange={onQueryChange}
          className="h-12 min-w-0 flex-1 bg-transparent text-base text-tinta outline-none placeholder:text-stone-400 md:text-sm"
        />
        {inputAction}
      </div>
      <Command.List
        aria-label={
          debouncedQuery.trim()
            ? "Materias ordenadas por relevancia"
            : "Materias ordenadas por relevancia general"
        }
        // cmdk envuelve los ítems en `[cmdk-list-sizer]`: el espaciado va ahí.
        className="min-h-0 flex-1 [scrollbar-gutter:stable] overflow-y-auto overscroll-contain p-2 [&_[cmdk-list-sizer]]:space-y-1.5"
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
            // Mismo estilo que la lista de cátedras: la seleccionada queda en relieve celeste.
            // `scroll-m*`: al navegar con flechas cmdk hace `scrollIntoView` y sin esto el ítem queda
            // pegado al borde, sin el `p-2` de la lista. Abajo suma los 3px del canto del relieve.
            className="flex cursor-pointer scroll-mt-2 scroll-mb-[calc(--spacing(2)+3px)] items-start gap-2 rounded-md border-[1.5px] border-transparent px-2.5 py-2 text-sm text-stone-600 select-none data-[selected=true]:relieve data-[selected=true]:bg-sky-50 data-[selected=true]:text-tinta data-[selected=true]:[--relieve:var(--color-fiuba-oscuro)]"
          >
            {/* Nombre primero y código apagado a la derecha, como en el encabezado de la materia. */}
            <TextoConMatches
              className="min-w-0 flex-1 font-acento text-base/6 wrap-break-word"
              texto={materia.nombre}
              indices={matches?.find((match) => match.key === "nombre")?.indices}
            />
            <TextoConMatches
              className="shrink-0 font-numeros text-sm/6 text-stone-400"
              texto={materia.codigo}
              indices={matches?.find((match) => match.key === "codigo")?.indices}
            />
          </Command.Item>
        ))}
        {restantes > 0 && (
          <p className="px-2 py-3 text-center text-xs text-stone-500">
            <span className="font-numeros tabular-nums">{restantes}</span>{" "}
            {restantes === 1 ? "materia más" : "materias más"} · seguí escribiendo para filtrar
          </p>
        )}
      </Command.List>
    </Command>
  );
}
