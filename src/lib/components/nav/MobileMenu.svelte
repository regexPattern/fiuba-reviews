<script lang="ts">
	import staticData from "$lib/staticData";
	import Link from "./Link.svelte";
	import { Dialog, DialogOverlay } from "@rgossiaux/svelte-headlessui";
	import GithubIcon from "~icons/bi/github";
	import CloseIcon from "~icons/ic/round-close";
	import MoonIcon from "~icons/lucide/moon";
	import SunIcon from "~icons/lucide/sun";
	import MenuIcon from "~icons/tabler/dots";

	export let modoOscuro: boolean;
	export let alternarModoOscuro: () => void;

	let menuAbierto = false;

	const abrirMenu = () => (menuAbierto = true);
	const cerrarMenu = () => (menuAbierto = false);
</script>

<div class="block sm:hidden">
	<button
		class="p-2 text-slate-600 dark:text-slate-200"
		on:click={abrirMenu}
		aria-controls="menu"
		aria-expanded={menuAbierto}
		aria-label="Abrir menu de navegacion"
	>
		<MenuIcon class="h-[20px] w-[20px]" />
	</button>

	<Dialog id="menu" open={menuAbierto} on:close={cerrarMenu}>
		<DialogOverlay
			class="fixed left-0 top-0 z-40 h-full w-full supports-[backdrop-filter]:bg-background/[0.65] supports-[backdrop-filter]:backdrop-blur"
		/>

		<div
			class="bordered-overlay fixed right-4 top-4 z-50 flex flex-col space-y-4 divide-y divide-slate-200 rounded-lg bg-white p-4 text-slate-500 shadow-lg dark:divide-slate-700 dark:bg-slate-800 dark:text-slate-400"
		>
			<div class="flex flex-row-reverse gap-6">
				<button class="self-start p-2 focus:outline-none" on:click={cerrarMenu}>
					<CloseIcon />
				</button>

				<nav class="with-ring mt-2 text-base font-medium">
					<ul class="space-y-4">
						{#each staticData.navLinks as link}
							<li><Link href={link.href} label={link.label} on:click={cerrarMenu} /></li>
						{/each}
					</ul>
				</nav>
			</div>

			<div
				class="flex items-center justify-center gap-4 px-2 pt-4 [&>*>svg]:h-[20px] [&>*>svg]:w-[20px]"
			>
				<button
					class="p-2"
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
					class="p-2"
					on:click={cerrarMenu}
					target="_blank"
				>
					<GithubIcon />
				</a>
			</div>
		</div>
	</Dialog>
</div>
