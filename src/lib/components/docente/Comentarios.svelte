<script lang="ts">
	import { twMerge } from "tailwind-merge";

	type Comentario = {
		codigo: string;
		cuatrimestre: string;
		contenido: string;
	};

	type Docente = {
		nombre: string;
		codigo: string;
		comentarios: Comentario[];
	};

	export let docente: Docente;

	let mostrarComentarios = true;
</script>

<div class="flex space-x-2 text-sm font-medium">
	<a
		href={`/docentes/${docente.codigo}`}
		target="_blank"
		class="flex items-center space-x-1 rounded-lg border border-fiuba bg-fiuba p-2 text-slate-50"
	>
		Agregar Reseña
	</a>

	{#if docente.comentarios.length}
		<button
			class="rounded-lg border border-fiuba p-2 text-fiuba"
			on:click={() => (mostrarComentarios = !mostrarComentarios)}
			aria-label={`muestra u oculta el listado de comentarios de docente ${docente.codigo}`}
			aria-controls={`comentarios-${docente.codigo}`}
			aria-expanded={mostrarComentarios}
		>
			{#if mostrarComentarios} Ocultar {:else} Mostrar {/if}
			{docente.comentarios.length}
			{#if docente.comentarios.length == 1} Comentario {:else} Comentarios {/if}
		</button>
	{/if}
</div>

{#if docente.comentarios.length}
	<div
		class={twMerge(
			"bordered-overlay rounded-lg bg-slate-50 dark:bg-slate-800",
			mostrarComentarios && "border"
		)}
	>
		<ul id={`comentarios-${docente.codigo}`} class="divide-y divide-border">
			{#if mostrarComentarios}
				{#each docente.comentarios as comentario (comentario.codigo)}
					<li class="p-2">
						<p>
							"{comentario.contenido}"
							<span class="text-sm font-medium text-slate-400 dark:text-slate-500">
								&dash; {comentario.cuatrimestre}
							</span>
						</p>
					</li>
				{/each}
			{/if}
		</ul>
	</div>
{/if}
