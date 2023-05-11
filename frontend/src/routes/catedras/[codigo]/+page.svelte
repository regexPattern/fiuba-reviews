<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;

	function ordenarAlfabeticamente() {
		data.docentes.sort((a, b) => a.nombre.localeCompare(b.nombre));
		data.docentes = data.docentes;
	}

	function ordenarPorPromedioDescendente() {
		data.docentes.sort((a, b) => b.promedio - a.promedio);
		data.docentes = data.docentes;
	}

	ordenarPorPromedioDescendente();
</script>

<h1>{data.catedra.nombre}</h1>
{#each data.docentes as docente}
	<article class="mb-12">
		<h3>{docente.promedio.toFixed(1)} - {docente.nombre}</h3>
		<ul>
			{#each docente.comentario as comentario}
				<li>{comentario.cuatrimestre} - {comentario.contenido}</li>
			{/each}
		</ul>
	</article>
{/each}
