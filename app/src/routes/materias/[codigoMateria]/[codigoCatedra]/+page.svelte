<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<div class="flex flex-col">
	{#each data.catedra.docentes as docente (docente.codigo)}
		<div class="mb-8">
			<p>{docente.nombre}</p>
			{#each docente.comentario as comentario}
				<p>{comentario.contenido}</p>
			{/each}
		</div>
	{/each}
</div>
