<script lang="ts">
	import type { PatchCatedra } from "$lib";
	import * as Card from "$lib/components/ui/card";
	import { Checkbox } from "$lib/components/ui/checkbox";
	import { Label } from "$lib/components/ui/label";

	interface Props {
		catedra: PatchCatedra;
		resoluciones: Map<string, string | undefined>;
	}

	let { catedra, resoluciones }: Props = $props();
</script>

<Card.Root>
	<Card.Content>
		<ul class="space-y-2">
			{#each catedra as docente (docente.nombre)}
				<li>
					<div class="flex gap-1.5">
						<Checkbox
							checked={resoluciones.get(docente.nombre) !== ""}
							onclick={(e) => e.preventDefault()}
						/>
						<Label>{docente.nombre}</Label>
					</div>
				</li>
			{/each}
		</ul>
	</Card.Content>
</Card.Root>
