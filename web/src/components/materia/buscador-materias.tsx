"use client";

import { useRouter } from "next/navigation";
import { type ReactNode, useMemo, useState } from "react";
import styles from "@/components/resultados-busqueda.module.css";
import { TextoConCoincidencias, rangosExactos } from "@/components/texto-con-coincidencias";
import { Command, CommandInput, CommandItem, CommandList } from "@/components/ui/command";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { crearBuscadorMaterias, type MateriaBuscador } from "@/lib/buscar-materias";

export function BuscadorMaterias({
  materias,
  onSelect,
  inputAction,
  className
}: {
  materias: MateriaBuscador[];
  onSelect?: () => void;
  inputAction?: ReactNode;
  className?: string;
}) {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedValue(query);
  const buscar = useMemo(() => crearBuscadorMaterias(materias), [materias]);
  const resultados = useMemo(() => buscar(debouncedQuery), [buscar, debouncedQuery]);

  return (
    <Command className={className} label="Buscador de materias" shouldFilter={false} loop>
      <CommandInput
        aria-label="Buscar materias por nombre o código"
        placeholder="Buscar materia por nombre o código…"
        value={query}
        onValueChange={setQuery}
        trailingAction={inputAction}
      />
      <CommandList aria-label="Materias ordenadas por cantidad de cátedras, priorizando cátedras propias">
        {resultados.length === 0 && (
          <p className="py-6 text-center text-sm">No se encontraron materias.</p>
        )}
        {resultados.map((materia) => (
          <CommandItem
            key={materia.codigo}
            className={styles.entrada}
            value={materia.codigo}
            onSelect={() => {
              onSelect?.();
              router.push(`/materia/${materia.codigo}`);
            }}
          >
            <span className="min-w-0 flex-1">
              <TextoConCoincidencias
                className="font-mono"
                texto={materia.codigo}
                rangos={rangosExactos(
                  materia.codigo,
                  debouncedQuery,
                  materia.coincidencias
                    ?.filter((match) => match.key === "codigo")
                    .flatMap((match) => match.indices)
                )}
              />
              {" · "}
              <TextoConCoincidencias
                texto={materia.nombre}
                rangos={rangosExactos(
                  materia.nombre,
                  debouncedQuery,
                  materia.coincidencias
                    ?.filter((match) => match.key === "nombre")
                    .flatMap((match) => match.indices)
                )}
              />
            </span>
          </CommandItem>
        ))}
      </CommandList>
    </Command>
  );
}
