import type { Metadata } from "next";
import Link from "next/link";
import { getMateriasVigentes } from "@/server/db/queries/materias-vigentes";

export const metadata: Metadata = {
  title: "FIUBA Reviews",
  description: "Calificaciones y comentarios de docentes de FIUBA, organizados por materia.",
  alternates: { canonical: "/" }
};

export default async function Page() {
  const materias = await getMateriasVigentes();

  return (
    <main className="mx-auto max-w-7xl space-y-4 p-4 md:p-6">
      <h1 className="text-2xl font-semibold">Materias</h1>
      <ul className="divide-y">
        {materias.map((materia) => (
          <li key={materia.codigo}>
            <Link
              href={`/materia/${materia.codigo}`}
              className="block rounded-md p-3 hover:bg-accent focus-visible:outline-ring"
            >
              {materia.codigo} · {materia.nombre}
            </Link>
          </li>
        ))}
      </ul>
    </main>
  );
}
