<script lang="ts">
	import type { ActionData, PageData } from "./$types";

	export let data: PageData;
	// export let form: ActionData | undefined;
</script>

<div class="flex flex-col gap-12">
	{#each data.docentes as docente (docente.codigo)}
		<div class="space-y-3">
			<div class="flex items-center gap-2">
				<span class="rounded-xl border p-1.5 font-medium">{docente.promedio.toFixed(1)}</span>
				<p class="text-2xl font-medium">{docente.nombre}</p>
			</div>

			<div
				class="flex flex-col gap-2 rounded-lg border border-slate-200 bg-slate-100 p-4 dark:border-slate-700 dark:bg-slate-900"
			>
				<p>{docente.descripcion}</p>
			</div>

			{#each docente.comentarios as comentario}
				<p>&bull; {comentario.contenido}</p>
			{/each}
		</div>
	{/each}
</div>
