<script lang="ts">
	import type { LayoutData } from "./$types";
	import { twMerge } from "tailwind-merge";
	import FlechaIcon from "~icons/iconamoon/arrow-up-2";

	export let data: LayoutData;

	let mostrarCatedras = false;
</script>

<div class="flex flex-1 flex-col lg:flex-row">
	<aside
		class="top-16 flex shrink-0 flex-col overflow-auto overflow-y-auto lg:sticky lg:max-h-[calc(100vh-4rem)] lg:w-80 border-r border-border"
	>
		<div class="sticky top-0 border-b border-border bg-background p-2 lg:border-b-0">
			<button
				class="with-ring flex w-full items-center justify-between p-2 text-left text-sm leading-6 focus:ring-fiuba lg:hover:cursor-default lg:focus:ring-0"
				on:click={() => (mostrarCatedras = !mostrarCatedras)}
				aria-label="muestra u oculta el listado de catedras"
				aria-controls="listado-catedras"
				aria-expanded={mostrarCatedras}
			>
				<p class="font-medium">
					<span class="font-semibold">{data.materia.codigo}</span> &bull; {data.materia
						.nombre}
				</p>
				<FlechaIcon class={twMerge("mr-2 lg:hidden", !mostrarCatedras && "rotate-180")} />
			</button>
		</div>

		<div id="listado-catedras" class={`flex-1 ${mostrarCatedras ? "block" : "hidden"} lg:block `}>
			<ul class="border-b border-border p-2 lg:border-b-0">
				{#each data.catedras as catedra}
					<li class="flex items-center gap-1 p-1 text-sm">
						<div class="flex h-6 w-6 items-center justify-center rounded font-medium text-fiuba">
							{catedra.promedio.toFixed(1)}
						</div>
						<a
							class="block flex-1 p-2 with-ring focus:ring-fiuba"
							href={`/materias/${data.materia.codigo}/${catedra.codigo}`}
							on:click={() => (mostrarCatedras = false)}>{catedra.nombre}</a
						>
					</li>
				{/each}
			</ul>
		</div>
	</aside>

	<main class="flex-1 lg:ml-2 p-4">
		<slot />
	</main>
</div>
