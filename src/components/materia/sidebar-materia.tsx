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
    <span className="flex items-center gap-2">
      <span className="min-w-0 truncate font-acento">{materia.nombre}</span>
      <span className="shrink-0 font-numeros font-normal text-stone-400">{materia.codigo}</span>
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
      <header className="shrink-0 space-y-1.5 px-4 pt-4 pb-2">
        <div className="flex items-start gap-2">
          <h1 className="min-w-0 flex-1 font-acento text-base/5 font-medium text-tinta">
            {materia.nombre}
          </h1>
          {/* Código como texto apagado, sin pill: el nombre es lo importante. `leading-5`: misma
              altura que el primer renglón del nombre. */}
          <span className="shrink-0 font-numeros text-sm/5 text-stone-400">{materia.codigo}</span>
        </div>
        <p className="text-sm text-stone-500">
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
          )}
        </p>
      </header>
      <ListaCatedras codigoMateria={codigo} catedras={catedras} />
    </>
  );
}
