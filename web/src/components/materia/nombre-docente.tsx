"use client";

import { Link as LinkIcon } from "lucide-react";
import { useEffect } from "react";

function scrollAlDocente(codigo: string) {
  document.getElementById(codigo)?.scrollIntoView({
    behavior: window.matchMedia("(prefers-reduced-motion: reduce)").matches ? "instant" : "smooth",
    block: "start"
  });
}

export function NombreDocente({
  codigo,
  nombre,
  catedra
}: {
  codigo: string;
  nombre: string;
  catedra: string;
}) {
  const href = `?catedra=${encodeURIComponent(catedra)}#${encodeURIComponent(codigo)}`;

  useEffect(() => {
    // El detalle puede llegar por streaming: recién acá existe el destino del hash.
    let frame: number | undefined;
    function navegarAlHash() {
      if (window.location.hash !== `#${encodeURIComponent(codigo)}`) return;
      if (frame !== undefined) cancelAnimationFrame(frame);
      frame = requestAnimationFrame(() => scrollAlDocente(codigo));
    }
    navegarAlHash();
    window.addEventListener("hashchange", navegarAlHash);
    return () => {
      if (frame !== undefined) cancelAnimationFrame(frame);
      window.removeEventListener("hashchange", navegarAlHash);
    };
  }, [codigo]);

  return (
    <h2 className="mb-2 text-3xl font-semibold tracking-tight">
      <a
        href={href}
        onClick={(event) => {
          if (
            event.button !== 0 ||
            event.metaKey ||
            event.ctrlKey ||
            event.shiftKey ||
            event.altKey
          ) {
            return;
          }
          event.preventDefault();
          // Actualiza la URL sin el salto instantáneo del ancla nativa.
          const url = new URL(window.location.href);
          url.searchParams.set("catedra", catedra);
          url.hash = codigo;
          if (url.href !== window.location.href) {
            window.history.pushState(null, "", url.href);
          }
          scrollAlDocente(codigo);
        }}
        className="inline-flex cursor-pointer items-baseline gap-2 rounded-sm text-left focus-visible:outline-2 focus-visible:outline-offset-4"
        aria-label={`Ir a ${nombre}`}
        title="Enlace al docente"
      >
        <span>{nombre}</span>
        <LinkIcon
          className="size-4 shrink-0 self-center text-muted-foreground"
          aria-hidden="true"
        />
      </a>
    </h2>
  );
}
