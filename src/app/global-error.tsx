"use client";

import "@fontsource-variable/geist";
import "@fontsource-variable/geist-mono";
import Link from "next/link";
import { useEffect } from "react";
import { PaginaError } from "@/components/error/pagina-error";
import "./globals.css";

// Reemplaza al layout raíz cuando este falla: no hay header, así que la página ocupa toda la
// pantalla (`--header-height: 0px`).
export default function GlobalError({ error, retry }: { error: unknown; retry: () => void }) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  return (
    <html lang="es">
      <body
        className="bg-white font-sans antialiased"
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
