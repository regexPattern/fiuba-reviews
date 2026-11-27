"use client";

import { Dialog } from "@base-ui/react/dialog";
import { Drawer } from "@base-ui/react/drawer";
import { Close, PenLine } from "@withicons/react/line";
import { type ElementType, type ReactNode, useRef, useState, useTransition } from "react";
import { enviarCalificacion } from "@/app/materias/actions";
import { BORRADOR_VACIO, type BorradorCalificacion } from "@/lib/calificaciones";
import { useMediaQuery } from "@/lib/media-query";
import type { Cuatrimestre } from "@/lib/ofertas";
import { FormularioCalificacion } from "./formulario-calificacion";

/**
 * Botón primario (celeste) que abre el panel para calificar al docente: un dialog en desktop y un
 * drawer desde abajo en móvil. El borrador sobrevive a cerrar el panel (se descarta solo al enviar),
 * así que cerrarlo sin querer no pierde lo cargado. Si el servidor rechaza la calificación, el panel
 * queda abierto con el error.
 */
export function BotonCalificar({
  docente,
  cuatrimestres,
  className
}: {
  docente: { codigo: string; nombre: string };
  /** Los que se pueden elegir, del más nuevo al más viejo. */
  cuatrimestres: Cuatrimestre[];
  className?: string;
}) {
  const botonRef = useRef<HTMLButtonElement>(null);
  const [abierto, setAbierto] = useState(false);
  const [borrador, setBorrador] = useState(BORRADOR_VACIO);
  const [error, setError] = useState<string | null>(null);
  const [enviando, startTransition] = useTransition();
  const enviado = useRef(false);
  // Mismo corte que el resto del sitio (`md`). Solo importa con el panel abierto, ya hidratado.
  const esDesktop = useMediaQuery("(min-width: 48rem)");

  function enviar(calificacion: BorradorCalificacion) {
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

  const raiz = {
    open: abierto,
    onOpenChange: setAbierto,
    // Se limpia al terminar la animación de cierre, así el formulario no se vacía mientras se va.
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
      onEnviar={enviar}
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
        <Dialog.Root {...raiz}>
          <Dialog.Portal>
            {/* `z-50`: por encima del header fijo (`z-40`). */}
            <Dialog.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black/20 transition-opacity duration-200 ease-out data-ending-style:opacity-0 data-starting-style:opacity-0 motion-reduce:transition-none" />
            <Dialog.Popup
              finalFocus={botonRef}
              className="fixed top-1/2 left-1/2 z-50 flex max-h-[calc(100dvh-4rem)] w-[calc(100%-2rem)] max-w-3xl -translate-1/2 flex-col overflow-hidden relieve bg-white text-stone-900 transition-[scale,opacity] duration-200 ease-out outline-none [--relieve:var(--color-stone-400)] data-ending-style:scale-97 data-ending-style:opacity-0 data-starting-style:scale-97 data-starting-style:opacity-0 motion-reduce:transition-none"
            >
              <Encabezado titulo={Dialog.Title} cerrar={Dialog.Close} teclado>
                {titulo}
              </Encabezado>
              {formulario()}
            </Dialog.Popup>
          </Dialog.Portal>
        </Dialog.Root>
      ) : (
        <Drawer.Root {...raiz}>
          {/* Lleva el campo enfocado a la vista cuando aparece el teclado. */}
          <Drawer.VirtualKeyboardProvider>
            <Drawer.Portal>
              {/* Se aclara a medida que se arrastra el drawer, como el de las cátedras. */}
              <Drawer.Backdrop className="fixed inset-0 z-50 min-h-dvh bg-black/20 opacity-[calc(1-var(--drawer-swipe-progress))] transition-opacity duration-300 ease-out data-ending-style:opacity-0 data-ending-style:duration-[calc(var(--drawer-swipe-strength)*300ms)] data-starting-style:opacity-0 data-swiping:duration-0 supports-[-webkit-touch-callout:none]:absolute motion-reduce:transition-none" />
              <Drawer.Viewport className="fixed inset-0 z-50 flex items-end">
                <Drawer.Popup
                  finalFocus={botonRef}
                  className="flex max-h-[calc(100dvh-var(--header-height))] w-full [transform:translateY(var(--drawer-swipe-movement-y))] flex-col rounded-t-xl border-t-[1.5px] border-stone-300 bg-white text-stone-900 transition-transform duration-300 ease-out outline-none data-ending-style:[transform:translateY(100%)] data-ending-style:duration-[calc(var(--drawer-swipe-strength)*300ms)] data-starting-style:[transform:translateY(100%)] data-swiping:duration-0 data-swiping:select-none motion-reduce:transition-none"
                >
                  {/* Manija: indica que se puede arrastrar para cerrar. */}
                  <div
                    aria-hidden="true"
                    className="mx-auto mt-2 h-1 w-10 shrink-0 rounded-full bg-stone-300"
                  />
                  <Encabezado titulo={Drawer.Title} cerrar={Drawer.Close} teclado={false}>
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

/**
 * Título y botón de cerrar, con las partes del dialog o del drawer según cuál esté abierto. El
 * cerrar es una tecla chica en relieve, como en el buscador: `Esc` en desktop y una cruz en móvil.
 */
function Encabezado({
  titulo: Titulo,
  cerrar: Cerrar,
  teclado,
  children
}: {
  titulo: ElementType<{ className?: string; children?: ReactNode }>;
  cerrar: ElementType<{ className?: string; "aria-label"?: string; children?: ReactNode }>;
  /** Si se muestra `Esc` (desktop) en vez de la cruz. */
  teclado: boolean;
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
        {teclado ? <span>Esc</span> : <Close size={14} aria-hidden="true" />}
      </Cerrar>
    </div>
  );
}
