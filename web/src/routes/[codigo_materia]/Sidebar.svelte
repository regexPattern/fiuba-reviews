<script lang="ts">
  import { ScrollArea } from "bits-ui";
  import Fuse from "fuse.js";

  interface Props {
    materia: {
      codigo: string;
      nombre: string;
      cuatrimestre: {
        numero: number;
        anio: number;
      } | null;
      equivalencias: {
        codigo: string;
        nombre: string;
      }[];
    };
    catedras: {
      codigo: string;
      codigoMateria: string;
      nombre: string;
      calificacion: number;
    }[];
  }

  let { materia, catedras }: Props = $props();

  const DEBOUNCE_TIMEOUT_MS = 300;

  let queryValue = $state("");
  let queryDebounced = $state("");
  let fuse = $derived(
    new Fuse(catedras, {
      ignoreDiacritics: true,
      shouldSort: false,
      includeScore: false,
      threshold: 0.5,
      keys: ["nombre"]
    })
  );
  let catedrasFiltradas = $derived.by(() => {
    if (queryDebounced.trim() === "") {
      return catedras;
    }
    return fuse.search(queryDebounced).map((result) => result.item);
  });

  $effect(() => {
    if (queryValue.trim() === "") {
      queryDebounced = "";
      return;
    }

    const handler = setTimeout(() => {
      queryDebounced = queryValue;
    }, DEBOUNCE_TIMEOUT_MS);

    return () => clearTimeout(handler);
  });
</script>

<aside class="relative flex h-full min-h-0 flex-col border-r">
  <div class="sticky top-0">
    {materia.codigo} - {materia.nombre}
    {#if materia.cuatrimestre}
      {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
    {:else}
      <ul>
        {#each materia.equivalencias as equivalencia (equivalencia.codigo)}
          <li>{equivalencia.codigo}</li>
        {/each}
      </ul>
    {/if}
  </div>

  <ScrollArea.Root class="min-h-0 flex-1 overflow-hidden">
    <ScrollArea.Viewport class="h-full">
      <ul>
        {#each catedrasFiltradas as catedra (catedra.codigo)}
          <li class="p-4">
            <a href={`${catedra.codigo}`}>
              {catedra.calificacion.toFixed(1)} - {catedra.nombre}
            </a>
          </li>
        {/each}
      </ul>
    </ScrollArea.Viewport>
    <ScrollArea.Scrollbar orientation="vertical">
      <ScrollArea.Thumb />
    </ScrollArea.Scrollbar>
    <ScrollArea.Corner />
  </ScrollArea.Root>

  <div class="sticky bottom-0">
    <input bind:value={queryValue} />
  </div>
</aside>
