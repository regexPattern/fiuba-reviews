<script lang="ts">
	import { FormItem, Label, Select, Validation } from "./ui/form";
	import { SelectContent, SelectItem, SelectTrigger, SelectValue } from "./ui/select";

	export let id: string;
	export let label: string;
</script>

<FormItem class="">
	<div class="flex flex-wrap justify-between gap-1">
		<Label for={id} class="flex shrink-0 items-center">{label}</Label>
		<Select>
			<SelectTrigger class="w-32">
				<SelectValue placeholder="Seleccionar" />
			</SelectTrigger>
			<SelectContent {id}>
				{#each Array.from({ length: 5 }, (_, i) => i + 1) as val}
					<SelectItem value={val}>{val}</SelectItem>
				{/each}
			</SelectContent>
		</Select>
	</div>
	<Validation />
</FormItem>
