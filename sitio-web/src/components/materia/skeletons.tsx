import { Skeleton } from "@/components/ui/skeleton";

export function SidebarSkeleton() {
  return (
    <div className="flex min-h-0 flex-1 flex-col divide-y divide-layout-border">
      <ul className="my-2 min-h-0 flex-1 overflow-hidden">
        {Array.from({ length: 10 }, (_, i) => (
          <li key={i} className="flex h-[50px] items-center gap-2 p-3">
            <Skeleton className="h-4 w-[3ch]" />
            <Skeleton className="h-4 flex-1" />
          </li>
        ))}
      </ul>
      <div className="h-[50px] shrink-0" />
    </div>
  );
}

export function CatedraSkeleton() {
  return (
    <div className="m-4 space-y-8 md:m-6" aria-hidden="true">
      <section className="space-y-3">
        <div>
          <Skeleton className="h-10 w-3/4 max-w-md" />
        </div>
        <div className="flex gap-2">
          <Skeleton className="h-[38px] w-40" />
          <Skeleton className="h-[38px] w-28" />
        </div>
        <div className="divide-y divide-layout-border">
          {Array.from({ length: 4 }, (_, i) => (
            <div key={i} className="py-3">
              <Skeleton className="h-[48px] w-full" />
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}
