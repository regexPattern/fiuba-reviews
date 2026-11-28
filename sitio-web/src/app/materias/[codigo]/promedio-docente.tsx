"use client";

import type { PromediosDocente } from "@/server/db/queries/catedras";
import { Popover } from "@base-ui/react/popover";
import { ChevronDown } from "@withicons/react/line";
import { Star } from "@withicons/react/plush";
import { CRITERIOS } from "@/lib/calificaciones";
import { formatearPuntaje } from "@/lib/utils";

export function PromedioDocente({
  promedios,
  cantidadCalificaciones
}: {
  promedios: PromediosDocente;
  cantidadCalificaciones: number;
}) {
  return (
    <Popover.Root>
      <Popover.Trigger
        aria-label={`Promedio ${formatearPuntaje(promedios.general)}. Ver desglose por criterio`}
        className="group boton-secundario"
      >
        <span className="font-mono tabular-nums">{formatearPuntaje(promedios.general)}</span>
        <Star size={18} aria-hidden="true" className="text-fiuba" />
        <ChevronDown
          size={14}
          aria-hidden="true"
          className="text-stone-400 transition-transform duration-150 group-data-popup-open:rotate-180 motion-reduce:transition-none"
        />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Positioner side="bottom" align="end" sideOffset={6} className="z-50">
          <Popover.Popup className="w-64 origin-(--transform-origin) relieve bg-white p-3 text-sm transition-[transform,opacity] duration-150 outline-none data-ending-style:scale-95 data-ending-style:opacity-0 data-starting-style:scale-95 data-starting-style:opacity-0">
            <Popover.Title className="border-b border-stone-200 pb-2 text-center rotulo text-stone-500">
              <span className="font-mono tabular-nums">{cantidadCalificaciones}</span>{" "}
              {cantidadCalificaciones === 1 ? "calificación" : "calificaciones"}
            </Popover.Title>
            <dl className="space-y-1.5 pt-2">
              {CRITERIOS.map(({ nombre, label }) => (
                <div key={nombre} className="flex items-center justify-between gap-4">
                  <dt className="text-stone-500">{label}</dt>
                  <dd
                    className={`flex items-center gap-1.5 ${COLOR_NUMERO[nivelPuntaje(promedios[nombre])]}`}
                  >
                    <span className="font-mono tabular-nums">
                      {formatearPuntaje(promedios[nombre])}
                    </span>
                    <Star size={16} aria-hidden="true" className="shrink-0" />
                  </dd>
                </div>
              ))}
            </dl>
          </Popover.Popup>
        </Popover.Positioner>
      </Popover.Portal>
    </Popover.Root>
  );
}

type NivelPuntaje = "muyBueno" | "bueno" | "regular" | "malo";

function nivelPuntaje(puntaje: number): NivelPuntaje {
  if (puntaje >= 4.5) return "muyBueno";
  if (puntaje >= 3.5) return "bueno";
  if (puntaje >= 2.5) return "regular";
  return "malo";
}

const COLOR_NUMERO: Record<NivelPuntaje, string> = {
  muyBueno: "text-green-700",
  bueno: "text-green-700",
  regular: "text-stone-700",
  malo: "text-red-700"
};
