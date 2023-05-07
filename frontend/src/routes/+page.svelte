<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;

	let searchQuery = "";
</script>

<input bind:value={searchQuery} placeholder="Buscar materia" />
<ul>
	{#each data.materias as materia}
		{#if materia.nombre.includes(searchQuery.toLowerCase()) || (searchQuery.length === 4 && materia.codigo === parseInt(searchQuery, 10))}
			<li class="uppercase">
				<a href={`/materias/${materia.codigo}`}>{materia.codigo} - {materia.nombre}</a>
			</li>
		{/if}
	{/each}
</ul>
