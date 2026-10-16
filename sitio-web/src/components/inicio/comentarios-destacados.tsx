import type { Inicio } from "@/server/queries/inicio";
import { cn } from "@/lib/utils";

type Comentario = Inicio["comentarios"][number];

function ComentarioCard({ comentario, className }: { comentario: Comentario; className?: string }) {
  return (
    <article className={cn("border border-button-border bg-button-background/50 p-4", className)}>
      <p className="line-clamp-4 before:content-['&quot;'] after:content-['&quot;'] md:line-clamp-none">
        {comentario.contenido}
      </p>
      <p className="mt-2 text-sm text-foreground/75">
        {comentario.nombreDocente} • {comentario.nombreMateria}
      </p>
    </article>
  );
}

function Fila({ comentarios, reverse = false }: { comentarios: Comentario[]; reverse?: boolean }) {
  const animacion = reverse
    ? "animate-scroll-horizontal sm:animate-scroll-horizontal-sm"
    : "animate-scroll-horizontal-reverse sm:animate-scroll-horizontal-reverse-sm";

  return (
    <div className="overflow-hidden">
      <div className={cn("flex w-max gap-4", animacion)}>
        {[...comentarios, ...comentarios].map((comentario, i) => (
          <ComentarioCard
            key={`${comentario.codigo}-${i}`}
            comentario={comentario}
            className="max-w-65 min-w-65 shrink-0"
          />
        ))}
      </div>
    </div>
  );
}

function Columna({ comentarios }: { comentarios: Comentario[] }) {
  return (
    <div className="h-157.5 overflow-hidden">
      <div className="flex animate-scroll-vertical flex-col gap-4 md:animate-scroll-vertical-md">
        {[...comentarios, ...comentarios].map((comentario, i) => (
          <ComentarioCard key={`${comentario.codigo}-${i}`} comentario={comentario} />
        ))}
      </div>
    </div>
  );
}

export function ComentariosDestacados({ comentarios }: { comentarios: Comentario[] }) {
  const mitad = Math.ceil(comentarios.length / 2);

  return (
    <section id="ultimos-comentarios" className="space-y-4 md:h-157.5">
      <div className="grid gap-4 md:hidden">
        <Fila comentarios={comentarios.filter((_, i) => i % 2 === 0)} />
        <Fila comentarios={comentarios.filter((_, i) => i % 2 === 1)} reverse />
      </div>

      <div className="hidden grid-cols-2 gap-4 md:grid md:h-157.5">
        <Columna comentarios={comentarios.slice(0, mitad)} />
        <Columna comentarios={comentarios.slice(mitad)} />
      </div>
    </section>
  );
}
