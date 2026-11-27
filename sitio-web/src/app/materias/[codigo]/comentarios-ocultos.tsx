"use client";

import { ChevronDown } from "@withicons/react/line";
import { type ReactNode, useId, useRef, useState } from "react";

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
            Mostrar <span className="font-mono tabular-nums">{cantidad}</span>{" "}
            {cantidad === 1 ? "comentario más" : "comentarios más"}
          </span>
        )}
      </button>
    </>
  );
}
