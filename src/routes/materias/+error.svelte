<script lang="ts">
	import { page } from "$app/stores";
</script>

{#if $page.status === 404}
	<h1>404 - Materia no encontrada</h1>
{:else}
	<h1>Otro error - Materia no encontrada</h1>
{/if}
