<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import { Textarea } from "$lib/components/ui/textarea";
	import { superForm } from "sveltekit-superforms/client";
	import SuperDebug from "sveltekit-superforms/client/SuperDebug.svelte";

	import type { PageData } from "./$types";

	export let data: PageData;
	const { form, enhance } = superForm(data.form);
</script>

<SuperDebug data={$form} />

<form method="POST" class="space-y-4 p-4" use:enhance>
	<div class="grid gap-3 grid-cols-3">
		<div>
			<Label for="acepta_critica">Acepta Crítica</Label>
			<Input
				id="acepta_critica"
				name="acepta_critica"
				type="number"
				step="0.1"
				bind:value={$form.acepta_critica}
			/>
		</div>

		<div>
			<Label for="asistencia">Acepta Crítica</Label>
			<Input
				id="asistencia"
				name="asistencia"
				type="number"
				step="0.1"
				bind:value={$form.asistencia}
			/>
		</div>
	</div>
	<div>
		<Label for="comentario">Comentario</Label>
		<Textarea id="comentario" name="comentario" bind:value={$form.comentario} />
	</div>
	<Button type="submit">Enviar</Button>
</form>
