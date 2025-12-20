<script lang="ts">
	import favicon from "$lib/assets/favicon.svg";
	import "./layout.css";

	let { data, children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="flex h-screen">
	<aside class="h-full w-md overflow-y-scroll border-r border-gray-300">
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
	</aside>

	<main class="h-full w-full">
		{@render children()}
	</main>
</div>
