"use client";

import Link, { useLinkStatus } from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useDeferredValue, useMemo, useState } from "react";
import type { CatedraListada } from "@/server/queries/catedras-materia";
import { ListFilter, Loader2, Star } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";
import { cn } from "@/lib/utils";
import { useMateriaShell } from "./materia-shell";

const normalizar = (texto: string) =>
  texto
    .normalize("NFD")
    .replace(/\p{Diacritic}/gu, "")
    .toLowerCase();

function filtrarCatedras(catedras: CatedraListada[], query: string) {
  const palabras = normalizar(query).split(/\s+/).filter(Boolean);

  if (palabras.length === 0) return catedras;

  return catedras.filter((catedra) => {
    const nombre = normalizar(catedra.nombre);
    return palabras.every((palabra) => nombre.includes(palabra));
  });
}

function Contenido({ nombre, activa }: { nombre: string; activa: boolean }) {
  const { pending } = useLinkStatus();

  return (
    <span
      className={cn("flex items-center gap-2 transition-opacity", pending && "opacity-60")}
      aria-busy={pending}
    >
      <span className={activa ? "text-fiuba" : ""}>{nombre}</span>
      {pending && <Loader2 className="size-3 shrink-0 animate-spin" aria-hidden="true" />}
    </span>
  );
}

function Items({
  codigoMateria,
  catedras,
  activo
}: {
  codigoMateria: string;
  catedras: CatedraListada[];
  activo: string | null;
}) {
  const { cerrarDrawer } = useMateriaShell();

  return (
    <ul className="my-2">
      {catedras.map((catedra) => {
        const calificacion = catedra.calificacion.toFixed(1);

        return (
          <li key={catedra.codigo}>
            <Link
              href={`/materia/${codigoMateria}?catedra=${catedra.codigo}`}
              scroll={false}
              prefetch
              onClick={cerrarDrawer}
              className="flex items-center gap-1.5 p-3 tabular-nums"
            >
              <span className="inline-block w-[3ch] shrink-0 text-center">
                {calificacion === "0.0" ? "–" : calificacion}
              </span>
              <Star className="size-3 shrink-0 fill-yellow-500 stroke-yellow-700 dark:stroke-yellow-400" />
              <Contenido nombre={catedra.nombre} activa={activo === catedra.codigo} />
            </Link>
          </li>
        );
      })}
    </ul>
  );
}

/** Lee `?catedra` (dinámico): si falta o no pertenece a la materia, la activa es la primera. */
function ItemsConActiva(props: {
  codigoMateria: string;
  catedras: CatedraListada[];
  todas: CatedraListada[];
}) {
  const param = useSearchParams().get("catedra");
  const activo = props.todas.some((c) => c.codigo === param)
    ? param
    : (props.todas[0]?.codigo ?? null);

  return <Items {...props} activo={activo} />;
}

export function ListaCatedras({
  codigoMateria,
  catedras
}: {
  codigoMateria: string;
  catedras: CatedraListada[];
}) {
  const [query, setQuery] = useState("");
  const deferida = useDeferredValue(query);
  const filtradas = useMemo(() => filtrarCatedras(catedras, deferida), [catedras, deferida]);

  return (
    <div className="flex min-h-0 flex-1 flex-col divide-y divide-layout-border">
      <ScrollArea className="min-h-0 flex-1">
        {catedras.length === 0 ? (
          <p className="p-3 text-sm text-muted-foreground">
            Esta materia no tiene cátedras todavía.
          </p>
        ) : (
          <Suspense
            fallback={<Items codigoMateria={codigoMateria} catedras={filtradas} activo={null} />}
          >
            <ItemsConActiva codigoMateria={codigoMateria} catedras={filtradas} todas={catedras} />
          </Suspense>
        )}
        {catedras.length > 0 && filtradas.length === 0 && (
          <p className="p-3 text-sm text-muted-foreground">Sin resultados.</p>
        )}
      </ScrollArea>

      <div className="relative shrink-0">
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Filtrar por nombre de docente"
          aria-label="Filtrar cátedras por nombre de docente"
          className="w-full bg-transparent py-3 pr-12 pl-3 outline-none placeholder:text-sm"
        />
        <span className="pointer-events-none absolute top-1/2 right-3 flex size-6.5 -translate-y-1/2 items-center justify-center rounded-full border text-foreground/50">
          <ListFilter className="size-3" aria-hidden="true" />
        </span>
      </div>
    </div>
  );
}
