"use client";

import type { MateriaBuscador } from "@/lib/buscar-materias";
import { SearchIcon } from "lucide-react";
import { Dialog } from "radix-ui";
import { useEffect, useRef, useState, useSyncExternalStore } from "react";
import { BuscadorMaterias } from "@/components/materia/buscador-materias";
import { Button } from "@/components/ui/button";

const subscribePlatform = () => () => {};
const getIsMac = () => /Mac|iPhone|iPad|iPod/i.test(navigator.platform);
const getServerIsMac = () => false;

export function CommandPalette({ materias }: { materias: MateriaBuscador[] }) {
  const contentRef = useRef<HTMLDivElement>(null);
  const [open, setOpen] = useState(false);
  const isMac = useSyncExternalStore(subscribePlatform, getIsMac, getServerIsMac);

  useEffect(() => {
    const mac = getIsMac();

    function handleKeyDown(event: KeyboardEvent) {
      if (
        event.key.toLowerCase() === "k" &&
        (mac ? event.metaKey : event.ctrlKey) &&
        !event.altKey &&
        !event.shiftKey &&
        !event.isComposing
      ) {
        event.preventDefault();
        if (!event.repeat) setOpen((previous) => !previous);
      }
    }

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <Dialog.Root open={open} onOpenChange={setOpen}>
      <Dialog.Trigger asChild>
        <Button
          variant="ghost"
          className="h-11 px-3 sm:h-9"
          aria-label="Buscar materias"
          aria-keyshortcuts={isMac ? "Meta+k" : "Control+k"}
        >
          <SearchIcon aria-hidden="true" />
          <span>
            Buscar<span className="hidden sm:inline"> materias</span>
          </span>
          <kbd className="hidden rounded-sm border px-1.5 py-0.5 text-xs text-muted-foreground sm:inline">
            {isMac ? "⌘ K" : "Ctrl K"}
          </kbd>
        </Button>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 bg-background/80 motion-safe:data-[state=closed]:animate-out motion-safe:data-[state=closed]:fade-out-0 motion-safe:data-[state=open]:animate-in motion-safe:data-[state=open]:fade-in-0" />
        <Dialog.Content
          ref={contentRef}
          onOpenAutoFocus={(event) => {
            event.preventDefault();
            contentRef.current?.querySelector<HTMLInputElement>("input")?.focus();
          }}
          className="command-palette fixed left-1/2 z-50 flex max-w-xl -translate-x-1/2 flex-col overflow-hidden rounded-md border bg-popover/75 text-popover-foreground shadow-glass backdrop-blur-surface outline-none"
        >
          <Dialog.Title className="sr-only">Buscar materias</Dialog.Title>
          <Dialog.Description className="sr-only">
            Buscá por nombre o código. Usá las flechas para elegir una materia y Enter para abrirla.
          </Dialog.Description>
          <BuscadorMaterias
            materias={materias}
            onSelect={() => setOpen(false)}
            inputAction={
              <Dialog.Close asChild>
                <Button variant="ghost" size="xs" className="border" aria-label="Cerrar buscador">
                  Esc
                </Button>
              </Dialog.Close>
            }
            className="min-h-0 border-0 bg-transparent [&_[cmdk-list]]:min-h-0"
          />
          <div
            aria-hidden="true"
            className="flex shrink-0 gap-4 border-t px-3 py-2 text-xs text-muted-foreground"
          >
            <span>↑ ↓ navegar</span>
            <span>↵ abrir</span>
            <span className="ml-auto">Esc cerrar</span>
          </div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
