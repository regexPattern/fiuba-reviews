import type { Cuatrimestre } from "@/lib/ofertas";
import type { ComentarioDocente, DocenteCatedra } from "@/server/db/queries/catedras";
import { MessageCircleQuestion, Quote, Sparkles } from "@withicons/react/plush";
import Image from "next/image";
import { notFound } from "next/navigation";
import dollyLogo from "@/assets/dolly-logo.png";
import { BotonCalificar } from "@/components/calificacion/boton-calificar";
import { ComentariosOcultos } from "@/components/materia/comentarios-ocultos";
import { PromedioDocente } from "@/components/materia/promedio-docente";
import { Numero } from "@/components/ui/numero";
import { getCatedra } from "@/server/db/queries/catedras";
import { getCuatrimestresCalificables } from "@/server/db/queries/cuatrimestres";

const COMENTARIOS_VISIBLES = 10;

export async function DetalleCatedra({ codigo }: { codigo: string }) {
  const [catedra, calificables] = await Promise.all([
    getCatedra(codigo),
    getCuatrimestresCalificables()
  ]);

  if (!catedra) notFound();

  const cuatrimestres = calificables.map(({ numero, anio }) => ({ numero, anio }));

  return (
    <div className="space-y-12 p-4 md:p-6">
      {catedra.docentes.map((docente) => (
        <Docente key={docente.codigo} docente={docente} cuatrimestres={cuatrimestres} />
      ))}
    </div>
  );
}

function Docente({
  docente,
  cuatrimestres
}: {
  docente: DocenteCatedra;
  cuatrimestres: Cuatrimestre[];
}) {
  return (
    <section
      id={docente.codigo}
      // Que los anchors a un docente no queden debajo del header (y de la barra sticky en móvil).
      className="scroll-mt-[calc(var(--header-height)+var(--materia-sticky-height,0px)+1rem)] space-y-4 md:scroll-mt-[calc(var(--header-height)+1.5rem)]"
    >
      <header className="flex flex-col gap-3 sm:flex-row sm:items-start sm:gap-2">
        <div className="min-w-0 flex-1 space-y-1">
          <h2 className="text-xl/7 font-medium tracking-tight text-stone-900">{docente.nombre}</h2>
          {(docente.nombreSiu || docente.rol) && (
            <div className="overflow-hidden">
              <p className="-ml-4 flex flex-wrap rotulo text-stone-500">
                {[docente.nombreSiu, docente.rol].filter(Boolean).map((parte) => (
                  <span
                    key={parte}
                    className="relative pl-4 before:absolute before:left-1.5 before:content-['·'_/_'']"
                  >
                    {parte}
                  </span>
                ))}
              </p>
            </div>
          )}
        </div>
        <div className="flex shrink-0 items-start gap-2">
          {docente.promedios ? (
            <PromedioDocente
              promedios={docente.promedios}
              cantidadCalificaciones={docente.cantidadCalificaciones}
            />
          ) : (
            <span className="shrink-0 relieve bg-stone-50 px-2.5 py-1 text-sm text-stone-400 select-none [--relieve:var(--color-stone-200)]">
              Sin calificaciones
            </span>
          )}
          <BotonCalificar
            docente={{ codigo: docente.codigo, nombre: docente.nombre }}
            cuatrimestres={cuatrimestres}
            className="-order-1 sm:order-none"
          />
        </div>
      </header>

      {docente.resumenComentarios && <ResumenIA resumen={docente.resumenComentarios} />}

      {docente.comentarios.length > 0 ? (
        <Comentarios id={`${docente.codigo}-comentarios`} comentarios={docente.comentarios} />
      ) : (
        <SinComentarios />
      )}
    </section>
  );
}

function ResumenIA({ resumen }: { resumen: string }) {
  // Los textos arrancan a ~26px del borde, igual que el de los comentarios (ícono de 18px + `gap-2`):
  // 1.5px de borde + `px-6`. Vertical: `py-5`, que con el interlineado de `text-base/relaxed` se ve
  // como el padding de los costados. El pie usa menos porque arranca con el ícono.
  return (
    <div className="overflow-hidden relieve bg-white">
      <p className="px-6 py-5 text-base/relaxed text-stone-700">{resumen}</p>
      <footer className="flex items-center gap-2 border-t border-stone-200 bg-stone-50 px-3 py-2 rotulo text-stone-500">
        <Sparkles size={18} aria-hidden="true" className="shrink-0" />
        Resumen generado con IA
      </footer>
    </div>
  );
}

function SinComentarios() {
  return (
    <p className="flex items-center gap-2 rounded-md border-[1.5px] border-dashed border-stone-300 px-3 py-3 text-sm text-stone-500">
      <MessageCircleQuestion size={18} aria-hidden="true" className="shrink-0 text-stone-400" />
      Todavía nadie comentó sobre este docente.
    </p>
  );
}

function Comentarios({ id, comentarios }: { id: string; comentarios: ComentarioDocente[] }) {
  const visibles = comentarios.slice(0, COMENTARIOS_VISIBLES);
  const ocultos = comentarios.slice(COMENTARIOS_VISIBLES);

  return (
    <div
      id={id}
      className="-mx-2 scroll-mt-[calc(var(--header-height)+var(--materia-sticky-height,0px)+1rem)] md:scroll-mt-[calc(var(--header-height)+1.5rem)]"
    >
      {visibles.map((comentario) => (
        <Comentario key={comentario.codigo} comentario={comentario} />
      ))}
      {ocultos.length > 0 && (
        <ComentariosOcultos cantidad={ocultos.length}>
          {ocultos.map((comentario) => (
            <Comentario key={comentario.codigo} comentario={comentario} />
          ))}
        </ComentariosOcultos>
      )}
    </div>
  );
}

function Comentario({ comentario }: { comentario: ComentarioDocente }) {
  return (
    <article className="flex gap-2 px-2 py-2">
      <Quote size={18} aria-hidden="true" className="mt-1 shrink-0 text-stone-400" />
      <p className="min-w-0 flex-1 text-base/relaxed text-stone-700">
        <span className="whitespace-pre-wrap">{comentario.contenido.trim()}</span>
        {"\u00A0 "}
        <span className="whitespace-nowrap text-stone-400 select-none">
          <span className="text-xs">
            <Numero>
              {comentario.cuatrimestre.numero}C{comentario.cuatrimestre.anio}
            </Numero>
          </span>
          {comentario.deDolly && (
            <>
              {/* Mismo espacio que entre el texto y el cuatrimestre (al tamaño del texto). */}
              {"\u00A0\u00A0"}
              <Image
                src={dollyLogo}
                alt="Dolly FIUBA"
                title="Comentario importado de Dolly FIUBA"
                draggable={false}
                className="inline-block h-4 w-auto align-[-3px] select-none"
              />
            </>
          )}
        </span>
      </p>
    </article>
  );
}
