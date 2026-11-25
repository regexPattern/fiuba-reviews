import Link from "next/link";
import { notFound } from "next/navigation";
import { getCatedrasMateria } from "@/server/db/queries/catedras";
import { getMateria } from "@/server/db/queries/materias";
import { ListaCatedras } from "./lista-catedras";

/** Título corto para la barra sticky móvil que abre el drawer de cátedras. */
export async function TituloMateria({ params }: Pick<LayoutProps<"/materias/[codigo]">, "params">) {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) notFound();

  return (
    // Misma tipografía que el título del sidebar: nombre a la izquierda y código a la derecha.
    <span className="flex items-baseline gap-2">
      <span className="min-w-0 flex-1 truncate font-acento text-base/5 font-medium">
        {materia.nombre}
      </span>
      <span className="shrink-0 font-numeros text-sm/5">{materia.codigo}</span>
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
        {/* Nombre a la izquierda y código a la derecha, del mismo color, como parte del título. Si
            el nombre baja, el código queda en la línea base del primer renglón. */}
        <h1 className="flex items-baseline gap-2 p-4 text-tinta">
          <span className="min-w-0 flex-1 font-acento text-base/5 font-medium">
            {materia.nombre}
          </span>
          <span className="shrink-0 font-numeros text-sm/5">{materia.codigo}</span>
        </h1>
        <p className="px-4 py-3 text-sm text-stone-500">
          {materia.cuatrimestre ? (
            <>
              Última actualización de cátedras en{" "}
              <span className="font-numeros text-stone-700 tabular-nums">
                {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
              </span>
              .
            </>
          ) : (
            <>
              Se muestran las cátedras de sus equivalencias:{" "}
              <span className="font-numeros text-stone-700 tabular-nums">
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
