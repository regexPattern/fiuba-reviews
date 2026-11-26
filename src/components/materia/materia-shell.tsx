"use client";

import { Drawer } from "@base-ui/react/drawer";
import { ChevronsUpDown, Close } from "@withicons/react/line";
import {
  type ReactNode,
  createContext,
  useContext,
  useLayoutEffect,
  useRef,
  useState
} from "react";

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
  const [drawerAbierto, setDrawerAbierto] = useState(false);

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
    setDrawerAbierto(false);
    // En móvil scrollea la ventana y en desktop el panel principal.
    mainRef.current?.scrollTo({ top: 0 });
    window.scrollTo({ top: 0 });
  }

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
              <Drawer.Root
                open={drawerAbierto}
                onOpenChange={setDrawerAbierto}
                swipeDirection="left"
              >
                <Drawer.Trigger className="flex w-full items-center gap-3 px-4 py-3 text-left">
                  <span className="min-w-0 flex-1 text-sm font-medium text-stone-900">
                    {titulo}
                  </span>
                  <ChevronsUpDown
                    size={18}
                    aria-hidden="true"
                    className="shrink-0 text-stone-500"
                  />
                  <span className="sr-only">Ver cátedras</span>
                </Drawer.Trigger>
                <Drawer.Portal>
                  {/* `z-50`: por encima del header fijo (`z-40`). Se aclara a medida que se arrastra el drawer. */}
                  <Drawer.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black opacity-[calc(0.2*(1-var(--drawer-swipe-progress)))] transition-opacity duration-300 ease-out data-ending-style:opacity-0 data-ending-style:duration-[calc(var(--drawer-swipe-strength)*300ms)] data-starting-style:opacity-0 data-swiping:duration-0 supports-[-webkit-touch-callout:none]:absolute motion-reduce:transition-none md:hidden" />
                  <Drawer.Viewport className="fixed inset-0 z-50 flex md:hidden">
                    <Drawer.Popup className="flex h-full w-[85vw] max-w-sm [transform:translateX(var(--drawer-swipe-movement-x))] flex-col border-r border-stone-300 bg-white pt-[env(safe-area-inset-top)] pb-[env(safe-area-inset-bottom)] transition-transform duration-300 ease-out outline-none data-ending-style:[transform:translateX(-100%)] data-ending-style:duration-[calc(var(--drawer-swipe-strength)*300ms)] data-starting-style:[transform:translateX(-100%)] data-swiping:duration-0 data-swiping:select-none motion-reduce:transition-none">
                      <Drawer.Content className="flex min-h-0 flex-1 flex-col">
                        <div className="flex shrink-0 items-center justify-between border-b border-stone-300 py-1 pr-1 pl-4">
                          <Drawer.Title className="text-base tracking-tight text-stone-900">
                            Cátedras
                          </Drawer.Title>
                          <Drawer.Close
                            aria-label="Cerrar"
                            className="flex size-11 items-center justify-center text-stone-500 hover:bg-stone-100 hover:text-stone-900"
                          >
                            <Close size={18} aria-hidden="true" />
                          </Drawer.Close>
                        </div>
                        <div className="flex min-h-0 flex-1 flex-col">{sidebar}</div>
                      </Drawer.Content>
                    </Drawer.Popup>
                  </Drawer.Viewport>
                </Drawer.Portal>
              </Drawer.Root>
            </div>

            <main
              ref={mainRef}
              className="min-w-0 flex-1 md:mt-(--header-height) md:min-h-0 md:overflow-y-auto"
            >
              {children}
            </main>
          </div>
        </div>
      </div>
    </ShellContext>
  );
}
