"use client";

import type { PromediosDocente } from "@/server/db/queries/catedra";
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
        <Button variant="outline" disabled={!promedios}>
          Promedio: {promedios ? promedios.general.toFixed(1) : "–"}
          {promedios && <ChevronDown />}
        </Button>
      </PopoverTrigger>
      {promedios && (
        <PopoverContent align="start" className="space-y-4 text-sm">
          <dl className="space-y-2">
            {CRITERIOS.map(({ name, label }) => (
              <div key={name} className="flex justify-between gap-4">
                <dt>{label}</dt>
                <dd className="tabular-nums">{promedios[name].toFixed(1)}</dd>
              </div>
            ))}
          </dl>
          <p className="text-muted-foreground">
            {cantidadCalificaciones} calificacion{cantidadCalificaciones === 1 ? "" : "es"}
          </p>
        </PopoverContent>
      )}
    </Popover>
  );
}
