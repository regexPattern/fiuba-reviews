<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

{#each data.comentarios as c}
	<div>{c.docente.nombre}</div>
{/each}
