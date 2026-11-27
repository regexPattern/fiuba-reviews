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

  return { title: `${materia.codigo} • ${materia.nombre} | FIUBA Reviews` };
}

/**
 * Las materias y cátedras solo cambian con un deploy, así que toda la ruta se prerenderiza
 * completa en el build. Una URL que no está en `generateStaticParams` espera el render completo
 * (sin fallback shell): si no existe devuelve un 404 real y no un 200 con `notFound()` streameado.
 *
 * Los `<Suspense>` siguen haciendo falta por Partial Prefetching: el App Shell se comparte entre
 * todas las URLs de la ruta, así que no puede incluir nada que dependa de `params`. Sus fallbacks
 * se ven al navegar con un link sin `prefetch={true}` (p. ej. desde el buscador).
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
