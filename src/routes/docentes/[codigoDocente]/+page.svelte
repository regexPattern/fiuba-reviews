<script lang="ts">
	import InputCalificacion from "$lib/components/input-calificacion.svelte";
	import Link from "$lib/components/link.svelte";
	import {
		Form,
		FormButton,
		FormField,
		FormItem,
		Label,
		Select,
		Textarea
	} from "$lib/components/ui/form";
	import { SelectContent, SelectItem, SelectTrigger, SelectValue } from "$lib/components/ui/select";
	import { Toaster } from "$lib/components/ui/sonner";
	import schema from "$lib/zod/schema";
	import type { FormOptions } from "formsnap";
	import { ChevronLeft, Loader2 } from "lucide-svelte";
	import { toast } from "svelte-sonner";
	import SuperDebug from "sveltekit-superforms/client/SuperDebug.svelte";

	import type { PageData } from "./$types";

	export let data: PageData;

	const options: FormOptions<typeof schema> = {
		onUpdated: ({ form }) => {
			if (form.valid) {
				toast.success(form.message);
			}
		},
		onError: "apply"
	};
</script>

<Toaster />

<main class="mx-auto max-w-screen-sm space-y-6 p-4 xs:space-y-8">
	<Link
		href={`/materias/${data.codigoMateria}/${data.codigoCatedra}`}
		class="flex items-center gap-1 underline"
	>
		<ChevronLeft class="w-4" />
		Ir a cátedra del docente
	</Link>

	<h1 class="text-5xl font-bold tracking-tight">{data.nombreDocente}</h1>

	<Form
		method="POST"
		form={data.form}
		{schema}
		{options}
		let:config
		let:formValues
		let:submitting
		class="space-y-6 xs:space-y-8"
	>
		<div class="space-y-4">
			<InputCalificacion name="acepta-critica" label="Acepta Crítica" {config} />
			<InputCalificacion name="asistencia" label="Asistencia" {config} />
			<InputCalificacion name="buen-trato" label="Buen Trato" {config} />
			<InputCalificacion name="claridad" label="Claridad" {config} />
			<InputCalificacion name="clase-organizada" label="Clase Organizada" {config} />
			<InputCalificacion name="cumple-horario" label="Cumple Horario" {config} />
			<InputCalificacion name="fomenta-participacion" label="Fomenta Participación" {config} />
			<InputCalificacion name="panorama-amplio" label="Panorama Amplio" {config} />
			<InputCalificacion name="responde-mails" label="Responde Mails" {config} />
		</div>

		<div class="space-y-4">
			<FormField {config} name="comentario">
				<FormItem>
					<Label for="comentario">Comentario (Opcional)</Label>
					<Textarea id="comentario" />
				</FormItem>
			</FormField>

			<FormField {config} name="cuatrimestre">
				<FormItem>
					<Label for="cuatrimestre">Cuatrimestre</Label>
					<Select disabled={!formValues.comentario}>
						<SelectTrigger>
							<SelectValue placeholder="Seleccionar" />
						</SelectTrigger>
						<SelectContent id="cuatrimestre">
							{#each data.cuatrimestres as cuatrimestre}
								<SelectItem value={cuatrimestre.nombre}>{cuatrimestre.nombre}</SelectItem>
							{/each}
						</SelectContent>
					</Select>
				</FormItem>
			</FormField>
		</div>

		<FormButton type="submit" class="items-center gap-1" disabled={submitting}>
			{#if submitting}
				<span>Enviando</span>
				<Loader2 class="h-4 w-4 animate-spin" />
			{:else}
				Enviar
			{/if}
		</FormButton>
	</Form>
</main>
