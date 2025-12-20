<script lang="ts">
	import type { PatchDocente } from "$lib";

	interface Props {
		docente: PatchDocente;
		resoluciones: Map<string, string | null>;
	}

	let { docente, resoluciones }: Props = $props();

	let nombreDb = $derived.by(() => {
		const primerNombre = docente.nombre.split(" ").at(0);
		if (!primerNombre) {
			return null;
		}
		return primerNombre.charAt(0).toUpperCase() + primerNombre.slice(1);
	});

	let codigoMatch = $state<string | null >(null);

	$effect(() => {
		resoluciones.set(docente.nombre, codigoMatch);
	});

	let resolucionJSON = $derived(JSON.stringify({ nombre_db: nombreDb, codigo_match: codigoMatch }));
</script>

<div class="space-y-3 rounded-xl border border-gray-300 p-3">
	<div class="space-x-1">
		<div class="text-xl font-medium">{docente.nombre}</div>
		<div class="text-sm text-gray-500">({docente.rol})</div>
	</div>

	<div>
		<label>
			<span class="text-sm">Nombre a mostrar</span>
			<input
				type="text"
				class="mt-0.5 w-full rounded-md border border-gray-300 p-1"
				bind:value={nombreDb}
			/>
		</label>
	</div>

	<div>
		<span class="text-sm">Matches propuestos</span>
		<div class="mt-0.5">
			{#each docente.matches as match (match.codigo)}
				<label class="flex items-center gap-1">
					<input type="radio" value={match.codigo} bind:group={codigoMatch} />
					<span>{match.nombre}</span><span class="text-sm text-gray-500"
						>(similitud {match.score.toFixed(2)})</span
					>
				</label>
			{/each}

			<label class="flex items-center gap-1">
				<input type="radio" value={null} bind:group={codigoMatch} />
				<span>Registrar nuevo docente</span>
			</label>
		</div>
	</div>

	<input type="hidden" name={`${docente.nombre}`} value={resolucionJSON} />
</div>
