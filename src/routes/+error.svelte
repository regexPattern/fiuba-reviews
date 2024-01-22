<script lang="ts">
	import { page } from "$app/stores";

	export let status: number | null = null;
	export let message = "Error interno.";
</script>

<div class="flex items-center justify-center gap-6 divide-x p-8">
	<h1 class="text-2xl font-medium tracking-tight">{status || $page.status}</h1>
	<p class="py-4 pl-6">{$page.error?.message || message}</p>
</div>
