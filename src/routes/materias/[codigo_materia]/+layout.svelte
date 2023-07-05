<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;

	let mostrarCatedras = false;
</script>

<div class="flex flex-1 flex-col lg:flex-row">
	<aside
		class="top-16 flex shrink-0 flex-col overflow-auto overflow-y-auto lg:sticky lg:max-h-[calc(100vh-4rem)] lg:w-80 divide-y"
	>
		<div class="sticky top-0 border-b border-border bg-background p-2 lg:border-b-0">
			<button
				class="w-full items-center p-2 text-left text-sm leading-6"
				on:click={() => (mostrarCatedras = !mostrarCatedras)}
				aria-label="muestra u oculta el listado de catedras"
				aria-controls="listado-catedras"
				aria-expanded={mostrarCatedras}
			>
				Catedras de materia {data.materia.codigo}<br/>
				<span class="uppercase"
					>{data.materia.nombre}</span
				>
				ordenadas por promedio 
			</button>
		</div>

		<div id="listado-catedras" class={`flex-1 ${mostrarCatedras ? "block" : "hidden"} lg:block `}>
			<ul class="border-b border-border p-2 lg:border-b-0">
				{#each data.catedras as catedra}
					<li class="flex items-center gap-1 p-2 text-sm">
						<div class="flex h-6 w-6 items-center justify-center rounded font-medium text-fiuba">
							{catedra.promedio.toFixed(1)}
						</div>
						<a
							class="block flex-1 p-1"
							href={`/materias/${data.materia.codigo}/${catedra.codigo}`}
							on:click={() => (mostrarCatedras = false)}>{catedra.nombre}</a
						>
					</li>
				{/each}
			</ul>
		</div>
	</aside>

	<main class="flex-1 p-4">
		<slot />
	</main>
</div>
