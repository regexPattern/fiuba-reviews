<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<main class="container mx-auto">
	<div class="mt-12 space-y-6">
		<h1 class="text-center text-7xl font-bold tracking-tighter">
			<span class="text-fiuba">FIUBA</span> Reviews
		</h1>
		<p class="mx-auto max-w-[48rem]">
			Lorem, ipsum dolor sit amet consectetur adipisicing elit. Non laborum numquam voluptas
			asperiores minima ab fugit debitis voluptatibus eos quod! Officiis ducimus dignissimos nemo
			impedit modi dolor quas aperiam natus?
		</p>
	</div>

	<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
		{#each data.materias as mat (mat.codigo)}
			<a href={`/materias/${mat.codigo}`}>
				<div class="rounded-lg border bg-secondary p-3 flex flex-col">
					<span>{mat.codigo} - {mat.nombre}</span>
					<span class="text-muted-foreground">{mat.cantidadCatedras} catedras</span>
					<span class="text-muted-foreground">{mat.cantidadComentarios} comentarios</span>
				</div>
			</a>
		{/each}
	</div>
</main>
