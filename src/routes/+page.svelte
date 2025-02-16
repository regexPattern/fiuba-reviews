<script lang="ts">
  import BuscadorMaterias from "$lib/components/buscador-materias.svelte";
  import Feature from "$lib/components/feature-fiuba-reviews.svelte";
  import Link from "$lib/components/link.svelte";
  import MateriasPopulares from "$lib/components/materias-populares.svelte";
  import type { PageData } from "./$types";
  import "@splidejs/svelte-splide/css";
  import { Cpu, Database, Paintbrush2 } from "lucide-svelte";

  export let data: PageData;
</script>

<svelte:head>
  <title>FIUBA Reviews</title>
</svelte:head>

<div
  class="sticky top-16 bg-fiuba p-3 text-center text-background bg-background/75 backdrop-blur-lg"
>
  <span class=""
    >Ayudame a actualizar los listados de cátedras y docentes para mejorar la
    info de FIUBA Reviews. <Link href="/planes" class="underline" external
      >Más Información.</Link
    ></span
  >
</div>

<main class="container my-12 flex max-w-screen-md flex-col items-center gap-12">
  <div class="space-y-6 text-center">
    <h1 class="text-6xl font-bold tracking-tighter xs:text-7xl">
      <span class="text-fiuba">FIUBA</span> Reviews
    </h1>

    <p class="mx-auto max-w-[40rem] text-lg xs:text-xl">
      Encontrá calificaciones y comentarios de los docentes de la facultad,
      subidos por otros estudiantes de la FIUBA. Basado en el legendario
      <Link href="http://dollyfiuba.com" class="underline after:content-link"
        >Dolly FIUBA</Link
      >.
    </p>
  </div>

  <BuscadorMaterias label="Buscar materias" materias={data.materias} />

  <section class="grid grid-cols-1 gap-10 md:grid-cols-3 md:gap-4 lg:gap-6">
    <Feature
      icon={Database}
      title="Mismos datos"
      desc="Tomamos los datos originales de Dolly para que podás acceder a los mismos comentarios recolectados durante años. Bajo el capó, toda la estructura de los datos fue adaptada a Postgres."
    />
    <Feature
      icon={Paintbrush2}
      title="Nuevo diseño"
      desc="Se reconstruyó completamente la interfaz de la página web y se le dió un estilo totalmente diferente, con modo claro y oscuro, y un estilo minimalista moderno."
    />
    <Feature
      icon={Cpu}
      title="Inteligencia Artificial"
      desc="Utilizando inteligencia artificial, generamos un resumen de lo que dicen los comentarios de los docentes más populares para que te ahorrés tiempo al evaluarlos."
    />
  </section>

  <section class="w-full space-y-4 overflow-x-hidden">
    <h2 class="text-center text-4xl font-semibold tracking-tight">
      Materias Más Populares
    </h2>
    <p class="text-center text-muted-foreground">
      Las {data.materiasMasPopulares.length} materias más cursadas de la facultad
      en base a la cantidad de planes de carrera en las que están presentes.
    </p>
    <MateriasPopulares materias={data.materiasMasPopulares.slice(0, 10)} />
    <MateriasPopulares materias={data.materiasMasPopulares.slice(10, 20)} />
  </section>

  <section class="space-y-4 text-center text-muted-foreground">
    <h2 class="text-4xl font-semibold tracking-tight text-foreground">
      Acerca del Proyecto
    </h2>
    <p>
      Desde que entré a la facultad he usado Dolly para saber en qué cátedras
      inscribirme para el cuatrimestre entrante. Sin embargo, me parecía que la
      página podría beneficiarse de algunos cambios para mejorar la experiencia.
      Así que decidí crear una nueva aplicación utilizando los datos que Dolly
      ya tenía
      <Link
        href="https://github.com/lugfi/dolly/issues/80"
        class="after:content-link"
      >
        (con el permiso de sus creadores)
      </Link>
      .
    </p>
    <p>
      Este proyecto no pretende ser un reemplazo a la aplicación original, sino
      más bien una propuesta de posibles cambios que me gustaría ver
      implementados o funcionalidades adicionales que me parece que resultan
      útiles para los estudiantes.
    </p>
  </section>

  <!-- Elemento que se renderiza tras el contenido de la página principal para
       mostrar el efecto de background blur con un blob de colores. -->
  <div
    class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
    aria-hidden="true"
  >
    <div
      class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-fiuba to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
      style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)"
    />

    <div
      class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
      aria-hidden="true"
    />
  </div>
</main>
