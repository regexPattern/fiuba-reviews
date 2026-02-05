<script lang="ts">
  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";
  import { calificarDocente } from "./data.remote";
  import { CircleCheck, CircleAlert, Loader } from "@lucide/svelte";
  import { Button, Label, Slider } from "bits-ui";
  import { mode } from "mode-watcher";
  import { Turnstile } from "svelte-turnstile";

  const LABELS_CAMPOS_CALIFICACION = [
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

  const VALOR_POR_DEFECTO = 2.5;

  let enviando = $state(false);

  for (const { name } of LABELS_CAMPOS_CALIFICACION) {
    const campo = calificarDocente.fields.calificaciones[name];
    if (campo.value() === null) {
      campo.set(VALOR_POR_DEFECTO);
    }
  }

  function resetearCalificaciones() {
    for (const { name } of LABELS_CAMPOS_CALIFICACION) {
      calificarDocente.fields.calificaciones[name].set(VALOR_POR_DEFECTO);
    }
  }
</script>

{#snippet campoCalificacion(
  nombreCampo: (typeof LABELS_CAMPOS_CALIFICACION)[number]["name"],
  label: string
)}
  {@const field = calificarDocente.fields.calificaciones[nombreCampo]}
  {@const value = field.value() ?? VALOR_POR_DEFECTO}
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
  {...calificarDocente.enhance(async ({ form, submit }) => {
    enviando = true;
    try {
      await submit();
      if (calificarDocente.result) {
        form.reset();
        resetearCalificaciones();
      }
    } catch (_) {
    } finally {
      enviando = false;
    }
  })}
  class="mx-auto flex w-full flex-col gap-12 md:w-fit lg:mx-0 lg:flex-row"
>
  <div class="w-full mx-auto max-w-[360px] sm:max-w-none space-y-8 sm:mx-0 md:w-fit">
    {#each LABELS_CAMPOS_CALIFICACION as { name, label } (name)}
      {@render campoCalificacion(name, label)}
    {/each}
  </div>

  <div class="">
    <Label.Root for="comentario" class="block">
      <span class="font-medium">Comentario</span>
    </Label.Root>

    <textarea
      {...calificarDocente.fields.comentario.as("text")}
      id="comentario"
      rows={5}
      class="mt-1 w-full border border-button-border bg-background p-2 dark:bg-background"
    >
    </textarea>

    <div class="mt-6 flex flex-col items-center justify-between gap-4 md:flex-row">
      <div class="h-[65px] w-[300px] overflow-hidden">
        <Turnstile
          siteKey={PUBLIC_TURNSTILE_SITE_KEY}
          responseFieldName="cfTurnstileResponse"
          language="es-es"
          theme={mode.current}
          on:callback={(e) => {
            calificarDocente.fields.cfTurnstileResponse.set(e.detail.token);
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
        {:else if calificarDocente.result?.success}
          Enviado
          <CircleCheck class="size-[16px]" />
        {:else}
          Enviar
        {/if}
      </Button.Root>
    </div>

    {#if calificarDocente.fields.allIssues() && calificarDocente.fields.allIssues()!.length > 0}
      <div class="mt-6 space-y-1 text-sm text-red-500 dark:text-red-400">
        {#each calificarDocente.fields.allIssues() as issue, i (i)}
          <p class="flex flex-wrap items-center gap-1">
            <CircleAlert class="size-[14px]" />{issue.message}
          </p>
        {/each}
      </div>
    {/if}
  </div>
</form>
