<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;
</script>

<div class="relative flex flex-1">
	<aside class="sticky top-0 flex max-h-[calc(100vh-4rem)] w-80 shrink-0 flex-col overflow-y-auto">
		<div class="sticky top-0 p-4 bg-white">{data.materia.codigo} - {data.materia.nombre}</div>

		<ul class="flex-1">
			{#each data.catedras as catedra}
				<li class="p-2">
					<a href={`/materias/${data.materia.codigo}/${catedra.codigo}`}
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
