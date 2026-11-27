import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { Suspense } from "react";
import { getCatedrasMateria } from "@/server/db/queries/catedras";
import { getMateriasVigentes } from "@/server/db/queries/materias";
import { DetalleCatedra } from "../../detalle-catedra";

export async function generateStaticParams() {
  const materias = await getMateriasVigentes();
  const catedras = await Promise.all(materias.map(({ codigo }) => getCatedrasMateria(codigo)));

  return materias.flatMap(({ codigo }, i) =>
    catedras[i].map((catedra) => ({ codigo, catedra: catedra.codigo }))
  );
}

// No se indexan las cátedras, solo las materias. Pero si se siguen sus links
// al resto de las cátedras y materias.
export const metadata: Metadata = {
  robots: { index: false, follow: true }
};

export default function Page({ params }: PageProps<"/materias/[codigo]/catedras/[catedra]">) {
  return (
    <Suspense fallback={<p className="p-4">Cargando cátedra…</p>}>
      <Catedra params={params} />
    </Suspense>
  );
}

async function Catedra({
  params
}: Pick<PageProps<"/materias/[codigo]/catedras/[catedra]">, "params">) {
  const { codigo, catedra } = await params;
  const catedras = await getCatedrasMateria(codigo);

  // La cátedra tiene que ser una de las que se muestran en esta materia.
  if (!catedras.some((c) => c.codigo === catedra)) notFound();

  return <DetalleCatedra codigo={catedra} />;
}
