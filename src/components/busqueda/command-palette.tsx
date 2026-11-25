"use client";

import { Close, Search } from "@withicons/react/pixel";
import { type ComponentProps, useEffect, useRef, useState, useSyncExternalStore } from "react";
import { BuscadorMaterias, type MateriaBuscador } from "./buscador-materias";

const subscribePlatform = () => () => {};
const getIsMac = () => /Mac|iPhone|iPad|iPod/i.test(navigator.platform);
const getServerIsMac = () => false;

export function CommandPalette({ materias }: { materias: MateriaBuscador[] }) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const [query, setQuery] = useState("");
  const isMac = useSyncExternalStore(subscribePlatform, getIsMac, getServerIsMac);

  function abrir() {
    const dialog = dialogRef.current;
    if (!dialog || dialog.open) return;
    dialog.showModal();
    dialog.querySelector("input")?.focus();
  }

  const cerrar = () => dialogRef.current?.close();

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
        if (event.repeat) return;
        if (dialogRef.current?.open) cerrar();
        else abrir();
      }
    }

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <>
      <BotonBuscar
        onClick={abrir}
        aria-keyshortcuts={isMac ? "Meta+k" : "Control+k"}
        atajo={isMac ? "⌘ K" : "Ctrl K"}
      />
      {/* oxlint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions -- Click en el backdrop; con teclado se cierra con Esc (nativo). */}
      <dialog
        ref={dialogRef}
        aria-label="Buscar materias"
        onClose={() => setQuery("")}
        // Click en el backdrop: el contenido ocupa todo el diálogo, así que el target solo es
        // el propio diálogo cuando se clickea afuera.
        onClick={(event) => event.target === event.currentTarget && cerrar()}
        className="inset-x-0 top-[max(calc(var(--header-height)+1rem),15dvh)] bottom-auto mx-auto max-h-[75dvh] w-[calc(100%-2rem)] max-w-xl scale-97 flex-col overflow-hidden border bg-white opacity-0 shadow-xl transition-all transition-discrete duration-200 ease-out backdrop:bg-black/0 backdrop:transition-all backdrop:transition-discrete backdrop:duration-200 open:flex open:scale-100 open:opacity-100 open:backdrop:bg-black/20 motion-reduce:transition-none starting:open:scale-97 starting:open:opacity-0 starting:open:backdrop:bg-black/0"
      >
        <BuscadorMaterias
          materias={materias}
          query={query}
          onQueryChange={setQuery}
          onSelect={cerrar}
          inputAction={
            <button
              type="button"
              onClick={cerrar}
              aria-label="Cerrar buscador"
              className="flex h-6 shrink-0 items-center border px-1.5 text-xs text-neutral-500 hover:bg-neutral-100"
            >
              <Close size={12} className="sm:hidden" aria-hidden="true" />
              <span className="hidden sm:inline">Esc</span>
            </button>
          }
        />
        <div
          aria-hidden="true"
          className="hidden shrink-0 gap-4 border-t px-3 py-2 text-xs text-neutral-500 sm:flex"
        >
          <span>↑ ↓ navegar</span>
          <span>↵ abrir</span>
          <span className="ml-auto">Esc cerrar</span>
        </div>
      </dialog>
    </>
  );
}

export function BotonBuscar({ atajo, ...props }: ComponentProps<"button"> & { atajo?: string }) {
  return (
    <button
      type="button"
      aria-label="Buscar materias"
      aria-haspopup="dialog"
      className="flex h-11 items-center gap-2 px-3 text-sm hover:bg-neutral-100 disabled:opacity-50 sm:h-9"
      {...props}
    >
      <Search size={16} aria-hidden="true" />
      <span>
        Buscar<span className="hidden sm:inline"> materias</span>
      </span>
      {atajo && (
        <kbd className="hidden border px-1.5 py-0.5 font-mono text-xs text-neutral-500 sm:inline">
          {atajo}
        </kbd>
      )}
    </button>
  );
}
