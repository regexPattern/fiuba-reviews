<script lang="ts">
  import { resolve } from "$app/paths";
  import { ScrollArea } from "bits-ui";
  import Catedra from "./components/Catedra.svelte";
  import Sidebar from "./components/Sidebar.svelte";

  let { data } = $props();

  let metaTitle = $derived(`${data.materia.codigo} • ${data.materia.nombre} | FIUBA Reviews`);
  let metaDescription = $derived(
    `Visitá la página de ${data.materia.nombre} para ver calificaciones y comentarios de las cátedras.`
  );

  let idxCatedra = $state(0);
  let tieneCatedras = $derived(data.catedras.length > 0);
  let viewportRef: HTMLDivElement | null = $state(null);

  const resetCatedraView = () => {
    viewportRef?.scrollTo({ top: 0 });
  };
</script>

<svelte:head>
  <title>{metaTitle}</title>
  <meta name="robots" content="index,follow,max-snippet:-1,max-image-preview:large" />
  <meta name="description" content={metaDescription} />
  <link rel="canonical" href={resolve(`/materia/${data.materia.codigo}`)} />

  <meta property="og:title" content={metaTitle} />
  <meta property="og:description" content={metaDescription} />
  <meta property="og:image" content={resolve(`/materia/${data.materia.codigo}/og.png`)} />
  <meta property="og:image:alt" content={`FIUBA Reviews Materia ${data.materia.codigo}`} />
  <meta name="twitter:card" content="summary_large_image" />
</svelte:head>

<div
  class="container mx-auto mt-[calc(-56px-env(safe-area-inset-top))] flex overflow-hidden"
  style="height: -webkit-fill-available; height: 100dvh"
>
  {#if tieneCatedras}
    <div class="hidden w-90 shrink-0 md:flex">
      <Sidebar
        materia={data.materia}
        catedras={data.catedras}
        bind:idxCatedra
        callback={resetCatedraView}
      />
    </div>
  {/if}

  <main class="min-h-0 w-full min-w-0">
    <ScrollArea.Root class="h-full min-h-0 overflow-hidden">
      <ScrollArea.Viewport
        bind:ref={viewportRef}
        class="h-full w-full pt-[calc(56px+env(safe-area-inset-top))]"
        data-scroll-container="main"
      >
        <button
          class="sticky top-0 z-200 flex w-full items-center justify-between border-b border-layout-border bg-background p-3 text-left font-serif text-lg font-medium md:hidden"
        >
          {data.materia.nombre}
        </button>

        {#if tieneCatedras}
          <Catedra catedra={data.catedras[idxCatedra]} />
        {:else}
          <div class="m-4 md:m-6">
            <p class="text-foreground-muted py-2 text-sm">
              Esta materia todavía no tiene cátedras cargadas.
            </p>
          </div>
        {/if}
      </ScrollArea.Viewport>

      <ScrollArea.Scrollbar orientation="vertical">
        <ScrollArea.Thumb />
      </ScrollArea.Scrollbar>
      <ScrollArea.Corner />
    </ScrollArea.Root>
  </main>
</div>
