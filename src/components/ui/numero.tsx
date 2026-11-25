import type { ReactNode } from "react";

/** Números con la fuente de números y ancho fijo, para que las columnas queden alineadas. */
export function Numero({ children }: { children: ReactNode }) {
  return <span className="font-numeros tabular-nums">{children}</span>;
}
