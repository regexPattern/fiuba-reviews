<script lang="ts">
	import type { PageData } from "./$types";

	import { Disclosure, DisclosureButton, DisclosurePanel } from "@rgossiaux/svelte-headlessui";

	export let data: PageData;
</script>

<article class="flex flex-col gap-4">
	{#each data.docentes as d}
		<Disclosure defaultOpen={true} class="rounded-lg border dark:border-slate-700 dark:bg-slate-800 shadow shadow-slate-700/50">
			<DisclosureButton class="w-full p-4 text-left font-semibold dark:text-slate-50">
				{d.promedio.toFixed(1)} - {d.nombre}
			</DisclosureButton>

			<DisclosurePanel>
				<ul class="divide-y divide-slate-700 text-slate-300">
					{#each d.comentario as c}
						<li class="px-4 py-3">
							<p class={`inline before:content-['"'] after:content-['"']`}>
								{c.contenido}
							</p>
							<span class="text-slate-500"> - {c.cuatrimestre}</span>
						</li>
					{/each}
				</ul>
			</DisclosurePanel>
		</Disclosure>
	{/each}
</article>
