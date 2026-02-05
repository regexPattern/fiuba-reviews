<script lang="ts">
  import { PUBLIC_TURNSTILE_SITE_KEY } from "$env/static/public";
  import { extraerMetadataOferta } from "$lib/parser-ofertas";
  import { enviarOferta } from "./data.remote";
  import { Check, CircleAlert, Loader } from "@lucide/svelte";
  import { Button } from "bits-ui";
  import { mode } from "mode-watcher";
  import { Turnstile } from "svelte-turnstile";

  let enviando = $state(false);
  let metadata = $derived(extraerMetadataOferta(enviarOferta.fields.contenido.value()));

  $effect(() => {
    if (metadata) {
      enviarOferta.fields.metadata.carrera.set(metadata.carrera);
      enviarOferta.fields.metadata.cuatrimestre.numero.set(metadata.cuatrimestre.numero);
      enviarOferta.fields.metadata.cuatrimestre.anio.set(metadata.cuatrimestre.anio);
    }
  });
</script>

<form
  {...enviarOferta.enhance(async ({ form, submit }) => {
    enviando = true;
    try {
      await new Promise((r) => setTimeout(r, 3000));
      await submit();
      if (enviarOferta.result) {
        form.reset();
      }
    } catch (_) {
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
        {...enviarOferta.fields.contenido.as("text")}
        rows={5}
        class="mt-1 w-full border border-button-border bg-background p-2 dark:bg-background"
      >
      </textarea>

      <input {...enviarOferta.fields.metadata.carrera.as("text")} hidden />
      <input {...enviarOferta.fields.metadata.cuatrimestre.numero.as("number")} hidden />
      <input {...enviarOferta.fields.metadata.cuatrimestre.anio.as("number")} hidden />
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
          enviarOferta.fields.cfTurnstileResponse.set(e.detail.token);
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
      {:else if enviarOferta.result?.success}
        Enviado
        <Check class="size-[16px]" />
      {:else}
        Enviar
      {/if}
    </Button.Root>
  </div>

  {#if enviarOferta.fields.allIssues() && enviarOferta.fields.allIssues()!.length > 0}
    <div class="space-y-1 text-sm text-red-500 dark:text-red-400">
      {#each enviarOferta.fields.allIssues() as issue, i (i)}
        <p class="flex flex-wrap items-center gap-1">
          <CircleAlert class="size-[14px]" />{issue.message}
        </p>
      {/each}
    </div>
  {/if}
</form>
