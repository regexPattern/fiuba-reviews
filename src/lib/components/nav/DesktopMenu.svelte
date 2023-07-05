<script lang="ts">
	import DesktopPageLink from "./DesktopPageLink.svelte";
	import GithubIcon from "~icons/bi/github";
	import MoonIcon from "~icons/lucide/moon";
	import SunIcon from "~icons/lucide/sun";

	export let modoOscuro: boolean;
	export let alternarModoOscuro: () => void;
</script>

<div class="hidden items-center space-x-3 divide-x divide-border sm:flex">
	<nav class="text-sm font-medium">
		<ul class="flex gap-2">
			<li><DesktopPageLink href="/" label="Inicio" /></li>
			<li><DesktopPageLink href="/materias" label="Materias" /></li>
			<li><DesktopPageLink href="/docentes" label="Docentes" /></li>
			<li><DesktopPageLink href="https://dollyfiuba.com/" label="Dolly" /></li>
		</ul>
	</nav>
	<div class="flex items-center gap-2 pl-3 text-slate-500 dark:text-slate-200">
		<button
			class="with-ring p-2 transition-colors duration-100 hover:text-slate-700 focus:text-slate-700 focus:ring-fiuba dark:hover:text-slate-100 dark:focus:text-slate-100"
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
			class="with-ring p-2 transition-colors duration-100 hover:text-slate-700 focus:text-slate-700 focus:ring-fiuba dark:hover:text-slate-100 dark:focus:text-slate-100"
			target="_blank"
		>
			<GithubIcon />
		</a>
	</div>
</div>
