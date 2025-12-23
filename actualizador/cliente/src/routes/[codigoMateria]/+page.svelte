<script lang="ts">
	import { enhance } from "$app/forms";
	import { Button } from "$lib/components/ui/button";
	import { ScrollArea } from "$lib/components/ui/scroll-area";
	import type { PageProps } from "./$types";
	import PatchCatedra from "./PatchCatedra.svelte";
	import PatchDocente from "./PatchDocente.svelte";
	import { SvelteMap } from "svelte/reactivity";

	let { data, form }: PageProps = $props();

	/* Todos los docentes inician sin resolución, esto para obligar al
	usuario a poner atención al momento de resolver los docentes. */
	let resoluciones = $derived.by(() => {
		const map = new SvelteMap<string, string | undefined>();
		for (const docente of data.patch.docentes_pendientes) {
			if (docente.matches.length === 0) {
				map.set(docente.nombre, undefined);
			}
		}
		return map;
	});

	let matchesYaAsignados = $state(new SvelteMap<string, string>());
</script>

{#if form?.docentesFaltantes}
	{JSON.stringify(form.docentesFaltantes)}
{/if}

<form method="POST" use:enhance>
	<header class="flex h-18 items-center justify-between border-b border-border px-6">
		<h1 class="text-3xl">
			<span class="font-mono font-semibold">{data.patch.codigo}</span><span class="mx-2">•</span
			><span>{data.patch.nombre}</span>
		</h1>

		<Button type="submit" disabled={resoluciones.values().some((x) => x === undefined)}
			>Aplicar cambios</Button
		>
	</header>

	<div class="mt-4 flex flex-1 space-x-4 px-6">
		<section class="flex w-4/12 flex-col">
			<h2 class="text-2xl font-semibold">Docentes</h2>

			<div class="py-4">
				<ScrollArea class="h-[calc(100vh-10rem)]">
					<div class="flex flex-col gap-4">
						{#each data.patch.docentes_pendientes as docente (docente.nombre)}
							<PatchDocente {docente} {resoluciones} {matchesYaAsignados} />
						{/each}
					</div>
				</ScrollArea>
			</div>
		</section>

		<section class="flex flex-1 flex-col">
			<h2 class="text-2xl font-semibold">Cátedras</h2>

			<div class="py-4">
				<ScrollArea class="h-[calc(100vh-10rem)]">
					<div class="grid grid-cols-2 gap-4">
						{#each data.patch.docentes_por_catedra as catedra, i (i)}
							<PatchCatedra {catedra} {resoluciones} />
						{/each}
					</div>
				</ScrollArea>
			</div>
		</section>
	</div>
</form>
