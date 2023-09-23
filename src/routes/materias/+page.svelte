<script lang="ts">
	import { Input } from "$lib/components/ui/input";

	import type { PageData } from "./$types";

	export let data: PageData;

	let searchQuery = "";
	let timeout: NodeJS.Timeout | undefined;

	function debounce(e: KeyboardEvent) {
		clearTimeout(timeout);
		timeout = setTimeout(() => {
			if (e.target instanceof HTMLInputElement) {
				searchQuery = e.target.value;
			}
		}, 300);
	}

	$: materias =
		searchQuery.length == 0
			? data.materias
			: data.materias.filter((m) => {
					return (
						m.codigo.includes(searchQuery) ||
						m.codigoEquivalencia?.includes(searchQuery) ||
						m.nombre.includes(searchQuery.toUpperCase())
					);
			  });
</script>

<div class="space-y-4 p-4">
	<h1 class="text-4xl font-bold tracking-tight lg:text-5xl">Materias</h1>
	<p class="leading-7">Buscá entre todas las materias de la facultad y sus equivalencias.</p>

	<Input
		placeholder="Código o nombre de materia"
		class="mx-auto mt-2 block w-full"
		on:keyup={debounce}
	/>

	<ul class="grid grid-cols-1 gap-2 md:grid-cols-2 md:gap-3 lg:grid-cols-3 lg:gap-4">
		{#each materias as m}
			{@const slug = m.codigoEquivalencia ?? m.codigo}

			<li class="flex flex-col rounded border [&>*]:p-4">
				<a href={`/materias/${slug}`} class="flex-1 rounded">
					<span class="font-medium">{m.codigo}</span>
					<span class="text-muted-foreground">&bullet;</span>
					<span>{m.nombre}</span>
				</a>
				<div class="text-sm text-muted-foreground">
					<span>
						{m.cantidadCatedras}
						{m.cantidadCatedras == 1 ? "cátedra" : "cátedras"}
					</span>
					{#if m.codigoEquivalencia}
						<span>&bullet;</span>
						<span>Equivalente a {m.codigoEquivalencia}</span>
					{/if}
				</div>
			</li>
		{/each}
	</ul>
</div>
