<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;
</script>

<div class="bg-red grid grid-cols-4">
	<aside>
		{#each data.catedras as c}
			<div>
				<a href={`/materias/${c.codigo_materia}/catedras/${c.codigo}`}
					>{c.promedio.toFixed(1)} - {c.nombre}</a
				>
			</div>
		{/each}
	</aside>
	<main class="col-span-3">
		<slot />
	</main>
</div>
