<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<h1>{data.materia.codigo} - {data.materia.nombre}</h1>
{#each data.materia.equivalencias as equivalencia}
	<p>{equivalencia.codigo} - {equivalencia.nombre}</p>
{/each}
