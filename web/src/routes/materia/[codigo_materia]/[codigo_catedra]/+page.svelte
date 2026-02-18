<script lang="ts">
  import { Info, MessageSquarePlus } from "@lucide/svelte";
  import Comentarios from "./Comentarios.svelte";
  import Promedios from "./Promedios.svelte";

  let { data } = $props();

  let nombreCatedra = $derived(
    data.catedras.find((catedra) => catedra.codigo === data.codigoCatedra)?.nombre ?? "Cátedra"
  );
</script>

<svelte:head>
  <title>FIUBA Reviews • {data.materia.codigo} • {nombreCatedra}</title>
</svelte:head>

<div class="m-4 space-y-8 md:m-6">
  {#each data.docentes as docente (docente.codigo)}
    <section id={docente.codigo} class="scroll-mt-[68px] space-y-3">
      <div>
        <h1 class="w-fit text-4xl font-semibold tracking-tight">
          {docente.nombre}
        </h1>
        {#if docente.rol}
          <small class="text-sm">({docente.rol})</small>
        {/if}
      </div>

      {#if docente.resumenComentario}
        <div class="divide-y divide-fiuba border border-fiuba bg-fiuba/45">
          <p class={`p-3 before:content-['"'] after:content-['"']`}>
            {docente.resumenComentario}
          </p>
          <div class="flex items-center gap-1 p-3 text-button-foreground select-none">
            <Info class="size-[16px]" />
            <span class="text-sm">Resumen generado con IA.</span>
          </div>
        </div>
      {/if}

      <div class="flex gap-2 text-sm text-button-foreground">
        <Promedios
          promedio={docente.promedioCalificaciones}
          cantidadCalificaciones={docente.cantidadCalificaciones}
        />

        <a
          href={`/calificar?docente=${docente.codigo}&catedra=${data.codigoCatedra}`}
          class="flex items-center gap-2 border border-button-border bg-button-background px-3 py-2 hover:bg-button-hover hover:transition-colors"
        >
          <span>Calificar</span>
          <MessageSquarePlus
            class="size-[16px] fill-fiuba/50 stroke-[#665889] dark:fill-[#D1BCE3]"
          />
        </a>
      </div>

      {#if docente.comentarios.length > 0}
        <Comentarios comentarios={docente.comentarios} />
      {:else}
        <p class="text-foreground-muted py-2 text-sm">Docente sin comentarios</p>
      {/if}
    </section>
  {/each}
</div>
