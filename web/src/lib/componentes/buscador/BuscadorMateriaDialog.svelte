<script lang="ts">
  import { Search } from "@lucide/svelte";
  import { Command, Dialog } from "bits-ui";
  import { buscadorMaterias } from "./materias.svelte";
  import { beforeNavigate } from "$app/navigation";

  beforeNavigate(() => {
    buscadorMaterias.cerrarBuscador();
    buscadorMaterias.limpiarQuery();
  });

  function handleKeyDownGlobal(evento: KeyboardEvent) {
    const target = evento.target as HTMLElement | null;
    const tag = target?.tagName?.toLowerCase();

    if (tag === "input" || tag === "textarea" || tag === "select" || target?.isContentEditable) {
      return;
    }

    if (evento.key.toLowerCase() === "k" && (evento.metaKey || evento.ctrlKey)) {
      evento.preventDefault();
      buscadorMaterias.abrirBuscador();
    }
  }
</script>

<svelte:document onkeydown={handleKeyDownGlobal} />

<Dialog.Root bind:open={buscadorMaterias.open}>
  <Dialog.Portal>
    <Dialog.Overlay
      class="fixed inset-0 z-500 bg-overlay-background backdrop-filter-(--backdrop-filter-overlay-blur) data-[state=closed]:animate-out data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:fade-in"
    />

    <Dialog.Content
      class="fixed top-[calc(56px+env(safe-area-inset-top)+16px)] left-1/2 z-501 max-h-[calc(100dvh-(56px+env(safe-area-inset-top)+32px))] w-full max-w-[min(560px,94vw)] -translate-x-1/2 border border-secondary-border bg-background shadow-xl data-[state=closed]:animate-out data-[state=closed]:zoom-out-95 data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:zoom-in-95 data-[state=open]:fade-in"
    >
      <Command.Root
        shouldFilter={false}
        vimBindings
        class="flex h-full w-full flex-col divide-y divide-secondary-border self-start overflow-hidden bg-background"
      >
        <div class="relative">
          <Search
            aria-hidden="true"
            class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-foreground/60"
          />

          <Command.Input
            value={buscadorMaterias.queryValue}
            oninput={(evento) => buscadorMaterias.setQueryValue(evento.currentTarget.value)}
            placeholder="Codigo o nombre de materia"
            class="w-full truncate p-3 pr-3 pl-10 text-base focus:ring-0 focus:outline-hidden md:pr-14"
          />

          <kbd
            class="pointer-events-none absolute top-1/2 right-3 hidden -translate-y-1/2 rounded-none border border-secondary-border px-1.5 py-0.5 font-mono text-[10px] leading-none text-foreground/60 md:inline-flex"
          >
            esc
          </kbd>
        </div>

        <Command.List class="max-h-[280px] overflow-x-hidden overflow-y-auto pb-2">
          <Command.Viewport>
            {#each buscadorMaterias.materiasFiltradas as materia (materia.codigo)}
              <Command.LinkItem
                href="/materia/{materia.codigo}"
                class="flex h-10 cursor-pointer items-center gap-3 p-3 text-sm outline-hidden select-none data-selected:text-fiuba"
              >
                <span class="text-xs text-foreground tabular-nums">
                  {materia.codigo}
                </span>
                <span>{materia.nombre}</span>
              </Command.LinkItem>
            {/each}

            <Command.Empty
              class="text-muted-foreground flex w-full items-center justify-center pt-8 pb-6 text-sm"
            >
              Sin resultados.
            </Command.Empty>
          </Command.Viewport>
        </Command.List>
      </Command.Root>
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
