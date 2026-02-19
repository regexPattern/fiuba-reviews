<script lang="ts">
  let { children, data } = $props();
</script>

<svelte:head>
  <meta name="robots" content="noindex,follow,max-snippet:-1,max-image-preview:large" />
  <link
    rel="canonical"
    href={`https://fiuba-reviews.com/materia/${encodeURIComponent(data.materia.codigo)}`}
  />
  <!-- og:image para páginas de cátedra (descomentar cuando esté implementado) -->
  <!-- <meta property="og:image" content={`${META_BASE_URL}/og/catedra/${data.materia.codigo}.png`} /> -->
  <!-- <meta name="twitter:image" content={`${META_BASE_URL}/og/catedra/${data.materia.codigo}.png`} /> -->
</svelte:head>

{@render children()}
