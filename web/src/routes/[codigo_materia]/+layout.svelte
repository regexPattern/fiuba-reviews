<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { ScrollArea, Dialog } from "bits-ui";
  import { afterNavigate } from "$app/navigation";

  let { children, data } = $props();

  let mobileDrawerOpen = $state(false);

  function onClickDrawer(event: MouseEvent) {
    const target = event.target as Element;
    if (target && target.closest("a[href]")) {
      mobileDrawerOpen = false;
    }
  }

  afterNavigate(() => {
    mobileDrawerOpen = false;
  });
</script>

<div class="container mx-auto flex h-screen overflow-hidden">
  {#if data.catedras.length > 0}
    <!-- Desktop sidebar -->
    <div class="hidden w-[280px] shrink-0 md:flex">
      <Sidebar materia={data.materia} catedras={data.catedras} />
    </div>

    <!-- Mobile drawer -->
    <Dialog.Root bind:open={mobileDrawerOpen}>
      <Dialog.Portal>
        <Dialog.Overlay
          class="fixed inset-0 z-300 bg-background/25 backdrop-blur-lg data-[state=closed]:animate-out data-[state=closed]:duration-200 data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:duration-200 data-[state=open]:fade-in md:hidden"
        />
        <Dialog.Content
          onclick={onClickDrawer}
          class="fixed inset-y-0 left-0 z-301 h-full w-4/5 overflow-hidden bg-background data-[state=closed]:animate-out data-[state=closed]:duration-200 data-[state=closed]:slide-out-to-left data-[state=open]:animate-in data-[state=open]:duration-200 data-[state=open]:slide-in-from-left md:hidden"
        >
          <Sidebar materia={data.materia} catedras={data.catedras} />
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  {/if}

  <main class="min-h-0 w-full min-w-0">
    <ScrollArea.Root class="h-full min-h-0 overflow-hidden">
      <ScrollArea.Viewport class="h-full w-full pt-[56px]" data-scroll-container="main">
        <button
          class="sticky top-0 z-200 w-full border-b border-border-muted bg-background p-3 text-left font-serif text-lg font-medium md:hidden"
          onclick={() => (mobileDrawerOpen = true)}
        >
          {data.materia.nombre}
        </button>

        {@render children()}
      </ScrollArea.Viewport>

      <ScrollArea.Scrollbar orientation="vertical">
        <ScrollArea.Thumb />
      </ScrollArea.Scrollbar>
      <ScrollArea.Corner />
    </ScrollArea.Root>
  </main>
</div>
