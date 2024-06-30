<script lang="ts">
import type { PageData } from "./$types";

	export let data: PageData;
</script>

<ul>
	{#each data.materias as mat}
		{@const url = mat.codigoEquivalencia || mat.codigo}
		<li><a href={`/materias/${url}`}>{mat.nombre}</a></li>
	{/each}
</ul>
