import { MateriaShell } from "@/components/materia/materia-shell";
import { CatedraSkeleton, SidebarSkeleton } from "@/components/materia/skeletons";
import { Skeleton } from "@/components/ui/skeleton";

export default function Loading() {
  return (
    <MateriaShell
      nombreMateria="Cargando materia…"
      sidebar={
        <aside className="flex h-full min-h-0 w-full flex-col">
          <output className="sr-only">Cargando materia…</output>
          <div aria-hidden="true" className="space-y-3 border-b p-4">
            <Skeleton className="h-6 w-3/4" />
            <Skeleton className="h-4 w-full" />
          </div>
          <SidebarSkeleton />
        </aside>
      }
    >
      <CatedraSkeleton />
    </MateriaShell>
  );
}
