"use client";

import type { CatedraListada } from "@/server/db/queries/catedras";
import { ListFilter } from "@withicons/react/line";
import { Star } from "@withicons/react/plush";
import Link from "next/link";
import { useParams } from "next/navigation";
import { useMemo, useState } from "react";
import { TextoConMatches } from "@/components/busqueda/texto-con-matches";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { crearBusqueda } from "@/lib/busqueda";
import { useMateriaShell } from "./materia-shell";

export function ListaCatedras({
  codigoMateria,
  catedras
}: {
  codigoMateria: string;
  catedras: CatedraListada[];
}) {
  const { alNavegar } = useMateriaShell();
  // Sin cátedra en la URL, la página muestra la primera.
  const activa = useParams<{ catedra?: string }>().catedra ?? catedras[0]?.codigo;
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedValue(query);
  const buscar = useMemo(() => crearBusqueda(catedras, ["nombre"]), [catedras]);
  const resultados = useMemo(() => buscar(debouncedQuery), [buscar, debouncedQuery]);

  if (catedras.length === 0) {
    return <p className="p-4 text-sm">Esta materia no tiene cátedras todavía.</p>;
  }

  return (
    <div className="flex min-h-0 flex-1 flex-col">
      {/* `px-2`, igual que la lista: el input queda alineado con el borde de la cátedra seleccionada. */}
      <div className="shrink-0 px-2 pt-3 pb-1">
        <label className="relative block">
          <span className="sr-only">Filtrar cátedras por docente</span>
          <input
            type="search"
            autoComplete="off"
            autoCorrect="off"
            autoCapitalize="off"
            spellCheck={false}
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Filtrar por docente"
            className="h-9 w-full campo pr-9 md:text-sm"
          />
          <ListFilter
            size={18}
            aria-hidden="true"
            className="pointer-events-none absolute top-1/2 right-2.5 -translate-y-1/2 text-stone-600"
          />
        </label>
      </div>
      <nav
        aria-label="Cátedras"
        className="min-h-0 flex-1 [scrollbar-gutter:stable] overflow-y-auto overscroll-contain"
      >
        <ul className="space-y-1.5 p-2">
          {resultados.map(({ item: catedra, matches }) => (
            <li key={catedra.codigo}>
              <Link
                href={`/materias/${codigoMateria}/catedras/${catedra.codigo}`}
                // Prefetch completo, no solo hasta el Suspense: el click es instantáneo.
                prefetch={true}
                // Scrollea la shell, así no queda debajo del header.
                scroll={false}
                onNavigate={alNavegar}
                aria-current={catedra.codigo === activa ? "page" : undefined}
                className="item-lista hover:bg-stone-100 hover:text-stone-900 aria-[current=page]:item-lista-activo"
              >
                <span className="flex h-6 w-12 shrink-0 items-center gap-1 font-mono text-stone-600 tabular-nums">
                  <Star size={16} aria-hidden="true" className="shrink-0 text-fiuba" />
                  {catedra.calificacion > 0 ? catedra.calificacion.toFixed(1) : "–"}
                </span>
                <TextoConMatches
                  className="min-w-0 flex-1 text-base/6 tracking-tight"
                  texto={catedra.nombre}
                  indices={matches?.[0]?.indices}
                />
              </Link>
            </li>
          ))}
        </ul>
        {resultados.length === 0 && <p className="p-4 text-sm text-stone-400">Sin resultados.</p>}
      </nav>
    </div>
  );
}
