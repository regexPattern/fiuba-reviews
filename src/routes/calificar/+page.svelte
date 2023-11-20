<script lang="ts">
	import { page } from "$app/stores";
	import {
		Select,
		SelectContent,
		SelectItem,
		SelectTrigger,
		SelectValue
	} from "$lib/components/ui/select";

	import type { PageData } from "./$types";

	export let data: PageData;

	$: codigoDocente = $page.url.searchParams.get("docente");

	$: console.log(codigoDocente);
</script>

<Select>
	<SelectTrigger  class="w-30">
		<SelectValue placeholder="Materia" />
	</SelectTrigger>
	<SelectContent>
		{#each data.materias as mat (mat.codigo)}
			<SelectItem value={mat.codigo}>{mat.codigo} - {mat.nombre}</SelectItem>
		{/each}
	</SelectContent>
</Select>
