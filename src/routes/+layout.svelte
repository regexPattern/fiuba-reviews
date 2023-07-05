<script lang="ts">
	import { browser } from "$app/environment";
	import "../app.css";

	let modoOscuro = false;

	function alternarModoOscuro() {
		modoOscuro = !modoOscuro;
		modoOscuro
			? document.documentElement.classList.add("dark")
			: document.documentElement.classList.remove("dark");
	}

	if (browser) {
		if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
			document.documentElement.classList.add("dark");
			modoOscuro = true;
		} else {
			document.documentElement.classList.remove("dark");
			modoOscuro = false;
		}
	}
</script>

<header
	class="border-border sticky top-0 z-30 h-16 border-b bg-background supports-[backdrop-filter]:bg-background/[0.75] supports-[backdrop-filter]:backdrop-blur">
	<button on:click={alternarModoOscuro}>Click</button>
</header>

<div class="flex flex-1 lg:container lg:mx-auto">
	<slot />
</div>

<footer class="h-16 border-t p-4">TESTING</footer>
