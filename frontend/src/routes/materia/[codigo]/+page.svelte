<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<h1 class="text-xl font-semibold">{data.codigo_materia} - {data.nombre_materia}</h1>
<p class="flex space-x-4">
	Equivalencias: 
	{#each data.codigos_equivalencias as codigo}
		<span class="p-2 bg-gray-200 rounded-sm">{codigo}, </span>
	{/each}
</p>

<ul>
	{#each data.catedras as catedra}
		<li><a href={`/catedra/${catedra.codigo}`}>{catedra.promedio.toFixed(2)} - {catedra.nombre}</a></li>
	{/each}
</ul>
