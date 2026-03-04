<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  let { data } = $props();

  let metaTitle = $derived(`${data.materia.codigo} • ${data.materia.nombre} | FIUBA Reviews`);
  let metaDescription = $derived(
    `Visitá la página de ${data.materia.nombre} para ver calificaciones y comentarios de las cátedras.`
  );
  let ogImageUrl = $derived(`https://fiuba-reviews.com/materia/${data.materia.codigo}/og.png`);
  let ogImageAlt = $derived(`FIUBA Reviews Materia ${data.materia.codigo}`);

  onMount(() => {
    if (data.catedras.length > 0) {
      goto(`/materia/${data.materia.codigo}/${data.catedras[0].codigo}`, { replaceState: true });
    }
  });
</script>

<svelte:head>
  <title>{metaTitle}</title>
  <meta name="robots" content="index,follow,max-snippet:-1,max-image-preview:large" />
  <meta name="description" content={metaDescription} />
  <link rel="canonical" href={`https://fiuba-reviews.com/materia/${data.materia.codigo}`} />

  <meta property="og:title" content={metaTitle} />
  <meta property="og:description" content={metaDescription} />
  <meta property="og:image" content={ogImageUrl} />
  <meta property="og:image:alt" content={ogImageAlt} />

  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content={metaTitle} />
  <meta name="twitter:description" content={metaDescription} />
  <meta name="twitter:image" content={ogImageUrl} />
  <meta name="twitter:image:alt" content={ogImageAlt} />
</svelte:head>
