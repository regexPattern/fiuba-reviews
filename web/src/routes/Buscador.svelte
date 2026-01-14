<script lang="ts">
  import { Command, Dialog } from "bits-ui";
  import Fuse from "fuse.js";
  import { Search } from "@lucide/svelte";

  interface Props {
    materias: {
      codigo: string;
      nombre: string;
    }[];
  }

  let { materias }: Props = $props();

  const FUZZY_SEARCH_THRESHOLD = 0.75;
  const FUZZY_SEARCH_DEBOUNCE_TIMEOUT_MS = 300;

  let dialogOpen = $state(false);
  let queryValue = $state("");
  let queryDebounced = $state("");

  let fuse = $derived(
    new Fuse(materias, {
      ignoreDiacritics: true,
      ignoreFieldNorm: true,
      shouldSort: true,
      includeScore: true,
      threshold: FUZZY_SEARCH_THRESHOLD,
      keys: ["codigo", "nombre"]
    })
  );

  let materiasFiltradas = $derived.by(() => {
    if (queryDebounced.trim() === "") {
      return materias;
    }
    return fuse
      .search(queryDebounced)
      .sort((a, b) => (a.score ?? 0) - (b.score ?? 0) || a.refIndex - b.refIndex)
      .map((r) => r.item);
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

  function handleKeydown(e: KeyboardEvent) {
    const target = e.target as HTMLElement | null;
    const tag = target?.tagName?.toLowerCase();
    const estaEscribiendo =
      tag === "input" || tag === "textarea" || tag === "select" || target?.isContentEditable;

    if (estaEscribiendo) {
      return;
    }

    if (e.key.toLowerCase() === "k" && (e.metaKey || e.ctrlKey)) {
      e.preventDefault();
      dialogOpen = true;
    }
  }
</script>

<svelte:document onkeydown={handleKeydown} />

<Dialog.Root bind:open={dialogOpen}>
  <Dialog.Trigger
    class="flex items-center gap-2 rounded-full border border-button-border bg-button-background px-3 py-2 text-sm text-button-foreground transition-colors hover:bg-button-hover"
  >
    <span class="hidden items-center gap-2 md:flex">
      <Search class="size-4" />
      <span>Buscar materias</span>
    </span>
    <span class="flex items-center gap-2 md:hidden">
      <span>Buscar</span>
      <Search class="size-4" />
    </span>
    <span class="hidden md:inline">⌘K</span>
  </Dialog.Trigger>

  <Dialog.Portal>
    <Dialog.Overlay
      class="fixed inset-0 z-500 bg-overlay-background backdrop-filter-(--backdrop-filter-overlay-blur) data-[state=closed]:animate-out data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:fade-in"
    />
    <Dialog.Content
      class="fixed top-1/2 left-1/2 z-501 w-full max-w-[min(560px,94vw)] -translate-x-1/2 -translate-y-1/2 border border-secondary-border bg-background shadow-xl data-[state=closed]:animate-out data-[state=closed]:zoom-out-95 data-[state=closed]:fade-out data-[state=open]:animate-in data-[state=open]:zoom-in-95 data-[state=open]:fade-in"
    >
      <Command.Root
        shouldFilter={false}
        class="flex h-full w-full flex-col divide-y divide-secondary-border self-start overflow-hidden bg-background"
      >
        <div class="relative">
          <Search
            aria-hidden="true"
            class="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-foreground/60"
          />

          <Command.Input
            bind:value={queryValue}
            autofocus
            placeholder="Código o nombre de materia"
            class="w-full truncate bg-transparent p-3 pr-3 pl-10 text-sm focus:ring-0 focus:outline-hidden md:pr-14"
          />

          <kbd
            class="pointer-events-none absolute top-1/2 right-3 hidden -translate-y-1/2 rounded-none border border-secondary-border px-1.5 py-0.5 font-mono text-[10px] leading-none text-foreground/60 md:inline-flex"
          >
            esc
          </kbd>
        </div>
        <Command.List class="max-h-[280px] overflow-x-hidden overflow-y-auto pb-2">
          <Command.Viewport>
            {#if materiasFiltradas.length === 0}
              <Command.Empty
                class="text-muted-foreground flex w-full items-center justify-center pt-8 pb-6 text-sm"
              >
                Sin resultados
              </Command.Empty>
            {:else}
              <Command.Group>
                <Command.GroupItems>
                  {#each materiasFiltradas as materia (materia.codigo)}
                    <Command.LinkItem
                      value={materia.codigo}
                      href={`/${materia.codigo}`}
                      onSelect={() => (dialogOpen = false)}
                      class="data-selected:bg-muted flex h-10 cursor-pointer items-center gap-3 p-3 text-sm outline-hidden select-none"
                    >
                      <span class="font-mono text-xs text-foreground/60 tabular-nums">
                        {materia.codigo}
                      </span>
                      <span>{materia.nombre}</span>
                    </Command.LinkItem>
                  {/each}
                </Command.GroupItems>
              </Command.Group>
            {/if}
          </Command.Viewport>
        </Command.List>
      </Command.Root>
    </Dialog.Content>
  </Dialog.Portal>
</Dialog.Root>
