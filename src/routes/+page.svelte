<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import "@splidejs/svelte-splide/css";

	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<main class="mx-auto space-y-12">
	<div class="flex flex-col items-center space-y-6 py-12 text-center">
		<h1 class="px-2 text-7xl font-bold tracking-tighter">
			<span class="text-fiuba">FIUBA</span> Reviews
		</h1>

		<p class="max-w-[40rem] px-8 text-lg xs:text-xl">
			Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros
			estudiantes de la FIUBA. Basado en el legendario
			<a href="http://dollyfiuba.com" target="_blank" rel="noopener noreferrer" class="underline"
				>Dolly FIUBA</a
			>
		</p>

		<Button variant="outline" class="w-64 flex justify-between p-2 font-normal text-muted-foreground">
			<span>Buscar materias</span>
			<kbd class="ml-4 rounded border bg-muted px-1.5 py-0.5 font-mono text-xs">
				<span class="mr-[3px]">⌘</span>K
			</kbd>
		</Button>
	</div>

	<div
		class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
		aria-hidden="true"
	>
		<div
			class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-fiuba to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
			style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
		/>

		<div
			class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
			aria-hidden="true"
		/>
	</div>
</main>
