import { Bot } from "@withicons/react/chrome";
import { cacheLife, cacheTag } from "next/cache";
import { notFound } from "next/navigation";
import { Card, CardContent, CardDescription, CardFooter } from "@/components/ui/card";
import { getCatedra } from "@/server/db/queries/catedra";
import { Comentarios } from "./comentarios";
import { Promedios } from "./promedios";

export async function DetalleCatedra({ uuid }: { uuid: string }) {
  "use cache";
  cacheLife("max");
  cacheTag(`catedra:${uuid}`);

  const catedra = await getCatedra(uuid);
  if (!catedra) notFound();

  return (
    <div className="space-y-8 p-4 md:p-6">
      {catedra.docentes.map((docente) => (
        <section key={docente.codigo} id={docente.codigo} className="space-y-4">
          <div>
            <h2 className="text-3xl font-semibold tracking-tight">{docente.nombre}</h2>
            {docente.nombreSiu && (
              <p className="text-sm text-muted-foreground">{docente.nombreSiu}</p>
            )}
            {docente.rol && <p className="text-sm text-muted-foreground">{docente.rol}</p>}
          </div>
          {docente.resumenComentario && (
            <Card className="gap-3 py-4 shadow-none">
              <CardContent className="px-4">
                <p>{docente.resumenComentario}</p>
              </CardContent>
              <CardFooter className="px-4">
                <CardDescription className="flex items-center gap-2 select-none">
                  <Bot size={18} className="shrink-0" aria-hidden="true" />
                  Resumen generado con IA
                </CardDescription>
              </CardFooter>
            </Card>
          )}
          <Promedios
            promedios={docente.promedios}
            cantidadCalificaciones={docente.cantidadCalificaciones}
          />
          {docente.comentarios.length > 0 ? (
            <Comentarios comentarios={docente.comentarios} />
          ) : (
            <p className="text-sm text-muted-foreground">Docente sin comentarios</p>
          )}
        </section>
      ))}
    </div>
  );
}
