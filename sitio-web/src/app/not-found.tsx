import type { Metadata } from "next";
import Link from "next/link";
import { PaginaError } from "@/components/error/pagina-error";

export const metadata: Metadata = { title: "Página no encontrada | FIUBA Reviews" };

export default function NotFound() {
  return (
    <PaginaError
      codigo={404}
      descripcion="Esta página se tomó otro colectivo. Revisá la dirección o volvé al inicio."
    >
      <Link href="/" className="boton-primario">
        Volver al inicio
      </Link>
    </PaginaError>
  );
}
