"use client";

import { Dialog } from "@base-ui/react/dialog";
import { Close, Search } from "@withicons/react/line";
import { type ComponentProps, useEffect, useRef, useState } from "react";
import { Tecla } from "@/components/ui/tecla";
import { usaCmd, usePlataforma } from "@/lib/plataforma";
import { BuscadorMaterias, type MateriaBuscador } from "./buscador-materias";

export function CommandPalette({ materias }: { materias: MateriaBuscador[] }) {
  const popupRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const [abierto, setAbierto] = useState(false);
  const [query, setQuery] = useState("");
  const isMac = usaCmd(usePlataforma());

  const cerrar = () => setAbierto(false);

  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (
        event.key.toLowerCase() === "k" &&
        (isMac ? event.metaKey : event.ctrlKey) &&
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
  }, [isMac]);

  return (
    <Dialog.Root
      open={abierto}
      onOpenChange={setAbierto}
      // Se limpia al terminar la animación de cierre, así el texto no desaparece mientras se va.
      onOpenChangeComplete={(abierto) => !abierto && setQuery("")}
    >
      <Dialog.Trigger
        render={
          <BotonBuscar
            aria-keyshortcuts={isMac ? "Meta+k" : "Control+k"}
            atajo={isMac ? "⌘ K" : "Ctrl K"}
          />
        }
      />
      <Dialog.Portal>
        {/* `z-50`: por encima del header fijo (`z-40`). */}
        <Dialog.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black/20 transition-opacity duration-200 ease-out data-ending-style:opacity-0 data-starting-style:opacity-0 supports-[-webkit-touch-callout:none]:absolute motion-reduce:transition-none" />
        <Dialog.Popup
          ref={popupRef}
          aria-label="Buscar materias"
          // Al input, salvo con touch: ahí iOS no abre el teclado con un foco por código, y con el input
          // enfocado el gesto scrollea la página en vez de la lista. Con touch se enfoca el popup.
          initialFocus={(tipo) => (tipo === "touch" ? popupRef.current : inputRef.current)}
          className="fixed inset-x-0 top-[max(calc(var(--header-height)+1rem),15dvh)] z-50 mx-auto flex max-h-[75dvh] w-[calc(100%-2rem)] max-w-xl flex-col overflow-hidden relieve bg-white text-stone-900 transition-[scale,opacity] duration-200 ease-out outline-none [--relieve:var(--color-stone-500)] data-ending-style:scale-97 data-ending-style:opacity-0 data-starting-style:scale-97 data-starting-style:opacity-0 motion-reduce:transition-none"
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
                // Tecla chica en relieve, como los atajos del pie.
                className="flex h-6 shrink-0 relieve-boton cursor-pointer items-center bg-white px-1.5 font-mono text-xs text-stone-500 outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fiuba"
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
              <Tecla chica>↑</Tecla>
              <Tecla chica>↓</Tecla>
              navegar
            </span>
            <span className="flex items-center gap-1.5">
              <Tecla chica>↵</Tecla>
              abrir
            </span>
            <span className="ml-auto flex items-center gap-1.5">
              <Tecla chica>Esc</Tecla>
              cerrar
            </span>
          </div>
        </Dialog.Popup>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

export function BotonBuscar({ atajo, ...props }: ComponentProps<"button"> & { atajo?: string }) {
  return (
    <button
      type="button"
      aria-label="Buscar materias"
      aria-haspopup="dialog"
      // Terciario, sin borde ni relieve: el atajo ya va en una tecla en relieve y dos cajas competirían.
      className="boton-terciario h-9 disabled:pointer-events-none disabled:opacity-50"
      {...props}
    >
      <Search size={18} aria-hidden="true" className="text-stone-500" />
      <span>
        Buscar<span className="hidden sm:inline"> materias</span>
      </span>
      {atajo && (
        <span className="ml-1 hidden sm:flex">
          <Tecla chica>{atajo}</Tecla>
        </span>
      )}
    </button>
  );
}
