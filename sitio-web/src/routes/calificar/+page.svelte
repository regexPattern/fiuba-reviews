<script lang="ts">
  import { resolve } from "$app/paths";
  import BackgroundBlob from "$ui/BackgroundBlob.svelte";
  import { ChevronLeft } from "@lucide/svelte";
  import Form from "./Form.svelte";

  let { data } = $props();

  const metaTitle = "Calificar docente | FIUBA Reviews";
  let metaDescription = $derived(
    `Deja tu calificación y comentario anónimo para el docente ${data.docente.nombre} de la materia ${data.docente.codigoMateria}.`
  );
  let ogImageUrl = "https://fiuba-reviews.com/calificar/og.png";
  let ogImageAlt = $derived(
    `FIUBA Reviews Calificar Docente ${data.docente.nombre} Materia ${data.docente.codigoMateria}`
  );
</script>

<svelte:head>
  <title>{metaTitle}</title>
  <meta name="robots" content="noindex,nofollow" />
  <meta name="description" content={metaDescription} />
  <link rel="canonical" href="https://fiuba-reviews.com/calificar" />

  <meta property="og:title" content={metaTitle} />
  <meta property="og:description" content={metaDescription} />
  <meta property="og:image" content={ogImageUrl} />
  <meta property="og:image:alt" content={ogImageAlt} />

  <meta name="twitter:title" content={metaTitle} />
  <meta name="twitter:description" content={metaDescription} />
  <meta name="twitter:image" content={ogImageUrl} />
  <meta name="twitter:image:alt" content={ogImageAlt} />
</svelte:head>

<div class="relative isolate">
  <BackgroundBlob upperLeft={true} lowerRight={false} />

  <main class="relative z-10 container mx-auto space-y-12 p-6">
    <div class="flex flex-col-reverse justify-between gap-4 md:flex-row md:items-center">
      <div class="flex flex-col gap-2">
        <h1 class="text-4xl font-medium">{data.docente.nombre}</h1>
        {#if data.docente.nombreSiu && data.docente.rol}
          <p class="text-sm">{data.docente.nombreSiu} • {data.docente.rol}</p>
        {/if}
      </div>

      <a
        href={resolve(`/materia/${data.docente.codigoMateria}`)}
        class="flex items-center text-sm underline"
      >
        <ChevronLeft class="size-[18px]" />
        Ir a materia del docente
      </a>
    </div>

    <Form cuatris={data.cuatris} />
  </main>
</div>
