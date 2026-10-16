import type { Metadata } from "next";
import { BuscadorMaterias } from "@/components/materia/buscador-materias";
import { getMateriasBuscador } from "@/server/db/queries/materias-buscador";

export const metadata: Metadata = {
  title: "FIUBA Reviews",
  description: "Calificaciones y comentarios de docentes de FIUBA, organizados por materia.",
  alternates: { canonical: "/" }
};

export default async function Page() {
  const materias = await getMateriasBuscador();

  return (
    <main className="mx-auto max-w-7xl space-y-4 p-4 md:p-6">
      <h1 className="text-2xl font-semibold">Materias</h1>
      <BuscadorMaterias materias={materias} />
    </main>
  );
}
