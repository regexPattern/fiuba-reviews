<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

{#each data.comentariosByDocente as doc}
	{doc.promedio.toFixed(2)} &bullet; {doc.nombre}
	{#each doc.comentarios as com}
		<p class="mb-4">{com.contenido}</p>
	{/each}
{/each}
