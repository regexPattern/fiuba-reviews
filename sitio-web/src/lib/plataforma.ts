import { useSyncExternalStore } from "react";

export type Plataforma = "pc" | "mac" | "android" | "ios";

/**
 * Detecta la plataforma del usuario según el user agent.
 * Los iPad se detectan como Mac, porque desde iPadOS 13 mandan el mismo user agent.
 */
export function detectarPlataforma(userAgent: string): Plataforma {
  if (/Android/i.test(userAgent)) return "android";
  if (/iPhone|iPod/i.test(userAgent)) return "ios";
  if (/Macintosh|iPad/i.test(userAgent)) return "mac";
  return "pc";
}

/** Las plataformas donde los atajos usan ⌘ en vez de Ctrl. */
export const usaCmd = (plataforma: Plataforma) => plataforma === "mac" || plataforma === "ios";

const subscribe = () => () => {};

const plataformaSnapshot = () => detectarPlataforma(navigator.userAgent);

const plataformaServerSnapshot = (): Plataforma => "pc";

/**
 * Hook para detectar la plataforma del usuario al hidratar la página. Durante
 * el SSR devuelve `"pc"` por defecto, así el primer frame ya tiene un valor y
 * no `undefined`.
 */
export function usePlataforma() {
  return useSyncExternalStore(subscribe, plataformaSnapshot, plataformaServerSnapshot);
}
