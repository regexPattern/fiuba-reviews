"use client";

import type { PromediosDocente } from "@/server/queries/catedra";
import { ChevronDown, Star } from "lucide-react";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { CRITERIOS } from "@/lib/criterios";
import { cn } from "@/lib/utils";

const MIN = 1;
const MAX = 5;

function Criterio({ label, valor }: { label: string; valor: number }) {
  const porcentaje = Math.max(0, Math.min(100, ((valor - MIN) / (MAX - MIN)) * 100));

  return (
    <div className="relative overflow-hidden px-2 py-1">
      <div className="absolute inset-0 bg-fiuba/5" />
      <div className="absolute inset-y-0 left-0 bg-fiuba/20" style={{ width: `${porcentaje}%` }} />
      <div className="relative z-10 grid grid-cols-[1fr_auto] items-center gap-3 text-foreground tabular-nums">
        <span>{label}</span>
        <span>{valor.toFixed(1)}</span>
      </div>
    </div>
  );
}

export function Promedios({
  promedios,
  cantidadCalificaciones
}: {
  promedios: PromediosDocente | null;
  cantidadCalificaciones: number;
}) {
  return (
    <Popover>
      <PopoverTrigger
        disabled={!promedios}
        className={cn(
          "flex items-center gap-2 border border-button-border bg-button-background px-3 py-2 transition-colors",
          promedios ? "hover:bg-button-hover" : "pointer-events-none"
        )}
      >
        <Star className="size-4 fill-yellow-500 stroke-yellow-700 dark:stroke-yellow-400" />
        <span>Promedio: {promedios ? promedios.general.toFixed(1) : "–"}</span>
        {promedios && <ChevronDown className="size-4" />}
      </PopoverTrigger>

      {promedios && (
        <PopoverContent align="start" sideOffset={6} className="w-56 p-3 text-sm">
          <div className="divide-y-2 divide-background">
            {CRITERIOS.map(({ name, label }) => (
              <Criterio key={name} label={label} valor={promedios[name]} />
            ))}
          </div>

          <div className="pt-3 text-center text-secondary-foreground">
            {cantidadCalificaciones} calificacion{cantidadCalificaciones === 1 ? "" : "es"}
          </div>
        </PopoverContent>
      )}
    </Popover>
  );
}
