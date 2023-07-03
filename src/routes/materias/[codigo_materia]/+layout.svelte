<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;

	let mostrarCatedras = false;
</script>

<div class="flex flex-1 flex-col lg:flex-row">
	<aside
		class="top-16 flex shrink-0 flex-col overflow-y-auto lg:sticky lg:max-h-[calc(100vh-4rem)] lg:w-80"
	>
		<div class="top-0 border-b bg-white p-2 dark:bg-black lg:sticky lg:border-b-0">
			<button
				class="w-full p-2 text-left"
				on:click={() => (mostrarCatedras = !mostrarCatedras)}
				aria-controls="listado-catedras"
				aria-expanded={mostrarCatedras}
			>
				<span>{data.materia.codigo}</span>
				<p>{data.materia.nombre}</p>
			</button>
		</div>

		<ul
			id="listado-catedras"
			class={`flex-1 border-b p-2 lg:block lg:border-b-0 ${mostrarCatedras ? "block" : "hidden"}`}
		>
			{#each data.catedras as catedra}
				<li class="flex items-center gap-1 p-2">
					<span>{catedra.promedio.toFixed(1)}</span>
					<a
						class="flex-1 py-1 pl-1 lg:py-0"
						href={`/materias/${data.materia.codigo}/${catedra.codigo}`}
						on:click={() => (mostrarCatedras = false)}>{catedra.nombre}</a
					>
				</li>
			{/each}
		</ul>
	</aside>

	<main class="flex-1 p-4 lg:p-8">
		<slot />
	</main>
</div>
