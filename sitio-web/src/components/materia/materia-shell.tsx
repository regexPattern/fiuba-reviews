"use client";

import { createContext, useContext, useState } from "react";
import { ChevronDown, ChevronUp } from "lucide-react";
import { Sheet, SheetContent, SheetDescription, SheetTitle } from "@/components/ui/sheet";

const ShellContext = createContext<{ cerrarDrawer: () => void }>({ cerrarDrawer: () => {} });

export const useMateriaShell = () => useContext(ShellContext);

export const SCROLL_CONTAINER_ID = "materia-scroll";

/**
 * Layout de dos columnas: sidebar fijo en desktop y drawer (Sheet) en mobile. El sidebar llega
 * como prop, así que sus Server Components siguen siendo parte de la shell estática.
 */
export function MateriaShell({
  nombreMateria,
  sidebar,
  children
}: {
  nombreMateria: string;
  sidebar: React.ReactNode;
  children: React.ReactNode;
}) {
  const [abierto, setAbierto] = useState(false);

  return (
    <ShellContext.Provider value={{ cerrarDrawer: () => setAbierto(false) }}>
      <div className="container mx-auto mt-[calc(-56px-env(safe-area-inset-top))] flex h-dvh overflow-hidden">
        <div className="hidden w-70 shrink-0 md:flex">{sidebar}</div>

        <Sheet open={abierto} onOpenChange={setAbierto}>
          <SheetContent side="left" showCloseButton={false} className="w-4/5 gap-0 md:hidden">
            <SheetTitle className="sr-only">{nombreMateria}</SheetTitle>
            <SheetDescription className="sr-only">Cátedras de la materia</SheetDescription>
            {sidebar}
          </SheetContent>
        </Sheet>

        <main className="min-h-0 w-full min-w-0">
          <div
            id={SCROLL_CONTAINER_ID}
            className="h-full overflow-y-auto overscroll-contain pt-[calc(56px+env(safe-area-inset-top))]"
          >
            <button
              type="button"
              onClick={() => setAbierto(true)}
              className="sticky top-0 z-20 flex w-full items-center justify-between border-b border-layout-border bg-background p-3 text-left font-serif text-lg font-medium md:hidden"
            >
              {nombreMateria}
              {abierto ? (
                <ChevronUp className="size-4.5 shrink-0" />
              ) : (
                <ChevronDown className="size-4.5 shrink-0" />
              )}
            </button>

            {children}
          </div>
        </main>
      </div>
    </ShellContext.Provider>
  );
}
