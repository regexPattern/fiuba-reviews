<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;
</script>

<div class="flex">
	<aside>
		<ul>
			{#each data.catedras as catedra (catedra.codigo)}
				<li>
					<a href={`/materias/${catedra.codigo_materia}/${catedra.codigo}`}>{catedra.nombre}</a>
				</li>
			{/each}
		</ul>
	</aside>

	<main class="flex-1">
		<slot />
	</main>
</div>
