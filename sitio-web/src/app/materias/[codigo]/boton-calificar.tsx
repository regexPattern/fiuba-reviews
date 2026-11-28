"use client";

import type { Cuatrimestre } from "@/lib/cuatrimestres";
import { Dialog } from "@base-ui/react/dialog";
import { Drawer } from "@base-ui/react/drawer";
import { Close, PenLine } from "@withicons/react/line";
import { type ElementType, type ReactNode, useRef, useState, useTransition } from "react";
import { enviarCalificacion } from "@/app/materias/actions";
import { useMediaQuery } from "@/hooks/use-media-query";
import { type BorradorCalificacion, FormularioCalificacion } from "./formulario-calificacion";

const BORRADOR_VACIO: BorradorCalificacion = {
  puntajes: {},
  cuatrimestre: null,
  comentario: ""
};

/**
 * Botón que abre el panel para calificar al docente: un dialog en desktop y un drawer en móvil. El
 * borrador se conserva al cerrar el panel y solo se borra después de enviarlo.
 */
export function BotonCalificar({
  docente,
  cuatrimestres,
  className
}: {
  docente: { codigo: string; nombre: string };
  cuatrimestres: Cuatrimestre[];
  className?: string;
}) {
  const botonRef = useRef<HTMLButtonElement>(null);
  const [abierto, setAbierto] = useState(false);
  const [borrador, setBorrador] = useState(BORRADOR_VACIO);
  const [error, setError] = useState<string | null>(null);
  const [enviando, startTransition] = useTransition();
  const enviado = useRef(false);
  /*
   * `48rem` es el breakpoint `md` de Tailwind, el mismo que usa el resto del sitio. El `false` del
   * SSR no se llega a ver, porque el panel se abre recién después de hidratar.
   */
  const esDesktop = useMediaQuery("(min-width: 48rem)");

  function handleSubmit(calificacion: BorradorCalificacion) {
    setError(null);
    startTransition(async () => {
      const resultado = await enviarCalificacion({
        codigoDocente: docente.codigo,
        ...calificacion
      });

      if (resultado.estado === "error") {
        setError(resultado.mensaje);
        return;
      }
      enviado.current = true;
      setAbierto(false);
    });
  }

  const rootProps = {
    open: abierto,
    onOpenChange: setAbierto,
    /*
     * Después de enviar, el borrador se borra al terminar de cerrarse, así el
     * formulario no se vacía durante la animación.
     */
    onOpenChangeComplete: (abierto: boolean) => {
      if (!abierto && enviado.current) {
        enviado.current = false;
        setBorrador(BORRADOR_VACIO);
        setError(null);
      }
    }
  };

  const titulo = `Calificar a ${docente.nombre}`;
  const formulario = (cuerpo?: ElementType) => (
    <FormularioCalificacion
      borrador={borrador}
      cuatrimestres={cuatrimestres}
      onChange={setBorrador}
      onSubmit={handleSubmit}
      enviando={enviando}
      error={error}
      cuerpo={cuerpo}
    />
  );

  return (
    <>
      <button
        ref={botonRef}
        type="button"
        aria-haspopup="dialog"
        onClick={() => setAbierto(true)}
        className={`boton-primario ${className ?? ""}`}
      >
        <PenLine size={16} aria-hidden="true" className="shrink-0" />
        Calificar
      </button>

      {esDesktop ? (
        <Dialog.Root {...rootProps}>
          <Dialog.Portal>
            <Dialog.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black/20 transition-opacity duration-200 ease-out data-ending-style:opacity-0 data-starting-style:opacity-0 motion-reduce:transition-none" />
            <Dialog.Popup
              finalFocus={botonRef}
              className="fixed top-1/2 left-1/2 z-50 flex max-h-[calc(100dvh-4rem)] w-[calc(100%-2rem)] max-w-3xl -translate-1/2 flex-col overflow-hidden relieve bg-white text-stone-900 transition-[scale,opacity] duration-200 ease-out outline-none [--relieve:var(--color-stone-400)] data-ending-style:scale-97 data-ending-style:opacity-0 data-starting-style:scale-97 data-starting-style:opacity-0 motion-reduce:transition-none"
            >
              <Encabezado titulo={Dialog.Title} cerrar={Dialog.Close} mostrarEsc>
                {titulo}
              </Encabezado>
              {formulario()}
            </Dialog.Popup>
          </Dialog.Portal>
        </Dialog.Root>
      ) : (
        <Drawer.Root {...rootProps}>
          <Drawer.VirtualKeyboardProvider>
            <Drawer.Portal>
              <Drawer.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black/20 opacity-[calc(1-var(--drawer-swipe-progress))] transition-opacity duration-300 ease-out data-ending-style:opacity-0 data-ending-style:duration-[calc(var(--drawer-swipe-strength)*300ms)] data-starting-style:opacity-0 data-swiping:duration-0 supports-[-webkit-touch-callout:none]:absolute motion-reduce:transition-none" />
              <Drawer.Viewport className="fixed inset-0 z-50 flex items-end">
                <Drawer.Popup
                  finalFocus={botonRef}
                  className="flex max-h-[calc(100dvh-var(--header-height))] w-full [transform:translateY(var(--drawer-swipe-movement-y))] flex-col rounded-t-xl border-t-[1.5px] border-stone-300 bg-white text-stone-900 transition-transform duration-300 ease-out outline-none data-ending-style:[transform:translateY(100%)] data-ending-style:duration-[calc(var(--drawer-swipe-strength)*300ms)] data-starting-style:[transform:translateY(100%)] data-swiping:duration-0 data-swiping:select-none motion-reduce:transition-none"
                >
                  <div
                    aria-hidden="true"
                    className="mx-auto mt-2 h-1 w-10 shrink-0 rounded-full bg-stone-300"
                  />
                  <Encabezado titulo={Drawer.Title} cerrar={Drawer.Close} mostrarEsc={false}>
                    {titulo}
                  </Encabezado>
                  {formulario(Drawer.Content)}
                </Drawer.Popup>
              </Drawer.Viewport>
            </Drawer.Portal>
          </Drawer.VirtualKeyboardProvider>
        </Drawer.Root>
      )}
    </>
  );
}

function Encabezado({
  titulo: Titulo,
  cerrar: Cerrar,
  mostrarEsc,
  children
}: {
  titulo: ElementType<{ className?: string; children?: ReactNode }>;
  cerrar: ElementType<{ className?: string; "aria-label"?: string; children?: ReactNode }>;
  mostrarEsc: boolean;
  children: ReactNode;
}) {
  return (
    <div className="flex shrink-0 items-center justify-between gap-3 border-b border-stone-200 px-4 py-3 md:px-6">
      <Titulo className="min-w-0 text-base font-medium tracking-tight text-stone-900 md:text-lg">
        {children}
      </Titulo>
      <Cerrar
        aria-label="Cerrar"
        className="flex h-6 shrink-0 relieve-boton cursor-pointer items-center bg-white px-1.5 font-mono text-xs text-stone-500 outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-fiuba"
      >
        {mostrarEsc ? <span>Esc</span> : <Close size={14} aria-hidden="true" />}
      </Cerrar>
    </div>
  );
}
