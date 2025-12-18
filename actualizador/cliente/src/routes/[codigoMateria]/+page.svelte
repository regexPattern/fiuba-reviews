<script lang="ts">
	import { enhance } from "$app/forms";
	import { SvelteMap } from "svelte/reactivity";
	import type { PageProps } from "./$types";
	import PatchDocente from "./PatchDocente.svelte";
	import PatchCatedra from "./PatchCatedra.svelte";

	let { data }: PageProps = $props();

	let resolucionesMatches = $derived.by(() => {
		const map = new SvelteMap<string, string | null>();
		for (const docente of data.patch.docentes)
			map.set(docente.nombre, docente.matches.at(0)?.codigo ?? null);
		return map;
	});

	$inspect(resolucionesMatches);
</script>

<form method="POST" use:enhance>
	<header class="mb-4 px-6 py-4 flex justify-between border-b border-gray-300">
		<h1 class="text-3xl">
			<span class="font-mono">{data.patch.codigo}</span><span class="mx-2">•</span><span
				>{data.patch.nombre}</span
			>
		</h1>

		<button type="submit" class="rounded-lg border border-gray-300 text-green-700 font-medium px-3"
			>Aplicar cambios</button
		>
	</header>

	<div class="mx-6 grid grid-cols-5 gap-8">
		<section class="col-span-2">
			<h2 class="text-2xl mb-3">Docentes</h2>
			<div class="h-full overflow-y-scroll space-y-3">
				{#each data.patch.docentes as docente (docente.nombre)}
					<PatchDocente {docente} resoluciones={resolucionesMatches} />
				{/each}
			</div>
		</section>

		<section class="col-span-3">
			<h2 class="text-2xl mb-3">Cátedras</h2>
			<div class="grid grid-cols-2 gap-3">
				{#each data.patch.catedras as catedra (catedra.codigo)}
					<PatchCatedra {catedra} resoluciones={resolucionesMatches} />
				{/each}
			</div>
		</section>
	</div>
</form>
