<script lang="ts">
	import { browser } from "$app/environment";
	import DesktopMenu from "./nav/DesktopMenu.svelte";
	import MobileMenu from "./nav/MobileMenu.svelte";
	import { onMount } from "svelte";

	let modoOscuro = false;

	function alternarModoOscuro() {
		modoOscuro = !modoOscuro;
		modoOscuro
			? document.documentElement.classList.add("dark")
			: document.documentElement.classList.remove("dark");
	}

	onMount(() => {
		if (browser) {
			if (window.matchMedia("(prefers-color-scheme: dark)").matches) {
				document.documentElement.classList.add("dark");
				modoOscuro = true;
			} else {
				document.documentElement.classList.remove("dark");
				modoOscuro = false;
			}
		}
	});
</script>

<header
	class="sticky top-0 z-30 h-16 border-b border-border bg-background supports-[backdrop-filter]:bg-background/[0.65] supports-[backdrop-filter]:backdrop-blur"
>
	<div class="flex h-full items-center justify-between px-4 lg:container">
		<h1 class="text-xl tracking-tighter">
			<span class="font-semibold text-fiuba">FIUBA</span>
			<span class="font-medium">Reviews</span>
		</h1>
		<div class="flex items-center">
			<MobileMenu {modoOscuro} {alternarModoOscuro} />
			<DesktopMenu {modoOscuro} {alternarModoOscuro} />
		</div>
	</div>
</header>
