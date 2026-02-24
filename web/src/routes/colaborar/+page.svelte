<script lang="ts">
  import Form from "./Form.svelte";
  import Instrucciones from "./Instrucciones.svelte";

  let { data } = $props();

  const metaTitle = "Colaborar | FIUBA Reviews";
  const metaDescription =
    "Ayudá a actualizar las ofertas de comisiones disponibles en FIUBA Reviews enviando tu oferta más actualizada.";
  const ogImageUrl = "https://fiuba-reviews.com/colaborar/og.png";
  const ogImageAlt = "FIUBA Reviews Colaborar";
</script>

<svelte:head>
  <title>{metaTitle}</title>
  <meta name="robots" content="index,follow,max-snippet:-1,max-image-preview:large" />
  <meta name="description" content={metaDescription} />
  <link rel="canonical" href="https://fiuba-reviews.com/colaborar" />

  <meta property="og:title" content={metaTitle} />
  <meta property="og:description" content={metaDescription} />
  <meta property="og:image" content={ogImageUrl} />
  <meta property="og:image:alt" content={ogImageAlt} />

  <meta name="twitter:card" content="summary" />
  <meta name="twitter:title" content={metaTitle} />
  <meta name="twitter:description" content={metaDescription} />
  <meta name="twitter:image" content={ogImageUrl} />
  <meta name="twitter:image:alt" content={ogImageAlt} />
</svelte:head>

<div class="relative isolate">
  <div class="fondo-decorativo fondo-superior" aria-hidden="true"></div>

  <main class="relative z-10 container mx-auto grid items-start lg:grid-cols-2">
    <section class="space-y-4 p-6">
      <h1 class="text-2xl font-semibold">Actualización de ofertas</h1>

      <p>
        La página de la oferta horaria de la que Dolly FIUBA solía obtener el listado de cátedras de
        cada materia dejó de actualizarse desde el 1C2024. Para solventar esto de manera
        automatizada se quiere de acceso a las ofertas de todas las carreras en cada cuatrimestre.
        Sin embargo, yo solo puedo acceder a la oferta para las materias de mi carrera, no de las
        materias que no están en mi plan.
      </p>
      <p>
        Para esto desde el 1C2025 FIUBA Reviews actualiza sus ofertas gracias a la colaboración de
        los alumnos de las diferentes carreras de la facultad que envían una copia de las ofertas de
        sus carreras para poder generar los listados correctos.
      </p>

      <table class="mx-auto caption-bottom text-center">
        <caption class="p-2 text-sm">
          Cuatrimestre de última actualización de la oferta de comisiones de cada carrera.
        </caption>
        <thead class="border border-button-border bg-button-background text-button-foreground">
          <tr class="divide-x divide-button-border">
            <th class="p-2 font-medium">Carrera</th>
            <th class="px-4 font-medium">Última actualización</th>
          </tr>
        </thead>
        <tbody class="border border-button-border">
          {#each data.actualizacionesOfertas as { carrera, cuatrimestre } (carrera.codigo)}
            <tr class="divide-x divide-button-border">
              <td class="px-4 py-1">{carrera.nombre}</td>
              <td>{cuatrimestre.numero}C{cuatrimestre.anio}</td>
            </tr>
          {/each}
        </tbody>
      </table>
    </section>

    <section class="space-y-6 p-6">
      <Instrucciones />
      <Form />
    </section>
  </main>
</div>
