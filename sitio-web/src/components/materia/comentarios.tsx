"use client";

import type { ComentarioDocente } from "@/server/queries/catedra";
import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";

const MAX_VISIBLES = 10;

function ComentarioItem({ comentario }: { comentario: ComentarioDocente }) {
  return (
    <div className="space-y-2 py-3">
      <p className="break-words whitespace-pre-wrap">{comentario.contenido.trim()}</p>

      <div className="inline-flex items-center gap-2 text-xs text-muted-foreground">
        <Badge variant="outline">
          {comentario.cuatrimestre.numero}C{comentario.cuatrimestre.anio}
        </Badge>
        {comentario.deDolly && <Badge variant="outline">Dolly FIUBA</Badge>}
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
      <div className="divide-y">
        {visibles.map((comentario) => (
          <ComentarioItem key={comentario.codigo} comentario={comentario} />
        ))}

        {ocultos.length > 0 && (
          <>
            <CollapsibleContent className="divide-y">
              {ocultos.map((comentario) => (
                <ComentarioItem key={comentario.codigo} comentario={comentario} />
              ))}
            </CollapsibleContent>

            <div className="w-full text-center">
              <CollapsibleTrigger asChild>
                <Button variant="ghost" className="mt-2">
                  {expandido
                    ? "Mostrar menos"
                    : `Mostrar ${ocultos.length} comentario${ocultos.length === 1 ? "" : "s"} más`}
                </Button>
              </CollapsibleTrigger>
            </div>
          </>
        )}
      </div>
    </Collapsible>
  );
}
