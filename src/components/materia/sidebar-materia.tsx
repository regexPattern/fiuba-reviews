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
      <span className="shrink-0 font-numeros">{materia.codigo}</span>
      <span className="min-w-0 truncate font-acento">{materia.nombre}</span>
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
      <header className="shrink-0 space-y-3 px-4 pt-4 pb-2">
        {/* Código a la izquierda y del mismo color que el nombre, como parte del título. Si el nombre
            baja, sigue alineado consigo mismo y el código queda en la línea base del primer renglón. */}
        <h1 className="flex items-baseline gap-2 text-tinta">
          <span className="shrink-0 font-numeros text-sm/5">{materia.codigo}</span>
          <span className="min-w-0 flex-1 font-acento text-base/5 font-medium">
            {materia.nombre}
          </span>
        </h1>
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
