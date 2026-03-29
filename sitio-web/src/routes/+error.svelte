<script lang="ts">
  import { PUBLIC_REPO_URL } from "$env/static/public";
  import { page } from "$app/state";
  import colectivo99 from "$lib/assets/colectivo-99.webp";
  import colectivo400 from "$lib/assets/colectivo-400.webp";
  import colectivo404 from "$lib/assets/colectivo-404.webp";
  import colectivo500 from "$lib/assets/colectivo-500.webp";
  import BackgroundBlob from "$ui/BackgroundBlob.svelte";
  import ExternalLink from "$ui/ExternalLink.svelte";

  function imagenColectivo(codigo: number) {
    if (codigo === 400) {
      return colectivo400;
    } else if (codigo === 404) {
      return colectivo404;
    } else if (codigo === 500) {
      return colectivo500;
    } else {
      return colectivo99;
    }
  }

  const codigoError = $derived(page.status);
</script>

<svelte:head>
  <meta name="robots" content="noindex,nofollow" />
</svelte:head>

<div class="relative isolate">
  <BackgroundBlob upperLeft={true} lowerRight={true} />

  {#if page.error}
    <div class="relative z-10 mx-auto space-y-6 p-6 sm:p-8">
      <img
        src={imagenColectivo(codigoError)}
        alt={`Error ${codigoError}`}
        class="mx-auto max-h-117"
      />
      <div class="mx-auto max-w-lg space-y-4 text-center">
        <div class="space-y-2">
          <h1 class="text-4xl font-semibold sm:text-6xl">Error {codigoError}</h1>
          <p>
            {page.error.message}
          </p>
        </div>
        <div class="mx-auto h-px w-18 bg-foreground/30" aria-hidden="true"></div>
        <div class="text-sm text-foreground/75">
          Si considerás que este error no debería haber ocurrido podes reportarlo en el
          <ExternalLink href={PUBLIC_REPO_URL} class="underline underline-offset-2">
            repositorio
          </ExternalLink>
          o escribirme un mail
          <a href="mailto:ccastillo@fi.uba.ar" class="underline underline-offset-2">acá</a>
          .
        </div>
      </div>
    </div>
  {/if}
</div>
