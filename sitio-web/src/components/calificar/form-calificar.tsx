"use client";

import { useRouter } from "next/navigation";
import { useRef, useTransition } from "react";
import type { CalificarInput } from "@/lib/calificar-schema";
import type { TurnstileInstance } from "@marsidev/react-turnstile";
import { zodResolver } from "@hookform/resolvers/zod";
import { Turnstile } from "@marsidev/react-turnstile";
import { Loader2 } from "lucide-react";
import { useTheme } from "next-themes";
import { Controller, useForm } from "react-hook-form";
import { toast } from "sonner";
import { calificarDocente } from "@/app/calificar/actions";
import { Button } from "@/components/ui/button";
import { Field, FieldError, FieldGroup, FieldLabel } from "@/components/ui/field";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select";
import { Slider } from "@/components/ui/slider";
import { Textarea } from "@/components/ui/textarea";
import { env } from "@/env";
import { calificarSchema } from "@/lib/calificar-schema";
import { CRITERIOS } from "@/lib/criterios";

const CALIFICACION_POR_DEFECTO = 2.5;

type Cuatri = { codigo: number; numero: number; anio: number };

export function FormCalificar({
  codigoDocente,
  codigoCatedra,
  codigoMateria,
  cuatris
}: {
  codigoDocente: string;
  codigoCatedra: string | null;
  codigoMateria: string | null;
  cuatris: Cuatri[];
}) {
  const router = useRouter();
  const { resolvedTheme } = useTheme();
  const turnstileRef = useRef<TurnstileInstance>(null);
  const [enviando, startTransition] = useTransition();

  const form = useForm<CalificarInput>({
    resolver: zodResolver(calificarSchema),
    defaultValues: {
      codigoDocente,
      codigoCatedra,
      codigoMateria,
      calificaciones: Object.fromEntries(
        CRITERIOS.map(({ name }) => [name, CALIFICACION_POR_DEFECTO])
      ) as CalificarInput["calificaciones"],
      comentario: "",
      cfTurnstileResponse: ""
    }
  });

  function onSubmit(values: CalificarInput) {
    startTransition(async () => {
      const resultado = await calificarDocente(values);

      if (resultado.ok) {
        toast.success("¡Gracias! Tu calificación fue enviada.");
        router.push(resultado.redirectTo as never);
        return;
      }

      toast.error(resultado.error);
      // El token de Turnstile es de un solo uso.
      form.setValue("cfTurnstileResponse", "");
      turnstileRef.current?.reset();
    });
  }

  return (
    <form
      // react-hook-form: falso positivo de react-hooks/refs
      // eslint-disable-next-line react-hooks/refs
      onSubmit={form.handleSubmit(onSubmit)}
      className="mx-auto flex w-full flex-col gap-12 md:w-fit lg:mx-0 lg:flex-row"
    >
      <FieldGroup className="mx-auto w-full max-w-[360px] gap-6 sm:mx-0 sm:max-w-none sm:gap-10 md:w-fit">
        {CRITERIOS.map(({ name, label }) => (
          <Controller
            key={name}
            control={form.control}
            name={`calificaciones.${name}`}
            render={({ field }) => (
              <Field className="sm:flex-row sm:items-center sm:justify-between sm:gap-8">
                <FieldLabel htmlFor={`criterio-${name}`} className="shrink-0 font-medium">
                  {label}
                </FieldLabel>
                <div className="flex w-full items-center gap-3 sm:w-[280px] md:ml-auto">
                  <Slider
                    id={`criterio-${name}`}
                    min={0}
                    max={5}
                    step={0.5}
                    value={[field.value]}
                    onValueChange={([valor]) => field.onChange(valor)}
                  />
                  <span className="w-[3ch] text-right text-sm tabular-nums">
                    {field.value.toFixed(1)}
                  </span>
                </div>
              </Field>
            )}
          />
        ))}
      </FieldGroup>

      <FieldGroup className="gap-6">
        <Controller
          control={form.control}
          name="comentario"
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel htmlFor="comentario" className="font-medium">
                Comentario
              </FieldLabel>
              <Textarea
                {...field}
                id="comentario"
                rows={5}
                aria-invalid={fieldState.invalid}
                className="bg-background"
              />
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Controller
          control={form.control}
          name="cuatrimestre"
          render={({ field, fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <FieldLabel htmlFor="cuatrimestre" className="font-medium">
                Cuatrimestre
              </FieldLabel>
              <Select
                value={field.value ? String(field.value) : ""}
                onValueChange={(v) => field.onChange(Number(v))}
              >
                <SelectTrigger
                  id="cuatrimestre"
                  aria-invalid={fieldState.invalid}
                  className="w-full bg-background"
                >
                  <SelectValue placeholder="Seleccioná un cuatrimestre" />
                </SelectTrigger>
                <SelectContent>
                  {cuatris.map((c) => (
                    <SelectItem key={c.codigo} value={String(c.codigo)}>
                      {c.numero}C{c.anio}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Controller
          control={form.control}
          name="cfTurnstileResponse"
          render={({ fieldState }) => (
            <Field data-invalid={fieldState.invalid}>
              <div className="h-[65px] w-[300px] overflow-hidden">
                <Turnstile
                  ref={turnstileRef}
                  siteKey={env.NEXT_PUBLIC_TURNSTILE_SITE_KEY}
                  options={{
                    language: "es",
                    theme: resolvedTheme === "dark" ? "dark" : "light"
                  }}
                  onSuccess={(token) =>
                    form.setValue("cfTurnstileResponse", token, { shouldValidate: true })
                  }
                  onExpire={() => form.setValue("cfTurnstileResponse", "")}
                />
              </div>
              <FieldError errors={[fieldState.error]} />
            </Field>
          )}
        />

        <Button type="submit" disabled={enviando} className="w-32 self-end rounded-full">
          {enviando ? (
            <>
              Enviando <Loader2 className="size-4 animate-spin" />
            </>
          ) : (
            "Enviar"
          )}
        </Button>
      </FieldGroup>
    </form>
  );
}
