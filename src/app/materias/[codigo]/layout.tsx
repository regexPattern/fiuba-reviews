import type { Metadata } from "next";
import { Suspense } from "react";
import { MateriaShell } from "@/components/materia/materia-shell";
import { SidebarMateria, TituloMateria } from "@/components/materia/sidebar-materia";
import { getMateria } from "@/server/db/queries/materias";

export async function generateMetadata({
  params
}: LayoutProps<"/materias/[codigo]">): Promise<Metadata> {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) return {};

  return { title: `${materia.codigo} • ${materia.nombre} | FIUBA Reviews` };
}

export default function Layout({ params, children }: LayoutProps<"/materias/[codigo]">) {
  return (
    <MateriaShell
      titulo={
        <Suspense fallback="Cargando materia…">
          <TituloMateria params={params} />
        </Suspense>
      }
      sidebar={
        <Suspense fallback={<p className="p-4">Cargando materia…</p>}>
          <SidebarMateria params={params} />
        </Suspense>
      }
    >
      {children}
    </MateriaShell>
  );
}
