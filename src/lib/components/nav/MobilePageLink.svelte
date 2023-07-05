<script lang="ts">
	import { twMerge } from "tailwind-merge";

	export let href: string;
	export let label: string;

	export let menuAbierto: boolean;

	export let external = href.startsWith("http");

	function cerrarMenu() {
		menuAbierto = false;
	}
</script>

<a
	{href}
	class={twMerge(
		"with-ring p-2 focus:ring-fiuba dark:focus:text-slate-200",
		external && "after:content-['↗']"
	)}
	on:click={cerrarMenu}
	target={external ? "_blank" : "_self"}>{label}</a
>
