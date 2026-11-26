import type { Metadata } from "next";
import { PasosOferta } from "@/components/planes/pasos-oferta";
import { TablaOfertas } from "@/components/planes/tabla-ofertas";
import { getOfertasCarreras } from "@/server/db/queries/ofertas";

export const metadata: Metadata = {
  title: "Ofertas de comisiones | FIUBA Reviews",
  description:
    "Ayudá a mantener al día las cátedras de FIUBA Reviews enviando la oferta de comisiones de tu carrera.",
  robots: { index: false, follow: false }
};

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
          <PasosOferta ofertas={ofertas} />
        </section>
      </div>
    </main>
  );
}
