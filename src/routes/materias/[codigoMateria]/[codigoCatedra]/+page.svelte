<script lang="ts">
  import type { PageData } from "./$types";

  import { ChevronDown, PlusCircle, Star } from "lucide-svelte";

  import * as Popover from "$lib/components/ui/popover";
  import { Button } from "$lib/components/ui/button";
  import TablaPromediosDocente from "./tabla-promedios-docente.svelte";
  import { Skeleton } from "$lib/components/ui/skeleton";

  export let data: PageData;
</script>

{#await data.docentes}
  <div class="space-y-4">
    <Skeleton class="h-10 w-72" />
    <div
      class="flex flex-col gap-2 xs:flex-row [&>*]:h-10 [&>*]:w-full [&>*]:xs:w-44">
      <Skeleton />
      <Skeleton />
    </div>
    <div class="space-y-2">
      {#each Array(10) as _}
        <Skeleton class="h-10" />
      {/each}
    </div>
  </div>
{:then cat}
  {#each cat.docentes as doc (doc.codigo)}
    <section id={doc.codigo} class="space-y-4">
      <h2 class="text-4xl font-bold tracking-tight">{doc.nombre}</h2>

      {#if doc.resumen_comentarios}
        <div
          class="divide-y divide-border rounded-lg border bg-secondary dark:divide-slate-700 dark:border-slate-700 [&>*]:p-3">
          <p
            class={`text-secondary-foreground before:content-['"'] after:content-['"']`}>
            {doc.resumen_comentarios}
          </p>
          <div class="text-sm text-slate-500">Resumen generado por IA.</div>
        </div>
      {/if}

      <div class="flex flex-col gap-2 xs:flex-row xs:items-center">
        {#if doc.calificaciones}
          <Popover.Root>
            <Popover.Trigger asChild let:builder>
              <Button
                builders={[builder]}
                variant="outline"
                class="items-center gap-1.5">
                <Star class="h-4 w-4 fill-current text-yellow-500" />
                <span
                  >Promedio: {doc.calificaciones.promedio_general.toFixed(
                    1,
                  )}</span>
                <ChevronDown class="h-[1.2rem] w-[1.2rem]" />
              </Button>
            </Popover.Trigger>
            <Popover.Content class="w-max">
              <TablaPromediosDocente
                cantidadCalificaciones={doc.cantidad_calificaciones}
                {...doc.calificaciones} />
            </Popover.Content>
          </Popover.Root>
        {:else}
          <Button variant="outline" class="items-center gap-1.5">
            <Star class="h-4 w-4 fill-none text-yellow-500" />
            <span>Sin calificaciones</span>
          </Button>
        {/if}

        <Button class="items-center gap-1.5" href={`/docentes/${doc.codigo}`}>
          Calificar <PlusCircle class="h-[1.2rem] w-[1.2rem]" />
        </Button>
      </div>

      <div class="flex flex-col gap-2 divide-y">
        {#each doc.comentarios as com (com.codigo)}
          <div class="pt-2 [&:first-child]:pt-0">
            <p class={`inline before:content-['"'] after:content-['"']`}>
              {com.contenido}
            </p>
            &dash;
            <span class="text-sm text-muted-foreground"
              >{com.cuatrimestre}</span>
          </div>
        {/each}
      </div>
    </section>
  {/each}
{/await}
