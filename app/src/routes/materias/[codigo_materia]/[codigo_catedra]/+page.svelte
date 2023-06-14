<script lang="ts">
	import type { PageData } from "./$types";

	import { Disclosure, DisclosureButton, DisclosurePanel } from "@rgossiaux/svelte-headlessui";

	export let data: PageData;
</script>

<article class="flex flex-col gap-4">
	{#each data.docentes as d}
		<Disclosure defaultOpen={true} class="rounded-lg border p-4">
			<DisclosureButton class="w-full text-left">
				{d.promedio.toFixed(1)} - {d.nombre}
			</DisclosureButton>
			<DisclosurePanel>
				<ul class="list-inside list-disc">
					{#each d.comentario as c}
						<li class="mt-2">{c.contenido}</li>
					{/each}
				</ul>
			</DisclosurePanel>
		</Disclosure>
	{/each}
</article>
