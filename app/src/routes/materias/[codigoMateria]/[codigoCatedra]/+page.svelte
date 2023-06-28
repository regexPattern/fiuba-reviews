<script lang="ts">
	import { enhance } from "$app/forms";
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

<div class="flex flex-col space-y-8 divide-y">
	{#each data.catedra.docentes as docente (docente.codigo)}
		<div>
			<p>{docente.promedio.toFixed(1)} - {docente.nombre}</p>

			<form method="POST" use:enhance>
				<input type="hidden" name="codigo" value={docente.codigo} />

				<select name="cuatrimestre">
					{#each data.cuatrimestres as c}
						<option value={c}>{c}</option>
					{/each}
				</select>

				<textarea name="comentario" class="border-2"/>
				<button type="submit">Enviar</button>
			</form>

			{#each docente.comentario as comentario}
				<p>{comentario.contenido}</p>
			{/each}
		</div>
	{/each}
</div>
