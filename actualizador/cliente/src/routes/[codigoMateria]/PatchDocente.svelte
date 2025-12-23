<script lang="ts">
	import type { PatchDocente } from "$lib";
	import * as Card from "$lib/components/ui/card";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import * as RadioGroup from "$lib/components/ui/radio-group";

	interface Props {
		docente: PatchDocente;
		resoluciones: Map<string, string | undefined>;
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

	// La lib de UI no permite configurar el valor como undefined manualmente.
	let codigoMatch = $state("");

	let resolucionJSON = $derived(JSON.stringify({ nombre_db: nombreDb, codigo_match: codigoMatch }));

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

<Card.Root>
	<Card.Header>
		<Card.Title>
			{docente.nombre}
		</Card.Title>
		<Card.Description>
			({docente.rol})
		</Card.Description>
	</Card.Header>
	<Card.Content class="space-y-5">
		<div class="space-y-2">
			<Label>Nombre a mostrar</Label>
			<Input bind:value={nombreDb} />
		</div>
		<div class="space-y-2">
			<Label>Matches propuestos</Label>
			<div class="space-y-2">
				<RadioGroup.Root bind:value={codigoMatch}>
					{#each docente.matches as match (match.codigo)}
						{@const currMatch = matchesYaAsignados.get(match.codigo)}

						<div class="flex gap-1.5">
							<RadioGroup.Item
								value={match.codigo}
								disabled={currMatch !== undefined && currMatch !== docente.nombre}
							/>
							<Label
								>{match.nombre}
								<span class="text-muted-foreground">(similitud {match.score.toFixed(2)})</span
								></Label
							>
						</div>
					{/each}

					<div class="flex gap-1">
						<RadioGroup.Item value={"__CREATE__"} />
						<Label>Registrar nuevo docente</Label>
					</div>
				</RadioGroup.Root>
			</div>
		</div>
	</Card.Content>
</Card.Root>

<input type="hidden" name={`${docente.nombre}`} value={resolucionJSON} />
