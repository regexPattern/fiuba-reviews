<script lang="ts">
	import favicon from "$lib/assets/favicon.svg";
	import { ScrollArea } from "$lib/components/ui/scroll-area";
	import "./layout.css";

	let { data, children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex">
	<aside class="w-md border-r border-border">
		<ScrollArea class="h-screen">
			<ul class="m-2">
				{#each data.patches as patch (patch.codigo)}
					<li class="p-2">
						<a href={`/${patch.codigo}`}>
							<span class="font-mono font-semibold">{patch.codigo}</span><span class="mx-1">
								•
							</span><span>{patch.nombre}</span>
						</a>
					</li>
				{/each}
			</ul>
		</ScrollArea>
	</aside>

	<main class="h-screen w-full">
		{@render children()}
	</main>
</div>
