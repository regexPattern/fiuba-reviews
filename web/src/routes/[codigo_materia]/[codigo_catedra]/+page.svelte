<script lang="ts">
  import Promedios from "./Promedios.svelte";
  import { Info, MessageSquarePlus } from "@lucide/svelte";
  import Comentarios from "./Comentarios.svelte";

  let { data } = $props();
</script>

<div class="m-4 space-y-8 md:m-6">
  {#each data.docentes as docente (docente.codigo)}
    <section id={docente.codigo} class="scroll-mt-[68px] space-y-3">
      <div>
        <h1 class="w-fit font-serif text-4xl font-semibold tracking-tight">
          {docente.nombre}
        </h1>
        {#if docente.rol}
          <small class="text-sm">({docente.rol})</small>
        {/if}
      </div>

      {#if docente.resumenComentario}
        <div class="divide-y divide-fiuba border border-fiuba bg-[#E6ADEC]/50">
          <p class={`p-3 before:content-['"'] after:content-['"']`}>
            {docente.resumenComentario}
          </p>
          <div class="flex items-center gap-1 p-3 text-[#495883] select-none">
            <Info class="size-[16px]" />
            <span class="text-sm">Resumen generado con IA.</span>
          </div>
        </div>
      {/if}

      <div class="flex gap-2 text-sm">
        <Promedios
          promedio={docente.promedioCalificaciones}
          cantidadCalificaciones={docente.cantidadCalificaciones}
        />

        <a
          href={`/calificar?docente=${docente.codigo}`}
          class="flex items-center gap-2 border border-foreground-muted bg-foreground-muted/50 px-3 py-2"
        >
          <span>Calificar</span>
          <MessageSquarePlus class="size-[16px] fill-fiuba/50 stroke-[#665889]" />
        </a>
      </div>

      {#if docente.comentarios.length > 0}
        <Comentarios comentarios={docente.comentarios} />
      {:else}
        <p class="py-2 text-sm text-foreground-muted">Docente sin comentarios</p>
      {/if}
    </section>
  {/each}
</div>
