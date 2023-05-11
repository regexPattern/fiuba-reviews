<script lang="ts">
	import type { PageData } from "./$types";

	export let data: PageData;
</script>

{#each data.docente as docente}{/each}
