<script lang="ts">
	import { SvelteMap } from "svelte/reactivity";
	import type { PageProps } from "./$types";

	let { data }: PageProps = $props();

	let docentesResueltos = $state(new SvelteMap<string, string>());

	const handleSubmit = async (
		e: SubmitEvent & { currentTarget: EventTarget & HTMLFormElement }
	) => {
		e.preventDefault();

		const resoluciones = new Map<string, string | null>();

		for (const [nombreSiu, codMatch] of docentesResueltos) {
			resoluciones.set(nombreSiu, codMatch);
		}

		for (const nombreSiu of data.docentesNuevos) {
			resoluciones.set(nombreSiu, null);
		}

		console.log(resoluciones);
	};
</script>

<form method="POST" onsubmit={handleSubmit}>
	<header class="mb-4 px-6 py-4 flex justify-between border-b border-gray-300">
		<h1 class="text-3xl">
			<span class="font-mono">{data.patch.codigo}</span><span class="mx-2">•</span><span
				>{data.patch.nombre}</span
			>
		</h1>

		<button class="rounded-lg border border-gray-300 text-green-700 font-medium px-3"
			>✅ Actualizar</button
		>
	</header>

	<div class="mx-6 grid grid-cols-5 gap-8">
		<section class="col-span-2">
			<h2 class="text-2xl mb-3">Docentes</h2>

			<div class="h-full overflow-y-scroll space-y-3">
				{#each data.patch.docentes as doc (doc.nombre)}
					<div class="p-3 border border-gray-300 rounded">
						<h3 class="space-x-1">
							<span>{doc.nombre}</span><span>•</span><span>{doc.rol}</span>
						</h3>

						{#if doc.matches.length > 0}
							<div class="flex flex-col">
								{#each doc.matches as match (match.codigo)}
									<label>
										<input
											type="radio"
											name={doc.nombre}
											value={match.codigo}
											onchange={() => {
												docentesResueltos.set(doc.nombre, match.codigo);
											}}
										/>
										<span>{match.nombre}</span><span>•</span><span
											>{match.similitud.toFixed(2)}</span
										>
									</label>
								{/each}
								<label>
									<input type="radio" />
									Registrar nuevo docente
								</label>
							</div>
						{:else}
							<label>
								<input type="radio" checked={true} />
								Registrar nuevo docente
							</label>
						{/if}
					</div>
				{/each}
			</div>
		</section>

		<section class="col-span-3">
			<h2 class="text-2xl mb-3">Cátedras</h2>

			<div class="grid grid-cols-2 gap-3">
				{#each data.patch.catedras as cat (cat.codigo)}
					<div class="border rounded border-gray-300 p-3">
						{#each cat.docentes as doc (doc.nombre)}
							<div>
								{#if docentesResueltos.has(doc.nombre)}
									✅
								{:else if data.docentesNuevos.has(doc.nombre)}
									🆕
								{:else}
									❓
								{/if}
								{doc.nombre}
							</div>
						{/each}
					</div>
				{/each}
			</div>
		</section>
	</div>
</form>
