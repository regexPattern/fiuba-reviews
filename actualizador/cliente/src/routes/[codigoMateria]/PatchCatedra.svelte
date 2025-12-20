<script lang="ts">
	import type { PatchCatedra } from "$lib";

	interface Props {
		catedra: PatchCatedra;
		resoluciones: Map<string, string | null>;
	}

	let { catedra, resoluciones }: Props = $props();
</script>

<div class="rounded-xl border border-gray-300 p-3">
	<ul>
		{#each catedra.docentes as docente (docente.nombre)}
			<li>
				{#if docente.codigo_ya_resuelto !== null}
					<input
						type="hidden"
						name={docente.nombre}
						value={JSON.stringify(docente.codigo_ya_resuelto)}
					/>
				{/if}
				<div class="space-x-1">
					<input
						type="checkbox"
						checked={resoluciones.has(docente.nombre)}
						onclick={(e) => e.preventDefault()}
					/><span>{docente.nombre}</span>
				</div>
			</li>
		{/each}
	</ul>
</div>
