<script lang="ts">
  import { goto } from "$app/navigation";
  import CommandStore from "$lib/command";
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

  export let label: string;

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

  $: if (!$CommandStore) {
    materiasFiltradas = [];
  }

  async function debounceSearch(e: Event) {
    clearTimeout(debounceTimeout);

    debounceTimeout = setTimeout(() => {
      if (e.target instanceof HTMLInputElement) {
        query = e.target.value;

        materiasFiltradas = fuse.search(query).map((r) => r.item);
      }
    }, 300);
  }
</script>

<Button
  class={cn("flex justify-between gap-1 px-3 py-2", className)}
  on:click={() => ($CommandStore = !$CommandStore)}
  {...$$restProps}>
  <span>{label}</span>
  <Search class="h-4 w-4" />
</Button>

<CommandDialog bind:open={$CommandStore} shouldFilter={false}>
  <CommandInput
    placeholder="Código o nombre de una materia"
    on:input={debounceSearch} />
  <CommandList>
    {#each materiasFiltradas as mat (mat.codigo)}
      {@const slug = mat.codigo}
      <CommandItem
        value={mat.codigo}
        onSelect={async () => {
          await goto(`/materias/${slug}`);
          $CommandStore = false;
        }}
        class="flex cursor-pointer items-start space-x-1.5">
        <span> {mat.nombre} </span>
      </CommandItem>
    {/each}
  </CommandList>
</CommandDialog>
