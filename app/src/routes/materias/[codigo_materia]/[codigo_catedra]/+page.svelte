<script lang="ts">
	import type { PageData } from "./$types";
	import { Disclosure, DisclosureButton, DisclosurePanel } from "@rgossiaux/svelte-headlessui";

	export let data: PageData;
</script>

<svelte:head>
	<title>Materia {data.codigo_materia} | {data.nombre_catedra}</title>
</svelte:head>

<div class="flex flex-col gap-4 md:gap-6 lg:gap-8">
	{#each data.docentes as d (`${data.codigo_catedra}-${d.codigo}`)}
		<Disclosure
			let:open
			defaultOpen={true}
			class="divide-color border-color divide-y rounded-lg border bg-slate-100 shadow shadow-slate-200 dark:bg-[#28374F] dark:shadow-slate-700/50"
		>
			<DisclosureButton
				class="flex w-full items-center gap-2 p-4 font-semibold tracking-tight text-slate-700 dark:text-slate-300"
			>
				<span>{d.promedio.toFixed(1)}</span>
				<span class="h-1 w-1 rounded-full bg-slate-700 dark:bg-slate-300" />
				<span class="text-md grow text-left lg:text-lg">{d.nombre}</span>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					stroke-width="1.5"
					stroke="currentColor"
					class={`${open ? "rotate-180" : "rotate-0"} h-[1.15rem] w-[1.15rem]`}
				>
					<path stroke-linecap="round" stroke-linejoin="round" d="M19.5 8.25l-7.5 7.5-7.5-7.5" />
				</svg>
			</DisclosureButton>

			<DisclosurePanel class="divide-color divide-y rounded-b-lg bg-white dark:bg-slate-800 ">
				{#if d.comentario.length > 0}
					{#each d.comentario as c}
						<div class="p-4">
							<p
								class={`inline text-slate-700 before:content-['"'] after:content-['"'] dark:text-slate-300`}
							>
								{c.contenido}
							</p>
							<span class="text-sm text-slate-500 dark:text-slate-400"> - {c.cuatrimestre}</span>
						</div>
					{/each}
				{:else}
					<div class="p-4 text-center text-slate-500">No hay comentarios</div>
				{/if}
			</DisclosurePanel>
		</Disclosure>
	{/each}
</div>
