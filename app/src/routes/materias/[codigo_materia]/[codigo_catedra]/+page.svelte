<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

{#each data.docentes as d}
	<div>{d.promedio.toFixed(1)} - {d.nombre}</div>
	{#each d.comentario as c}
		<p>{c.contenido}</p>
	{/each}
	<br />
	<br />
	<br />
	<br />
	<br />
	<br />
{/each}
