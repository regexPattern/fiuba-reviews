<script lang="ts">
	import type { PatchCatedra } from "$lib";
	import { CirclePlus, Link } from "@lucide/svelte";
	import { Checkbox, Label } from "bits-ui";

	interface Props {
		catedra: PatchCatedra;
		resoluciones: Map<string, string | undefined>;
	}

	let { catedra, resoluciones }: Props = $props();
</script>

<div class="rounded-xl border bg-card p-4">
	<div>
		<ul class="space-y-2">
			{#each catedra.docentes as docente, i (i)}
				<li>
					<div class="flex items-center gap-2">
						<Checkbox.Root
							id={`${docente.nombre}_catedra_${i}`}
							checked={resoluciones.get(docente.nombre) !== ""}
							onclick={(e) => e.preventDefault()}
							class="data-[state=unchecked]:border-border-input inline-flex size-[22px] items-center justify-center rounded-md border border-muted bg-foreground transition-all duration-150 ease-in-out active:scale-[0.98] data-[state=unchecked]:bg-background"
						>
							{#snippet children({ checked })}
								{@const res = resoluciones.get(docente.nombre)}
								<div class="inline-flex items-center justify-center text-background">
									{#if checked && res === "__CREATE__"}
										<CirclePlus class="size-[14px]" />
									{:else if checked && res !== ""}
										<Link class="size-[14px]" />
									{/if}
								</div>
							{/snippet}
						</Checkbox.Root>
						<Label.Root for={`${docente.nombre}_catedra_${i}`} class="text-base font-normal"
							>{docente.nombre}</Label.Root
						>
					</div>
				</li>
			{/each}
		</ul>
	</div>
</div>
