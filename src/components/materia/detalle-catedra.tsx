import type { ComentarioDocente, DocenteCatedra } from "@/server/db/queries/catedras";
import { notFound } from "next/navigation";
import { CRITERIOS } from "@/lib/criterios";
import { getCatedra } from "@/server/db/queries/catedras";

const COMENTARIOS_VISIBLES = 10;

export async function DetalleCatedra({ codigo }: { codigo: string }) {
  const catedra = await getCatedra(codigo);

  if (!catedra) notFound();

  return (
    <div className="space-y-10 p-4 md:p-6">
      {catedra.docentes.map((docente) => (
        <Docente key={docente.codigo} docente={docente} />
      ))}
    </div>
  );
}

function Docente({ docente }: { docente: DocenteCatedra }) {
  const visibles = docente.comentarios.slice(0, COMENTARIOS_VISIBLES);
  const ocultos = docente.comentarios.slice(COMENTARIOS_VISIBLES);

  return (
    <section id={docente.codigo} className="space-y-4">
      <div>
        <h2 className="text-2xl font-semibold">{docente.nombre}</h2>
        {docente.nombreSiu && <p className="text-sm text-neutral-500">{docente.nombreSiu}</p>}
        {docente.rol && <p className="text-sm text-neutral-500">{docente.rol}</p>}
      </div>

      {docente.resumenComentarios && (
        <div className="rounded border p-4">
          <p>{docente.resumenComentarios}</p>
          <p className="mt-2 text-sm text-neutral-500">Resumen generado con IA</p>
        </div>
      )}

      {docente.promedios ? (
        <details>
          <summary className="cursor-pointer">
            Promedio: <span className="font-mono">{docente.promedios.general.toFixed(1)}</span> (
            {docente.cantidadCalificaciones} calificaciones)
          </summary>
          <dl className="mt-2 grid max-w-xs grid-cols-[1fr_auto] gap-x-4 text-sm">
            {CRITERIOS.map(({ nombre, etiqueta }) => (
              <div key={nombre} className="contents">
                <dt>{etiqueta}</dt>
                <dd className="font-mono">{docente.promedios![nombre].toFixed(1)}</dd>
              </div>
            ))}
          </dl>
        </details>
      ) : (
        <p>Sin calificaciones</p>
      )}

      {docente.comentarios.length > 0 ? (
        <div className="divide-y">
          {visibles.map((comentario) => (
            <Comentario key={comentario.codigo} comentario={comentario} />
          ))}
          {ocultos.length > 0 && (
            <details className="divide-y">
              <summary className="cursor-pointer py-3">
                Mostrar {ocultos.length} comentarios más
              </summary>
              {ocultos.map((comentario) => (
                <Comentario key={comentario.codigo} comentario={comentario} />
              ))}
            </details>
          )}
        </div>
      ) : (
        <p className="text-sm text-neutral-500">Docente sin comentarios</p>
      )}
    </section>
  );
}

function Comentario({ comentario }: { comentario: ComentarioDocente }) {
  return (
    <article className="space-y-1 py-3">
      <p className="whitespace-pre-wrap">{comentario.contenido.trim()}</p>
      <p className="font-mono text-xs text-neutral-500">
        {comentario.cuatrimestre.numero}C{comentario.cuatrimestre.anio}
        {comentario.deDolly && " · Dolly"}
      </p>
    </article>
  );
}
