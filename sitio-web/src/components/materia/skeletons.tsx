import { Skeleton } from "@/components/ui/skeleton";

export function SidebarSkeleton() {
  return (
    <div className="min-h-0 flex-1 overflow-hidden p-3">
      <output className="sr-only">Cargando cátedras…</output>
      <div aria-hidden="true" className="space-y-3">
        {Array.from({ length: 8 }, (_, index) => (
          <Skeleton key={index} className="h-9 w-full" />
        ))}
      </div>
    </div>
  );
}

export function CatedraSkeleton() {
  return (
    <div className="p-4 md:p-6">
      <output className="sr-only">Cargando docentes…</output>
      <div aria-hidden="true" className="space-y-4">
        <Skeleton className="h-7 w-2/3 max-w-sm" />
        <Skeleton className="h-4 w-24" />
        <Skeleton className="h-9 w-36" />
        <div className="divide-y">
          {Array.from({ length: 4 }, (_, index) => (
            <div key={index} className="space-y-2 py-3">
              <Skeleton className="h-4 w-full" />
              <Skeleton className="h-4 w-3/4" />
              <Skeleton className="h-5 w-20" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
