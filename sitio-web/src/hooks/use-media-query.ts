import { useCallback, useSyncExternalStore } from "react";

/**
 * Hook para saber si la pantalla cumple la media query `query`. Durante el
 * SSR devuelve `false`, así que no conviene usarlo para lo que se ve al
 * cargar la página, porque cambiaría al hidratar.
 */
export function useMediaQuery(query: string) {
  const suscribir = useCallback(
    (onChange: () => void) => {
      const lista = window.matchMedia(query);
      lista.addEventListener("change", onChange);
      return () => lista.removeEventListener("change", onChange);
    },
    [query]
  );

  return useSyncExternalStore(
    suscribir,
    () => window.matchMedia(query).matches,
    () => false
  );
}
