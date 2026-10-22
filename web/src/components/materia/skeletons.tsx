import styles from "./skeletons.module.css";

function StarsSkeleton({ label, className }: { label: string; className: string }) {
  return (
    <div className={className}>
      <output className="sr-only">{label}</output>
      <div aria-hidden="true" className={styles.stars} />
    </div>
  );
}

export function MateriaHeaderSkeleton() {
  return (
    <StarsSkeleton
      label="Cargando información de la materia…"
      className="pointer-events-none relative h-36 shrink-0 overflow-hidden border-b"
    />
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
