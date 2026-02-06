<script lang="ts">
  import { Database, LayersPlus } from "@lucide/svelte";
  import type { Component } from "svelte";

  let { data } = $props();
</script>

{#snippet featureCard(Icono: Component, descripcion: string)}
  <article class="flex items-start gap-3 p-4">
    <div class="border p-2 select-none">
      <Icono class="size-[22px]" />
    </div>
    <p>{descripcion}</p>
  </article>
{/snippet}

<main class="container mx-auto grid md:grid-cols-2">
  <section id="hero">
    <div class="text-center">
      <h1 class="xs:text-7xl text-6xl">
        <span class="font-bold tracking-tight text-fiuba">FIUBA</span>
        <span class="font-sans font-semibold tracking-tighter">Reviews</span>
      </h1>

      <p>
        Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros
        estudiantes de la FIUBA. Basado en el legendario Dolly FIUBA.
      </p>
    </div>

    <div class="grid grid-cols-2 gap-4">
      {@render featureCard(
        Database,
        `Usamos los datos originales de Dolly para que podás acceder a los comentarios recolectados durante años.`
      )}
      {@render featureCard(
        LayersPlus,
        `Se agregaron los listados de materias de los nuevos planes.`
      )}
    </div>
  </section>

  <section id="ultimos-comentarios">
    {#each data.comentarios as com (com.codigo)}
      <div>{com.contenido}</div>
    {/each}
  </section>

  <section id="acerca-del-proyecto">
    <h2>Acerca del proyecto</h2>
    <p>
      Cuando entré a la facultad había una aplicación llamada Dolly FIUBA que recolectaba para saber
      en qué cátedras inscribirme para el cuatrimestre entrante, pero me parecía que la página
      podría beneficiarse de algunos cambios para mejorar la experiencia. Así que decidí crear una
      nueva aplicación utilizando los datos que Dolly ya tenía (con el permiso de sus creadores) .
      Ahora que Dolly ya no está en funcionamiento, me parece necesario mantener activa una
      plataforma donde los alumnos puedan comentar sobre sus experiencias con las diferente
      cátedras.
    </p>
  </section>
</main>
