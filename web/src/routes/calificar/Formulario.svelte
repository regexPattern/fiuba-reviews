<script lang="ts">
  import { CircleAlert, CircleCheck, Loader } from "@lucide/svelte";
  import { Button, Label, Select, Slider } from "bits-ui";
  import { mode } from "mode-watcher";
  import { Turnstile } from "svelte-turnstile";
  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";
  import { submitForm } from "./form.remote";

  const CAMPOS_CALIFICACION = [
    { name: "aceptaCritica", label: "Acepta crítica" },
    { name: "asistencia", label: "Asistencia" },
    { name: "buenTrato", label: "Buen trato" },
    { name: "claridad", label: "Claridad" },
    { name: "claseOrganizada", label: "Clase organizada" },
    { name: "cumpleHorarios", label: "Cumple horarios" },
    { name: "fomentaParticipacion", label: "Fomenta participación" },
    { name: "panoramaAmplio", label: "Panorama amplio" },
    { name: "respondeMails", label: "Responde mails" }
  ] as const;

  const CALIFICACION_POR_DEFECTO = 2.5;

  for (const { name } of CAMPOS_CALIFICACION) {
    const campo = submitForm.fields.calificaciones[name];
    if (campo.value() === null) {
      campo.set(CALIFICACION_POR_DEFECTO);
    }
  }

  type Cuatri = { codigo: number; numero: number; anio: number };

  interface Props {
    cuatris: Cuatri[];
  }

  let { cuatris }: Props = $props();

  let cuatriSeleccionado = $state<Cuatri | null>(null);

  let enviando = $state(false);

  function resetearCalificaciones() {
    for (const { name } of CAMPOS_CALIFICACION) {
      submitForm.fields.calificaciones[name].set(CALIFICACION_POR_DEFECTO);
    }
  }
</script>

