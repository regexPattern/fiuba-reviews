<script lang="ts">
	export let contenido: string;
	export let cuatrimestre: string;
	export let docente: { nombre: string };
</script>

<article class="flex flex-col rounded-md border p-4">
	<p class="line-clamp-6 flex-1">{contenido}</p>
	<p class="mt-2 text-right">
		{docente.nombre} - {cuatrimestre}
	</p>
</article>
