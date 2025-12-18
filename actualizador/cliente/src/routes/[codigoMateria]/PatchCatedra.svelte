<script lang="ts">
	import type { PatchCatedra } from "$lib";

	interface Props {
		catedra: PatchCatedra;
		resoluciones: Map<string, string | null>;
	}

	let { catedra, resoluciones }: Props = $props();
</script>

<div class="border">
	<ul>
		{#each catedra.docentes as docente (docente.nombre)}
			<li>
				{#if resoluciones.has(docente.nombre)}
					✅
				{:else}
					❌
				{/if}
				<span>{docente.nombre}</span>
			</li>
		{/each}
	</ul>
</div>
