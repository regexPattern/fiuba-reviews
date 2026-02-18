<script lang="ts">
  import type { Component } from "svelte";
  import { Database, GraduationCap, HatGlasses, LayersPlus } from "@lucide/svelte";
  import BuscadorMateriaTrigger from "$lib/componentes/buscador/BuscadorMateriaTrigger.svelte";
  import "@fontsource/noto-serif/600.css";

  let { data } = $props();
</script>

<div class="relative isolate overflow-hidden">
  <div class="pointer-events-none fixed inset-0 z-0">
    <div
      class="absolute -top-24 -left-28 h-112 w-180 rounded-full bg-fiuba/25 blur-[140px]"
      aria-hidden="true"
    ></div>
  </div>

  <main class="container mx-auto grid gap-12 p-6 lg:grid-cols-2 lg:flex-row">
    <section id="hero" class="space-y-8 self-center">
      <div class="mx-auto space-y-4 text-center lg:max-w-[512px]">
        <h1 class="font-noto-serif text-5xl sm:text-6xl font-semibold tracking-tight">
          <span class="text-fiuba">FIUBA</span>
          <span>Reviews</span>
        </h1>
        <p class="mx-auto max-w-[468px] font-medium">
          Encontrá calificaciones y comentarios de los docentes de la facultad, subidos por otros
          estudiantes. Basado en el legendario Dolly FIUBA.
        </p>
      </div>

      <div class="flex justify-center">
        <BuscadorMateriaTrigger variante="hero" />
      </div>

      <div class="mx-auto grid max-w-[620px] gap-6 sm:grid-cols-2">
        {#snippet tarjetaFeature(Icono: Component, titulo: string, descripcion: string)}
          <article class="flex flex-1 items-start gap-3">
            <div class="p-2 select-none">
              <Icono class="size-[22px] stroke-fiuba" />
            </div>
            <div class="space-y-2">
              <p class="font-medium">{titulo}</p>
              <p>{descripcion}</p>
            </div>
          </article>
        {/snippet}

        {@render tarjetaFeature(
          HatGlasses,
          "Reviews anónimas",
          `Las calificaciones y comentarios agregados son totalmente anónimos.`
        )}
        {@render tarjetaFeature(
          Database,
          "Mismos datos de Dolly",
          `Usamos los datos originales de Dolly recolectados durante más de 10 años.`
        )}
        {@render tarjetaFeature(
          LayersPlus,
          "Nuevos planes",
          `Se agregaron las materias de los nuevos planes manteniendo comentarios de sus equivalencias.`
        )}
        {@render tarjetaFeature(
          GraduationCap,
          "Todas las ingenierías",
          `Están disponibles todas las materias de las 11 carreras de ingeniería.`
        )}
      </div>
    </section>

    <section id="ultimos-comentarios" class="space-y-4 md:h-[630px]">
      {#snippet filaComentarios(comentarios: typeof data.comentarios, claseAnimacion: string)}
        <div class="overflow-hidden">
          <div class={`flex w-max gap-4 ${claseAnimacion} will-change-transform`}>
            {#each [...comentarios, ...comentarios] as com, i (`fila-${com.codigo}-${i}`)}
              <article
                class="max-w-[260px] min-w-[260px] shrink-0 border border-button-border bg-button-background/50 p-4"
              >
                <p
                  class={`comentario-contenido line-clamp-4 before:content-['"'] after:content-['"'] md:line-clamp-none`}
                >
                  {com.contenido}
                </p>
                <p class="text-muted-foreground mt-2 text-sm text-foreground/75">
                  {com.nombreDocente} • {com.nombreMateria}
                </p>
              </article>
            {/each}
          </div>
        </div>
      {/snippet}

      {#snippet columnaComentarios(comentarios: typeof data.comentarios)}
        <div class="h-[630px] overflow-hidden">
          <div
            class="flex animate-[scroll-vertical_50s_linear_infinite] flex-col gap-4 will-change-transform md:animate-[scroll-vertical_30s_linear_infinite]"
          >
            {#each [...comentarios, ...comentarios] as com, i (`${com.codigo}-${i}`)}
              <article class="border border-button-border bg-button-background/50 p-4">
                <p
                  class={`comentario-contenido line-clamp-4 before:content-['"'] after:content-['"'] md:line-clamp-none`}
                >
                  {com.contenido}
                </p>
                <p class="text-muted-foreground mt-2 text-sm text-foreground/75">
                  {com.nombreDocente} • {com.nombreMateria}
                </p>
              </article>
            {/each}
          </div>
        </div>
      {/snippet}

      <div class="grid gap-4 md:hidden">
        {@render filaComentarios(
          data.comentarios.filter((_, i) => i % 2 === 0),
          "animate-[scroll-horizontal_68s_linear_infinite] sm:animate-[scroll-horizontal_56s_linear_infinite]"
        )}
        {@render filaComentarios(
          data.comentarios.filter((_, i) => i % 2 === 1),
          "animate-[scroll-horizontal-reverse_64s_linear_infinite] sm:animate-[scroll-horizontal-reverse_52s_linear_infinite]"
        )}
      </div>

      <div class="hidden grid-cols-2 gap-4 md:grid md:h-[630px]">
        {@render columnaComentarios(data.comentarios.slice(0, data.comentarios.length / 2))}
        {@render columnaComentarios(data.comentarios.slice(data.comentarios.length / 2))}
      </div>
    </section>

    <style>
      @keyframes scroll-horizontal {
        0% {
          transform: translateX(0);
        }
        100% {
          transform: translateX(-50%);
        }
      }

      @keyframes scroll-horizontal-reverse {
        0% {
          transform: translateX(-50%);
        }
        100% {
          transform: translateX(0);
        }
      }

      @keyframes scroll-vertical {
        0% {
          transform: translateY(0);
        }
        100% {
          transform: translateY(-50%);
        }
      }
    </style>

    <section id="materias-populares" class="space-y-4 text-center">
      <h2 class="text-2xl font-semibold">Materias más populares</h2>
      <p>
        Las {data.materiasPopulares.length} más populares según la cantidad carreras que la cursan. Con
        los nuevos planes se unificaron bastantes materias que antes eran distinguidas por diferente código,
        pero realmente eran comunes a varias carreras.
      </p>

      <div class="mt-4 flex flex-wrap justify-center gap-2">
        {#each data.materiasPopulares as materia (materia.codigo)}
          <a
            href={`/materia/${materia.codigo}`}
            class="rounded-full border border-button-border bg-button-background px-3 py-1 text-sm"
          >
            {materia.nombre} • {materia.codigo}
          </a>
        {/each}
      </div>
    </section>

    <section id="acerca-del-proyecto" class="space-y-4 text-center">
      <h2 class="text-2xl font-semibold">Acerca del proyecto</h2>
      <p>
        <a href="https://github.com/lugfi/dolly">Dolly FIUBA</a> era el sitio original en donde los estudiantes
        de FIUBA publicaban calificaciones y comentarios de los docentes con los que cursaban. Como alumno,
        desde que entré a la facultad fue un recurso invaluable al momento de elegir cátedras al iniciar
        cada cuatrimestre.
      </p>
      <p>
        Ahora que Dolly ya no está en funcionamiento, me parece necesario mantener activa una
        plataforma donde los alumnos puedan comentar sobre sus experiencias con las diferente
        cátedras.
      </p>
    </section>
  </main>
</div>
