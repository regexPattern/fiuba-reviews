import type { Metadata } from "next";
import { Suspense } from "react";
import { SidebarMateria } from "@/components/materia/sidebar-materia";
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
    <div className="mx-auto flex max-w-6xl flex-col md:h-[calc(100dvh-var(--header-height))] md:flex-row">
      <aside className="flex shrink-0 flex-col border-b md:w-80 md:border-r md:border-b-0">
        <Suspense fallback={<p className="p-4">Cargando materia…</p>}>
          <SidebarMateria params={params} />
        </Suspense>
      </aside>
      <main className="min-w-0 flex-1 md:overflow-y-auto">{children}</main>
    </div>
  );
}
