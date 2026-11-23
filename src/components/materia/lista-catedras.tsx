"use client";

import type { CatedraListada } from "@/server/db/queries/catedras";
import Link from "next/link";
import { useParams } from "next/navigation";

export function ListaCatedras({
  codigoMateria,
  catedras
}: {
  codigoMateria: string;
  catedras: CatedraListada[];
}) {
  // En `/materias/[codigo]` no hay cátedra en la URL y se muestra la primera.
  const activa = useParams<{ catedra?: string }>().catedra ?? catedras[0]?.codigo;

  if (catedras.length === 0) {
    return <p className="p-4 text-sm">Esta materia no tiene cátedras todavía.</p>;
  }

  return (
    <nav aria-label="Cátedras" className="min-h-0 flex-1 overflow-y-auto">
      <ul className="p-2">
        {catedras.map((catedra) => (
          <li key={catedra.codigo}>
            <Link
              href={`/materias/${codigoMateria}/catedras/${catedra.codigo}`}
              // Trae la página prerenderizada de la cátedra antes del click.
              prefetch={true}
              aria-current={catedra.codigo === activa ? "page" : undefined}
              className="flex gap-3 rounded p-2 hover:bg-neutral-100 aria-[current=page]:bg-neutral-200"
            >
              <span className="w-8 shrink-0 text-right font-mono">
                {catedra.calificacion > 0 ? catedra.calificacion.toFixed(1) : "–"}
              </span>
              <span>{catedra.nombre}</span>
            </Link>
          </li>
        ))}
      </ul>
    </nav>
  );
}
