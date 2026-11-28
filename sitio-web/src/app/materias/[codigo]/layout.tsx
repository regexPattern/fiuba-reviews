import type { Metadata } from "next";
import { Suspense } from "react";
import { getMateria } from "@/server/db/queries/materias";
import { MateriaShell } from "./materia-shell";
import { SidebarMateria, TituloMateria } from "./sidebar-materia";

export async function generateMetadata({
  params
}: LayoutProps<"/materias/[codigo]">): Promise<Metadata> {
  const { codigo } = await params;
  const materia = await getMateria(codigo);

  if (!materia) return {};

  return { title: `${materia.codigo} • ${materia.nombre}` };
}

/**
 * Exige que toda la ruta se genere en el build. Los datos se actualizan
 * invalidando por tag (ver `server/cache.ts`).
 */
export const ensureStatic = "navigation";

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
