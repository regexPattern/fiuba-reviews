<script lang="ts">
  import Formulario from "./Formulario.svelte";
  import { ChevronLeft } from "@lucide/svelte";

  let { data } = $props();
</script>

<main class="container mx-auto space-y-12 p-6">
  <div class="flex flex-col-reverse justify-between gap-4 md:flex-row md:items-center">
    <div class="flex flex-col gap-2">
      <h1 class="text-4xl font-medium">{data.docente.nombre}</h1>
      {#if data.docente.nombreSiu && data.docente.rol}
        <p class="text-sm">{data.docente.nombreSiu} • {data.docente.rol}</p>
      {/if}
    </div>

    <a href="/materia/{data.docente.codigoMateria}" class="flex items-center text-sm underline">
      <ChevronLeft class="size-[18px]" />
      {#if data.codigoCatedra}
        Ir a cátedra del docente
      {:else}
        Ir a materia del docente
      {/if}
    </a>
  </div>

  <Formulario cuatris={data.cuatris} />
</main>
