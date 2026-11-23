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
    <>
      {materia.codigo} · {materia.nombre}
    </>
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
      <header className="shrink-0 space-y-2 border-b p-4">
        <h1 className="font-pixel text-lg">
          {materia.codigo} · {materia.nombre}
        </h1>
        <p className="text-sm">
          {materia.cuatrimestre ? (
            <>
              Última actualización de cátedras en{" "}
              <span className="font-mono">
                {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
              </span>
              .
            </>
          ) : (
            <>
              Se muestran las cátedras de sus equivalencias:{" "}
              <span className="font-mono">
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
