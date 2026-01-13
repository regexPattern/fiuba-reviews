<script lang="ts">
  import Promedios from "./Promedios.svelte";
  import { Info, MessageSquarePlus } from "@lucide/svelte";
  import Comentarios from "./Comentarios.svelte";

  let { data } = $props();
</script>

<div class="mx-4 space-y-6 py-4 md:mx-6">
  {#each data.docentes as docente (docente.codigo)}
    <section id={docente.codigo} class="scroll-mt-[68px] space-y-3">
      <div>
        <a href={`#${docente.codigo}`}>
          <h1 class="font-serif text-4xl font-semibold tracking-tight">{docente.nombre}</h1>
        </a>
        {#if docente.rol}
          <small class="text-sm">({docente.rol})</small>
        {/if}
      </div>

      {#if docente.resumenComentario}
        <div class="divide-y divide-fiuba border border-fiuba bg-[#E6ADEC]/50">
          <p class={`p-3 before:content-['"'] after:content-['"']`}>
            {docente.resumenComentario}
          </p>
          <div class="flex items-center gap-1 p-2 text-[#495883] select-none">
            <Info class="size-[16px]" />
            <span class="text-sm">Resumen generado con IA.</span>
          </div>
        </div>
      {/if}

      <div class="flex gap-2 text-sm">
        <Promedios promedio={docente.promedioCalificaciones} />

        <a
          href={`/calificar?docente=${docente.codigo}`}
          class="flex items-center gap-2 border border-[#AB9E9C] bg-[#AB9E9C]/50 px-3 py-2"
        >
          <span>Calificar</span>
          <MessageSquarePlus class="size-[16px] fill-fiuba/50 stroke-[#665889]" />
        </a>
      </div>

      <Comentarios comentarios={docente.comentarios} />
    </section>
  {/each}
</div>
