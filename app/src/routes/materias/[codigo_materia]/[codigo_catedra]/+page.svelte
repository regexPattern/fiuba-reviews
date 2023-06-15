<script lang="ts">
	import type { PageData } from "./$types";

	import { Disclosure, DisclosureButton, DisclosurePanel } from "@rgossiaux/svelte-headlessui";

	export let data: PageData;
</script>

<article class="flex flex-col gap-4 text-slate-700 dark:text-slate-300">
	{#each data.docentes as d (`${data.codigo_catedra}-${d.codigo}`)}
		<Disclosure
			defaultOpen={true}
			class="divide-y divide-slate-300 rounded-lg border border-slate-300 bg-slate-50/30 dark:divide-slate-700 dark:border-slate-700 dark:bg-slate-800"
		>
			<DisclosureButton
				class="divide-2 w-full divide-x divide-slate-50 p-4 text-left font-semibold text-slate-800 dark:text-slate-50"
			>
				{d.promedio.toFixed(1)} - {d.nombre}
			</DisclosureButton>

			<DisclosurePanel class="divide-y divide-inherit">
				{#each d.comentario as c}
					<div class="p-4">
						<p class={`inline before:content-['"'] after:content-['"']`}>
							{c.contenido}
						</p>
						<span class="text-sm text-slate-500 dark:text-slate-400"> - {c.cuatrimestre}</span>
					</div>
				{/each}
			</DisclosurePanel>
		</Disclosure>
	{/each}
</article>