{#snippet campoCalificacion(
  nombreCampo: (typeof CAMPOS_CALIFICACION)[number]["name"],
  label: string
)}
  {@const field = submitForm.fields.calificaciones[nombreCampo]}
  {@const value = field.value() ?? CALIFICACION_POR_DEFECTO}
  {@const inputComponent = field.as("number")}

  <div class="flex flex-col gap-6 sm:flex-row sm:items-center sm:justify-between sm:gap-8">
    <Label.Root for={inputComponent.name} class="shrink-0">
      <span class="font-medium">{label}</span>
    </Label.Root>

    <input
      id={inputComponent.name}
      type="hidden"
      name={inputComponent.name}
      {value}
      aria-invalid={inputComponent["aria-invalid"]}
    />

    <Slider.Root
      type="single"
      min={0}
      max={5}
      step={0.5}
      {value}
      onValueChange={(valor) => field.set(valor)}
      trackPadding={3}
      class="relative flex w-full touch-none items-center select-none sm:w-[280px] md:ml-auto"
    >
      {#snippet children({ tickItems })}
        <span
          class="relative h-4 w-full grow cursor-pointer overflow-hidden rounded-full bg-button-background"
        >
          <Slider.Range class="absolute h-full bg-fiuba/50" />
        </span>

        <Slider.Thumb
          index={0}
          class="z-5 size-[22px] rounded-full border border-neutral-400 bg-white transition-transform active:scale-[1.15] active:focus:outline-none"
        />

        {#each tickItems as { index, value } (index)}
          <Slider.Tick {index} class="size-[2px] rounded-full bg-neutral-400" />

          {#if Number.isInteger(value)}
            <Slider.TickLabel {index} position="top" class="pb-1 text-xs text-neutral-600">
              {value}
            </Slider.TickLabel>
          {/if}
        {/each}
      {/snippet}
    </Slider.Root>
  </div>
{/snippet}

<form
  {...submitForm.enhance(async ({ form, submit }) => {
    enviando = true;
    try {
      await submit();
      if (submitForm.result) {
        form.reset();
        resetearCalificaciones();
      }
    } finally {
      enviando = false;
    }
  })}
  class="mx-auto flex w-full flex-col gap-12 md:w-fit lg:mx-0 lg:flex-row"
>
  <div class="mx-auto w-full max-w-[360px] space-y-6 sm:mx-0 sm:max-w-none sm:space-y-10 md:w-fit">
    {#each CAMPOS_CALIFICACION as { name, label } (name)}
      {@render campoCalificacion(name, label)}
    {/each}
  </div>

  <div>
    <Label.Root for="comentario" class="block">
      <span class="font-medium">Comentario</span>
    </Label.Root>

    <textarea
      {...submitForm.fields.comentario.as("text")}
      id="comentario"
      rows={5}
      class="mt-1 w-full border border-button-border bg-background p-2 dark:bg-background"
    >
    </textarea>

    <div class="mt-4">
      <Label.Root for="cuatrimestre" class="block">
        <span class="font-medium">Cuatrimestre</span>
      </Label.Root>

      <input
        type="hidden"
        name={submitForm.fields.cuatrimestre.as("number").name}
        value={cuatriSeleccionado?.codigo ?? ""}
        aria-invalid={submitForm.fields.cuatrimestre.as("number")["aria-invalid"]}
      />

      <Select.Root
        type="single"
        onValueChange={(v) => {
          cuatriSeleccionado = cuatris.find((c) => `${c.codigo}` === v) || null;
          if (cuatriSeleccionado) {
            submitForm.fields.cuatrimestre.set(cuatriSeleccionado.codigo);
          }
        }}
      >
        <Select.Trigger
          id="cuatrimestre"
          class="mt-1 inline-flex w-full items-center rounded-none border border-button-border bg-background px-2 py-2 text-sm text-neutral-900 data-placeholder:text-neutral-500 dark:text-neutral-100 dark:data-placeholder:text-neutral-400"
          aria-label="Seleccionar cuatrimestre"
        >
          <span class="truncate">
            {#if cuatriSeleccionado}
              {cuatriSeleccionado.numero}C{cuatriSeleccionado.anio}
            {:else}
              Seleccióna un cuatrimestre
            {/if}
          </span>
        </Select.Trigger>

        <Select.Portal>
          <Select.Content
            class="z-50 w-(--bits-select-anchor-width) rounded-none border border-button-border bg-background text-neutral-900 shadow-md outline-hidden dark:text-neutral-100"
            sideOffset={6}
          >
            <Select.Viewport class="py-1">
              {#each cuatris as cuatri (cuatri.codigo)}
                {@const label = `${cuatri.numero}C${cuatri.anio}`}

                <Select.Item
                  value={`${cuatri.codigo}`}
                  {label}
                  class="flex w-full items-center px-2 py-2 text-sm text-neutral-900 select-none data-highlighted:bg-neutral-100 data-highlighted:text-neutral-900 dark:text-neutral-100 dark:data-highlighted:bg-neutral-800 dark:data-highlighted:text-neutral-100"
                >
                  {label}
                </Select.Item>
              {/each}
            </Select.Viewport>
          </Select.Content>
        </Select.Portal>
      </Select.Root>
    </div>

    <div class="mt-6 flex flex-col items-center justify-between gap-4 md:flex-row">
      <div class="h-[65px] w-[300px] overflow-hidden">
        <Turnstile
          siteKey={PUBLIC_TURNSTILE_SITE_KEY}
          responseFieldName="cfTurnstileResponse"
          language="es-es"
          theme={mode.current}
          on:callback={(e) => {
            submitForm.fields.cfTurnstileResponse.set(e.detail.token);
          }}
        />
      </div>

      <Button.Root
        type="submit"
        disabled={enviando}
        class="flex w-32 shrink-0 items-center justify-center gap-1 rounded-full border border-green-700 bg-[#65eb95] py-2.5 text-sm font-medium text-green-800 transition-colors hover:bg-green-400 disabled:border-slate-400 disabled:bg-[#C4D8E2] disabled:text-slate-400 dark:disabled:border-slate-600 dark:disabled:bg-slate-900"
      >
        {#if enviando}
          Enviando
          <Loader class="size-[16px] animate-spin" />
        {:else if submitForm.result?.success}
          Enviado
          <CircleCheck class="size-[16px]" />
        {:else}
          Enviar
        {/if}
      </Button.Root>
    </div>

    {#if submitForm.fields.allIssues() && submitForm.fields.allIssues()!.length > 0}
      <div class="mt-6 space-y-1 text-sm text-red-500 dark:text-red-400">
        {#each submitForm.fields.allIssues() as issue, i (i)}
          <p class="flex flex-wrap items-center gap-1">
            <CircleAlert class="size-[14px]" />{issue.message}
          </p>
        {/each}
      </div>
    {/if}
  </div>
</form>
