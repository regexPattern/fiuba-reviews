<script lang="ts">
  import PromedioCalificaciones from "./PromedioCalificaiones.svelte";

  let { data } = $props();
</script>

{#each data.docentes as docente (docente.codigo)}
  <section id={docente.codigo}>
    <div>
      <h1>{docente.nombre}</h1>
      <small>{docente.rol}</small>
    </div>

    {#if docente.resumenComentario}
      <div class="border">
        {docente.resumenComentario}
      </div>
    {/if}

    <div>
      <PromedioCalificaciones promedio={docente.promedioCalificaciones} />
    </div>

    <div>
      {#each docente.comentarios as comentario (comentario.codigo)}
        <div>
          ({comentario.cuatrimestre.anio}C{comentario.cuatrimestre.numero})
          {comentario.contenido}
          {#if comentario.esDeDolly}(Dolly){/if}
        </div>
      {/each}
    </div>
  </section>
{/each}
