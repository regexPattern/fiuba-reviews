<script lang="ts">
  import { page } from "$app/state";
</script>

{#if page.error}
  <h1>TESTING {page.error.message}</h1>
{/if}
