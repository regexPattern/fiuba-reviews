"use client";

import { Star } from "@withicons/react/chrome";
import { Plus, X } from "lucide-react";
import { Dialog } from "radix-ui";
import { useId, useState } from "react";
import { Button } from "@/components/ui/button";
import { CRITERIOS, type CriterioName } from "@/lib/criterios";

type ValoresCalificacion = Record<CriterioName, number>;

const VALORES_INICIALES = Object.fromEntries(
  CRITERIOS.map(({ name }) => [name, 3])
) as ValoresCalificacion;

export function AgregarCalificacion({ nombreDocente }: { nombreDocente: string }) {
  const id = useId();
  const [valores, setValores] = useState<ValoresCalificacion>(VALORES_INICIALES);
  const [comentario, setComentario] = useState("");

  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>
        <Button variant="outline" className="shadow-none">
          <Plus aria-hidden="true" /> Agregar calificación
        </Button>
      </Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-50 bg-background/80 backdrop-blur-sm" />
        <Dialog.Content
          aria-describedby={undefined}
          style={{ maxHeight: "calc(100dvh - 2rem)", width: "calc(100% - 2rem)" }}
          className="fixed top-1/2 left-1/2 z-50 flex max-w-4xl -translate-x-1/2 -translate-y-1/2 flex-col overflow-hidden rounded-lg border bg-background shadow-glass outline-none"
        >
          <header className="relative shrink-0 border-b px-4 py-3 pr-14 sm:px-6 sm:pr-14">
            <Dialog.Title className="text-lg font-semibold">
              Calificá a {nombreDocente}
            </Dialog.Title>
            <Dialog.Close asChild>
              <Button
                variant="ghost"
                size="icon"
                className="absolute top-1/2 right-3 -translate-y-1/2"
                aria-label="Cerrar formulario"
              >
                <X aria-hidden="true" />
              </Button>
            </Dialog.Close>
          </header>
          <div className="min-h-0 overflow-y-auto overscroll-contain p-4 sm:p-6">
            <form onSubmit={(event) => event.preventDefault()} className="space-y-6">
              <div className="grid gap-x-3 gap-y-6 md:grid-cols-2 md:gap-y-2">
                <section
                  aria-labelledby={`${id}-puntajes-label`}
                  aria-describedby={`${id}-puntajes-ayuda`}
                  className="grid min-w-0 gap-2 md:row-span-3 md:grid-rows-subgrid"
                >
                  <p id={`${id}-puntajes-label`} className="text-sm font-medium">
                    Puntajes
                  </p>
                  <p id={`${id}-puntajes-ayuda`} className="text-sm text-muted-foreground">
                    Mové las barritas para ajustar cada puntaje del 1 al 5.
                  </p>
                  <div>
                    {CRITERIOS.map(({ name, label }) => (
                      <div
                        key={name}
                        className="relative isolate min-w-0 bg-muted/50 focus-within:z-10 focus-within:ring-2 focus-within:ring-ring focus-within:ring-inset"
                      >
                        <span
                          aria-hidden="true"
                          className="pointer-events-none absolute inset-y-0 left-0 -z-10 bg-fiuba/25"
                          style={{ width: `${(valores[name] / 5) * 100}%` }}
                        />
                        <div className="pointer-events-none flex min-h-12 items-center gap-2 px-3 py-2 text-sm select-none">
                          <label htmlFor={`${id}-${name}`} className="mr-auto font-medium">
                            {label}
                          </label>
                          <output
                            htmlFor={`${id}-${name}`}
                            className="flex shrink-0 items-center gap-2 font-mono text-base font-semibold tracking-tighter tabular-nums"
                          >
                            {valores[name].toFixed(1)}
                            <Star size={18} className="gold-star shrink-0" aria-hidden="true" />
                          </output>
                        </div>
                        <input
                          id={`${id}-${name}`}
                          name={name}
                          type="range"
                          min={1}
                          max={5}
                          step={0.1}
                          value={valores[name]}
                          aria-valuetext={`${valores[name]} de 5`}
                          aria-describedby={`${id}-puntajes-ayuda`}
                          title="Arrastrá hacia los lados para ajustar la nota"
                          onChange={(event) =>
                            setValores((previous) => ({
                              ...previous,
                              [name]: Number(event.target.value)
                            }))
                          }
                          className="absolute inset-0 m-0 h-full w-full cursor-ew-resize opacity-0"
                        />
                      </div>
                    ))}
                  </div>
                </section>
                <div className="grid min-w-0 gap-2 md:row-span-3 md:grid-rows-subgrid">
                  <label htmlFor={`${id}-comentario`} className="text-sm font-medium">
                    Comentario <span className="font-normal text-muted-foreground">(opcional)</span>
                  </label>
                  <p id={`${id}-comentario-ayuda`} className="text-sm text-muted-foreground">
                    Agregá un comentario para contar cómo fueron las clases y qué destacarías del
                    docente.
                  </p>
                  <textarea
                    id={`${id}-comentario`}
                    name="comentario"
                    aria-describedby={`${id}-comentario-ayuda`}
                    rows={4}
                    value={comentario}
                    onChange={(event) => setComentario(event.target.value)}
                    placeholder="Compartí tu experiencia…"
                    className="block min-h-48 w-full resize-y rounded-md border border-input bg-transparent px-3 py-2 text-base placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/50 focus-visible:outline-none md:text-sm"
                  />
                </div>
              </div>
            </form>
          </div>
          <footer className="flex shrink-0 flex-wrap items-center justify-end gap-3 border-t p-4">
            <p className="mr-auto max-w-sm text-xs text-muted-foreground">
              Vista previa: la publicación todavía no está habilitada. Los datos no se guardan al
              recargar.
            </p>
            <Dialog.Close asChild>
              <Button variant="outline">Cerrar</Button>
            </Dialog.Close>
            <Button disabled>Publicar calificación</Button>
          </footer>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
