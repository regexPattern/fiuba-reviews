<script lang="ts">
	// TODO: https://dev.to/willkre/persistent-theme-switch-dark-mode-with-svelte-sveltekit-tailwind-1b9g
	let darkMode = false;

	function handleSwitchDarkMode() {
		darkMode = !darkMode;

		darkMode
			? document.documentElement.classList.add("dark")
			: document.documentElement.classList.remove("dark");
	}
</script>

<header class="sticky top-0 h-16 border-b border-border bg-background shadow-sm">
	<button on:click={handleSwitchDarkMode}>Toggle</button>
</header>
