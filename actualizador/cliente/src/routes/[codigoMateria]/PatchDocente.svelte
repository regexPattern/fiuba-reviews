<script lang="ts">
	import type { PatchDocente } from "$lib";
	import { Label, RadioGroup } from "bits-ui";

	interface Props {
		docente: PatchDocente;
		resoluciones: Map<string, string>;
		matchesYaAsignados: Map<string, string>;
	}

	let { docente, resoluciones, matchesYaAsignados }: Props = $props();

	/*
	Se sugiere el apellido del docente como posible nombre en la base de
	datos. Muchos docentes van a requirir algún tipo de intervención
	manual para ajustar este valor, ya sea por falta de tilde o por la
	complejidad con apellidos compuestos.
	*/
	let nombreDb = $derived.by(() => {
		const primerNombre = docente.nombre.split(" ").at(0);
		if (!primerNombre) {
			return null;
		}
		return primerNombre.charAt(0).toUpperCase() + primerNombre.slice(1);
	});

	let codigoMatch = $state("");

	let resolucionJSON = $derived(
		JSON.stringify({
			rol: docente.rol,
			nombre_db: nombreDb,
			codigo_match: codigoMatch
		})
	);

	$effect(() => {
		const prevMatch = resoluciones.get(docente.nombre);
		if (prevMatch) {
			matchesYaAsignados.delete(prevMatch);
		}
		resoluciones.set(docente.nombre, codigoMatch);
		if (codigoMatch !== "__CREATE__") {
			matchesYaAsignados.set(codigoMatch, docente.nombre);
		}
	});
</script>

<div class="space-y-4 rounded-xl border bg-card p-4">
	<div class="flex flex-col">
		<span class="text-lg font-semibold">{docente.nombre}</span>
		<span class="text-sm text-muted-foreground">({docente.rol})</span>
	</div>
	<div class="">
		<Label.Root for={`${nombreDb}_nombre`} class="text-sm text-muted-foreground">
			Nombre a mostrar
		</Label.Root>
		<input
			id={`${nombreDb}_nombre`}
			bind:value={nombreDb}
			class="border-border-input w-full rounded-md border bg-background px-2 py-1"
		/>
	</div>
	<div class="space-y-2">
		<Label.Root for={`${nombreDb}_match`} class="text-sm text-muted-foreground"
			>Matches propuestos</Label.Root
		>
		<div class="space-y-2">
			<RadioGroup.Root id={`${nombreDb}_match`} bind:value={codigoMatch}>
				{#each docente.matches as match (match.codigo)}
					{@const id = `${docente.nombre}-${match.codigo}`}
					{@const currMatch = matchesYaAsignados.get(match.codigo)}

					<div class="flex gap-1.5">
						<RadioGroup.Item
							{id}
							value={match.codigo}
							disabled={currMatch !== undefined && currMatch !== docente.nombre}
							class="border-border-input hover:border-dark-40 size-5 shrink-0 cursor-default rounded-full border bg-background transition-all duration-100 ease-in-out data-[state=checked]:border-6 data-[state=checked]:border-foreground"
						/>
						<label for={id}>
							{match.nombre}
							<span class="text-sm text-muted-foreground">(similitud {match.score.toFixed(2)})</span
							>
						</label>
					</div>
				{/each}

				<div class="flex items-center gap-1.5">
					<RadioGroup.Item
						id={`${docente.nombre}-__CREATE__`}
						value={"__CREATE__"}
						class="border-border-input hover:border-dark-40 size-5 shrink-0 cursor-default rounded-full border bg-background transition-all duration-100 ease-in-out data-[state=checked]:border-6 data-[state=checked]:border-foreground"
					/>
					<Label.Root for={`${docente.nombre}-__CREATE__`}>Registrar nuevo docente</Label.Root>
				</div>
			</RadioGroup.Root>
		</div>
	</div>
</div>

<input type="hidden" name={`${docente.nombre}`} value={resolucionJSON} />
