<script lang="ts">
	import type { PatchCatedra } from "$lib";

	interface Props {
		catedra: PatchCatedra;
		resolucionesActuales: Map<string, string | null>;
	}

	let { catedra, resolucionesActuales: _ }: Props = $props();
</script>

<div class="border">
	<ul>
		{#each catedra.docentes as docente (docente.nombre)}
			<li>
				{#if docente.codigo_ya_resuelto !== null}
					✅
					<input
						type="hidden"
						name={docente.nombre}
						value={JSON.stringify(docente.codigo_ya_resuelto)}
					/>
				{/if}
				<span>{docente.nombre}</span>
			</li>
		{/each}
	</ul>
</div>
