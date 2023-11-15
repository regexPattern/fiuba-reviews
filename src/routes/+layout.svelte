<script lang="ts">
	import { goto } from "$app/navigation";
	import { Button } from "$lib/components/ui/button";
	import { CommandDialog, CommandInput, CommandList } from "$lib/components/ui/command";
	import CommandItem from "$lib/components/ui/command/command-item.svelte";
	import {
		DropdownMenu,
		DropdownMenuContent,
		DropdownMenuItem,
		DropdownMenuTrigger
	} from "$lib/components/ui/dropdown-menu";
	import "@fontsource-variable/inter";
	import { onMount } from "svelte";

	import "../app.css";
	import type { LayoutData } from "./$types";

	export let data: LayoutData;

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

	$: filtered = data.materias.filter(
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

<header class="sticky top-0 z-20 space-x-2 border-b bg-background/70 p-3 backdrop-blur-lg">
	<Button
		variant="outline"
		on:click={() => (open = !open)}
		class="p-2 font-normal text-muted-foreground"
	>
		Buscar materias
		<kbd class="ml-4 flex gap-1 rounded border bg-muted px-1.5 py-0.5 font-mono text-xs">
			<span>⌘</span> K
		</kbd>
	</Button>

	<DropdownMenu positioning={{ placement: "bottom-end" }}>
		<DropdownMenuTrigger asChild let:builder>
			<Button builders={[builder]} variant="outline" size="icon">T</Button>
		</DropdownMenuTrigger>
		<DropdownMenuContent>
			<DropdownMenuItem>Claro</DropdownMenuItem>
			<DropdownMenuItem>Oscuro</DropdownMenuItem>
			<DropdownMenuItem>Dispositivo</DropdownMenuItem>
		</DropdownMenuContent>
	</DropdownMenu>
</header>

<CommandDialog bind:open shouldFilter={false}>
	<CommandInput placeholder="Código o nombre de una materia" on:input={debounceSearch} />
	<CommandList>
		{#each filtered as mat}
			{@const slug = mat.codigoEquivalencia || mat.codigo}
			<CommandItem
				value={mat.codigo}
				onSelect={async () => {
					await goto(`/materias/${slug}`);
					open = false;
				}}
			>
				{mat.nombre}
			</CommandItem>
		{/each}
	</CommandList>
</CommandDialog>

<slot />
