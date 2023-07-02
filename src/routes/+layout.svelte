<script lang="ts">
	import PageHeader from "$components/PageHeader.svelte";
	import "../app.css";
</script>

<PageHeader />

<div class="container mx-auto flex-1">
	<slot />
</div>
