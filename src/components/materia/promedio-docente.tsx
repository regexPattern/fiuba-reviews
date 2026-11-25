"use client";

import type { PromediosDocente } from "@/server/db/queries/catedras";
import { Popover } from "@base-ui/react/popover";
import { ChevronDown } from "@withicons/react/line";
import { Star } from "@withicons/react/plush";
import { Numero } from "@/components/ui/numero";
import { CRITERIOS } from "@/lib/criterios";

/** Botón con el promedio general del docente que abre el desglose por criterio. */
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
        aria-label={`Promedio ${promedios.general.toFixed(1)}. Ver desglose por criterio`}
        className="group boton-secundario"
      >
        <Numero>{promedios.general.toFixed(1)}</Numero>
        <Star size={18} aria-hidden="true" className="text-fiuba" />
        {/* Indica que abre el desglose. Del set `line` y en gris: al lado de la estrella, otro
            ícono a color competiría con la nota. Gira cuando el popover está abierto. */}
        <ChevronDown
          size={14}
          aria-hidden="true"
          className="text-stone-400 transition-transform duration-150 group-data-popup-open:rotate-180 motion-reduce:transition-none"
        />
      </Popover.Trigger>
      <Popover.Portal>
        <Popover.Positioner side="bottom" align="end" sideOffset={6} className="z-50">
          <Popover.Popup className="w-64 origin-(--transform-origin) relieve bg-white p-3 text-sm transition-[transform,opacity] duration-150 outline-none data-ending-style:scale-95 data-ending-style:opacity-0 data-starting-style:scale-95 data-starting-style:opacity-0">
            <Popover.Title className="flex items-baseline justify-between gap-4 border-b border-stone-200 pb-2 text-tinta">
              <span className="font-acento text-base/5">Promedios</span>
              <span className="rotulo text-stone-500">
                <Numero>{cantidadCalificaciones}</Numero>{" "}
                {cantidadCalificaciones === 1 ? "calificación" : "calificaciones"}
              </span>
            </Popover.Title>
            <dl className="space-y-1.5 pt-2">
              {CRITERIOS.map(({ nombre, etiqueta }) => (
                <div key={nombre} className="flex items-center justify-between gap-4">
                  <dt className="text-stone-500">{etiqueta}</dt>
                  <dd className={colorPuntaje(promedios[nombre])}>
                    <Numero>{promedios[nombre].toFixed(1)}</Numero>
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

/** Verde para lo bueno, rojo para lo malo, como los `+`/`−` de un diff. */
function colorPuntaje(puntaje: number) {
  if (puntaje >= 4) return "text-green-700";
  if (puntaje < 3) return "text-red-700";
  return "text-stone-700";
}
