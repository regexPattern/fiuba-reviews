import type { ReactNode } from "react";

/**
 * Tecla en relieve, como los botones. `chica` es la versión compacta y apagada para atajos
 * secundarios (botón de buscar, pie del buscador).
 *
 * Resetea peso, tracking y mayúsculas para no heredarlos cuando va dentro de un `rotulo`.
 */
export function Tecla({ children, chica = false }: { children: ReactNode; chica?: boolean }) {
  return (
    <kbd
      className={`inline-flex items-center justify-center relieve bg-white font-mono text-xs font-normal tracking-normal normal-case [--relieve-alto:2px] ${
        chica ? "h-5 min-w-5 px-1 text-stone-500" : "h-6 min-w-6 px-1.5 text-stone-700"
      }`}
    >
      {children}
    </kbd>
  );
}
