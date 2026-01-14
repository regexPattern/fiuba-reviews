<script lang="ts">
  import { ListFilter, Star } from "@lucide/svelte";
  import { ScrollArea } from "bits-ui";
  import Fuse from "fuse.js";
  import { page } from "$app/state";

  const FUZZY_SEARCH_THRESHOLD = 0.15;
  const FUZZY_SEARCH_DEBOUNCE_TIMEOUT_MS = 300;

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

  let queryValue = $state("");
  let queryDebounced = $state("");
  let fuse = $derived(
    new Fuse(catedras, {
      ignoreDiacritics: true,
      shouldSort: false,
      includeScore: false,
      threshold: FUZZY_SEARCH_THRESHOLD,
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
    }, FUZZY_SEARCH_DEBOUNCE_TIMEOUT_MS);

    return () => clearTimeout(handler);
  });
</script>

<aside
  class="flex h-full min-h-0 flex-col divide-y divide-layout-border border-r border-layout-border md:pt-[56px]"
>
  <div class="shrink-0 space-y-2 border-b p-3 font-serif text-lg font-medium">
    {materia.nombre}
  </div>

  <div class="shrink-0 p-3 text-sm tabular-nums">
    {#if materia.cuatrimestre}
      <p>
        Última actualización de cátedras en <span class="tracking-tight">
          {materia.cuatrimestre.numero}C{materia.cuatrimestre.anio}
        </span>. Si ves que la oferta está desactualizada, podés colaborar
        <a href={`/colaborar?materia=${materia.codigo}`} class="text-fiuba underline">
          enviándola acá
        </a>.
      </p>
    {:else}
      La oferta de esta materia no está actualizada todavía. Se muestran las ofertas de sus
      equivalencias en los planes anteriores: <span class="tracking-tight">
        {materia.equivalencias.map((e) => e.codigo).join(", ")}
      </span>. Si tenés la oferta actualizada, podés colaborar
      <a href={`/colaborar?materia=${materia.codigo}`} class="text-fiuba underline">
        enviándola acá
      </a>.
    {/if}
  </div>

  <ScrollArea.Root class="min-h-0 flex-1 overflow-hidden">
    <ScrollArea.Viewport class="h-full">
      <ul class="my-2">
        {#each catedrasFiltradas as catedra (catedra.codigo)}
          {@const calificacion = catedra.calificacion.toFixed(1)}

          <li class="p-3">
            <a href={`${catedra.codigo}`} class="flex items-center gap-1.5 tabular-nums">
              {calificacion === "0.0" ? "–" : calificacion}
              <Star class="size-[12px] shrink-0 fill-yellow-500 stroke-yellow-700" />
              <span class={page.params.codigo_catedra === catedra.codigo ? "text-fiuba" : ""}
                >{catedra.nombre}</span
              >
            </a>
          </li>
        {/each}
      </ul>
    </ScrollArea.Viewport>
    <ScrollArea.Scrollbar orientation="vertical">
      <ScrollArea.Thumb />
    </ScrollArea.Scrollbar>
  </ScrollArea.Root>

  <div class="relative shrink-0">
    <input
      bind:value={queryValue}
      placeholder="Filtrar por nombre de docente"
      class="w-full py-3 pr-12 pl-3 outline-none placeholder:text-sm"
    />
    <span
      class="pointer-events-none absolute top-1/2 right-3 flex size-[26px] -translate-y-1/2 items-center justify-center rounded-full border border-border text-foreground/50"
    >
      <ListFilter class="size-[12px]" aria-hidden="true" />
    </span>
  </div>
</aside>
