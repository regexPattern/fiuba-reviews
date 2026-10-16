"use client";

import type { ComentarioDocente } from "@/server/queries/catedra";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";

const MAX_VISIBLES = 10;

function ComentarioItem({ comentario }: { comentario: ComentarioDocente }) {
  return (
    <div className="space-x-0.5 py-3">
      <p className="inline w-fit before:content-['&quot;'] after:content-['&quot;']">
        {comentario.contenido.trim()}
      </p>

      <div className="inline-flex items-center gap-2 text-xs text-button-foreground select-none">
        <Badge variant="outline" className="border-fiuba/30 bg-fiuba/10">
          {comentario.cuatrimestre.numero}C{comentario.cuatrimestre.anio}
        </Badge>
        {comentario.esDeDolly && (
          <Badge variant="outline" className="border-fiuba/30 bg-fiuba/10">
            Dolly FIUBA
          </Badge>
        )}
      </div>
    </div>
  );
}

export function Comentarios({ comentarios }: { comentarios: ComentarioDocente[] }) {
  const [expandido, setExpandido] = useState(false);

  const visibles = comentarios.slice(0, MAX_VISIBLES);
  const ocultos = comentarios.slice(MAX_VISIBLES);

  return (
    <Collapsible open={expandido} onOpenChange={setExpandido}>
      <div className="divide-y divide-layout-border">
        {visibles.map((comentario) => (
          <ComentarioItem key={comentario.codigo} comentario={comentario} />
        ))}

        {ocultos.length > 0 && (
          <>
            <CollapsibleContent className="divide-y divide-layout-border">
              {ocultos.map((comentario) => (
                <ComentarioItem key={comentario.codigo} comentario={comentario} />
              ))}
            </CollapsibleContent>

            <div className="w-full text-center">
              <CollapsibleTrigger className="mt-2 w-full cursor-pointer p-2 text-sm text-secondary-foreground hover:underline md:w-fit">
                {expandido
                  ? "Mostrar menos"
                  : `Mostrar ${ocultos.length} comentario${ocultos.length === 1 ? "" : "s"} más`}
              </CollapsibleTrigger>
            </div>
          </>
        )}
      </div>
    </Collapsible>
  );
}
