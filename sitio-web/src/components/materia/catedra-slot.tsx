import { Suspense } from "react";
import { UUID_RE } from "@/lib/utils";
import { getCatedrasMateria } from "@/server/queries/catedras-materia";
import { DetalleCatedra } from "./detalle-catedra";
import { ScrollTop } from "./scroll-top";
import { CatedraSkeleton } from "./skeletons";

/**
 * Única parte dinámica de la página: lee `?catedra`. Después todo sale del caché. Si el parámetro
 * es inválido o la cátedra no pertenece a la materia, se usa la primera del listado.
 */
export async function CatedraSlot({
  codigo,
  searchParams
}: {
  codigo: string;
  searchParams: Promise<{ catedra?: string | string[] }>;
}) {
  const [{ catedra }, catedras] = await Promise.all([searchParams, getCatedrasMateria(codigo)]);

  if (catedras.length === 0) {
    return (
      <div className="m-4 space-y-2 md:m-6">
        <h1 className="text-3xl font-semibold tracking-tight">Sin cátedras</h1>
        <p className="text-muted-foreground">
          Todavía no hay cátedras con calificaciones o comentarios para esta materia.
        </p>
      </div>
    );
  }

  const uuid =
    typeof catedra === "string" &&
    UUID_RE.test(catedra) &&
    catedras.some((c) => c.codigo === catedra)
      ? catedra
      : catedras[0].codigo;

  return (
    <>
      <ScrollTop key={`scroll-${uuid}`} />
      <Suspense key={uuid} fallback={<CatedraSkeleton />}>
        <DetalleCatedra uuid={uuid} />
      </Suspense>
    </>
  );
}
