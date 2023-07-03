<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;

	let mostrarCatedras = false;

	function mostrarOOcultarCatedras() {
		mostrarCatedras = !mostrarCatedras;
	}
</script>

<div class="flex flex-1 flex-col lg:flex-row">
	<aside
		class="top-16 flex shrink-0 flex-col overflow-y-auto lg:sticky lg:max-h-[calc(100vh-4rem)] lg:w-80"
	>
		<div class="top-0 border-b bg-white p-2 lg:sticky lg:border-b-0">
			<button
				class="w-full p-2 text-left"
				on:click={mostrarOOcultarCatedras}
				aria-controls="listado-catedras"
				aria-expanded={mostrarCatedras}
			>
				{data.materia.codigo} - {data.materia.nombre}
			</button>
		</div>

		<ul
			id="listado-catedras"
			class={`flex-1 border-b p-2 lg:block lg:border-b-0 ${mostrarCatedras ? "block" : "hidden"}`}
		>
			{#each data.catedras as catedra}
				<li class="flex items-center gap-1 p-2">
					<span>{catedra.promedio.toFixed(1)}</span>
					<a class="flex-1 pl-1" href={`/materias/${data.materia.codigo}/${catedra.codigo}`}
						>{catedra.nombre}</a
					>
				</li>
			{/each}
		</ul>
	</aside>

	<main class="flex-1 p-4">
		<slot />
	</main>
</div>
