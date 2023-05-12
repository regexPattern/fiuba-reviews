<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<ul>
	{#each data.materias as materia}
		<li class="uppercase">
			<a href={`/materias/${materia.codigo_equivalencia || materia.codigo}`}
				>{materia.codigo} - {materia.nombre}</a
			>
		</li>
	{/each}
</ul>
