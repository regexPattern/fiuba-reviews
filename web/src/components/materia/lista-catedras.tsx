"use client";

import type { CatedraListada } from "@/server/db/queries/catedras-materia";
import { Star } from "@withicons/react/chrome";
import Fuse from "fuse.js";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useMemo, useState } from "react";
import resultadosStyles from "@/components/resultados-busqueda.module.css";
import {
  TextoConCoincidencias,
  rangosExactos,
  type Rangos
} from "@/components/texto-con-coincidencias";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useDebouncedValue } from "@/hooks/use-debounced-value";
import styles from "./lista-catedras.module.css";
import { useMateriaShell } from "./materia-shell";

type Coincidencias = Map<string, Rangos>;

function Items({
  codigoMateria,
  catedras,
  coincidencias,
  activo
}: {
  codigoMateria: string;
  catedras: CatedraListada[];
  coincidencias: Coincidencias;
  activo: string | null;
}) {
  const { cerrarDrawer, navegarCatedra, catedraPendiente } = useMateriaShell();
  const seleccionado = catedraPendiente ?? activo;

  return (
    <ul className="space-y-1 p-2">
      {catedras.map((catedra) => {
        const calificacion = catedra.calificacion.toFixed(1);

        return (
          <li key={catedra.codigo} className={resultadosStyles.entrada}>
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
                <TextoConCoincidencias
                  texto={catedra.nombre}
                  rangos={coincidencias.get(catedra.codigo)}
                />
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
  coincidencias: Coincidencias;
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
  const debouncedQuery = useDebouncedValue(query);
  const indice = useMemo(
    () =>
      new Fuse(catedras, {
        keys: ["nombre"],
        ignoreDiacritics: true,
        ignoreLocation: true,
        threshold: 0.35,
        shouldSort: false,
        includeMatches: true
      }),
    [catedras]
  );
  const { filtradas, coincidencias } = useMemo(() => {
    const palabras = debouncedQuery.trim().split(/\s+/).filter(Boolean);
    const coincidencias: Coincidencias = new Map();
    if (palabras.length === 0) return { filtradas: catedras, coincidencias };

    // Todas las palabras deben coincidir, sin exigir orden ni perder el ranking por calificación.
    const resultados = palabras.map(
      (palabra) =>
        new Map(indice.search(palabra).map((resultado) => [resultado.item.codigo, resultado]))
    );
    const filtradas = catedras.filter((catedra) =>
      resultados.every((resultado) => resultado.has(catedra.codigo))
    );
    for (const catedra of filtradas) {
      coincidencias.set(
        catedra.codigo,
        resultados.flatMap((resultado, i) =>
          rangosExactos(
            catedra.nombre,
            palabras[i]!,
            (resultado.get(catedra.codigo)?.matches ?? []).flatMap((match) => match.indices)
          )
        )
      );
    }
    return { filtradas, coincidencias };
  }, [catedras, indice, debouncedQuery]);

  return (
    <div className="flex min-h-0 flex-1 flex-col divide-y">
      <nav aria-label="Cátedras" className={`min-h-0 flex-1 overflow-y-scroll ${styles.lista}`}>
        {catedras.length === 0 ? (
          <p className="p-3 text-sm text-muted-foreground">
            Esta materia no tiene cátedras todavía.
          </p>
        ) : (
          <Suspense
            fallback={
              <Items
                codigoMateria={codigoMateria}
                catedras={filtradas}
                coincidencias={coincidencias}
                activo={null}
              />
            }
          >
            <ItemsConActiva
              codigoMateria={codigoMateria}
              catedras={filtradas}
              todas={catedras}
              coincidencias={coincidencias}
            />
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
