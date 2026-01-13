<script lang="ts">
  import Sidebar from "./Sidebar.svelte";
  import { ScrollArea, Dialog } from "bits-ui";

  let { children, data } = $props();

  let mobileDrawerOpen = $state(false);
</script>

<div class="container mx-auto flex h-[calc(100vh-56px)] overflow-hidden">
  {#if data.catedras.length > 0}
    <!-- Desktop sidebar -->
    <div class="hidden md:flex">
      <Sidebar materia={data.materia} catedras={data.catedras} />
    </div>

    <!-- Mobile drawer -->
    <Dialog.Root bind:open={mobileDrawerOpen}>
      <Dialog.Portal>
        <Dialog.Overlay
          class="fixed inset-0 z-99 bg-black/80 data-[state=closed]:animate-out data-[state=closed]:duration-200 data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:duration-200 data-[state=open]:fade-in md:hidden"
        />
        <Dialog.Content
          class="fixed inset-y-0 left-0 z-100 h-full w-[280px] overflow-hidden data-[state=closed]:animate-out data-[state=closed]:duration-200 data-[state=closed]:slide-out-to-left data-[state=open]:animate-in data-[state=open]:duration-200 data-[state=open]:slide-in-from-left md:hidden"
        >
          <Sidebar materia={data.materia} catedras={data.catedras} />
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  {/if}

  <main class="min-h-0 w-full min-w-0">
    <ScrollArea.Root class="h-full min-h-0 overflow-hidden">
      <ScrollArea.Viewport class="h-full w-full">
        <button
          class="bg-background sticky top-0 z-10 w-full border-b p-4 text-left md:hidden"
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
