<script lang="ts">
	import { browser } from "$app/environment";
	import DesktopMenu from "./nav/DesktopMenu.svelte";
	import MobileMenu from "./nav/MobileMenu.svelte";

	let modoOscuro = false;

	if (browser) {
		if (localStorage.theme === "dark") {
			document.documentElement.classList.add("dark");
			modoOscuro = true;
		} else {
			document.documentElement.classList.remove("dark");
			modoOscuro = false;
		}
	}

	function alternarModoOscuro() {
		modoOscuro = !modoOscuro;

		localStorage.setItem("theme", modoOscuro ? "dark" : "light");

		modoOscuro
			? document.documentElement.classList.add("dark")
			: document.documentElement.classList.remove("dark");
	}
</script>

<header
	class="sticky top-0 z-30 h-16 border-b border-border bg-background supports-[backdrop-filter]:bg-background/[0.65] supports-[backdrop-filter]:backdrop-blur"
>
	<div class="flex h-full items-center justify-between px-4 md:container">
		<h1 class="text-xl tracking-tighter">
			<span class="font-semibold text-fiuba">FIUBA</span>
			<span class="font-medium">Reviews</span>
		</h1>

		<div class="contents">
			<MobileMenu {modoOscuro} {alternarModoOscuro} />
			<DesktopMenu {modoOscuro} {alternarModoOscuro} />
		</div>
	</div>
</header>
