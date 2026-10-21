import { Star } from "@withicons/react/chrome";
import styles from "./skeletons.module.css";

function StarsSkeleton({ label, className }: { label: string; className: string }) {
  return (
    <div className={className}>
      <output className="sr-only">{label}</output>
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

export function SidebarSkeleton() {
  return (
    <StarsSkeleton
      label="Cargando cátedras…"
      className="pointer-events-none relative h-full min-h-0 flex-1 overflow-hidden"
    />
  );
}

export function CatedraSkeleton() {
  return <StarsSkeleton label="Cargando información de la cátedra…" className={styles.canvas} />;
}
