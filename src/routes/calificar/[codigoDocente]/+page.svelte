<script lang="ts">
	import AnchorTag from "$lib/components/AnchorTag.svelte";
	import InputCalificacion from "$lib/components/InputCalificacion.svelte";
	import {
		Form,
		FormButton,
		FormField,
		FormItem,
		Label,
		Textarea,
		Validation
	} from "$lib/components/ui/form";
	import { ChevronLeft, Loader2 } from "lucide-svelte";

	import type { PageData } from "./$types";
	import schema from "./schema";

	export let data: PageData;
</script>

<main class="mx-auto max-w-screen-sm space-y-6 p-4">
	<AnchorTag
		href={`/materias/${data.codigoMateria}/${data.codigoCatedra}`}
		class="flex items-center gap-1 underline"
		><ChevronLeft class="w-4" />Ir a cátedra del docente</AnchorTag
	>
	<h1 class="text-5xl font-bold tracking-tight">{data.nombreDocente}</h1>
	<Form
		{schema}
		form={data.form}
		options={{ delayMs: 1000, timeoutMs: 1000 }}
		let:config
		let:delayed
		class="space-y-6"
	>
		<div class="space-y-3">
			<FormField {config} name="acepta_critica">
				<InputCalificacion id="acepta-critica" label="Acepta Crítica" />
			</FormField>
			<FormField {config} name="asistencia">
				<InputCalificacion id="asistencia" label="Asistencia" />
			</FormField>
			<FormField {config} name="buen_trato">
				<InputCalificacion id="buen-trato" label="Buen Trato" />
			</FormField>
			<FormField {config} name="claridad">
				<InputCalificacion id="claridad" label="Claridad" />
			</FormField>
			<FormField {config} name="clase_organizada">
				<InputCalificacion id="clase-organizada" label="Clase Organizada" />
			</FormField>
			<FormField {config} name="cumple_horario">
				<InputCalificacion id="cumple-horario" label="Cumple Horario" />
			</FormField>
			<FormField {config} name="fomenta_participacion">
				<InputCalificacion id="fomenta-participacion-horario" label="Fomenta Participación" />
			</FormField>
			<FormField {config} name="panorama_amplio">
				<InputCalificacion id="panorama-amplio" label="Panorama Amplio" />
			</FormField>
			<FormField {config} name="responde_mails">
				<InputCalificacion id="responde-mails" label="Responde Mails" />
			</FormField>
		</div>
		<FormField {config} name="comentario">
			<FormItem>
				<Label for="comentario">Comentario (Opcional)</Label>
				<Textarea id="comentario" />
				<Validation />
			</FormItem>
		</FormField>
		<FormButton type="submit" class="items-center gap-1" disabled={delayed}>
			<span>Enviar</span>
			{#if delayed}
				<Loader2 class="mr-2 h-4 w-4 animate-spin" />
			{/if}
		</FormButton>
	</Form>
</main>
