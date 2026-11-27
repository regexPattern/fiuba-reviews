"use client";

import "@fontsource-variable/geist";
import "@fontsource-variable/geist-mono";
import Link from "next/link";
import { useEffect } from "react";
import { PaginaError } from "@/components/error/pagina-error";
import "./globals.css";

export default function GlobalError({ error, retry }: { error: unknown; retry: () => void }) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <html lang="es-AR">
      <body
        className="bg-white font-sans antialiased"
        /* Como el header esta inyectado en el layout global de la app, para
          que el error inesperado se muestre solo en pantalla tenemos que
          reescribir esta variable CSS para esconder el header. */
        style={{ "--header-height": "0px" } as React.CSSProperties}
      >
        <title>Algo salió mal | FIUBA Reviews</title>
        <PaginaError descripcion="Algo se rompió en el camino. Probá de nuevo en un rato.">
          <button type="button" onClick={() => retry()} className="boton-primario">
            Reintentar
          </button>
          <Link href="/" className="boton-secundario">
            Volver al inicio
          </Link>
        </PaginaError>
      </body>
    </html>
  );
}
