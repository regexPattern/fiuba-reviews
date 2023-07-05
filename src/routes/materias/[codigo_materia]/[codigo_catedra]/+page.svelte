<script lang="ts">
	import type { ActionData, PageData } from "./$types";

	export let data: PageData;
	// export let form: ActionData | undefined;
</script>

<div class="flex flex-col gap-12">
	{#each data.docentes as docente (docente.codigo)}
		<div class="space-y-3">
			<div class="flex items-center gap-2 font-semibold">
				<span class="text-sm text-fiuba">{docente.promedio.toFixed(1)}</span>
				<p class="text-2xl tracking-tight">
					{docente.nombre}
				</p>
			</div>

			<div
				class="-bg-slate-100 flex flex-col gap-2 rounded-lg border p-3 dark:border-slate-700 dark:bg-slate-800"
			>
				<p>{docente.descripcion}</p>
			</div>

			<ul class="list-inside list-disc space-y-4">
				{#each docente.comentarios as comentario}
					<li>{comentario.contenido}</li>
				{/each}
			</ul>
		</div>
	{/each}
</div>
