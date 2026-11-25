"use client";

import { Close, Search } from "@withicons/react/plush";
import {
  type ComponentProps,
  type ReactNode,
  useEffect,
  useRef,
  useState,
  useSyncExternalStore
} from "react";
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
        className="inset-x-0 top-[max(calc(var(--header-height)+1rem),15dvh)] bottom-auto mx-auto max-h-[75dvh] w-[calc(100%-2rem)] max-w-xl scale-97 flex-col overflow-hidden relieve bg-white text-tinta opacity-0 transition-all transition-discrete duration-200 ease-out backdrop:bg-black/0 backdrop:transition-all backdrop:transition-discrete backdrop:duration-200 open:flex open:scale-100 open:opacity-100 open:backdrop:bg-black/20 motion-reduce:transition-none starting:open:scale-97 starting:open:opacity-0 starting:open:backdrop:bg-black/0"
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
              // Tecla chica en relieve, como los atajos del pie.
              className="flex h-6 shrink-0 relieve-boton cursor-pointer items-center bg-white px-1.5 font-numeros text-xs text-stone-500 outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fiuba"
            >
              <Close size={14} className="sm:hidden" aria-hidden="true" />
              <span className="hidden sm:inline">Esc</span>
            </button>
          }
        />
        {/* Pie apagado, como el del resumen con IA. */}
        <div
          aria-hidden="true"
          className="hidden shrink-0 items-center gap-4 border-t border-stone-200 bg-stone-50 px-3 py-2 rotulo text-stone-500 sm:flex"
        >
          <span className="flex items-center gap-1.5">
            <Tecla>↑</Tecla>
            <Tecla>↓</Tecla>
            navegar
          </span>
          <span className="flex items-center gap-1.5">
            <Tecla>↵</Tecla>
            abrir
          </span>
          <span className="ml-auto flex items-center gap-1.5">
            <Tecla>Esc</Tecla>
            cerrar
          </span>
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
          <Tecla>{atajo}</Tecla>
        </span>
      )}
    </button>
  );
}

/** Tecla chica en relieve bajo, para atajos de teclado. */
function Tecla({ children }: { children: ReactNode }) {
  return (
    <kbd className="inline-flex h-5 min-w-5 items-center justify-center relieve bg-white px-1 font-numeros text-xs font-normal tracking-normal text-stone-500 normal-case [--relieve-alto:2px]">
      {children}
    </kbd>
  );
}
