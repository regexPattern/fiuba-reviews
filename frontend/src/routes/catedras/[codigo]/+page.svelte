<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<h1>{data.catedra.nombre}</h1>
{#each data.docentes as docente}
	<article class="mb-12">
		<h3>{docente.nombre}</h3>
		<ul>
		{#each docente.comentarios as comentario}
			<li>{comentario.cuatrimestre} - {comentario.contenido}</li>
		{/each}
		</ul>
	</article>
{/each}
