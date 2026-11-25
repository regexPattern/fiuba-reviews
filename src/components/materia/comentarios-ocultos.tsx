"use client";

import { ChevronDown } from "@withicons/react/line";
import { type ReactNode, useId, useRef, useState } from "react";
import { Numero } from "@/components/ui/numero";

/**
 * Comentarios que arrancan ocultos, con un botón al final para mostrarlos y volver a ocultarlos.
 * Usa `boton-terciario` (ver globals.css): es una acción secundaria y no tiene que competir con
 * "Calificar" ni con el contenido.
 */
export function ComentariosOcultos({
  cantidad,
  children
}: {
  cantidad: number;
  children: ReactNode;
}) {
  const [abiertos, setAbiertos] = useState(false);
  const botonRef = useRef<HTMLButtonElement>(null);
  const id = useId();

  function alternar() {
    setAbiertos(!abiertos);

    // Al ocultar, lo de arriba se achica y el botón sube: lo traemos a la vista para no dejar al
    // usuario en medio del docente siguiente.
    if (abiertos) {
      requestAnimationFrame(() => botonRef.current?.scrollIntoView({ block: "nearest" }));
    }
  }

  return (
    <>
      {/* En el HTML aunque esté oculto, así el contenido se prerenderiza igual. */}
      <div id={id} hidden={!abiertos}>
        {children}
      </div>
      <button
        ref={botonRef}
        type="button"
        aria-expanded={abiertos}
        aria-controls={id}
        onClick={alternar}
        // Centrado bajo los comentarios. `scroll-*`: para el `scrollIntoView` al ocultar.
        className="mx-auto mt-1 boton-terciario scroll-mt-[calc(var(--header-height)+var(--materia-sticky-height,0px)+1rem)] scroll-mb-4"
      >
        <ChevronDown
          size={18}
          aria-hidden="true"
          className={`shrink-0 transition-transform duration-150 motion-reduce:transition-none ${abiertos ? "rotate-180" : ""}`}
        />
        {abiertos ? (
          "Mostrar menos comentarios"
        ) : (
          <span>
            Mostrar <Numero>{cantidad}</Numero>{" "}
            {cantidad === 1 ? "comentario más" : "comentarios más"}
          </span>
        )}
      </button>
    </>
  );
}
