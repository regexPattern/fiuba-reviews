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

			{#if docente.descripcion}
				<div
					class="divide divide-y divide-slate-300 rounded-lg border border-slate-300 dark:divide-slate-700 dark:border-slate-700"
				>
					<p class="rounded-t-lg bg-slate-100 p-3 dark:bg-slate-800">
						{docente.descripcion}
					</p>
					<small
						class="block rounded-b-lg bg-slate-100 p-3 text-right text-xs dark:bg-slate-800 text-slate-500 dark:text-slate-400"
					>
						Resumen generado por IA. <a href="/" class="underline underline-offset-2 after:content-['↗']"
							>Más información</a
						>.
					</small>
				</div>
			{/if}

			<ul class="list-inside list-disc space-y-4">
				{#each docente.comentarios as comentario}
					<li>{comentario.contenido}</li>
				{/each}
			</ul>
		</div>
	{/each}
</div>
