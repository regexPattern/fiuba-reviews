<script lang="ts">
	import type { LayoutData } from "./$types";
	import { page } from "$app/stores";

	export let data: LayoutData;

	let showList = false;

	const collapseList = () => (showList = false);
	const toggleList = () => (showList = !showList);
</script>

<div class="flex flex-col lg:flex-row">
	<div class="sticky z-10 h-full w-full text-sm lg:top-24 lg:w-96">
		<div class="border-color border-b p-4 lg:hidden">
			<button
				aria-controls="catedras-list"
				aria-expanded={showList}
				class="flex w-full items-center gap-1 text-left text-slate-700 dark:text-slate-100"
				on:click={toggleList}
				><span>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						stroke-width="1.5"
						stroke="currentColor"
						class={`${showList ? "rotate-0" : "rotate-[270deg]"} h-4 w-4`}
					>
						<path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
					</svg>
				</span>Catedras</button
			>
		</div>
		<ul
			id="catedras-list"
			class={`${
				showList ? "block" : "hidden"
			} border-color border-b px-4 py-3 text-slate-700 dark:text-slate-400 lg:block lg:border-none lg:py-0 lg:pr-0`}
		>
			{#each data.catedras as c}
				<li class="flex w-full gap-4 whitespace-nowrap py-3 lg:py-2">
					<a
						href={`/materias/${c.codigo_materia}/${c.codigo}`}
						title={c.nombre}
						class="w-full truncate hover:text-slate-900 dark:hover:text-slate-300 lg:w-min"
						on:click={collapseList}
					>
						<span class="mr-1 w-6 text-center font-semibold text-fiuba">
							{c.promedio.toFixed(1)}
						</span>
						<span
							class={`${
								$page.params.codigo_catedra === c.codigo
									? "font-semibold text-fiuba"
									: "font-normal"
							}`}>{c.nombre}</span
						></a
					>
				</li>
			{/each}
		</ul>
	</div>

	<main class="mt-4 flex-1 px-4 lg:mt-0">
		<slot />
	</main>
</div>
