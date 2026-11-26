import { useSyncExternalStore } from "react";

export type Plataforma = "windows" | "mac" | "android" | "ios";

/** Lo que hace falta del `navigator` para detectar la plataforma (aparte, para poder testearlo). */
export type DatosNavegador = {
  userAgent: string;
  /** `navigator.userAgentData.platform` si existe (Chromium), si no `navigator.platform`. */
  platform: string;
  maxTouchPoints: number;
};

/**
 * Plataforma del usuario según el navegador. El iPad con iPadOS 13+ se hace pasar por una Mac
 * (`MacIntel`), así que una "Mac" con pantalla táctil se toma como iOS.
 */
export function detectarPlataforma({
  userAgent,
  platform,
  maxTouchPoints
}: DatosNavegador): Plataforma {
  if (/Android/i.test(userAgent)) return "android";
  if (/iPhone|iPad|iPod/i.test(userAgent) || /iPhone|iPad|iPod/i.test(platform)) return "ios";
  if (/Mac/i.test(platform)) return maxTouchPoints > 1 ? "ios" : "mac";
  return "windows";
}

/** Las plataformas donde los atajos usan ⌘ en vez de Ctrl (el iPad con teclado incluido). */
export const usaCmd = (plataforma: Plataforma) => plataforma === "mac" || plataforma === "ios";

const suscribir = () => () => {};

function plataformaDelNavegador() {
  const { userAgentData } = navigator as Navigator & { userAgentData?: { platform: string } };

  return detectarPlataforma({
    userAgent: navigator.userAgent,
    platform: userAgentData?.platform || navigator.platform,
    maxTouchPoints: navigator.maxTouchPoints
  });
}

/** En el servidor no se sabe: se asume Windows/Linux y se corrige al hidratar. */
const plataformaDelServidor = (): Plataforma => "windows";

/** Plataforma detectada del navegador. */
export function usePlataforma() {
  return useSyncExternalStore(suscribir, plataformaDelNavegador, plataformaDelServidor);
}
