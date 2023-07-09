<script lang="ts">
	import { page } from "$app/stores";
	import type { LayoutData } from "./$types";
	import { setContext } from "svelte";
	import { get, writable } from "svelte/store";
	import { twMerge } from "tailwind-merge";
	import ChevronIcon from "~icons/iconamoon/arrow-up-2";
	import StarIcon from "~icons/material-symbols/star";

	export let data: LayoutData;

	const mostrarCatedras = writable(false);
	setContext("mostrarCatedras", mostrarCatedras);
</script>

<div class="flex flex-1 flex-col lg:flex-row">
	<aside
		class="top-16 flex shrink-0 flex-col overflow-auto overflow-y-auto border-border lg:sticky lg:max-h-[calc(100vh-4rem)] lg:w-80 lg:border-r"
	>
		<div class="sticky top-0 border-b border-border bg-background p-2 lg:border-b-0">
			<button
				class="flex w-full items-center justify-between p-2 text-left text-sm leading-6 lg:hover:cursor-default lg:focus:ring-0"
				on:click={() => mostrarCatedras.set(!get(mostrarCatedras))}
				aria-label="muestra u oculta el listado de catedras"
				aria-controls="listado-catedras"
				aria-expanded={get(mostrarCatedras)}
			>
				<div class="flex align-middle font-medium">
					<span class="font-semibold text-fiuba">{data.materia.codigo}</span>
					<span class="mr-1.5 rotate-90 text-slate-400" aria-hidden="true"><ChevronIcon /></span>
					<p>{data.materia.nombre}</p>
				</div>

				<ChevronIcon
					class={twMerge("mr-2 text-slate-400 lg:hidden", !get(mostrarCatedras) && "rotate-180")}
				/>
			</button>
		</div>

		<div
			id="listado-catedras"
			class={twMerge("flex-1 lg:block", get(mostrarCatedras) ? "block" : "hidden")}
		>
			<ul class="space-y-1.5 border-b border-border p-4 lg:border-b-0">
				{#each data.catedras as catedra}
					<li class="flex items-center space-x-2 text-sm font-medium">
						<div class="flex space-x-2">
							<StarIcon class="text-yellow-500" />
							<div class="flex items-center justify-center rounded text-fiuba">
								{catedra.promedio.toFixed(1)}
							</div>
						</div>

						<a
							class={twMerge(
								"block flex-1 py-1.5",
								$page.params.codigo_catedra === catedra.codigo && "text-fiuba"
							)}
							href={`/materias/${data.materia.codigo}/${catedra.codigo}`}
							on:click={() => mostrarCatedras.set(false)}>{catedra.nombre}</a
						>
					</li>
				{/each}
			</ul>
		</div>
	</aside>

	<main class="flex-1 p-4 lg:ml-2">
		<slot />
	</main>
</div>
