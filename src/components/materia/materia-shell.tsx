"use client";

import { ChevronsUpDown, Close } from "@withicons/react/line";
import { type ReactNode, createContext, useContext, useLayoutEffect, useRef } from "react";

const ShellContext = createContext<{ alNavegar: () => void }>({ alNavegar: () => {} });
export const useMateriaShell = () => useContext(ShellContext);

/**
 * Layout de una materia. En desktop, sidebar fijo y panel principal con su propio scroll.
 * En móvil scrollea el documento, así la barra flotante de Safari (iOS 26) deja ver el contenido
 * por debajo. El sidebar pasa a un drawer que se abre desde una barra sticky.
 */
export function MateriaShell({
  titulo,
  sidebar,
  children
}: {
  titulo: ReactNode;
  sidebar: ReactNode;
  children: ReactNode;
}) {
  const shellRef = useRef<HTMLDivElement>(null);
  const barraRef = useRef<HTMLDivElement>(null);
  const mainRef = useRef<HTMLElement>(null);
  const drawerRef = useRef<HTMLDialogElement>(null);

  useLayoutEffect(() => {
    const shell = shellRef.current;
    const barra = barraRef.current;
    if (!shell || !barra) return;

    // Para el `scroll-margin` de los anchors a docentes. Cambia con el ancho, el texto y las fuentes.
    const actualizar = () =>
      shell.style.setProperty(
        "--materia-sticky-height",
        `${barra.getBoundingClientRect().height}px`
      );
    actualizar();
    const observer = new ResizeObserver(actualizar);
    observer.observe(barra);
    return () => observer.disconnect();
  }, []);

  function alNavegar() {
    drawerRef.current?.close();
    // En móvil scrollea la ventana y en desktop el panel principal.
    mainRef.current?.scrollTo({ top: 0 });
    window.scrollTo({ top: 0 });
  }

  const cerrarDrawer = () => drawerRef.current?.close();

  return (
    <ShellContext value={{ alNavegar }}>
      {/* Fondo blanco a todo el ancho, igual que el header. */}
      <div className="min-h-[calc(100dvh-var(--header-height))] bg-white md:-mt-(--header-height) md:min-h-0">
        <div ref={shellRef} className="mx-auto flex max-w-6xl md:h-dvh">
          <aside className="hidden min-h-0 w-80 shrink-0 flex-col border-r border-stone-300 pt-(--header-height) md:flex">
            {sidebar}
          </aside>

          <div className="flex min-w-0 flex-1 flex-col">
            <div
              ref={barraRef}
              className="sticky top-(--header-height) z-30 border-b border-stone-300 bg-white/75 backdrop-blur-xl md:hidden"
            >
              <button
                type="button"
                aria-haspopup="dialog"
                onClick={() => drawerRef.current?.showModal()}
                className="flex w-full items-center gap-3 px-4 py-3 text-left"
              >
                <span className="min-w-0 flex-1 text-sm font-medium text-stone-900">{titulo}</span>
                <ChevronsUpDown size={18} aria-hidden="true" className="shrink-0 text-stone-500" />
                <span className="sr-only">Ver cátedras</span>
              </button>
            </div>

            <main
              ref={mainRef}
              className="min-w-0 flex-1 md:mt-(--header-height) md:min-h-0 md:overflow-y-auto"
            >
              {children}
            </main>
          </div>

          {/* oxlint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions -- Click en el backdrop; con teclado se cierra con Esc (nativo). */}
          <dialog
            ref={drawerRef}
            aria-label="Cátedras"
            onClick={(event) => event.target === event.currentTarget && cerrarDrawer()}
            className="inset-y-0 right-auto left-0 m-0 h-dvh max-h-none w-[85vw] max-w-sm [translate:-100%] flex-col border-r border-stone-300 bg-white pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)] transition-all transition-discrete duration-300 ease-out backdrop:bg-black/0 backdrop:transition-all backdrop:transition-discrete backdrop:duration-200 open:flex open:[translate:0] open:backdrop:bg-black/20 motion-reduce:transition-none md:hidden starting:open:[translate:-100%] starting:open:backdrop:bg-black/0"
          >
            <div className="flex shrink-0 items-center justify-between border-b border-stone-300 py-1 pr-1 pl-4">
              <h2 className="text-base tracking-tight text-stone-900">Cátedras</h2>
              <button
                type="button"
                onClick={cerrarDrawer}
                aria-label="Cerrar"
                className="flex size-11 items-center justify-center text-stone-500 hover:bg-stone-100 hover:text-stone-900"
              >
                <Close size={18} aria-hidden="true" />
              </button>
            </div>
            <div className="flex min-h-0 flex-1 flex-col">{sidebar}</div>
          </dialog>
        </div>
      </div>
    </ShellContext>
  );
}
