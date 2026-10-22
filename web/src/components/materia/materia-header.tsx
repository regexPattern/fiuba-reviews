import { cacheLife, cacheTag } from "next/cache";
import { getMateria } from "@/server/db/queries/materia";

export async function MateriaHeader({ codigo }: { codigo: string }) {
  "use cache";
  cacheLife("max");
  cacheTag(`materia:${codigo}`);

  const materia = await getMateria(codigo);

  if (!materia) return null;

  return (
    <>
      <h1 className="shrink-0 border-b p-4 text-lg font-semibold">
        <span className="font-mono">{materia.codigo}</span> · {materia.nombre}
      </h1>

      <div className="shrink-0 border-b p-3 text-sm">
        {materia.cuatrimestre ? (
          <p>
            Última actualización de cátedras en{" "}
            <span className="font-mono tabular-nums">
              {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
            </span>
            .
          </p>
        ) : (
          <p>
            La oferta de esta materia no está actualizada todavía. Se muestran las ofertas de sus
            equivalencias en los planes anteriores:{" "}
            <span className="font-mono">
              {materia.equivalencias.map((e) => e.codigo).join(", ")}
            </span>
            .
          </p>
        )}
      </div>
    </>
  );
}
