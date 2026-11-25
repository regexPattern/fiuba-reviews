import type { ReactNode } from "react";

/** Rótulo chico en versalitas y en relieve, como las cátedras del listado (códigos, roles, marcas). */
export function Etiqueta({ children, className }: { children: ReactNode; className?: string }) {
  return (
    <span
      className={`shrink-0 relieve bg-white px-1.5 py-px rotulo whitespace-nowrap text-stone-600 [--relieve-alto:2px] ${className ?? ""}`}
    >
      {children}
    </span>
  );
}
