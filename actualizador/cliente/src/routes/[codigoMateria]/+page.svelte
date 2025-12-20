<script lang="ts">
	import { enhance } from "$app/forms";
	import type { PageProps } from "./$types";
	import PatchCatedra from "./PatchCatedra.svelte";
	import PatchDocente from "./PatchDocente.svelte";
	import { SvelteMap } from "svelte/reactivity";

	let { data }: PageProps = $props();

	let resoluciones = $derived.by(() => {
		const map = new SvelteMap<string, string | null>();
		for (const docente of data.patch.docentes_sin_resolver) {
			if (docente.matches.length === 0) {
				map.set(docente.nombre, null);
			}
		}
		return map;
	});
</script>

<form method="POST" use:enhance>
	<header class="flex h-18 items-center justify-between border-b border-gray-300 px-6">
		<h1 class="text-3xl">
			<span class="font-mono font-semibold">{data.patch.codigo}</span><span class="mx-2">•</span
			><span>{data.patch.nombre}</span>
		</h1>

		<button
			type="submit"
			class="rounded-lg border border-green-700/50 px-3 py-1 font-medium text-green-700 transition-colors hover:cursor-pointer hover:bg-green-200 focus:bg-green-200"
			>Aplicar cambios</button
		>
	</header>

	<div class="flex space-x-6 px-6 pt-4" style="height: calc(100vh - 4.5rem);">
		<section class="flex h-full w-4/12 flex-col">
			<h2 class="mb-3 text-2xl font-semibold">Docentes</h2>
			<div class="flex flex-1 flex-col gap-3 overflow-y-auto pb-3">
				{#each data.patch.docentes_sin_resolver as docente (docente.nombre)}
					<PatchDocente {docente} {resoluciones} />
				{/each}
			</div>
		</section>

		<section class="flex flex-1 flex-col">
			<h2 class="mb-3 text-2xl font-semibold">Cátedras</h2>
			<div class="grid grid-cols-2 gap-3 overflow-y-auto pb-3">
				{#each data.patch.catedras as catedra, i (i)}
					<PatchCatedra {catedra} {resoluciones} />
				{/each}
			</div>
		</section>
	</div>
</form>
