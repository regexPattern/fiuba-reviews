<script lang="ts">
	import { page } from "$app/stores";

	let status: number = $page.status || 500;
	let message = $page.error?.message || "Error interno.";
</script>

<svelte:head>
	<title>{status} | Che FIUBA</title>
</svelte:head>

<div class="flex items-center justify-center gap-6 divide-x p-8">
	<h1 class="text-2xl font-medium tracking-tight">{status}</h1>
	<p class="py-4 pl-6">{message}</p>
</div>
