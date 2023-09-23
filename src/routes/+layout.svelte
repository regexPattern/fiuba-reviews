<script lang="ts">
	import * as Command from "$lib/components/ui/command";
	import "@fontsource-variable/inter";
	import { onMount } from "svelte";

	import "../app.css";
	import type { LayoutData } from "./$types";

	export let data: LayoutData;
	let open = false;

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

<header class="sticky top-0 z-20 h-16 border-b bg-background/70 backdrop-blur-md">
	<button
		on:click={() => {
			open = !open;
		}}
	>
		Buscar materias...
		<kbd>⌘K</kbd>
	</button>
</header>

<Command.Dialog bind:open>
	<Command.Input placeholder="Buscar materias..." />
	<Command.List>
    <Command.Empty>Sin resultados.</Command.Empty>
		<Command.Group heading="Materias">
      {#each data.materias as mat}
        <Command.Item>{mat.codigo} &bullet; {mat.nombre}</Command.Item>
      {/each}
    </Command.Group>
	</Command.List>
</Command.Dialog>
