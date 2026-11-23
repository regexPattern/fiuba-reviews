import { notFound } from "next/navigation";
import { Suspense } from "react";
import { DetalleCatedra } from "@/components/materia/detalle-catedra";
import { getCatedrasMateria } from "@/server/db/queries/catedras";
import { getMateriasVigentes } from "@/server/db/queries/materias";

export async function generateStaticParams() {
  const materias = await getMateriasVigentes();
  const catedras = await Promise.all(materias.map(({ codigo }) => getCatedrasMateria(codigo)));

  return materias.flatMap(({ codigo }, i) =>
    catedras[i].map((catedra) => ({ codigo, catedra: catedra.codigo }))
  );
}

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
