<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;

	let queryFiltro = "";
	$: materiasMostradas = materiasFiltradas(queryFiltro);

	function materiasFiltradas(queryFiltro: string) {
		return data.materias.filter(
			(m) => m.nombre.includes(queryFiltro.toLowerCase()) || m.codigo === parseInt(queryFiltro, 10),
		);
	}
</script>

<input bind:value={queryFiltro} placeholder="Buscar materia" use:debouce />
<ul>
	{#each materiasMostradas as materia}
		<li class="uppercase">
			<a href={`/materia/${materia.codigo_equivalencia || materia.codigo}`}
				>{materia.codigo} - {materia.nombre}</a
			>
		</li>
	{/each}
</ul>
