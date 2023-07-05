<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;

	let filtro: string;
</script>

<main class="my-4">
	<h1 class="scroll-m-20 text-4xl font-bold tracking-tight lg:text-5xl">Materias</h1>
	<p class="text-muted-foreground mt-2 text-xl">Busca entre todas las materias de la FIUBA.</p>

	<form class="mt-2 flex flex-col gap-2 shadow-sm sm:flex-row" on:submit|preventDefault={() => {}}>
		<input type="text" placeholder="Código o nombre de materia" bind:value={filtro} />
		<button type="submit" class="font-semibold">Buscar</button>
	</form>

	<ul class="mt-4 flex flex-col divide-y divide-border rounded border border-border shadow-sm">
		{#each data.materias as m}
			<li class="bg-card p-1 [&:first-child]:rounded-t">
				<a href={`/materias/${m.codigoEquivalencia || m.codigo}`}>{m.codigo} - {m.nombre}</a>
			</li>
		{/each}
	</ul>
</main>
