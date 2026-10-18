"use client";

import { useEffect } from "react";
import { SCROLL_CONTAINER_ID } from "./materia-shell";

/** Al montarse (con `key` distinto por cátedra) lleva el panel principal arriba de todo. */
export function ScrollTop() {
  useEffect(() => {
    document.getElementById(SCROLL_CONTAINER_ID)?.scrollTo({ top: 0 });
  }, []);

  return null;
}
