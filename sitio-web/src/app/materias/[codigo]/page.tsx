import { Suspense } from "react";
import { getCatedrasMateria } from "@/server/db/queries/catedras";
import { getMateriasVigentes } from "@/server/db/queries/materias";
import { DetalleCatedra } from "./detalle-catedra";

export async function generateStaticParams() {
  const materias = await getMateriasVigentes();

  return materias.map(({ codigo }) => ({ codigo }));
}

/** Sin cátedra elegida se muestra la primera del listado. */
export default function Page({ params }: PageProps<"/materias/[codigo]">) {
  return (
    <Suspense fallback={<p className="p-4">Cargando cátedra…</p>}>
      <PrimeraCatedra params={params} />
    </Suspense>
  );
}

async function PrimeraCatedra({ params }: Pick<PageProps<"/materias/[codigo]">, "params">) {
  const { codigo } = await params;
  const [primera] = await getCatedrasMateria(codigo);

  if (!primera) {
    return <p className="p-4">Todavía no hay cátedras para esta materia.</p>;
  }

  return <DetalleCatedra codigo={primera.codigo} />;
}
