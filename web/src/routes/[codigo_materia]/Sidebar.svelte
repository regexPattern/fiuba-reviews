<script lang="ts">
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

<aside class="h-full border-r">
  <!-- <input bind:value={queryValue} /> -->

  <div>
    {materia.codigo} - {materia.nombre}
    {#if materia.cuatrimestre}
      {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
    {:else}
      {#each materia.equivalencias as equivalencia (equivalencia.codigo)}
        {equivalencia.codigo}
      {/each}
    {/if}
  </div>

  <ul class="overflow-y-auto">
    {#each catedrasFiltradas as catedra (catedra.codigo)}
      <li class="p-4">
        <a href={`${catedra.codigo}`}>
          {catedra.calificacion.toFixed(1)} - {catedra.nombre}
        </a>
      </li>
    {/each}
  </ul>
</aside>
