<script lang="ts">
	import { enhance } from "$app/forms";
	import type { PageProps } from "./$types";
	import PatchCatedra from "./PatchCatedra.svelte";
	import PatchDocente from "./PatchDocente.svelte";
	import { ChevronRight } from "@lucide/svelte";
	import { Button, ScrollArea } from "bits-ui";
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
				class="rounded-md bg-foreground px-4 py-2 text-sm font-medium text-background ring-primary hover:bg-foreground/95 focus:ring-2 active:scale-[0.98] active:transition-all disabled:bg-foreground/50 disabled:active:scale-[1]"
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
			<h2 class="border-b p-4 text-2xl font-medium tracking-tight">Cátedras</h2>

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
