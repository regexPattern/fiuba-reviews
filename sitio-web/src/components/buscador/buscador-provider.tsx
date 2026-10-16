"use client";

import type { MateriaBuscador } from "@/server/queries/materias-vigentes";
import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";

type BuscadorContext = {
  materias: MateriaBuscador[];
  abierto: boolean;
  setAbierto: (abierto: boolean) => void;
  abrir: () => void;
};

const Context = createContext<BuscadorContext | null>(null);

export function useBuscador() {
  const ctx = useContext(Context);
  if (!ctx) throw new Error("useBuscador debe usarse dentro de <BuscadorProvider>");
  return ctx;
}

/**
 * Los puntos de acceso al buscador (navbar, hero) son solo triggers; el diálogo se monta una sola
 * vez y comparte el estado a través de este contexto.
 */
export function BuscadorProvider({
  materias,
  children
}: {
  materias: MateriaBuscador[];
  children: React.ReactNode;
}) {
  const [abierto, setAbierto] = useState(false);
  const abrir = useCallback(() => setAbierto(true), []);

  useEffect(() => {
    function onKeyDown(evento: KeyboardEvent) {
      if (evento.key.toLowerCase() === "k" && (evento.metaKey || evento.ctrlKey)) {
        evento.preventDefault();
        setAbierto((prev) => !prev);
      }
    }

    document.addEventListener("keydown", onKeyDown);
    return () => document.removeEventListener("keydown", onKeyDown);
  }, []);

  const value = useMemo(
    () => ({ materias, abierto, setAbierto, abrir }),
    [materias, abierto, abrir]
  );

  return <Context.Provider value={value}>{children}</Context.Provider>;
}
