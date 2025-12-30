<script lang="ts">
	import { page } from "$app/state";
	import img400 from "$lib/assets/400.webp";
	import img404 from "$lib/assets/404.webp";
	import img500 from "$lib/assets/500.webp";
</script>

<div class="flex h-screen w-full flex-col items-center justify-center space-y-8">
	{#if page.status === 404}
		<img width={520} alt="Colectivo del AMBA línea 404" src={img404} />
	{:else if page.status >= 500}
		<img width={520} alt="Colectivo del AMBA línea 500" src={img500} />
	{:else}
		<img width={520} alt="Colectivo del AMBA línea 400" src={img400} />
	{/if}
	<div class="space-y-4 text-center">
		<h1 class="font-mono text-8xl font-semibold">{page.status}</h1>
		<p class="w-72">{page.error?.message}</p>
	</div>
</div>
