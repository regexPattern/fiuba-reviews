"use client";

import { ChevronsUpDown, Close } from "@withicons/react/pixel";
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
      <div ref={shellRef} className="mx-auto flex max-w-6xl md:-mt-(--header-height) md:h-dvh">
        <aside className="hidden min-h-0 w-80 shrink-0 flex-col border-r pt-(--header-height) md:flex">
          {sidebar}
        </aside>

        <div className="flex min-w-0 flex-1 flex-col">
          <div
            ref={barraRef}
            className="sticky top-(--header-height) z-30 border-b bg-white/75 backdrop-blur-xl md:hidden"
          >
            <button
              type="button"
              aria-haspopup="dialog"
              onClick={() => drawerRef.current?.showModal()}
              className="flex w-full items-center gap-3 p-4 text-left"
            >
              <span className="min-w-0 flex-1 font-pixel text-lg wrap-break-word">{titulo}</span>
              <ChevronsUpDown size={16} aria-hidden="true" className="shrink-0 text-neutral-500" />
              <span className="sr-only">Ver cátedras</span>
            </button>
          </div>

          <main
            ref={mainRef}
            className="min-w-0 flex-1 md:min-h-0 md:overflow-y-auto md:pt-(--header-height)"
          >
            {children}
          </main>
        </div>

        {/* oxlint-disable-next-line jsx-a11y/click-events-have-key-events, jsx-a11y/no-noninteractive-element-interactions -- Click en el backdrop; con teclado se cierra con Esc (nativo). */}
        <dialog
          ref={drawerRef}
          aria-label="Cátedras"
          onClick={(event) => event.target === event.currentTarget && cerrarDrawer()}
          className="modal-drawer h-dvh max-h-none w-[85vw] max-w-sm flex-col border-r bg-white pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)] open:flex md:hidden"
        >
          <div className="flex shrink-0 items-center justify-between border-b px-4 py-3">
            <h2 className="font-pixel text-lg">Cátedras</h2>
            <button
              type="button"
              onClick={cerrarDrawer}
              aria-label="Cerrar"
              className="flex size-11 items-center justify-center rounded hover:bg-neutral-100"
            >
              <Close size={16} aria-hidden="true" />
            </button>
          </div>
          <div className="flex min-h-0 flex-1 flex-col">{sidebar}</div>
        </dialog>
      </div>
    </ShellContext>
  );
}
