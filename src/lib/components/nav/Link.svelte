<script lang="ts">
	import { createEventDispatcher } from "svelte";

	export let href: string;
	export let label: string;

	let classes = "";
	export { classes as class };

	const external = href.startsWith("http");
	const dispatch = createEventDispatcher();
</script>

<a {href} target={external ? "_blank" : "_self"} class={classes} on:click={() => dispatch("click")}
	>{label}</a
>
