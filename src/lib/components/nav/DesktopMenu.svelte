<script lang="ts">
	import staticData from "$lib/staticData";
	import Link from "./Link.svelte";
	import GithubIcon from "~icons/bi/github";
	import MoonIcon from "~icons/lucide/moon";
	import SunIcon from "~icons/lucide/sun";

	export let modoOscuro: boolean;
	export let alternarModoOscuro: () => void;
</script>

<div class="hidden items-center space-x-4 divide-x divide-border sm:flex">
	<nav class="text-sm font-semibold text-slate-600 dark:text-slate-200">
		<ul class="flex space-x-4 [&>li]:px-1">
			{#each staticData.navLinks as link}
				<li>
					<Link href={link.href} label={link.label} class="transition-colors hover:text-fiuba" />
				</li>
			{/each}
		</ul>
	</nav>

	<div
		class="flex items-center space-x-5 pl-5 text-slate-500 dark:text-slate-400 [&>*>svg]:h-[20px] [&>*>svg]:w-[20px]"
	>
		<button
			class="transition-colors hover:text-slate-700 dark:hover:text-slate-100"
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
			class="transition-colors hover:text-slate-700 dark:hover:text-slate-100"
			target="_blank"
		>
			<GithubIcon />
		</a>
	</div>
</div>
