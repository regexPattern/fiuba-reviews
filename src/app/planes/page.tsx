import type { Metadata } from "next";
import type { ReactNode } from "react";
import { ArrowUpRight, ChevronRight } from "@withicons/react/line";
import { FormularioOferta } from "@/components/planes/formulario-oferta";
import { TablaOfertas } from "@/components/planes/tabla-ofertas";
import { getOfertasCarreras } from "@/server/db/queries/ofertas";

export const metadata: Metadata = {
  title: "Ofertas de comisiones | FIUBA Reviews",
  description:
    "Ayudá a mantener al día las cátedras de FIUBA Reviews enviando la oferta de comisiones de tu carrera.",
  robots: { index: false, follow: false }
};

const URL_OFERTA_SIU = "https://guaraniautogestion.fi.uba.ar/g3w/oferta_comisiones";

export default async function Page() {
  const { ofertas, ultimoCuatrimestre } = await getOfertasCarreras();

  return (
    <main className="mx-auto max-w-6xl px-4 py-8 md:px-6 md:py-12">
      <header className="max-w-2xl space-y-3">
        <p className="rotulo text-fiuba">Colaborar</p>
        <h1 className="font-acento text-3xl font-medium text-tinta">Ofertas de comisiones</h1>
        <p className="text-base/relaxed text-stone-700">
          Las cátedras de cada materia salen de la oferta de comisiones del SIU, pero cada alumno
          solo puede ver la de su propia carrera. Por eso, desde el 1C2025 las ofertas se actualizan
          gracias a alumnos de cada carrera que comparten la suya.
        </p>
        <p className="text-base/relaxed text-stone-500">
          Antes se usaba la oferta horaria pública de la que se nutría Dolly FIUBA, que dejó de
          actualizarse en el 1C2024.
        </p>
      </header>

      <div className="mt-10 grid items-start gap-10 lg:grid-cols-2 lg:gap-12">
        <section aria-labelledby="estado" className="space-y-4">
          <h2 id="estado" className="font-acento text-xl font-medium text-tinta">
            Estado por carrera
          </h2>
          <TablaOfertas ofertas={ofertas} ultimoCuatrimestre={ultimoCuatrimestre} />
        </section>

        <section aria-labelledby="enviar" className="space-y-4">
          <h2 id="enviar" className="font-acento text-xl font-medium text-tinta">
            Enviar la oferta de tu carrera
          </h2>
          <ol className="space-y-5">
            <Paso numero={1} titulo="Abrí la oferta de comisiones en el SIU">
              <a
                href={URL_OFERTA_SIU}
                target="_blank"
                rel="noopener noreferrer"
                className="boton-secundario w-fit"
              >
                Reportes
                <ChevronRight size={14} aria-hidden="true" className="shrink-0 text-stone-400" />
                Oferta de comisiones
                <ArrowUpRight size={16} aria-hidden="true" className="shrink-0 text-stone-500" />
              </a>
            </Paso>
            <Paso numero={2} titulo="Seleccioná toda la página">
              <Atajo tecla="A" />
            </Paso>
            <Paso numero={3} titulo="Copiala">
              <Atajo tecla="C" />
            </Paso>
            <Paso numero={4} titulo="Pegala acá abajo y enviala">
              <FormularioOferta ofertas={ofertas} />
            </Paso>
          </ol>
        </section>
      </div>
    </main>
  );
}

/** Paso numerado: el número en un círculo y una línea que lo une con el siguiente. */
function Paso({
  numero,
  titulo,
  children
}: {
  numero: number;
  titulo: string;
  children: ReactNode;
}) {
  return (
    <li className="relative flex gap-3 pb-1 not-last:after:absolute not-last:after:top-8 not-last:after:-bottom-4 not-last:after:left-3 not-last:after:w-px not-last:after:-translate-x-1/2 not-last:after:bg-stone-200">
      <span
        aria-hidden="true"
        className="flex size-6 shrink-0 items-center justify-center rounded-full border-[1.5px] border-stone-300 bg-white font-numeros text-xs text-stone-600"
      >
        {numero}
      </span>
      <div className="min-w-0 flex-1 space-y-2.5">
        <h3 className="text-base/6 font-medium text-tinta">{titulo}</h3>
        {children}
      </div>
    </li>
  );
}

/** El mismo atajo en Windows/Linux y en macOS. */
function Atajo({ tecla }: { tecla: string }) {
  return (
    <p className="flex flex-wrap items-center gap-x-3 gap-y-1.5 text-sm text-stone-500">
      <span className="flex items-center gap-1">
        <Tecla>Ctrl</Tecla>
        <Tecla>{tecla}</Tecla>
      </span>
      <span>o</span>
      <span className="flex items-center gap-1">
        <Tecla>⌘</Tecla>
        <Tecla>{tecla}</Tecla>
        <span className="ml-1">en macOS</span>
      </span>
    </p>
  );
}

/** Tecla en relieve, como los botones. */
function Tecla({ children }: { children: ReactNode }) {
  return (
    <kbd className="inline-flex h-6 min-w-6 items-center justify-center relieve bg-white px-1.5 font-mono text-xs text-stone-700 [--relieve-alto:2px]">
      {children}
    </kbd>
  );
}
