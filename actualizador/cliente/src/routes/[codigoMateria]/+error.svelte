<script lang="ts">
	import { page } from "$app/state";
	import imagen404 from "$lib/assets/404.webp";
	import imagen500 from "$lib/assets/500.webp";
	import imagenUnknown from "$lib/assets/unknown.webp";
</script>

<div class="flex h-screen w-full flex-col items-center justify-center space-y-8">
	{#if page.status === 404}
		<img width={520} alt="Colectivo del AMBA ruta 404" src={imagen404} />
	{:else if page.status >= 500}
		<img width={520} alt="Colectivo del AMBA ruta 500" src={imagen500} />
	{:else}
		<img width={520} alt="Colectivo del AMBA ruta desconocida" src={imagenUnknown} />
	{/if}
	<div class="space-y-4 text-center">
		<h1 class="font-mono text-8xl font-semibold">{page.status}</h1>
		<p class="w-72">{page.error?.message}</p>
	</div>
</div>
