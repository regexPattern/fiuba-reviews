<script lang="ts">
	import MobilePageLink from "./MobilePageLink.svelte";
	import { Dialog, DialogOverlay } from "@rgossiaux/svelte-headlessui";
	import GithubIcon from "~icons/bi/github";
	import CloseIcon from "~icons/ic/round-close";
	import MoonIcon from "~icons/lucide/moon";
	import SunIcon from "~icons/lucide/sun";
	import MenuIcon from "~icons/tabler/dots";

	export let modoOscuro: boolean;
	export let alternarModoOscuro: () => void;

	let menuAbierto = false;

	function cerrarMenu() {
		menuAbierto = false;
	}
</script>

<div class="block sm:hidden">
	<button
		class="with-ring p-1.5 focus:ring-fiuba"
		on:click={() => (menuAbierto = true)}
		aria-controls="menu"
		aria-expanded={menuAbierto}
		aria-label="Abrir menu de navegacion"
	>
		<MenuIcon class="text-slate-400 dark:text-slate-300" />
	</button>

	<Dialog id="menu" open={menuAbierto} on:close={cerrarMenu}>
		<DialogOverlay
			class="fixed left-0 top-0 z-40 h-full w-full supports-[backdrop-filter]:bg-background/[0.65] supports-[backdrop-filter]:backdrop-blur-lg"
		/>

		<div
			class="fixed right-4 top-4 z-50 flex flex-col rounded-lg border border-slate-200 bg-slate-50 p-4 text-slate-500 shadow-lg dark:border-slate-700 dark:bg-slate-800 dark:text-slate-400"
		>
			<button class="self-end p-2 focus:outline-none" on:click={cerrarMenu}>
				<CloseIcon />
			</button>

			<nav class="with-ring mt-2 text-base font-medium">
				<ul class="space-y-4">
					<li><MobilePageLink href="/" label="Inicio" bind:menuAbierto /></li>
					<li><MobilePageLink href="/materias" label="Materias" bind:menuAbierto /></li>
					<li><MobilePageLink href="/docentes" label="Docentes" bind:menuAbierto /></li>
					<li><MobilePageLink href="https://dollyfiuba.com/" label="Dolly" bind:menuAbierto /></li>
				</ul>
			</nav>

			<div class="mt-4 flex space-x-6 border-t border-slate-300 px-2 pt-4 dark:border-slate-700">
				<button
					class="with-ring p-2 focus:ring-fiuba dark:focus:text-slate-200"
					on:click={alternarModoOscuro}
					aria-label="Alternar entre modo oscuro y claro"
				>
					{#if modoOscuro}
						<MoonIcon />
					{:else}
						<SunIcon />
					{/if}
				</button>
				<a
					href="https://github.com/regexPattern/fiuba-reviews"
					class="with-ring p-2 focus:ring-fiuba dark:focus:text-slate-200"
					target="_blank"
				>
					<GithubIcon />
				</a>
			</div>
		</div>
	</Dialog>
</div>
