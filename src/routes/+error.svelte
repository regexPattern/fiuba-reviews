<script lang="ts">
  import { page } from "$app/state";
  import bondi404 from "$lib/assets/404.webp";
  import ricardoFort500 from "$lib/assets/500.gif";
</script>

<svelte:head>
  <title>{page.status} | FIUBA Reviews</title>
</svelte:head>

<div class="space-y-6 p-4 text-center sm:mt-6">
  {#if page.status == 404}
    <img
      alt="Colectivo 404."
      src={bondi404}
      class="mx-auto"
      height={350}
      width={525} />
  {:else if page.status == 500}
    <img
      alt="Ricardo Fort bailando."
      src={ricardoFort500}
      class="mx-auto"
      height={350}
      width={525} />
  {/if}
  <p class="divide-x-2 text-3xl font-medium">
    Error <span class="font-mono">{page.status}</span>
  </p>
  {#if page.status == 404}
    <p>No encontramos lo que andabas buscando.</p>
  {:else if page.status == 500}
    <p>Error interno del servidor.</p>
  {/if}
</div>
