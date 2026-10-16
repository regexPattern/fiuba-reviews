import Link from "next/link";
import { notFound } from "next/navigation";
import { Suspense } from "react";
import type { Metadata } from "next";
import { ChevronLeft } from "lucide-react";
import { BackgroundBlob } from "@/components/background-blob";
import { FormCalificar } from "@/components/calificar/form-calificar";
import { Skeleton } from "@/components/ui/skeleton";
import { UUID_RE } from "@/lib/utils";
import { getCuatrimestresCalificables, getDocenteParaCalificar } from "@/server/queries/calificar";

type SearchParams = Promise<{ docente?: string; catedra?: string; materia?: string }>;

export const metadata: Metadata = {
  title: "Calificar docente | FIUBA Reviews",
  robots: { index: false, follow: false },
  alternates: { canonical: "/calificar" }
};

async function Contenido({ searchParams }: { searchParams: SearchParams }) {
  const { docente: codigoDocente, catedra, materia } = await searchParams;

  if (!codigoDocente || !UUID_RE.test(codigoDocente)) notFound();

  const codigoCatedra = catedra && UUID_RE.test(catedra) ? catedra : null;
  const codigoMateria = materia && /^[0-9A-Za-z.-]+$/.test(materia) ? materia : null;

  const [docente, cuatris] = await Promise.all([
    getDocenteParaCalificar(codigoDocente),
    getCuatrimestresCalificables()
  ]);

  if (!docente || cuatris.length === 0) notFound();

  const materiaVolver = codigoMateria ?? docente.codigoMateria;
  const hrefVolver = codigoCatedra
    ? (`/materia/${materiaVolver}?catedra=${codigoCatedra}` as const)
    : (`/materia/${materiaVolver}` as const);

  return (
    <>
      <div className="flex flex-col-reverse justify-between gap-4 md:flex-row md:items-center">
        <div className="flex flex-col gap-2">
          <h1 className="text-4xl font-medium">{docente.nombre}</h1>
          {docente.nombreSiu && docente.rol && (
            <p className="text-sm">
              {docente.nombreSiu} • {docente.rol}
            </p>
          )}
        </div>

        <Link href={hrefVolver} className="flex items-center text-sm underline">
          <ChevronLeft className="size-[18px]" />
          {codigoCatedra ? "Ir a cátedra del docente" : "Ir a materia del docente"}
        </Link>
      </div>

      <FormCalificar
        codigoDocente={docente.codigo}
        codigoCatedra={codigoCatedra}
        codigoMateria={codigoMateria}
        cuatris={cuatris}
      />
    </>
  );
}

export default function Page({ searchParams }: { searchParams: SearchParams }) {
  return (
    <div className="relative isolate">
      <BackgroundBlob upperLeft lowerRight={false} />

      <main className="relative z-10 container mx-auto space-y-12 p-6">
        <Suspense fallback={<Skeleton className="h-96 w-full" />}>
          <Contenido searchParams={searchParams} />
        </Suspense>
      </main>
    </div>
  );
}
