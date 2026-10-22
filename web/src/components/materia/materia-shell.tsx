"use client";

import { Menu } from "lucide-react";
import { useRouter } from "next/navigation";
import { createContext, useContext, useOptimistic, useState, useTransition } from "react";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger
} from "@/components/ui/sheet";
import { CatedraSkeleton } from "./skeletons";

const ShellContext = createContext<{
  cerrarDrawer: () => void;
  navegarCatedra: (codigoMateria: string, uuid: string) => void;
  catedraPendiente: string | null;
}>({ cerrarDrawer: () => {}, navegarCatedra: () => {}, catedraPendiente: null });
export const useMateriaShell = () => useContext(ShellContext);
export const SCROLL_CONTAINER_ID = "materia-scroll";

/** El sidebar se recibe del servidor para conservar la shell estática. */
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
  const router = useRouter();
  const [navegando, startTransition] = useTransition();
  const [catedraPendiente, setCatedraPendiente] = useOptimistic<string | null>(null);

  function navegarCatedra(codigoMateria: string, uuid: string) {
    setAbierto(false);
    document.getElementById(SCROLL_CONTAINER_ID)?.scrollTo({ top: 0 });
    window.scrollTo({ top: 0 });
    // El panel responde al clic sin esperar el payload RSC. La transición termina
    // cuando Next confirma la navegación; no dependemos de timers ni de efectos sobre la URL.
    startTransition(() => {
      setCatedraPendiente(uuid);
      router.push(`/materia/${codigoMateria}?catedra=${uuid}`, { scroll: false });
    });
  }

  return (
    <ShellContext.Provider
      value={{ cerrarDrawer: () => setAbierto(false), navegarCatedra, catedraPendiente }}
    >
      <div className="mx-auto -mt-(--header-height) flex max-w-7xl md:h-[calc(100%+var(--header-height))]">
        <div className="hidden w-80 shrink-0 border-r pt-(--header-height) md:flex">{sidebar}</div>
        <div className="flex min-w-0 flex-1 flex-col">
          <div className="border-b p-3 pt-[calc(var(--header-height)+0.75rem)] md:hidden">
            <Sheet open={abierto} onOpenChange={setAbierto}>
              <SheetTrigger asChild>
                <Button variant="outline">
                  <Menu /> Ver cátedras
                </Button>
              </SheetTrigger>
              <SheetContent
                side="left"
                className="gap-0"
                onOpenAutoFocus={(event) => event.preventDefault()}
              >
                <SheetHeader>
                  <SheetTitle>Cátedras</SheetTitle>
                  <SheetDescription>{nombreMateria}</SheetDescription>
                </SheetHeader>
                <div className="min-h-0 flex-1">{sidebar}</div>
              </SheetContent>
            </Sheet>
          </div>
          <main
            id={SCROLL_CONTAINER_ID}
            className="flex-1 md:min-h-0 md:overflow-y-auto md:pt-(--header-height)"
          >
            {navegando ? <CatedraSkeleton /> : children}
          </main>
        </div>
      </div>
    </ShellContext.Provider>
  );
}
