import { notFound } from "next/navigation";
import { OG_SIZE, ogMateria } from "@/lib/og";
import { getMateria } from "@/server/db/queries/materia";
import { getMateriasVigentes } from "@/server/db/queries/materias-vigentes";

export const size = OG_SIZE;
export const contentType = "image/png";
export const alt = "FIUBA Reviews";

export async function generateStaticParams() {
  const materias = await getMateriasVigentes();

  return materias.map(({ codigo }) => ({ codigo }));
}

export default async function Image({ params }: { params: Promise<{ codigo: string }> }) {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) notFound();

  return ogMateria(materia);
}
