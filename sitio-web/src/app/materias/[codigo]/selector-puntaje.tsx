"use client";

import { Field } from "@base-ui/react/field";
import { Fieldset } from "@base-ui/react/fieldset";
import { Radio } from "@base-ui/react/radio";
import { RadioGroup } from "@base-ui/react/radio-group";
import { Star } from "@withicons/react/plush";
import { useState } from "react";
import { Numero } from "@/components/ui/numero";
import {
  PUNTAJE_MAXIMO,
  PUNTAJE_MINIMO,
  type Puntaje,
  formatearPuntaje
} from "@/lib/calificaciones";

const ESTRELLAS = Array.from({ length: PUNTAJE_MAXIMO }, (_, i) => i + 1);

/**
 * Puntaje de 1 a 5 con estrellas. Cada estrella tiene dos mitades clickeables (media y entera),
 * salvo la primera, que entera vale 1 (no hay 0,5).
 * Es un radio group de Base UI: con las flechas se mueve de a medio punto, y con el mouse se
 * previsualiza el puntaje antes de elegirlo.
 */
export function SelectorPuntaje({
  etiqueta,
  valor,
  onChange,
  invalido
}: {
  etiqueta: string;
  valor: Puntaje | undefined;
  onChange: (valor: Puntaje) => void;
  invalido?: boolean;
}) {
  const [resaltado, setResaltado] = useState<Puntaje | null>(null);
  const mostrado = resaltado ?? valor ?? 0;

  // El grupo es un `Fieldset` con `RadioGroup` adentro y la leyenda como etiqueta, como indica la
  // guía de formularios de Base UI.
  return (
    <Field.Root invalid={invalido} className="group">
      <Fieldset.Root
        render={
          <RadioGroup
            value={valor ?? null}
            onValueChange={(valor) => onChange(valor as Puntaje)}
            onPointerLeave={() => setResaltado(null)}
          />
        }
        className="flex items-center gap-3 py-0.5 md:gap-4 md:py-0.5"
      >
        <Fieldset.Legend className="min-w-0 flex-1 text-sm text-stone-700 group-data-invalid:text-red-700">
          {etiqueta}
        </Fieldset.Legend>
        <div className="flex shrink-0">
          {ESTRELLAS.map((estrella) => (
            <span
              key={estrella}
              // `size-8`: área táctil cómoda (dos mitades de 16px) aunque la estrella sea más chica.
              className="relative flex size-8 items-center justify-center rounded-md has-focus-visible:outline-2 has-focus-visible:-outline-offset-2 has-focus-visible:outline-fiuba"
            >
              <EstrellaPuntaje relleno={Math.min(Math.max(mostrado - (estrella - 1), 0), 1)} />
              {[estrella - 0.5, estrella]
                .filter((puntaje) => puntaje >= PUNTAJE_MINIMO)
                .map((puntaje, _, mitades) => (
                  <Radio.Root
                    key={puntaje}
                    value={puntaje}
                    aria-label={`${formatearPuntaje(puntaje)} de ${PUNTAJE_MAXIMO}`}
                    onPointerEnter={(event) =>
                      event.pointerType === "mouse" && setResaltado(puntaje)
                    }
                    className={`absolute inset-y-0 cursor-pointer outline-none ${
                      mitades.length === 1
                        ? "inset-x-0"
                        : puntaje % 1 === 0
                          ? "right-0 w-1/2"
                          : "left-0 w-1/2"
                    }`}
                  />
                ))}
            </span>
          ))}
        </div>
        {/* Con medias estrellas, el número evita tener que adivinar. */}
        <span aria-hidden="true" className="w-6 shrink-0 text-right text-sm text-stone-500">
          <Numero>{mostrado > 0 ? formatearPuntaje(mostrado) : "–"}</Numero>
        </span>
      </Fieldset.Root>
    </Field.Root>
  );
}

/**
 * La misma estrella `plush` dos veces: abajo vacía (rellenos en blanco, conserva el contorno y las
 * sombras) y arriba la dorada, recortada a la mitad si es media estrella.
 */
function EstrellaPuntaje({ relleno }: { relleno: number }) {
  return (
    <span aria-hidden="true" className="pointer-events-none relative size-6">
      <Star
        size={24}
        className="absolute inset-0 [--with-plush-c1:var(--color-white)] [--with-plush-c2:var(--color-white)]"
      />
      {relleno > 0 && (
        <Star
          size={24}
          style={{ clipPath: `inset(0 ${(1 - relleno) * 100}% 0 0)` }}
          className="absolute inset-0 [--with-plush-c2:var(--color-amber-400)]"
        />
      )}
    </span>
  );
}
