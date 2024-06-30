<script lang="ts">
	import { page } from "$app/stores";
	import { Sheet, SheetTrigger } from "$lib/components/ui/sheet";
	import SheetContent from "$lib/components/ui/sheet/sheet-content.svelte";
	import { Skeleton } from "$lib/components/ui/skeleton";
	import { cn } from "$lib/utils";
	import { ChevronDown, Star } from "lucide-svelte";

	import type { LayoutData } from "./$types";

	export let data: LayoutData;

	let open = false;
</script>

<svelte:head>
	<title>{data.materia.codigo} - {data.materia.nombre} | Che FIUBA</title>
</svelte:head>

<div class="relative md:container md:mx-auto md:flex-row">
	<div class="sticky top-16 z-30 md:top-auto">
		<aside
			class="fixed hidden h-[calc(100vh-4rem)] w-80 shrink-0 overflow-y-auto border-r bg-background md:block"
		>
			<div
				class="sticky top-0 flex w-full items-start gap-1.5 border-b bg-background p-3 text-center font-medium"
			>
				{data.materia.codigo}
				<span class="font-bold">&bullet;</span>
				{data.materia.nombre}
			</div>

			<ul class="space-y-1.5 py-2">
				{#await data.streamed.catedras}
					{#each Array(10) as _}
						<li class="px-2 py-0.5">
							<Skeleton class="h-10" />
						</li>
					{/each}
				{:then catedras}
					{#each catedras as cat (cat.codigo)}
						<li class="flex items-center gap-1.5 px-5 py-2 md:pl-2 md:pr-4">
							<span class={`w-[2.5ch] shrink-0 font-medium ${!cat.promedio ? "text-center" : ""}`}>
								{cat.promedio?.toFixed(1) || "-"}
							</span>
							<Star class="h-3 w-3 shrink-0 fill-current pr-0.5 text-yellow-500" />
							<a
								href={`/materias/${$page.params.codigoMateria}/${cat.codigo}`}
								class={cn($page.params.codigoCatedra === cat.codigo && "text-fiuba")}
							>
								{cat.nombre}
							</a>
						</li>
					{/each}
				{/await}
			</ul>
		</aside>

		<Sheet bind:open>
			<SheetTrigger asChild>
				<button
					class="z-20 flex w-full items-center justify-between gap-3 border-b bg-background p-3 text-left font-medium md:hidden"
					on:click={() => {
						open = !open;
						window.scroll({ top: 0, behavior: "instant" });
					}}
				>
					<span class="flex items-start gap-1">
						{data.materia.codigo}
						<span class="font-bold">&bullet;</span>
						{data.materia.nombre}
					</span>
					<ChevronDown class="shrink-0" />
				</button>
			</SheetTrigger>
			<SheetContent class="z-[120] p-0 pt-8" side="left">
				<ul class="h-full space-y-1.5 overflow-y-scroll py-2">
					{#await data.streamed.catedras}
						{#each Array(10) as _}
							<li class="px-2 py-0.5">
								<Skeleton class="h-10" />
							</li>
						{/each}
					{:then catedras}
						{#each catedras as cat (cat.codigo)}
							<li class="flex items-center gap-1.5 px-5 py-2 md:pl-2 md:pr-4">
								<span class={`w-[3ch] shrink-0 font-medium ${!cat.promedio ? "text-center" : ""}`}
									>{cat.promedio?.toFixed(1) || "-"}</span
								>
								<Star class="h-3 w-3 shrink-0 fill-current text-yellow-500" />
								<a
									href={`/materias/${$page.params.codigoMateria}/${cat.codigo}`}
									class={cn($page.params.codigoCatedra === cat.codigo && "text-fiuba")}
									on:click={() => (open = !open)}
								>
									{cat.nombre}
								</a>
							</li>
						{/each}
					{/await}
				</ul>
			</SheetContent>
		</Sheet>
	</div>

	<main class="flex flex-col gap-12 p-4 md:ml-80 md:min-h-[calc(100vh-4rem)] md:p-6">
		<slot />
	</main>
</div>
