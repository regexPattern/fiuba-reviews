<script lang="ts">
	import { goto } from "$app/navigation";
	import { Button } from "$lib/components/ui/button";
	import {
		CommandDialog,
		CommandGroup,
		CommandInput,
		CommandItem,
		CommandList
	} from "$lib/components/ui/command";
	import { cn } from "$lib/utils";
	import { onMount } from "svelte";

	export let materias: { nombre: string; codigo: string; codigoEquivalencia: string | null }[];

	let className = "";
	export { className as class };

	let open = false;
	let search = "";
	let debounceTimeout: NodeJS.Timeout | undefined;

	function debounceSearch(e: Event) {
		clearTimeout(debounceTimeout);
		debounceTimeout = setTimeout(() => {
			if (e.target instanceof HTMLInputElement) {
				search = e.target.value;
			}
		}, 300);
	}

	$: filtered = materias.filter(
		(m) =>
			m.codigo.includes(search) ||
			m.codigoEquivalencia?.includes(search) ||
			m.nombre.includes(search.toUpperCase())
	);

	$: if (search.length === 0 || open == false) {
		filtered = [];
	}

	onMount(() => {
		function handleKeydown(e: KeyboardEvent) {
			if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
				e.preventDefault();
				open = !open;
			}
		}
		document.addEventListener("keydown", handleKeydown);

		return () => {
			document.removeEventListener("keydown", handleKeydown);
		};
	});
</script>

<Button
	variant="outline"
	class={cn("flex justify-between p-2 font-normal text-muted-foreground", className)}
	on:click={() => (open = !open)}
>
	<span>Buscar materias</span>
	<kbd class="ml-4 rounded border bg-muted px-1.5 py-0.5 font-mono text-xs">
		<span class="mr-[3px]">⌘</span>K
	</kbd>
</Button>

<CommandDialog bind:open shouldFilter={false}>
	<CommandInput placeholder="Código o nombre de una materia" on:input={debounceSearch} />
  <CommandGroup heading="Materias">
    <CommandList>
      {#each filtered as mat (mat.codigo)}
        {@const slug = mat.codigoEquivalencia || mat.codigo}
        <CommandItem
          value={mat.codigo}
          onSelect={async () => {
            await goto(`/materias/${slug}`);
            open = false;
          }}
          class="flex items-start space-x-1.5"
          >
          <span class="font-mono font-semibold">{mat.codigo}</span>
          <span class="font-bold">&bullet;</span>
          <span>
            {mat.nombre}
            {#if mat.codigoEquivalencia}
              (Equivalente a {mat.codigoEquivalencia})
            {/if}
          </span>
        </CommandItem>
      {/each}
    </CommandList>
  </CommandGroup>
</CommandDialog>
