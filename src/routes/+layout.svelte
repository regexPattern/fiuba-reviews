<script lang="ts">
	import PageHeader from "$components/PageHeader.svelte";
	import "../app.css";
</script>

<PageHeader />

<div class="flex flex-1 lg:container lg:mx-auto">
	<slot />
</div>

<footer class="h-16 border-t border-border">
  <div class="container">footer</div>
</footer>
