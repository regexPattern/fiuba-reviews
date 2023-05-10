<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<h1 class="text-xl font-semibold">{data.codigo_materia} - {data.nombre_materia}</h1>
<ul>
	{#each data.catedras as catedra}
		<li><a href={`/catedra/${catedra.codigo}`}>{catedra.codigo} - {catedra.nombre}</a></li>
	{/each}
</ul>
