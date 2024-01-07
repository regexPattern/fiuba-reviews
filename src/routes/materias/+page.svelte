<script lang="ts">
	import Link from "$lib/components/link.svelte";
	import Input from "$lib/components/ui/input/input.svelte";

	import type { PageData } from "./$types";

	export let data: PageData;

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

	$: if (search.length === 0) {
		filtered = data.materias;
	}
</script>

<main class="space-y-3 p-4 md:container md:mx-auto">
	<Input
		bind:value={search}
		on:input={debounceSearch}
		placeholder="Código o nombre de una materia..."
		class="mx-auto max-w-xl"
	/>
	<ul class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
		{#each filtered as mat}
			{@const url = mat.codigoEquivalencia ?? mat.codigo}
			<li
				class="text-muted-foreground rounded-lg border bg-slate-50 text-center text-sm dark:bg-slate-900"
			>
				<Link href={`/materias/${url}`} class="flex flex-col p-3">
					<span class="text-foreground font-medium">{mat.codigo} &bull; {mat.nombre}</span>
					{#if mat.codigoEquivalencia}
						<span>Equivalente a {mat.codigoEquivalencia}</span>
					{/if}
					<span>{mat.cantidadCatedras} cátedra{mat.cantidadCatedras > 1 ? "s" : ""}</span>
				</Link>
			</li>
		{/each}
	</ul>
</main>
