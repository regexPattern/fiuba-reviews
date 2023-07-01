<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;
</script>

<div class="relative flex h-full">
	<aside class="relative w-80 overflow-y-scroll">
		<ul class="h-full">
			{#each data.catedras as catedra (catedra.codigo)}
				<li class="py-1">
					<a href={`/materias/${catedra.codigo_materia}/${catedra.codigo}`}
						>{catedra.promedio.toFixed(1)} - {catedra.nombre}</a
					>
				</li>
			{/each}
		</ul>
	</aside>

	<main class="flex-1">
		<slot />
	</main>
</div>
