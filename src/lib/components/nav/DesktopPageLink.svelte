<script lang="ts">
	import { twMerge } from "tailwind-merge";

	export let href: string;
	export let label: string;

	export let external = href.startsWith("http");
</script>

<a
	{href}
	class={twMerge(
		"with-ring p-2 text-slate-500 transition-colors duration-100 hover:text-slate-700 focus:text-slate-700 focus:ring-fiuba dark:text-slate-200 dark:hover:text-slate-100 dark:focus:text-slate-100",
		external && "after:content-['↗']"
	)}
	target={external ? "_blank" : "_self"}>{label}</a
>
