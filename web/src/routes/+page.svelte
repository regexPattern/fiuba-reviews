<script lang="ts">
  import { Database, LayersPlus } from "@lucide/svelte";
  import type { Options as SplideOptions } from "@splidejs/splide";
  import { AutoScroll } from "@splidejs/splide-extension-auto-scroll";
  import { Splide, SplideSlide } from "@splidejs/svelte-splide";
  import "@splidejs/svelte-splide/css";
  import type { Component } from "svelte";

  let { data } = $props();

  const mitadComentarios = $derived(Math.ceil(data.comentarios.length / 2));
  const comentariosIzquierda = $derived(data.comentarios.slice(0, mitadComentarios));
  const comentariosDerecha = $derived(data.comentarios.slice(mitadComentarios));

  const SPLIDE_OPTS: SplideOptions = {
    type: "loop",
    direction: "ttb",
    height: "22rem",
    autoHeight: true,
    perPage: 4,
    perMove: 1,
    gap: "0.75rem",
    arrows: false,
    pagination: false,
    autoScroll: { speed: 0.3, autoStart: true, pauseOnHover: false, pauseOnFocus: false }
  };
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
    <div class="grid grid-cols-2 gap-4">
      <div class="relative">
        <Splide
          options={SPLIDE_OPTS}
          extensions={{ AutoScroll }}
          aria-label="Últimos comentarios izquierda"
        >
          {#each comentariosIzquierda as comentario (comentario.codigo)}
            <SplideSlide>
              <article class="border border-button-border bg-button-background p-2">
                <p class="line-clamp-4 text-sm leading-5">{comentario.contenido}</p>
                {comentario.nombreDocente}
              </article>
            </SplideSlide>
          {/each}
        </Splide>
        <div class="pointer-events-none absolute inset-0">
          <div
            class="absolute top-0 right-0 left-0 h-16 bg-linear-to-b from-background to-transparent"
          ></div>
          <div
            class="absolute right-0 bottom-0 left-0 h-16 bg-linear-to-t from-background to-transparent"
          ></div>
        </div>
      </div>

      <div class="relative">
        <Splide
          options={SPLIDE_OPTS}
          extensions={{ AutoScroll }}
          aria-label="Últimos comentarios derecha"
        >
          {#each comentariosDerecha as comentario (comentario.codigo)}
            <SplideSlide>
              <article class="translate-y-12 border border-button-border bg-button-background p-2">
                <p class="line-clamp-4 text-sm leading-5">{comentario.contenido}</p>
              </article>
            </SplideSlide>
          {/each}
        </Splide>
        <div class="pointer-events-none absolute inset-0">
          <div
            class="absolute top-0 right-0 left-0 h-8 bg-gradient-to-b from-background to-transparent"
          />
          <div
            class="absolute right-0 bottom-0 left-0 h-8 bg-gradient-to-t from-background to-transparent"
          />
        </div>
      </div>
    </div>
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
