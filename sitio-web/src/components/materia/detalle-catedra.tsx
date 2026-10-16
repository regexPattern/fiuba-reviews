import { cacheLife, cacheTag } from "next/cache";
import { notFound } from "next/navigation";
import { Card, CardContent, CardDescription, CardHeader } from "@/components/ui/card";
import { getCatedra } from "@/server/queries/catedra";
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
            <h2 className="text-xl font-semibold">{docente.nombre}</h2>
            {docente.rol && <p className="text-sm text-muted-foreground">{docente.rol}</p>}
          </div>
          {docente.resumenComentario && (
            <Card>
              <CardHeader>
                <CardDescription>Resumen generado con IA</CardDescription>
              </CardHeader>
              <CardContent>
                <p>{docente.resumenComentario}</p>
              </CardContent>
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
