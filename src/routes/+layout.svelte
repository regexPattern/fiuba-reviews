<script lang="ts">
	import Header from "$components/Header.svelte";
	import "../app.css";
</script>

<Header />

<div class="flex flex-1 md:container md:mx-auto">
	<slot />
</div>
