"use client";

import type { CatedraListada } from "@/server/db/queries/catedras-materia";
import { Star } from "@withicons/react/chrome";
import Fuse from "fuse.js";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useMemo, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { mozillaHeadline } from "@/lib/fonts";
import { useMateriaShell } from "./materia-shell";

function Items({
  codigoMateria,
  catedras,
  activo
}: {
  codigoMateria: string;
  catedras: CatedraListada[];
  activo: string | null;
}) {
  const { cerrarDrawer, navegarCatedra, catedraPendiente } = useMateriaShell();
  const seleccionado = catedraPendiente ?? activo;

  return (
    <ul className="space-y-1 p-2">
      {catedras.map((catedra) => {
        const calificacion = catedra.calificacion.toFixed(1);

        return (
          <li key={catedra.codigo}>
            <Button
              asChild
              variant={seleccionado === catedra.codigo ? "secondary" : "ghost"}
              className="h-auto w-full justify-start text-left whitespace-normal"
            >
              <Link
                aria-current={seleccionado === catedra.codigo ? "page" : undefined}
                href={`/materia/${codigoMateria}?catedra=${catedra.codigo}`}
                scroll={false}
                onNavigate={(event) => {
                  event.preventDefault();
                  cerrarDrawer();
                  if (seleccionado !== catedra.codigo) {
                    navegarCatedra(codigoMateria, catedra.codigo);
                  }
                }}
              >
                <span className="inline-flex w-12 shrink-0 items-center justify-end gap-1 tabular-nums">
                  {calificacion === "0.0" ? "–" : calificacion}
                  <Star size={16} className="gold-star" aria-hidden="true" />
                </span>
                <span className={mozillaHeadline.className}>{catedra.nombre}</span>
              </Link>
            </Button>
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
  const indice = useMemo(
    () =>
      new Fuse(catedras, {
        keys: ["nombre"],
        ignoreDiacritics: true,
        ignoreLocation: true,
        threshold: 0.35,
        shouldSort: false
      }),
    [catedras]
  );
  const filtradas = useMemo(() => {
    const palabras = query.trim().split(/\s+/).filter(Boolean);
    if (palabras.length === 0) return catedras;

    // Todas las palabras deben coincidir, sin exigir orden ni perder el ranking por calificación.
    const coincidencias = palabras.map(
      (palabra) => new Set(indice.search(palabra).map(({ item }) => item.codigo))
    );
    return catedras.filter((catedra) => coincidencias.every((set) => set.has(catedra.codigo)));
  }, [catedras, indice, query]);

  return (
    <div className="flex min-h-0 flex-1 flex-col divide-y">
      <nav aria-label="Cátedras" className="min-h-0 flex-1 overflow-y-auto">
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
      </nav>

      <div className="shrink-0 p-3">
        <Input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Filtrar por nombre de docente"
          aria-label="Filtrar cátedras por nombre de docente"
        />
      </div>
    </div>
  );
}
