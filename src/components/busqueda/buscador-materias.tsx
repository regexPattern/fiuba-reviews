"use client";

import { Search } from "@withicons/react/pixel";
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
      <div className="flex shrink-0 items-center gap-2 border-b px-3">
        <Search size={16} aria-hidden="true" className="shrink-0 text-neutral-500" />
        {/* 16px en móvil evita el zoom de iOS al enfocar. */}
        <Command.Input
          aria-label="Buscar materias por nombre o código"
          placeholder="Buscar materia por nombre o código…"
          value={query}
          onValueChange={onQueryChange}
          className="h-12 min-w-0 flex-1 bg-transparent text-base outline-none placeholder:text-neutral-500 md:text-sm"
        />
        {inputAction}
      </div>
      <Command.List
        aria-label={
          debouncedQuery.trim()
            ? "Materias ordenadas por relevancia"
            : "Materias ordenadas por relevancia general"
        }
        className="min-h-0 flex-1 overflow-y-auto overscroll-contain p-1"
      >
        {resultados.length === 0 && (
          <p className="py-6 text-center text-sm">No se encontraron materias.</p>
        )}
        {resultados.map(({ item: materia, matches }) => (
          <Command.Item
            key={materia.codigo}
            value={materia.codigo}
            onSelect={() => {
              onSelect?.();
              router.push(`/materias/${materia.codigo}`);
            }}
            className="flex cursor-pointer items-center gap-3 rounded px-2 py-3 text-sm select-none data-[selected=true]:bg-neutral-100"
          >
            <span className="min-w-0 flex-1">
              <TextoConMatches
                className="font-mono"
                texto={materia.codigo}
                indices={matches?.find((match) => match.key === "codigo")?.indices}
              />
              {" · "}
              <TextoConMatches
                texto={materia.nombre}
                indices={matches?.find((match) => match.key === "nombre")?.indices}
              />
            </span>
          </Command.Item>
        ))}
        {restantes > 0 && (
          <p className="px-2 py-3 text-center text-xs text-neutral-500">
            {restantes === 1 ? "1 materia más" : `${restantes} materias más`} · seguí escribiendo
            para filtrar
          </p>
        )}
      </Command.List>
    </Command>
  );
}
