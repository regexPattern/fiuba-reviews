<script lang="ts">
  import Sidebar from "./Sidebar.svelte";

  let { children, data } = $props();
</script>

<div class="container mx-auto flex">
  {#if data.catedras.length > 0}
    <Sidebar materia={data.materia} catedras={data.catedras} />
  {/if}

  <main class="w-full">
    {@render children()}
  </main>
</div>
