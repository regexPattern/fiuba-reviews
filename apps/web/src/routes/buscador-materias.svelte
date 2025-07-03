<script lang="ts">
  import { goto } from "$app/navigation";
  import { Button } from "$lib/components/ui/button";
  import {
    CommandDialog,
    CommandInput,
    CommandItem,
    CommandList,
  } from "$lib/components/ui/command";
  import { cn } from "$lib/utils";
  import Fuse from "fuse.js";
  import { Search } from "lucide-svelte";
  import { onMount } from "svelte";
  import { writable } from "svelte/store";

  const activo = writable(false);

  let className = "";
  export { className as class };

  type Materia = { codigo: string; nombre: string };
  export let materias: Materia[];
  let materiasFiltradas: Materia[] = materias;

  const fuse = new Fuse(materias, {
    keys: ["nombre"],
    ignoreDiacritics: true,
    minMatchCharLength: 3,
    threshold: 0.2,
  });

  let query = "";
  let debounceTimeout: ReturnType<typeof setTimeout> | undefined;

  $: if (!$activo) {
    materiasFiltradas = [];
  }

  async function debounceSearch(e: Event) {
    clearTimeout(debounceTimeout);

    if (e.target instanceof HTMLInputElement) {
      query = e.target.value;
      debounceTimeout = setTimeout(
        () => {
          materiasFiltradas = fuse.search(query).map((r) => r.item);
        },
        query === "" ? 0 : 300,
      );
    }
  }

  onMount(() => {
    function manejarAtajo(e: KeyboardEvent) {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        $activo = !$activo;
      }
    }

    document.addEventListener("keydown", manejarAtajo);
    return () => {
      document.removeEventListener("keydown", manejarAtajo);
    };
  });
</script>

<Button
  class={cn("flex justify-between gap-2 p-2", className)}
  on:click={() => ($activo = !$activo)}
  {...$$restProps}>
  <Search class="h-4 w-4" />
  <span>
    <span>Buscar</span> <span class="hidden sm:inline">Materias</span>
  </span>
  <kbd
    class="hidden rounded border px-1 py-0.5 font-mono text-sm tracking-widest sm:inline">
    Ctrl+K
  </kbd>
</Button>

<CommandDialog bind:open={$activo} shouldFilter={false}>
  <CommandInput placeholder="Nombre de una materia" on:input={debounceSearch} />
  <CommandList>
    {#each materiasFiltradas as mat (mat.codigo)}
      {@const slug = mat.codigo}
      <CommandItem
        value={mat.codigo}
        onSelect={async () => {
          await goto(`/materias/${slug}`);
          $activo = false;
        }}
        class="flex cursor-pointer items-start space-x-1.5">
        <span> {mat.nombre} </span>
      </CommandItem>
    {/each}
  </CommandList>
</CommandDialog>
