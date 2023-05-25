<script lang="ts">
	import { Button } from "$components/ui/button";
	import { Input } from "$components/ui/input";
	import MateriaLink from "$components/ui/materiaLink/MateriaLink.svelte";
	import type { PageData } from "./$types";

	export let data: PageData;

	let filtro: string;

	function buscarMateria() {}
</script>

<main class="my-4">
	<h1 class="scroll-m-20 text-4xl font-bold tracking-tight lg:text-5xl">Materias</h1>
	<p class="mt-2 text-xl text-muted-foreground">Busca entre todas las materias de la FIUBA.</p>

	<form
		class="mt-2 flex flex-col gap-2 shadow-sm sm:flex-row"
		on:submit|preventDefault={buscarMateria}
	>
		<Input type="text" placeholder="Código o nombre de materia" bind:value={filtro} />
		<Button type="submit" class="font-semibold">Buscar</Button>
	</form>

	<ul class="mt-4 flex flex-col divide-y divide-border rounded border border-border shadow-sm">
		{#each data.materias as materia}
			<li class="bg-card p-1 [&:first-child]:rounded-t">
				<MateriaLink {materia} />
			</li>
		{/each}
	</ul>
</main>
