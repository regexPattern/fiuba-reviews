<script lang="ts">
  import { ChevronLeft } from "@lucide/svelte";
  import Formulario from "./Formulario.svelte";
  import { resolve } from "$app/paths";

  let { data } = $props();
</script>

<svelte:head>
  <title>FIUBA Reviews • Calificar</title>
</svelte:head>

<div class="relative isolate">
  <div class="fondo-decorativo fondo-superior" aria-hidden="true"></div>

  <main class="relative z-10 container mx-auto space-y-12 p-6">
    <div class="flex flex-col-reverse justify-between gap-4 md:flex-row md:items-center">
      <div class="flex flex-col gap-2">
        <h1 class="text-4xl font-medium">{data.docente.nombre}</h1>
        {#if data.docente.nombreSiu && data.docente.rol}
          <p class="text-sm">{data.docente.nombreSiu} • {data.docente.rol}</p>
        {/if}
      </div>

      <!-- TODO: hacer que funcione el scroll del docente y poner el codigo de la catedra -->
      <a
        href={resolve(
          data.codigoCatedra
            ? `/materia/${data.docente.codigoMateria}/${data.codigoCatedra}`
            : `/materia/${data.docente.codigoMateria}`
        )}
        class="flex items-center text-sm underline"
      >
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
</div>
