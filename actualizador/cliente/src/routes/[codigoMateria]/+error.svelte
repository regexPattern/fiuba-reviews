<script lang="ts">
	import { page } from "$app/state";
	import imagen404 from "$lib/assets/404.webp";
</script>

<div class="flex h-screen w-full flex-col items-center justify-center space-y-8">
	<img width={520} alt="Colectivo del AMBA ruta 404" src={imagen404} />
	<div class="space-y-4 text-center">
		<h1 class="font-mono text-8xl font-semibold">{page.status}</h1>
		<p class="w-72">{page.error?.message}</p>
	</div>
</div>
