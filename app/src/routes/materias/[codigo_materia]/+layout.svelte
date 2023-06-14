<script lang="ts">
	import type { LayoutData } from "./$types";

	import { page } from "$app/stores";

	export let data: LayoutData;

	let showList = false;

	function collapseList() {
		showList = false;
	}

	function toggleList() {
		showList = !showList;
	}
</script>

<div class="flex flex-col lg:flex-row">
	<aside class="static top-24 z-10 h-full w-full lg:sticky lg:w-80 xl:w-96">
		<div class="block border-b p-4 lg:hidden">
			<button
				aria-controls="catedras-list"
				aria-expanded={showList}
				class="w-full text-left"
				on:click={toggleList}>Catedras</button
			>
		</div>

		<ul id="catedras-list" class={`${showList ? "block" : "hidden"} px-2 lg:block`}>
			{#each data.catedras as c}
				<li
					class={`${
						$page.params.codigo_catedra === c.codigo ? "font-semibold" : "font-normal"
					} overflow-hidden overflow-ellipsis whitespace-nowrap px-2 py-1`}
				>
					<a
						href={`/materias/${c.codigo_materia}/${c.codigo}`}
						title={c.nombre}
						on:click={collapseList}>{c.promedio.toFixed(1)} - {c.nombre}</a
					>
				</li>
			{/each}
		</ul>
	</aside>

	<main class="flex-1 px-4 lg:px-0">
		<slot />
	</main>
</div>
