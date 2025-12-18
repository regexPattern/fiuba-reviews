<script lang="ts">
	import type { PatchDocente } from "$lib";

	interface Props {
		docente: PatchDocente;
		resoluciones: Map<string, string | null>;
	}

	let { docente, resoluciones }: Props = $props();

	let nombreDb = $derived.by(() => {
		const primerNombre = docente.nombre.split(" ").at(0);
		if (!primerNombre) return null;
		return primerNombre.charAt(0).toUpperCase() + primerNombre.slice(1);
	});

	let codigoMatch = $derived.by<string | null>(() => docente.matches.at(0)?.codigo ?? null);

	$effect(() => {
		resoluciones.set(docente.nombre, codigoMatch);
	});

	let resolucionJSON = $derived(JSON.stringify({ nombre_db: nombreDb, codigo_match: codigoMatch }));
</script>

<div class="p-3 border border-gray-300 rounded">
	<h3 class="space-x-1">
		<span>{docente.nombre}</span><span>•</span><span>{docente.rol}</span>
	</h3>

	<input type="text" class="border" bind:value={nombreDb} />

	{#each docente.matches as match (match.codigo)}
		<label class="flex">
			<input type="radio" value={match.codigo} bind:group={codigoMatch} />
			<span>{match.nombre}</span><span>•</span><span>{match.similitud.toFixed(2)}</span>
		</label>
	{/each}

	<label class="flex">
		<input type="radio" value={null} bind:group={codigoMatch} />
		<span>Registrar nuevo docente</span>
	</label>

	<input type="hidden" name={`${docente.nombre}`} value={resolucionJSON} />
</div>
