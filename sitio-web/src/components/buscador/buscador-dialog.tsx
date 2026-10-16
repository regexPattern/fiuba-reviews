"use client";

import { useRouter } from "next/navigation";
import {
  CommandDialog,
  CommandEmpty,
  CommandInput,
  CommandItem,
  CommandList
} from "@/components/ui/command";
import { useBuscador } from "./buscador-provider";

const normalizar = (texto: string) =>
  texto
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .toLowerCase();

/** Todas las palabras buscadas tienen que aparecer (sin tildes ni mayúsculas) en código o nombre. */
function filtrar(value: string, search: string) {
  const texto = normalizar(value);
  return normalizar(search)
    .split(/\s+/)
    .filter(Boolean)
    .every((palabra) => texto.includes(palabra))
    ? 1
    : 0;
}

export function BuscadorDialog() {
  const { materias, abierto, setAbierto } = useBuscador();
  const router = useRouter();

  return (
    <CommandDialog
      open={abierto}
      onOpenChange={setAbierto}
      title="Buscar materias"
      description="Buscá una materia por código o nombre."
      showCloseButton={false}
      className="top-[calc(56px+env(safe-area-inset-top)+16px)] max-w-[min(560px,94vw)] translate-y-0 sm:max-w-[min(560px,94vw)]"
      commandProps={{ filter: filtrar }}
    >
      <CommandInput placeholder="Código o nombre de materia" />
      <CommandList className="max-h-[280px]">
        <CommandEmpty>Sin resultados.</CommandEmpty>
        {materias.map((materia) => (
          <CommandItem
            key={materia.codigo}
            value={`${materia.codigo} ${materia.nombre}`}
            onSelect={() => {
              setAbierto(false);
              router.push(`/materia/${materia.codigo}`);
            }}
            className="gap-3 data-[selected=true]:text-fiuba"
          >
            <span className="text-xs tabular-nums">{materia.codigo}</span>
            <span>{materia.nombre}</span>
          </CommandItem>
        ))}
      </CommandList>
    </CommandDialog>
  );
}
