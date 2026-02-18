<script lang="ts">
  import { page } from "$app/state";
  import { PUBLIC_GITHUB_URL } from "$env/static/public";
  import colectivo99 from "$lib/assets/colectivo-99.webp";
  import colectivo400 from "$lib/assets/colectivo-400.webp";
  import colectivo404 from "$lib/assets/colectivo-404.webp";
  import colectivo500 from "$lib/assets/colectivo-500.webp";

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

<div class="relative isolate overflow-hidden">
  <div class="pointer-events-none fixed inset-0 z-0">
    <div
      class="absolute -top-24 -left-28 h-112 w-180 rounded-full bg-fiuba/25 blur-[140px]"
      aria-hidden="true"
    ></div>
    <div
      class="absolute -right-16 -bottom-16 h-56 w-56 rounded-full bg-fiuba/20 blur-[110px]"
      aria-hidden="true"
    ></div>
  </div>

  {#if page.error}
    <div class="relative z-10 mx-auto space-y-6 p-6 sm:p-8">
      <img
        src={imagenColectivo(codigoError)}
        alt={`Error ${codigoError}`}
        class="mx-auto max-h-[468px]"
      />
      <div class="mx-auto max-w-[512px] space-y-4 text-center">
        <div class="space-y-2">
          <h1 class="text-4xl font-semibold sm:text-6xl">Error {codigoError}</h1>
          <p>
            {page.error.message}
          </p>
        </div>
        <div class="mx-auto h-px w-[72px] bg-foreground/30" aria-hidden="true"></div>
        <div class="text-sm text-foreground/75">
          Si considerás que este error no debería haber ocurrido podes reportarlo en el
          <a
            href={PUBLIC_GITHUB_URL}
            target="_blank"
            rel="noreferrer"
            class="underline underline-offset-2"
          >
            repositorio
          </a>
          o escribirme un mail
          <a href="mailto:ccastillo@fi.uba.ar" class="underline underline-offset-2">acá</a>
          .
        </div>
      </div>
    </div>
  {/if}
</div>
