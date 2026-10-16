"use client";

import { Search } from "lucide-react";
import { cn } from "@/lib/utils";
import { useBuscador } from "./buscador-provider";

export function BuscadorTrigger({
  variante = "navbar",
  className
}: {
  variante?: "navbar" | "hero";
  className?: string;
}) {
  const { abrir } = useBuscador();

  return (
    <button
      type="button"
      onClick={abrir}
      aria-label="Abrir buscador de materias"
      className={cn(
        "flex items-center gap-2 rounded-full border border-button-border bg-button-background text-sm text-button-foreground transition-colors hover:bg-button-hover",
        variante === "hero" ? "px-4 py-2.5" : "px-3 py-2",
        className
      )}
    >
      {variante === "hero" ? (
        <>
          <Search className="size-4" />
          <span>Buscar materias</span>
        </>
      ) : (
        <>
          <span className="hidden items-center gap-2 md:flex">
            <Search className="size-4" />
            <span>Buscar materias</span>
          </span>
          <span className="flex items-center gap-2 md:hidden">
            <span>Buscar</span>
            <Search className="size-4" />
          </span>
        </>
      )}
      <span className="hidden md:inline">⌘K</span>
    </button>
  );
}
