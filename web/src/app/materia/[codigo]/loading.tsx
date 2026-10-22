import { MateriaShell } from "@/components/materia/materia-shell";
import {
  CatedraSkeleton,
  MateriaHeaderSkeleton,
  SidebarSkeleton
} from "@/components/materia/skeletons";

export default function Loading() {
  return (
    <MateriaShell
      nombreMateria="Cargando materia…"
      sidebar={
        <aside className="flex h-full min-h-0 w-full flex-col">
          <MateriaHeaderSkeleton />
          <SidebarSkeleton />
        </aside>
      }
    >
      <CatedraSkeleton />
    </MateriaShell>
  );
}
