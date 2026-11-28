"use client";

import { Dialog } from "@base-ui/react/dialog";
import { Close, Search } from "@withicons/react/line";
import { type ComponentProps, useEffect, useRef, useState } from "react";
import { usaCmd, usePlataforma } from "@/lib/plataforma";
import { BuscadorMaterias, type MateriaBuscador } from "./buscador-materias";

export function CommandPalette({ materias }: { materias: MateriaBuscador[] }) {
  const popupRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [abierto, setAbierto] = useState(false);
  const [query, setQuery] = useState("");
  const conCmd = usaCmd(usePlataforma());

  const cerrar = () => setAbierto(false);

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (
        event.key.toLowerCase() === "k" &&
        (conCmd ? event.metaKey : event.ctrlKey) &&
        !event.altKey &&
        !event.shiftKey &&
        !event.isComposing
      ) {
        event.preventDefault();
        if (event.repeat) return;
        setAbierto((abierto) => !abierto);
      }
    }

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [conCmd]);

  return (
    <Dialog.Root
      open={abierto}
      onOpenChange={setAbierto}
      /*
       * La búsqueda se borra al terminar de cerrarse, así el texto no desaparece durante la
       * animación.
       */
      onOpenChangeComplete={(abierto) => !abierto && setQuery("")}
    >
      <Dialog.Trigger
        render={
          <BotonBuscar
            aria-keyshortcuts={conCmd ? "Meta+k" : "Control+k"}
            shortcut={conCmd ? "⌘ K" : "Ctrl K"}
          />
        }
      />
      <Dialog.Portal>
        <Dialog.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black/20 transition-opacity duration-200 ease-out data-ending-style:opacity-0 data-starting-style:opacity-0 supports-[-webkit-touch-callout:none]:absolute motion-reduce:transition-none" />
        <Dialog.Popup
          ref={popupRef}
          aria-label="Buscar materias"
          /*
           * Con mouse o teclado, el foco va al input para escribir directo.
           * Con touch va al popup, porque en iOS el teclado no se abre si el
           * input se enfoca solo al abrir el buscador, y además se rompe el
           * scroll de la lista.
           */
          initialFocus={(tipo) => (tipo === "touch" ? popupRef.current : inputRef.current)}
          className="fixed inset-x-0 top-[max(calc(var(--header-height)+1rem),15dvh)] z-50 mx-auto flex max-h-[75dvh] w-[calc(100%-2rem)] max-w-xl flex-col overflow-hidden relieve bg-white text-stone-900 transition-[scale,opacity] duration-200 ease-out outline-none [--relieve:var(--color-stone-400)] data-ending-style:scale-97 data-ending-style:opacity-0 data-starting-style:scale-97 data-starting-style:opacity-0 motion-reduce:transition-none"
        >
          <BuscadorMaterias
            materias={materias}
            query={query}
            onQueryChange={setQuery}
            onSelect={cerrar}
            inputRef={inputRef}
            inputAction={
              <Dialog.Close
                aria-label="Cerrar buscador"
                className="tecla shrink-0 relieve-boton cursor-pointer text-stone-500 outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fiuba"
              >
                <Close size={14} className="sm:hidden" aria-hidden="true" />
                <span className="hidden sm:inline">Esc</span>
              </Dialog.Close>
            }
          />
          <div
            aria-hidden="true"
            className="hidden shrink-0 items-center gap-4 border-t border-stone-200 bg-stone-50 px-3 py-2 rotulo text-stone-500 sm:flex"
          >
            <span className="flex items-center gap-1.5">
              <kbd className="tecla tecla-chica relieve [--relieve-alto:2px]">↑</kbd>
              <kbd className="tecla tecla-chica relieve [--relieve-alto:2px]">↓</kbd>
              navegar
            </span>
            <span className="flex items-center gap-1.5">
              <kbd className="tecla tecla-chica relieve [--relieve-alto:2px]">↵</kbd>
              abrir
            </span>
            <span className="ml-auto flex items-center gap-1.5">
              <kbd className="tecla tecla-chica relieve [--relieve-alto:2px]">Esc</kbd>
              cerrar
            </span>
          </div>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

export function BotonBuscar({
  shortcut,
  ...props
}: ComponentProps<"button"> & { shortcut?: string }) {
  return (
    <button
      type="button"
      aria-label="Buscar materias"
      aria-haspopup="dialog"
      className="boton-terciario h-9 disabled:pointer-events-none disabled:opacity-50"
      {...props}
    >
      <Search size={18} aria-hidden="true" className="text-stone-500" />
      <span>
        Buscar<span className="hidden sm:inline"> materias</span>
      </span>
      {shortcut && (
        <span className="ml-1 hidden sm:flex">
          <kbd className="tecla tecla-chica relieve [--relieve-alto:2px]">{shortcut}</kbd>
        </span>
      )}
    </button>
  );
}
