"use client";

import Link from "next/link";
import { useEffect } from "react";
import { PaginaError } from "@/components/error/pagina-error";

export default function Error({ error, retry }: { error: unknown; retry: () => void }) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  // Los errores de Server Components llegan con `digest`; el resto es desconocido (el 99).
  const delServidor = error instanceof Error && "digest" in error;

  return (
    <>
      <title>Algo salió mal | FIUBA Reviews</title>
      <PaginaError
        codigo={delServidor ? 500 : undefined}
        descripcion="Algo se rompió en el camino. Probá de nuevo en un rato."
      >
        <button type="button" onClick={() => retry()} className="boton-primario">
          Reintentar
        </button>
        <Link href="/" className="boton-secundario">
          Volver al inicio
        </Link>
      </PaginaError>
    </>
  );
}
