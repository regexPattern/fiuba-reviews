<script lang="ts">
	import PageHeader from "$components/PageHeader.svelte";
	import "../app.css";
</script>

<PageHeader />

<div class="flex flex-1 lg:container lg:mx-auto">
	<slot />
</div>
