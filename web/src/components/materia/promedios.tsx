"use client";

import type { PromediosDocente } from "@/server/db/queries/catedra";
import { Star } from "@withicons/react/chrome";
import { ChevronDown } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { CRITERIOS } from "@/lib/criterios";

export function Promedios({
  promedios,
  cantidadCalificaciones
}: {
  promedios: PromediosDocente | null;
  cantidadCalificaciones: number;
}) {
  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline" className="shadow-none select-none" disabled={!promedios}>
          Promedio: {promedios ? promedios.general.toFixed(1) : "–"}
          <Star size={16} className="gold-star" aria-hidden="true" />
          {promedios && <ChevronDown />}
        </Button>
      </PopoverTrigger>
      {promedios && (
        <PopoverContent
          align="start"
          className="w-60 space-y-3 bg-popover/75 text-xs shadow-glass backdrop-blur-surface select-none"
        >
          <dl className="space-y-1">
            {CRITERIOS.map(({ name, label }) => (
              <div
                key={name}
                className="relative isolate flex justify-between gap-3 overflow-hidden rounded-sm border border-border bg-muted/50 px-2 py-1.5"
              >
                <span
                  aria-hidden="true"
                  className="absolute inset-y-0 left-0 -z-10 bg-fiuba/25"
                  style={{ width: `${Math.min(100, Math.max(0, (promedios[name] / 5) * 100))}%` }}
                />
                <dt className="font-medium">{label}</dt>
                <dd className="flex shrink-0 items-center gap-1 font-semibold tabular-nums">
                  {promedios[name].toFixed(1)}
                  <Star size={12} className="gold-star shrink-0" aria-hidden="true" />
                </dd>
              </div>
            ))}
          </dl>
          <p className="text-center text-muted-foreground">
            {cantidadCalificaciones} calificacion{cantidadCalificaciones === 1 ? "" : "es"}
          </p>
        </PopoverContent>
      )}
    </Popover>
  );
}
