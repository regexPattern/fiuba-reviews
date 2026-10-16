"use client";

import { useRouter } from "next/navigation";
import { useMemo, useState } from "react";
import {
  Command,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList
} from "@/components/ui/command";
import { crearBuscadorMaterias, type MateriaBuscador } from "@/lib/buscar-materias";

export function BuscadorMaterias({ materias }: { materias: MateriaBuscador[] }) {
  const router = useRouter();
  const [query, setQuery] = useState("");
  const buscar = useMemo(() => crearBuscadorMaterias(materias), [materias]);
  const resultados = useMemo(() => buscar(query), [buscar, query]);

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
            value={materia.codigo}
            onSelect={() => router.push(`/materia/${materia.codigo}`)}
          >
            <span className="min-w-0 flex-1">
              {materia.codigo} · {materia.nombre}
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
