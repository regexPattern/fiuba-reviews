<script lang="ts">
	import { goto } from "$app/navigation";
	import PlaceholderCatedra from "$lib/components/placeholder-catedra.svelte";
	import { onMount } from "svelte";

	import type { PageData } from "./$types";

	export let data: PageData;

	onMount(async () => {
		const catedras = await data.streamed.catedras;
		goto(`/materias/${data.materia.codigo}/${catedras[0].codigo}`);
	});
</script>

<PlaceholderCatedra />
