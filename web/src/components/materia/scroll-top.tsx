"use client";

import { useEffect } from "react";
import { SCROLL_CONTAINER_ID } from "./materia-shell";

/** Al montarse (con `key` distinto por cátedra) lleva el panel principal arriba de todo. */
export function ScrollTop() {
  useEffect(() => {
    // El detalle resuelve el hash cuando termina de llegar por streaming.
    if (window.location.hash) return;
    document.getElementById(SCROLL_CONTAINER_ID)?.scrollTo({ top: 0 });
    window.scrollTo({ top: 0 });
  }, []);

  return null;
}
