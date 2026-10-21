import { Star } from "@withicons/react/chrome";
import { Skeleton } from "@/components/ui/skeleton";
import styles from "./skeletons.module.css";

export function SidebarSkeleton() {
  return (
    <div className="min-h-0 flex-1 overflow-hidden p-3">
      <output className="sr-only">Cargando cátedras…</output>
      <div aria-hidden="true" className="space-y-3">
        {Array.from({ length: 8 }, (_, index) => (
          <div key={index} className="flex h-9 items-center gap-2 px-3">
            <span className="inline-flex w-12 shrink-0 items-center justify-end gap-1">
              <Skeleton className="h-4 w-6" />
              <Star size={16} className="gold-star shrink-0" />
            </span>
            <Skeleton className="h-4 flex-1" />
          </div>
        ))}
      </div>
    </div>
  );
}

export function CatedraSkeleton() {
  return (
    <div className={styles.canvas}>
      <output className="sr-only">Cargando información de la cátedra…</output>
      <div aria-hidden="true" className={styles.stars}>
        {Array.from({ length: 480 }, (_, index) => (
          <span key={index} className={styles.star}>
            <Star size={18} className="gold-star" />
          </span>
        ))}
      </div>
    </div>
  );
}
