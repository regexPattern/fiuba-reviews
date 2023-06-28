<script lang="ts">
	import { enhance } from "$app/forms";
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<div class="flex flex-col space-y-8 divide-y">
	{#each data.catedra.docentes as docente (docente.codigo)}
		<div>
			<p>{docente.nombre}</p>

			<form method="POST" use:enhance>
				<input type="hidden" name="codigo" value={docente.codigo} />
				<textarea name="comentario" class="border" />
				<button type="submit">Enviar</button>
			</form>

			<!-- {#each docente.comentario as comentario}
				<p>{comentario.contenido}</p>
			{/each} -->
		</div>
	{/each}
</div>
