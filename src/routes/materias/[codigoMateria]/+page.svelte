<script lang="ts">
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

	import type { PageData } from "./$types";
	import { Skeleton } from "$lib/components/ui/skeleton";

	export let data: PageData;

	onMount(async () => {
		const catedras = data.catedras;
		goto(`/materias/${data.materia.codigo}/${catedras[0].codigo}`);
	});
</script>

<div class="space-y-4">
	<Skeleton class="h-10 w-72" />
	<div class="flex flex-col gap-2 xs:flex-row [&>*]:h-10 [&>*]:w-full [&>*]:xs:w-44">
		<Skeleton />
		<Skeleton />
	</div>
	<div class="space-y-2">
		{#each Array(10) as _}
			<Skeleton class="h-10" />
		{/each}
	</div>
</div>
