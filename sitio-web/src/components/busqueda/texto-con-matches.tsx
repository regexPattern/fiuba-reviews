import type { Rango } from "@/lib/busqueda";

export function TextoConMatches({
  texto,
  indices = [],
  className
}: {
  texto: string;
  indices?: readonly Rango[];
  className?: string;
}) {
  const partes = [];
  let cursor = 0;
  for (const [inicio, fin] of indices) {
    partes.push(texto.slice(cursor, inicio));
    partes.push(
      <span key={inicio} className="underline decoration-fiuba decoration-2 underline-offset-2">
        {texto.slice(inicio, fin + 1)}
      </span>
    );
    cursor = fin + 1;
  }
  partes.push(texto.slice(cursor));

  return <span className={className}>{partes}</span>;
}
