import { cacheLife, cacheTag } from "next/cache";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Info, MessageSquarePlus } from "lucide-react";
import { getCatedra } from "@/server/queries/catedra";
import { Comentarios } from "./comentarios";
import { Promedios } from "./promedios";

export async function DetalleCatedra({
  uuid,
  codigoMateria
}: {
  uuid: string;
  codigoMateria: string;
}) {
  "use cache";
  cacheLife("max");
  cacheTag(`catedra:${uuid}`);

  const catedra = await getCatedra(uuid);

  if (!catedra) notFound();

  return (
    <div className="m-4 space-y-8 md:m-6">
      {catedra.docentes.map((docente) => (
        <section key={docente.codigo} id={docente.codigo} className="scroll-mt-[68px] space-y-3">
          <div>
            <h1 className="w-fit text-4xl font-semibold tracking-tight">{docente.nombre}</h1>
            {docente.rol && <small className="text-sm">({docente.rol})</small>}
          </div>

          {docente.resumenComentario && (
            <div className="divide-y divide-fiuba border border-fiuba bg-fiuba/45">
              <p className="p-3 before:content-['&quot;'] after:content-['&quot;']">
                {docente.resumenComentario}
              </p>
              <div className="flex items-center gap-1 p-3 text-button-foreground select-none">
                <Info className="size-4" />
                <span className="text-sm">Resumen generado con IA.</span>
              </div>
            </div>
          )}

          <div className="flex gap-2 text-sm text-button-foreground">
            <Promedios
              promedios={docente.promedios}
              cantidadCalificaciones={docente.cantidadCalificaciones}
            />

            <Link
              href={{
                pathname: "/calificar",
                query: { docente: docente.codigo, catedra: catedra.codigo, materia: codigoMateria }
              }}
              className="flex items-center gap-2 border border-button-border bg-button-background px-3 py-2 transition-colors hover:bg-button-hover"
            >
              <span>Calificar</span>
              <MessageSquarePlus className="size-4 fill-fiuba/50 stroke-[#665889] dark:fill-[#D1BCE3]" />
            </Link>
          </div>

          {docente.comentarios.length > 0 ? (
            <Comentarios comentarios={docente.comentarios} />
          ) : (
            <p className="py-2 text-sm text-muted-foreground">Docente sin comentarios</p>
          )}
        </section>
      ))}
    </div>
  );
}
