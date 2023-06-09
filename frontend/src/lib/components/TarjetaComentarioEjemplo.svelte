<script lang="ts">
	export let contenido: string;
	export let cuatrimestre: string;
	export let docente: { nombre: string };
</script>

<article
	class="flex-col rounded-md border border-light-border p-4 dark:border-dark-border sm:flex"
>
	<p class="line-clamp-6 flex-1 text-dark dark:text-light">
		{contenido}
	</p>
	<p class="mt-3 text-right text-sm font-medium text-dark dark:text-slate-400">
		{docente.nombre} - {cuatrimestre}
	</p>
</article>
