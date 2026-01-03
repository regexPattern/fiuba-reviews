<script lang="ts">
	import { enhance } from "$app/forms";
	import type { PageProps } from "./$types";
	import PatchCatedra from "./PatchCatedra.svelte";
	import PatchDocente from "./PatchDocente.svelte";
	import { ChevronRight, TriangleAlert } from "@lucide/svelte";
	import { Button, ScrollArea, Tooltip } from "bits-ui";
	import { SvelteMap } from "svelte/reactivity";

	let { data }: PageProps = $props();

	let resoluciones = $derived.by(() => {
		const map = new SvelteMap<string, string>();
		for (const docente of data.patch.docentes_pendientes) {
			if (docente.matches.length === 0) {
				map.set(docente.nombre, "");
			}
		}
		return map;
	});

	let matchesYaAsignados = $state(new SvelteMap<string, string>());

	let posibleErrorScraper = $state.raw(
		(() => {
			const map: Record<string, number> = {};
			for (const catedra of data.patch.catedras) {
				for (const docente of catedra.docentes) {
					const n = map[docente.nombre];
					if (!n) {
						map[docente.nombre] = 1;
					} else if (n === 2) {
						return true;
					} else {
						map[docente.nombre]++;
					}
				}
			}
			return false;
		})()
	);
</script>

<form method="POST" use:enhance>
	<header class="flex h-24 flex-col divide-y border-b border-border">
		<div class="flex items-center justify-between px-6 py-3">
			<h1 class="flex items-center gap-2">
				<span class="text-2xl font-medium tracking-tight tabular-nums">{data.patch.codigo}</span>
				<ChevronRight class="text-muted-foreground/50" />
				<span class="text-xl">{data.patch.nombre}</span>
			</h1>

			<Button.Root
				type="submit"
				disabled={resoluciones.values().some((x) => x === "")}
				class="rounded-md border border-[#33b5f9] bg-primary/85 px-4 py-2 text-sm transition-all hover:bg-primary focus:ring-2 active:bg-primary disabled:border-border disabled:bg-card disabled:text-muted-foreground"
			>
				Aplicar cambios
			</Button.Root>
		</div>

		<div class="flex flex-1 items-center px-6">
			<p class="text-sm text-muted-foreground">
				{data.patch.carrera} • {data.patch.cuatrimestre.numero}C{data.patch.cuatrimestre.anio}
			</p>
		</div>
	</header>

	<div class="flex divide-x">
		<section class="flex h-[calc(100vh-6rem)] w-1/2 flex-col xl:w-1/3">
			<h2 class="border-b p-4 text-2xl font-medium tracking-tight">Docentes</h2>

			<div class="flex-1 overflow-hidden">
				<ScrollArea.Root class="h-full">
					<ScrollArea.Viewport class="h-full p-4">
						<div class="flex flex-col gap-4">
							{#each data.patch.docentes_pendientes as docente (docente.nombre)}
								<PatchDocente {docente} {resoluciones} {matchesYaAsignados} />
							{/each}
						</div>
					</ScrollArea.Viewport>
					<ScrollArea.Scrollbar orientation="vertical">
						<ScrollArea.Thumb />
					</ScrollArea.Scrollbar>
					<ScrollArea.Corner />
				</ScrollArea.Root>
			</div>

			<div class="border-t px-4 py-5">
				<p>
					docentes faltantes: {Array.from(resoluciones.values()).filter((r) => r === "").length}
				</p>
			</div>
		</section>

		<section class="flex h-[calc(100vh-6rem)] flex-1 flex-col">
			<div class="flex items-center justify-between border-b p-4">
				<h2 class="text-2xl font-medium tracking-tight">Cátedras</h2>
				{#if posibleErrorScraper}
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<TriangleAlert class="size-[20px] text-yellow-500" />
							</Tooltip.Trigger>
							<Tooltip.Content
								side="left"
								class="border-border-input mr-2 rounded-lg border bg-card px-3 py-1.5 text-yellow-500"
							>
								Materia tiene posible error en el scraper.
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				{/if}
			</div>

			<div class="flex-1 overflow-hidden">
				<ScrollArea.Root class="h-full">
					<ScrollArea.Viewport class="h-full p-4">
						<div class="grid gap-4 xl:grid-cols-2">
							{#each data.patch.catedras as catedra, i (i)}
								<PatchCatedra {catedra} {resoluciones} />
							{/each}
						</div>
					</ScrollArea.Viewport>
					<ScrollArea.Scrollbar orientation="vertical">
						<ScrollArea.Thumb />
					</ScrollArea.Scrollbar>
					<ScrollArea.Corner />
				</ScrollArea.Root>
			</div>

			<div class="border-t px-4 py-5">
				<svelte:boundary>
					{@const catedrasYaExistentes = data.patch.catedras.filter((c) => c.ya_existente).length}
					<p class="space-x-1.5">
						<span>cátedras nuevas: {data.patch.catedras.length - catedrasYaExistentes}</span>
						<span>•</span>
						<span>cátedras ya existentes: {catedrasYaExistentes}</span>
					</p>
				</svelte:boundary>
			</div>
		</section>
	</div>
</form>
