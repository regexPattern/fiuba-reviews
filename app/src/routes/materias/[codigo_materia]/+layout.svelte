<script lang="ts">
	import type { LayoutData } from "./$types";

	import { page } from "$app/stores";

	export let data: LayoutData;

	let showList = false;

	const collapseList = () => (showList = false);
	const toggleList = () => (showList = !showList);
</script>

<div class="flex flex-col lg:flex-row">
	<aside class="static top-24 z-10 h-full w-full lg:sticky lg:w-96">
		<div
			class="mb-4 block border-b border-slate-300 p-4 text-slate-800 dark:border-slate-700 dark:text-slate-100 lg:hidden"
		>
			<button
				aria-controls="catedras-list"
				aria-expanded={showList}
				class="w-full text-left"
				on:click={toggleList}>Catedras</button
			>
		</div>

		<ul
			id="catedras-list"
			class={`${showList ? "block" : "hidden"} px-2 text-slate-500 dark:text-slate-400 lg:block`}
		>
			{#each data.catedras as c}
				<li class="group flex gap-2 overflow-hidden overflow-ellipsis whitespace-nowrap p-2">
					<a
						href={`/materias/${c.codigo_materia}/${c.codigo}`}
						title={c.nombre}
						on:click={collapseList}
					>
						<span class="w-6 text-center font-medium text-fiuba">
							{c.promedio.toFixed(1)}
						</span>
						&#x2022;
						<span
							class={`${
								$page.params.codigo_catedra === c.codigo
									? "font-semibold text-fiuba"
									: "font-normal group-hover:text-slate-800 dark:group-hover:text-slate-100"
							}`}>{c.nombre}</span
						></a
					>
				</li>
			{/each}
		</ul>
	</aside>

	<main class="flex-1 px-4 lg:px-0">
		<slot />
	</main>
</div>
