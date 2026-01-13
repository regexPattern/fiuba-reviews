<script lang="ts">
  import { SquareStop, Star } from "@lucide/svelte";
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
    return fuse.search(queryDebounced).map((r) => r.item);
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

<aside
  class="flex h-full min-h-0 flex-col divide-y divide-border-muted border-r border-border-muted md:pt-[56px]"
>
  <div class="shrink-0 space-y-2 border-b p-3 font-serif text-lg font-medium">
    {materia.nombre}
  </div>

  <div class="shrink-0 p-3 text-sm tabular-nums">
    {#if materia.cuatrimestre}
      <p>
        Última actualización de cátedras en <span class="tracking-tight">
          {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
        </span>. Si ves que la oferta está desactualizada, podés colaborar enviandola
        <a href={`/ofertas?materia=${materia.codigo}`} class="text-fiuba underline">acá</a>.
      </p>
    {:else}
      La oferta de esta materia no está actualizada todavía. Se muestran las ofertas de sus
      equivalencias en los planes anteriores: <span class="tracking-tight">
        {materia.equivalencias.map((e) => e.codigo).join(", ")}
      </span>. Si tenés la oferta actualizada, podés colaborar enviandola
      <a href={`/ofertas?materia=${materia.codigo}`} class="text-fiuba underline">acá</a>.
    {/if}
  </div>

  <ScrollArea.Root class="min-h-0 flex-1 overflow-hidden">
    <ScrollArea.Viewport class="h-full">
      <ul class="my-2">
        {#each catedrasFiltradas as catedra (catedra.codigo)}
          <li class="p-3">
            <a
              href={`${catedra.codigo}`}
              class="flex items-center gap-1.5 font-serif font-medium tabular-nums"
            >
              {catedra.calificacion.toFixed(1)}
              <Star class="size-[12px] shrink-0 fill-yellow-500 stroke-yellow-700" />
              {catedra.nombre}
            </a>
          </li>
        {/each}
      </ul>
    </ScrollArea.Viewport>
    <ScrollArea.Scrollbar orientation="vertical">
      <ScrollArea.Thumb />
    </ScrollArea.Scrollbar>
  </ScrollArea.Root>

  <div class="shrink-0">
    <input
      bind:value={queryValue}
      placeholder="Filtrar por nombre de docente"
      class="w-full p-3 placeholder:text-sm"
    />
  </div>
</aside>
