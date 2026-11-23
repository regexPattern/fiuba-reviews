import { notFound } from "next/navigation";
import { getCatedrasMateria } from "@/server/db/queries/catedras";
import { getMateria } from "@/server/db/queries/materias";
import { ListaCatedras } from "./lista-catedras";

export async function SidebarMateria({
  params
}: Pick<LayoutProps<"/materias/[codigo]">, "params">) {
  const { codigo } = await params;
  const [materia, catedras] = await Promise.all([getMateria(codigo), getCatedrasMateria(codigo)]);

  if (!materia) notFound();

  return (
    <>
      <header className="space-y-2 border-b p-4">
        <h1 className="text-lg font-semibold">
          <span className="font-mono">{materia.codigo}</span> · {materia.nombre}
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
