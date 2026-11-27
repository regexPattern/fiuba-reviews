import { useCallback, useSyncExternalStore } from "react";

/**
 * Si la media query coincide. En el servidor no se sabe, así que da `false` hasta hidratar: usarlo
 * solo para lo que no se ve en el primer render (por ejemplo, el contenido de un popup cerrado).
 */
export function useMediaQuery(query: string) {
  const suscribir = useCallback(
    (alCambiar: () => void) => {
      const lista = window.matchMedia(query);
      lista.addEventListener("change", alCambiar);
      return () => lista.removeEventListener("change", alCambiar);
    },
    [query]
  );

  return useSyncExternalStore(
    suscribir,
    () => window.matchMedia(query).matches,
    () => false
  );
}
