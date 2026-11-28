"use client";

import Link from "next/link";
import { useEffect } from "react";
import { PaginaError } from "@/components/error/pagina-error";

export default function Error({ error, retry }: { error: unknown; retry: () => void }) {
  useEffect(() => {
    console.error(error);
  }, [error]);

  /*
   * Los errores de Server Components llegan con `digest`, un hash que genera
   * Next para no mandarle al navegador el mensaje real del error. Con ese hash
   * se lo encuentra en los logs del servidor. Para esos errores la página dice
   * "Error 500", y para los del navegador dice solo "Error".
   */
  const esDelServidor = error instanceof Error && "digest" in error;

  return (
    <>
      <title>Algo salió mal | FIUBA Reviews</title>
      <PaginaError
        codigo={esDelServidor ? 500 : undefined}
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
