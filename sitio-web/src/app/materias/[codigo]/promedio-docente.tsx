"use client";

import type { PromediosDocente } from "@/server/db/queries/catedras";
import { Popover } from "@base-ui/react/popover";
import { ChevronDown } from "@withicons/react/line";
import { Star } from "@withicons/react/plush";
import { CRITERIOS, nivelPuntaje, type NivelPuntaje } from "@/lib/calificaciones";

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
        <span className="font-mono tabular-nums">{promedios.general.toFixed(1)}</span>
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
            <Popover.Title className="border-b border-stone-200 pb-2 text-center rotulo text-stone-500">
              <span className="font-mono tabular-nums">{cantidadCalificaciones}</span>{" "}
              {cantidadCalificaciones === 1 ? "calificación" : "calificaciones"}
            </Popover.Title>
            <dl className="space-y-1.5 pt-2">
              {CRITERIOS.map(({ nombre, etiqueta }) => (
                <div key={nombre} className="flex items-center justify-between gap-4">
                  <dt className="text-stone-500">{etiqueta}</dt>
                  <dd
                    className={`flex items-center gap-1.5 ${COLOR_NUMERO[nivelPuntaje(promedios[nombre])]}`}
                  >
                    <span className="font-mono tabular-nums">{promedios[nombre].toFixed(1)}</span>
                    {/* Decorativa: el número ya lo dice. */}
                    <Star
                      size={16}
                      aria-hidden="true"
                      className={`shrink-0 ${COLOR_ESTRELLA[nivelPuntaje(promedios[nombre])]}`}
                    />
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
const COLOR_NUMERO: Record<NivelPuntaje, string> = {
  muyBueno: "text-green-700",
  bueno: "text-green-700",
  regular: "text-stone-700",
  malo: "text-red-700"
};

/** Relleno de la estrella, de dorado a rojo a medida que baja el nivel. El relleno del set
    `plush` es `--with-plush-c2`, así que lo pisamos en cada nivel. */
const COLOR_ESTRELLA: Record<NivelPuntaje, string> = {
  muyBueno: "[--with-plush-c2:var(--color-amber-400)]",
  bueno: "[--with-plush-c2:var(--color-yellow-300)]",
  regular: "[--with-plush-c2:var(--color-orange-500)]",
  malo: "[--with-plush-c2:var(--color-red-600)]"
};
