import { useSyncExternalStore } from "react";

export type Plataforma = "pc" | "mac" | "android" | "ios";

/**
 * Plataforma del usuario según el user agent.
 * El iPad se toma como Mac: en modo escritorio (lo predeterminado desde iPadOS 13) su UA ya es el de
 * una Mac, y así se trata igual cuando pide la versión móvil.
 */
export function detectarPlataforma(userAgent: string): Plataforma {
  if (/Android/i.test(userAgent)) return "android";
  if (/iPhone|iPod/i.test(userAgent)) return "ios";
  if (/Macintosh|iPad/i.test(userAgent)) return "mac";
  return "pc";
}

/** Las plataformas donde los atajos usan ⌘ en vez de Ctrl. */
export const usaCmd = (plataforma: Plataforma) => plataforma === "mac" || plataforma === "ios";

const suscribir = () => () => {};

const plataformaDelNavegador = () => detectarPlataforma(navigator.userAgent);

/** En el servidor no se sabe así que se corrige al hidratar. */
const plataformaDelServidor = (): Plataforma => "pc";

/** Plataforma detectada del navegador. */
export function usePlataforma() {
  return useSyncExternalStore(suscribir, plataformaDelNavegador, plataformaDelServidor);
}
