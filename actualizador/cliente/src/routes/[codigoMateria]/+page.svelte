<script lang="ts">
	import imagen204 from "$lib/assets/204.webp";
	import type { PageProps } from "./$types";
	import FormularioPatch from "./FormularioPatch.svelte";

	let { data }: PageProps = $props();
</script>

{#if data.patch}
	<FormularioPatch patch={data.patch} />
{:else}
	<img width={520} alt="Colectivo del AMBA ruta 204" src={imagen204} />
	Patch materia ya resuelto con exito
{/if}
