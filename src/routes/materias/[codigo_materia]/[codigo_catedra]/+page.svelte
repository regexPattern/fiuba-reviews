<script lang="ts">
	import Calificaciones from "$components/docente/Calificaciones.svelte";
	import Comentarios from "$components/docente/Comentarios.svelte";
	import Descripcion from "$components/docente/Descripcion.svelte";
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<div class="flex flex-col space-y-16">
	{#each data.docentes as docente (docente.codigo)}
		<article
			id={encodeURIComponent(docente.nombre.replaceAll(" ", "-").toLowerCase())}
			class="scroll-mt-[4.5rem] space-y-4"
		>
			<h2 class="text-4xl font-bold tracking-tight">{docente.nombre}</h2>
			<Calificaciones {docente} />

			{#if docente.descripcion}
        <Descripcion descripcion={docente.descripcion} />
			{/if}

			<Comentarios {docente} />
		</article>
	{/each}
</div>
