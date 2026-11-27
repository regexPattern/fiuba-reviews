"use client";

import type { Cuatrimestre } from "@/lib/ofertas";
import { Button } from "@base-ui/react/button";
import { Field } from "@base-ui/react/field";
import { Fieldset } from "@base-ui/react/fieldset";
import { Form } from "@base-ui/react/form";
import { Select } from "@base-ui/react/select";
import { Check, ChevronsUpDown, LoaderCircle, Send } from "@withicons/react/line";
import { type ElementType, type FormEvent, type ReactNode, useMemo, useRef, useState } from "react";
import { SelectorPuntaje } from "@/components/calificacion/selector-puntaje";
import { Numero } from "@/components/ui/numero";
import {
  CRITERIOS,
  LARGO_MAXIMO_COMENTARIO,
  type BorradorCalificacion,
  criteriosFaltantes
} from "@/lib/calificaciones";

export function FormularioCalificacion({
  borrador,
  cuatrimestres: cuatrimestresCalificables,
  onChange,
  onEnviar,
  enviando = false,
  error = null,
  cuerpo: Cuerpo = "div"
}: {
  borrador: BorradorCalificacion;
  /** Los que se pueden elegir, del más nuevo al más viejo. */
  cuatrimestres: Cuatrimestre[];
  onChange: (borrador: BorradorCalificacion) => void;
  onEnviar: (borrador: BorradorCalificacion) => void;
  enviando?: boolean;
  /** Error del servidor al enviar. */
  error?: string | null;
  /** Contenedor del cuerpo: en el drawer, `Drawer.Content`. */
  cuerpo?: ElementType<{ className?: string; children?: ReactNode }>;
}) {
  const formRef = useRef<HTMLFormElement>(null);
  const selectRef = useRef<HTMLButtonElement>(null);
  // Los errores se muestran recién después del primer intento de enviar.
  const [intentoEnviar, setIntentoEnviar] = useState(false);
  const cuatrimestres = useMemo(
    () =>
      cuatrimestresCalificables.map((cuatrimestre) => ({
        value: claveCuatrimestre(cuatrimestre),
        // Igual que en los comentarios: `1C2026`.
        label: `${cuatrimestre.numero}C${cuatrimestre.anio}`,
        cuatrimestre
      })),
    [cuatrimestresCalificables]
  );

  const faltantes = criteriosFaltantes(borrador);
  const faltaCuatrimestre = borrador.cuatrimestre === null;
  const largoComentario = borrador.comentario.length;

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (enviando) return;
    setIntentoEnviar(true);

    if (faltantes.length > 0) {
      formRef.current
        ?.querySelector<HTMLElement>(`[data-criterio="${faltantes[0]}"] [role="radio"]`)
        ?.focus();
      return;
    }
    if (faltaCuatrimestre) {
      selectRef.current?.focus();
      return;
    }
    onEnviar({ ...borrador, comentario: borrador.comentario.trim() });
  }

  const errorCriterios =
    intentoEnviar && faltantes.length > 0
      ? faltantes.length === 1
        ? "Falta calificar 1 criterio."
        : `Faltan calificar ${faltantes.length} criterios.`
      : null;

  return (
    <Form ref={formRef} onSubmit={handleSubmit} className="flex min-h-0 flex-1 flex-col">
      <Cuerpo className="min-h-0 flex-1 space-y-6 overflow-y-auto overscroll-contain p-4 md:grid md:grid-cols-[minmax(0,1fr)_minmax(0,1fr)] md:gap-6 md:space-y-0 md:px-6 md:py-6">
        <Fieldset.Root>
          <Fieldset.Legend className="mb-1 rotulo text-stone-500 md:mb-2">
            Criterios
          </Fieldset.Legend>
          {/* En desktop sobra lugar: sin divisores y con aire entre criterios, como la página de
              calificar del sitio anterior. */}
          <div className="divide-y divide-stone-100 md:space-y-1 md:divide-y-0">
            {CRITERIOS.map(({ nombre, etiqueta }) => (
              <div key={nombre} data-criterio={nombre}>
                <SelectorPuntaje
                  etiqueta={etiqueta}
                  valor={borrador.puntajes[nombre]}
                  invalido={intentoEnviar && faltantes.includes(nombre)}
                  onChange={(puntaje) =>
                    onChange({ ...borrador, puntajes: { ...borrador.puntajes, [nombre]: puntaje } })
                  }
                />
              </div>
            ))}
          </div>
        </Fieldset.Root>

        <div className="flex flex-col gap-6 md:border-l md:border-stone-200 md:pl-6">
          <Field.Root
            invalid={intentoEnviar && faltaCuatrimestre}
            className="space-y-1.5 md:space-y-2"
          >
            {/* No es un `<label>` nativo: el control es un botón. */}
            <Field.Label nativeLabel={false} render={<div />} className="rotulo text-stone-500">
              Cuatrimestre que cursaste
            </Field.Label>
            <Select.Root
              items={cuatrimestres}
              value={borrador.cuatrimestre ? claveCuatrimestre(borrador.cuatrimestre) : null}
              onValueChange={(valor) =>
                onChange({
                  ...borrador,
                  cuatrimestre: cuatrimestres.find((c) => c.value === valor)?.cuatrimestre ?? null
                })
              }
            >
              <Select.Trigger
                ref={selectRef}
                className="flex h-10 w-full cursor-pointer items-center justify-between gap-2 campo text-left text-sm text-stone-900 data-invalid:border-red-400 data-popup-open:border-fiuba"
              >
                <Select.Value
                  placeholder="Elegí un cuatrimestre"
                  className="truncate font-mono tabular-nums data-placeholder:font-sans data-placeholder:text-stone-400"
                />
                <Select.Icon>
                  <ChevronsUpDown size={16} aria-hidden="true" className="text-stone-500" />
                </Select.Icon>
              </Select.Trigger>
              <Select.Portal>
                {/* `z-60`: por encima del dialog y del drawer (`z-50`). */}
                {/* Sin `alignItemWithTrigger`: con él, el popup tapa el botón y se corre a la izquierda
                    para alinear el texto de la opción con el del botón. Con dos opciones, un
                    desplegable común abajo y del mismo ancho se entiende mejor. */}
                <Select.Positioner
                  alignItemWithTrigger={false}
                  sideOffset={6}
                  className="z-60 outline-none select-none"
                >
                  <Select.Popup className="w-(--anchor-width) origin-(--transform-origin) relieve bg-white p-1 text-sm text-stone-900 transition-[scale,opacity] duration-100 ease-out outline-none data-ending-style:scale-97 data-ending-style:opacity-0 data-starting-style:scale-97 data-starting-style:opacity-0 motion-reduce:transition-none">
                    <Select.List className="max-h-(--available-height) overflow-y-auto">
                      {cuatrimestres.map(({ value, label }) => (
                        <Select.Item
                          key={value}
                          value={value}
                          className="grid cursor-pointer grid-cols-[1rem_1fr] items-center gap-2 rounded-md py-2 pr-4 pl-2 outline-none data-highlighted:bg-stone-100"
                        >
                          <Select.ItemIndicator className="col-start-1">
                            <Check size={16} aria-hidden="true" className="text-fiuba" />
                          </Select.ItemIndicator>
                          <Select.ItemText className="col-start-2 font-mono tabular-nums">
                            {label}
                          </Select.ItemText>
                        </Select.Item>
                      ))}
                    </Select.List>
                  </Select.Popup>
                </Select.Positioner>
              </Select.Portal>
            </Select.Root>
            {/* `match`: la validez la maneja el formulario, no la validación nativa. */}
            {intentoEnviar && faltaCuatrimestre && (
              <Field.Error match className="text-sm text-red-700">
                Elegí el cuatrimestre en que cursaste.
              </Field.Error>
            )}
          </Field.Root>

          <Field.Root className="flex flex-1 flex-col gap-1.5 md:gap-2">
            <div className="flex items-baseline justify-between gap-2">
              <Field.Label className="rotulo text-stone-500">
                Comentario <span className="tracking-normal normal-case">(opcional)</span>
              </Field.Label>
              {/* Solo cerca del límite, para no distraer. */}
              {largoComentario > LARGO_MAXIMO_COMENTARIO * 0.8 && (
                <span className="text-xs text-stone-400">
                  <Numero>
                    {largoComentario}/{LARGO_MAXIMO_COMENTARIO}
                  </Numero>
                </span>
              )}
            </div>
            <Field.Control
              render={<textarea rows={4} maxLength={LARGO_MAXIMO_COMENTARIO} />}
              value={borrador.comentario}
              onValueChange={(comentario) => onChange({ ...borrador, comentario })}
              placeholder="¿Cómo fue cursar con este docente?"
              className="block min-h-24 w-full flex-1 resize-none campo py-2 text-stone-900 md:min-h-44 md:text-sm"
            />
            <Field.Description className="text-xs text-stone-500">
              Se publica de forma anónima.
            </Field.Description>
          </Field.Root>
        </div>
      </Cuerpo>

      <footer className="flex shrink-0 items-center justify-end gap-3 border-t border-stone-200 bg-stone-50 px-4 pt-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] md:px-6">
        <p aria-live="polite" className="min-w-0 flex-1 text-sm text-red-700">
          {errorCriterios ?? error}
        </p>
        <Button type="submit" disabled={enviando} focusableWhenDisabled className="boton-primario">
          {enviando ? (
            <LoaderCircle size={16} aria-hidden="true" className="shrink-0 animate-spin" />
          ) : (
            <Send size={16} aria-hidden="true" className="shrink-0" />
          )}
          {enviando ? "Enviando…" : "Enviar"}
        </Button>
      </footer>
    </Form>
  );
}

function claveCuatrimestre({ numero, anio }: Cuatrimestre) {
  return `${anio}-${numero}`;
}
