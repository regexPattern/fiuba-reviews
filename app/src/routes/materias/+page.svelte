<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<main>
	<ul>
		{#each data.materias as materia}
			<li><a href={`/materias/${materia.codigo}`}>{materia.nombre}</a></li>
		{/each}
	</ul>
</main>
