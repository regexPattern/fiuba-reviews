<script lang="ts">
	import type schema from "$lib/zod/schema";
	import type { Form, FormFieldName } from "formsnap/dist/internal/types";

	import { FormField, Label, Select, Validation } from "./ui/form";
	import { SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select";

	export let name: FormFieldName<typeof schema>;
	export let label: string;
	export let config: Form<typeof schema>;

	const id = String(name);
</script>

<FormField {config} {name}>
	<div class="grid grid-cols-2">
		<Label for={id} class="flex items-center">{label}</Label>
		<Select>
			<SelectTrigger>
				<SelectValue placeholder="Seleccionar" />
			</SelectTrigger>
			<SelectContent {id}>
				{#each Array.from({ length: 5 }, (_, i) => i + 1) as val}
					<SelectItem value={val}>{val}</SelectItem>
				{/each}
			</SelectContent>
		</Select>
	</div>
</FormField>
