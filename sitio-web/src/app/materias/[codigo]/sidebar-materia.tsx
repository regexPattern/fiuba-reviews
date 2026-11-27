import Link from "next/link";
import { notFound } from "next/navigation";
import { getCatedrasMateria } from "@/server/db/queries/catedras";
import { getMateria } from "@/server/db/queries/materias";
import { ListaCatedras } from "./lista-catedras";

export async function TituloMateria({ params }: Pick<LayoutProps<"/materias/[codigo]">, "params">) {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) notFound();

  return (
    <span className="flex items-baseline gap-2">
      <span className="min-w-0 flex-1 truncate text-base/5 font-medium tracking-tight">
        {materia.nombre}
      </span>
      <span className="shrink-0 font-mono text-sm/5">{materia.codigo}</span>
    </span>
  );
}

export async function SidebarMateria({
  params
}: Pick<LayoutProps<"/materias/[codigo]">, "params">) {
  const { codigo } = await params;
  const [materia, catedras] = await Promise.all([getMateria(codigo), getCatedrasMateria(codigo)]);

  if (!materia) notFound();

  return (
    <>
      <header className="shrink-0 divide-y divide-stone-200 border-b border-stone-200">
        <h1 className="flex items-baseline gap-2 p-4 text-stone-900">
          <span className="min-w-0 flex-1 text-base/5 font-medium tracking-tight">
            {materia.nombre}
          </span>
          <span className="shrink-0 font-mono text-sm/5">{materia.codigo}</span>
        </h1>
        <p className="px-4 py-3 text-sm text-stone-500">
          {materia.cuatrimestre ? (
            <>
              Última actualización de cátedras en{" "}
              <span className="font-mono text-stone-700 tabular-nums">
                {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
              </span>
              .
            </>
          ) : (
            <>
              Se muestran las cátedras de sus equivalencias:{" "}
              <span className="font-mono text-stone-700 tabular-nums">
                {materia.equivalencias.map((e) => e.codigo).join(", ")}
              </span>
              .
            </>
          )}{" "}
          Si ves que la oferta está desactualizada, podés colaborar{" "}
          <Link
            href="/planes"
            className="text-fiuba underline decoration-fiuba/40 underline-offset-2 hover:decoration-fiuba"
          >
            enviándola acá
          </Link>
          .
        </p>
      </header>
      <ListaCatedras codigoMateria={codigo} catedras={catedras} />
    </>
  );
}
