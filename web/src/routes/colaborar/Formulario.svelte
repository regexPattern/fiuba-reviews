<script lang="ts">
  import { Check, CircleAlert, Loader } from "@lucide/svelte";
  import { Button } from "bits-ui";
  import { mode } from "mode-watcher";
  import { Turnstile } from "svelte-turnstile";
  import { submitForm } from "./form.remote";
  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";
  import { extraerMetadataOferta } from "$lib/ofertas";

  let enviando = $state(false);
  let metadata = $derived(extraerMetadataOferta(submitForm.fields.contenido.value()));

  $effect(() => {
    if (metadata) {
      submitForm.fields.metadata.carrera.set(metadata.carrera);
      submitForm.fields.metadata.cuatrimestre.numero.set(metadata.cuatrimestre.numero);
      submitForm.fields.metadata.cuatrimestre.anio.set(metadata.cuatrimestre.anio);
    }
  });
</script>

<form
  {...submitForm.enhance(async ({ form, submit }) => {
    enviando = true;
    try {
      await submit();
      if (submitForm.result) {
        form.reset();
      }
    } finally {
      enviando = false;
    }
  })}
  class="w-full space-y-6"
>
  <div class="space-y-1">
    <label class="block">
      <span class="font-medium">Contenido copiado del SIU</span>
      <textarea
        {...submitForm.fields.contenido.as("text")}
        rows={5}
        class="mt-1 w-full border border-button-border bg-background p-2 dark:bg-background"
      >
      </textarea>

      <input {...submitForm.fields.metadata.carrera.as("text")} hidden />
      <input {...submitForm.fields.metadata.cuatrimestre.numero.as("number")} hidden />
      <input {...submitForm.fields.metadata.cuatrimestre.anio.as("number")} hidden />
    </label>

    <p class="text-sm">
      {#if metadata}
        Detectada oferta de <span class="underline">{metadata.carrera}</span> para
        <span class="underline">{metadata.cuatrimestre.numero}C{metadata.cuatrimestre.anio}</span>.
      {:else}
        Pegá el contenido copiado del SIU.
      {/if}
    </p>
  </div>

  <div class="flex flex-col items-center justify-between gap-4 md:flex-row">
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
      disabled={!metadata || enviando}
      class="flex w-32 shrink-0 items-center justify-center gap-1 rounded-full border border-green-700 bg-[#65eb95] py-2.5 text-sm font-medium text-green-800 transition-colors hover:bg-green-400 disabled:border-slate-400 disabled:bg-[#C4D8E2] disabled:text-slate-400 dark:disabled:border-slate-600 dark:disabled:bg-slate-900"
    >
      {#if enviando}
        Enviando
        <Loader class="size-[16px] animate-spin" />
      {:else if submitForm.result?.success}
        Enviado
        <Check class="size-[16px]" />
      {:else}
        Enviar
      {/if}
    </Button.Root>
  </div>

  {#if submitForm.fields.allIssues() && submitForm.fields.allIssues()!.length > 0}
    <div class="space-y-1 text-sm text-red-500 dark:text-red-400">
      {#each submitForm.fields.allIssues() as issue, i (i)}
        <p class="flex flex-wrap items-center gap-1">
          <CircleAlert class="size-[14px]" />{issue.message}
        </p>
      {/each}
    </div>
  {/if}
</form>
