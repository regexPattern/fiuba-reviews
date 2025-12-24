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
			{#each catedra.docentes as docente, i (i)}
				<li>
					<div class="flex items-center gap-2">
						<Checkbox
							checked={resoluciones.get(docente.nombre) !== ""}
							onclick={(e) => e.preventDefault()}
						/>
						<Label class="text-base font-normal">{docente.nombre}</Label>
					</div>
				</li>
			{/each}
		</ul>
	</Card.Content>
</Card.Root>
