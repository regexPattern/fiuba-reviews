<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;

	function ordenarPromedioDescendente() {
		data.materia.catedra.sort((a, b) => b.promedio - a.promedio);
		data.materia.catedra = data.materia.catedra;
	}

	function ordenarPromedioAscendente() {
		data.materia.catedra.sort((a, b) => a.promedio - b.promedio);
		data.materia.catedra = data.materia.catedra;
	}

	ordenarPromedioDescendente();
</script>

<h1 class="text-xl font-semibold">{data.materia.codigo} - {data.materia.nombre}</h1>
<button on:click={ordenarPromedioAscendente}>Ascendente</button>
<button on:click={ordenarPromedioDescendente}>Descendente</button>
<div>
	{#each data.materia.other_materia as equivalencia}
		<span>{equivalencia.codigo}, </span>
	{/each}
</div>
<ul>
	{#each data.materia.catedra as catedra}
		<li>
			<a href={`/catedras/${catedra.codigo}`}>{catedra.promedio.toFixed(1)} - {catedra.nombre}</a>
		</li>
	{/each}
</ul>
