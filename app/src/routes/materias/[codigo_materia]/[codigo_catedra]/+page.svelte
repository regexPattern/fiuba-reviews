<script lang="ts">
	import { enhance } from "$app/forms";
	import type { ActionData, PageData } from "./$types";

	export let data: PageData;
	export let form: ActionData | undefined;
</script>

<div class="flex flex-col space-y-8 divide-y">
	{#each data.catedra.docentes as docente (docente.codigo)}
		<div>
			<p>{docente.promedio.toFixed(1)} - {docente.nombre}</p>

			<div class="border-2 p-4">
				{#if form?.issues}
					{#each form.issues as issue}
						<li>{issue.path} - {issue.message}</li>
					{/each}
				{/if}
				<form method="POST" use:enhance class="flex flex-col">
					<input type="hidden" name="codigo" value={docente.codigo} />

					<select name="cuatrimestre">
						{#each data.cuatrimestres as c}
							<option value={c}>{c}</option>
						{/each}
					</select>

					<textarea name="comentario" class="border-2" />
					<button type="submit">Enviar</button>
				</form>
			</div>

			{#each docente.comentario as comentario}
				<p>{comentario.contenido}</p>
			{/each}
		</div>
	{/each}
</div>
