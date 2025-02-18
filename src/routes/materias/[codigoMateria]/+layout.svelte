<script lang="ts">
  import { page } from "$app/stores";
  import Link from "$lib/components/link.svelte";
  import { Alert } from "$lib/components/ui/alert";
  import { Sheet, SheetTrigger } from "$lib/components/ui/sheet";
  import SheetContent from "$lib/components/ui/sheet/sheet-content.svelte";
  import { cn } from "$lib/utils";
  import { ChevronDown, Star } from "lucide-svelte";

  import type { LayoutData } from "./$types";

  export let data: LayoutData;

  let open = false;
</script>

<svelte:head>
  <title>{data.materia.nombre} | FIUBA Reviews</title>
</svelte:head>

<div class="relative md:container md:mx-auto md:flex-row">
  <div class="sticky top-16 z-30 md:top-auto">
    <aside
      class="fixed hidden h-[calc(100vh-4rem)] w-80 shrink-0 overflow-y-auto border-r bg-background md:block">
      <div
        class="sticky top-0 flex w-full items-start gap-1.5 border-b bg-background p-3 text-center font-medium">
        {data.materia.nombre}
      </div>

      <ul class="space-y-1.5 py-2">
        {#each data.catedras as cat (cat.codigo)}
          <li class="flex items-center gap-1.5 px-5 py-2 md:pl-2 md:pr-4">
            <span
              class={`w-[2.5ch] shrink-0 font-medium ${
                !cat.promedio ? "text-center" : ""
              }`}>
              {cat.promedio?.toFixed(1) || "-"}
            </span>
            <Star
              class="h-3 w-3 shrink-0 fill-current pr-0.5 text-yellow-500" />
            <a
              href={`/materias/${$page.params.codigoMateria}/${cat.codigo}`}
              class={cn(
                $page.params.codigoCatedra === cat.codigo && "text-fiuba",
              )}>
              {cat.nombre}
            </a>
          </li>
        {/each}
      </ul>
    </aside>

    <Sheet bind:open>
      <SheetTrigger asChild>
        <button
          class="z-20 flex w-full items-center justify-between gap-3 border-b bg-background p-3 text-left font-medium md:hidden"
          on:click={() => {
            open = !open;
            window.scroll({ top: 0, behavior: "instant" });
          }}>
          <span class="flex items-start gap-1">
            {data.materia.nombre}
          </span>
          <ChevronDown class="shrink-0" />
        </button>
      </SheetTrigger>
      <SheetContent class="z-[120] p-0 pt-8" side="left">
        <ul class="h-full space-y-1.5 overflow-y-scroll py-2">
          {#each data.catedras as cat (cat.codigo)}
            <li class="flex items-center gap-1.5 px-5 py-2 md:pl-2 md:pr-4">
              <span
                class={`w-[3ch] shrink-0 font-medium ${
                  !cat.promedio ? "text-center" : ""
                }`}>{cat.promedio?.toFixed(1) || "-"}</span>
              <Star class="h-3 w-3 shrink-0 fill-current text-yellow-500" />
              <a
                href={`/materias/${$page.params.codigoMateria}/${cat.codigo}`}
                class={cn(
                  $page.params.codigoCatedra === cat.codigo && "text-fiuba",
                )}
                on:click={() => (open = !open)}>
                {cat.nombre}
              </a>
            </li>
          {/each}
        </ul>
      </SheetContent>
    </Sheet>
  </div>

  <main class="space-y-6 p-4 md:ml-80 md:min-h-[calc(100vh-4rem)] md:p-6">
    <Alert class="border-cyan-300 bg-fiuba text-background dark:border-cyan-600"
      >⚠️ Debido a la migración de los datos hacia los nuevos planes por ahora
      se están mostrando las cátedras de las materias equivalentes a {data
        .materia.nombre} en los planes anteriores: {#each data.equivalencias as eq, i (eq.codigo)}
        <span class="font-mono font-medium"
          >{eq.codigo}{#if i < data.equivalencias.length - 1},
          {/if}</span>
      {/each}. <Link href="/planes" class="font-medium underline" external
        >Más Información</Link
      >.
    </Alert>
    <div class="flex flex-col gap-12">
      <slot />
    </div>
  </main>
</div>
