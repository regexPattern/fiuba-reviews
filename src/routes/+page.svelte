<script lang="ts">
  import TarjetaFeature from "./tarjeta-feature.svelte";
  import Link from "$lib/components/link.svelte";
  import MateriasPopulares from "$lib/components/materias-populares.svelte";
  import "@splidejs/svelte-splide/css";
  import { BookCheck, Database } from "lucide-svelte";
  import Annotation from "svelte-rough-notation";

  import type { PageData } from "./$types";
  import { onMount } from "svelte";

  export let data: PageData;

  let rn;
  let anotacionVisible = false;

  onMount(() => {
    setTimeout(() => {
      anotacionVisible = true;
    }, 250);
  });
</script>

<svelte:head>
  <title>FIUBA Reviews</title>
</svelte:head>

<div
  class="sticky top-16 z-50 bg-background/75 bg-fiuba p-3 text-center text-background backdrop-blur-lg">
  <span class=""
    >⚠️ Estamos actualizando los listados de cátedras y docentes para las
    materias de los planes nuevos.
    <Link href="/planes" class="font-semibold underline" external
      >Más Información.</Link
    ></span>
</div>

<main class="container my-12 flex max-w-screen-md flex-col items-center gap-12">
  <div class="space-y-6 text-center">
    <h1 class="text-6xl xs:text-7xl">
      <span class="font-serif font-bold tracking-tight text-fiuba">FIUBA</span>
      <Annotation
        bind:visible={anotacionVisible}
        bind:this={rn}
        type="highlight"
        color="#4EACD4"
        ><span class="font-semibold tracking-tighter">Reviews</span
        ></Annotation>
    </h1>

    <p class="mx-auto max-w-[40rem] text-lg xs:text-xl">
      Encontrá calificaciones y comentarios de los docentes de la facultad,
      subidos por otros estudiantes de la FIUBA. Basado en el legendario
      <Link
        href="https://github.com/lugfi/dolly"
        class="underline after:content-link">Dolly FIUBA</Link
      >.
    </p>
  </div>

  <section class="grid grid-cols-1 gap-10 md:grid-cols-2 md:gap-4 lg:gap-6">
    <TarjetaFeature
      icon={Database}
      title="Mismos Datos"
      desc="Usamos los datos originales de Dolly para que podás acceder a los comentarios recolectados durante años. Bajo el capó, toda la estructura de los datos fue reorganizada." />
    <TarjetaFeature
      icon={BookCheck}
      title="Nuevos Planes"
      desc="Se agregaron los listados de materias de los nuevos planes. Poco a poco se van a ir agregando las cátedras correspondientes a estas nuevas materias." />
  </section>

  <section class="w-full space-y-4 overflow-x-hidden">
    <h2 class="text-center text-4xl font-semibold tracking-tight">
      Materias Más Populares
    </h2>
    <p class="text-center text-muted-foreground">
      Las {data.materiasMasPopulares.length} materias más cursadas de la facultad
      en base a la cantidad de planes de carrera en las que están presentes. Con
      los nuevos planes se unificaron bastantes materias que antes eran distinguidas
      por diferente código, pero realmente eran comunes a varias carreras.
    </p>
    <MateriasPopulares materias={data.materiasMasPopulares.slice(0, 10)} />
    <MateriasPopulares materias={data.materiasMasPopulares.slice(10, 20)} />
    <MateriasPopulares materias={data.materiasMasPopulares.slice(20, 30)} />
  </section>

  <section class="space-y-4">
    <h2 class="text-center text-4xl font-semibold tracking-tight">
      Últimos Comentarios
    </h2>
    <p class="text-center text-muted-foreground">
      Los comentarios están hechos por los mismos alumnos de manera anónima,
      para recomendarte o advertirte cuando vayas a cursar una materia.
    </p>
    <div class="grid gap-4 sm:grid-cols-2">
      <div class="flex flex-col gap-4">
        {#each data.ultimosComentarios.slice(0, 2) as com (com.codigo)}
          <div class="divide-y rounded-lg border [&>*]:p-3">
            <p class="font-medium text-fiuba">{com.nombreDocente}</p>
            <p class={`before:content-['"'] after:content-['"']`}>
              {com.contenido}
            </p>
          </div>
        {/each}
      </div>
      <div class="flex flex-col gap-4">
        {#each data.ultimosComentarios.slice(2, 4) as com (com.codigo)}
          <div class="divide-y rounded-lg border [&>*]:p-3">
            <p class="font-medium text-fiuba">{com.nombreDocente}</p>
            <p class={`before:content-['"'] after:content-['"']`}>
              {com.contenido}
            </p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <section class="space-y-4">
    <h2 class="text-center text-4xl font-semibold tracking-tight">
      Estadísticas
    </h2>
    <p class="text-center text-muted-foreground">
      Continuamos el legado de Dolly aprovechando los casi 10 años de datos
      recolectados por la aplicación original mientras estuvo en actividad.
    </p>
    <div
      class="grid w-full grid-cols-1 gap-3 2xs:grid-cols-2 sm:grid-cols-3 md:gap-4">
      <div class="rounded-lg border p-4">
        <p class="font-mono text-4xl font-bold md:text-5xl">
          >{data.cantidadCatedras}
        </p>
        <p>Cátedras</p>
      </div>
      <div class="rounded-lg border p-4">
        <p class="font-mono text-4xl font-bold md:text-5xl">
          >{data.cantidadComentarios}
        </p>
        <p>Comentarios</p>
      </div>
      <div
        class="mx-auto w-full rounded-lg border p-4 2xs:col-span-2 2xs:min-w-min 2xs:max-w-[50%] sm:col-span-1 sm:w-full sm:max-w-full">
        <p class="font-mono text-4xl font-bold md:text-5xl">
          >{data.cantidadCalificaciones}
        </p>
        <p>Calificaciones</p>
      </div>
    </div>
  </section>

  <section class="space-y-4 text-center text-muted-foreground">
    <h2 class="text-4xl font-semibold tracking-tight text-foreground">
      Acerca del Proyecto
    </h2>
    <p>
      Desde que entré a la facultad usé Dolly para saber en qué cátedras
      inscribirme para el cuatrimestre entrante, pero me parecía que la página
      podría beneficiarse de algunos cambios para mejorar la experiencia. Así
      que decidí crear una nueva aplicación utilizando los datos que Dolly ya
      tenía
      <Link
        href="https://github.com/lugfi/dolly/issues/80"
        class="after:content-link">
        (con el permiso de sus creadores)
      </Link>
      .
    </p>
    <p>
      Ahora que Dolly ya no está en funcionamiento, me parece necesario mantener
      activa una plataforma donde los alumnos puedan comentar sobre sus
      experiencias con las diferente cátedras.
    </p>
  </section>

  <!-- Elemento que se renderiza tras el contenido de la página principal para
       mostrar el efecto de background blur con un blob de colores. -->
  <div
    class="absolute inset-x-0 -top-40 -z-10 transform-gpu overflow-hidden blur-3xl sm:-top-80"
    aria-hidden="true">
    <div
      class="relative left-[calc(50%-11rem)] aspect-[1155/678] w-[36.125rem] -translate-x-1/2 rotate-[30deg] bg-gradient-to-tr from-fiuba to-[#9089fc] opacity-30 sm:left-[calc(50%-30rem)] sm:w-[72.1875rem]"
      style="clip-path: polygon(74.1% 44.1%, 100% 61.6%, 97.5% 26.9%, 85.5% 0.1%, 80.7% 2%, 72.5% 32.5%, 60.2% 62.4%, 52.4% 68.1%, 47.5% 58.3%, 45.2% 34.5%, 27.5% 76.7%, 0.1% 64.9%, 17.9% 100%, 27.6% 76.8%, 76.1% 97.7%, 74.1% 44.1%)" />

    <div
      class="absolute inset-x-0 top-[calc(100%-13rem)] -z-10 transform-gpu overflow-hidden blur-3xl sm:top-[calc(100%-30rem)]"
      aria-hidden="true" />
  </div>
</main>
