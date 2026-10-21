"use client";

import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import styles from "@/components/resultados-busqueda.module.css";
import { TextoConCoincidencias, rangosExactos } from "@/components/texto-con-coincidencias";
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList
} from "@/components/ui/command";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import { crearBuscadorMaterias, type MateriaBuscador } from "@/lib/buscar-materias";

export function BuscadorMaterias({ materias }: { materias: MateriaBuscador[] }) {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const debouncedQuery = useDebouncedValue(query);
  const buscar = useMemo(() => crearBuscadorMaterias(materias), [materias]);
  const resultados = useMemo(() => buscar(debouncedQuery), [buscar, debouncedQuery]);

  return (
    <Command label="Buscador de materias" shouldFilter={false} loop>
      <CommandInput
        aria-label="Buscar materias por nombre o código"
        placeholder="Buscar materia por nombre o código…"
        value={query}
        onValueChange={setQuery}
      />
      <CommandList aria-label="Materias ordenadas por cantidad de comentarios">
        <CommandEmpty>No se encontraron materias.</CommandEmpty>
        {resultados.map((materia) => (
          <CommandItem
            key={materia.codigo}
            className={styles.entrada}
            value={materia.codigo}
            onSelect={() => router.push(`/materia/${materia.codigo}`)}
          >
            <span className="min-w-0 flex-1">
              <TextoConCoincidencias
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
            <span className="shrink-0 text-xs text-muted-foreground">
              {materia.cantidadComentarios}{" "}
              {materia.cantidadComentarios === 1 ? "comentario" : "comentarios"}
            </span>
          </CommandItem>
        ))}
      </CommandList>
    </Command>
  );
}
