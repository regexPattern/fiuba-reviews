<script lang="ts">
	import { page } from "$app/state";
	import favicon from "$lib/assets/favicon.svg";
	import { cn } from "$lib/utils";
	import "./layout.css";
	import "@fontsource-variable/inter";
	import { ScrollArea } from "bits-ui";

	let { data, children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<aside class="fixed top-0 left-0 flex h-screen w-96 flex-col border-r border-border">
	<div class="flex-1 overflow-hidden">
		<ScrollArea.Root class="h-full">
			<ScrollArea.Viewport class="h-full py-2">
				{#each data.patches as patch (patch.codigo)}
					<div class="p-2">
						<a href={`/${patch.codigo}`} class="flex items-center gap-2">
							<span class="rounded-md border bg-secondary px-1.5 py-1 text-xs tabular-nums"
								>{patch.codigo}</span
							>
							<span
								class={cn(
									page.params.codigoMateria === patch.codigo
										? "text-foreground"
										: "text-muted-foreground/50"
								)}>{patch.nombre}</span
							>
						</a>
					</div>
				{/each}
			</ScrollArea.Viewport>
			<ScrollArea.Scrollbar orientation="vertical">
				<ScrollArea.Thumb />
			</ScrollArea.Scrollbar>
			<ScrollArea.Corner />
		</ScrollArea.Root>
	</div>
	<div class="border-t p-6">
		<p>materias faltantes: {data.patches.length}</p>
		<p>
			materias con solo cambios en cátedras: {data.patches.filter((p) => p.docentes === 0).length}
		</p>
	</div>
</aside>

<main class="ml-96 h-screen">
	{@render children()}
</main>
