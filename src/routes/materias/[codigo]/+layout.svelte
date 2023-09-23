<script lang="ts">
	import type { LayoutData } from "./$types";

	export let data: LayoutData;

	let openMenu = false;
</script>

<div class="flex flex-col md:container md:mx-auto md:flex-row">
	<div class="sticky top-16 md:top-auto">
		<div class="block md:hidden">
			<button class="w-full bg-white p-2" on:click={() => (openMenu = !openMenu)}
				>Abrir/Cerrar</button
			>
		</div>
		<aside
			class={`max-h-64 overflow-y-scroll md:fixed md:h-[calc(100vh-4rem)] md:max-h-full md:w-80 md:shrink-0 ${
				openMenu ? "block" : "hidden"
			} md:block`}
		>
			<ul class="overflow-y-scroll">
				{#each data.catedras as cat}
					<li class="flex items-center gap-1.5 p-1.5">
						<span class="w-10 shrink-0">{cat.promedio.toFixed(2)}</span>
						<span>&bullet;</span>
						<a href={`/materias/${data.materia.codigo}/${cat.codigo}`} class="rounded p-1.5"
							>{cat.nombre}</a
						>
					</li>
				{/each}
			</ul>
		</aside>
	</div>

	<main class="border-4 border-blue-700 bg-blue-400 md:ml-80 md:min-h-[calc(100vh-4rem)]">
		<slot />
	</main>
</div>
