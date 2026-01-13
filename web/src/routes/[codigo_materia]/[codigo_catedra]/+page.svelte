<script lang="ts">
  let { data } = $props();
</script>

{#each data.docentes as docente (docente.codigo)}
  <div>
    <div>
      {docente.promedioCalificaciones?.general.toFixed(1) ?? "—"} - {docente.nombre}
    </div>
    <div>{docente.rol}</div>
    <div>{docente.cantidadCalificaciones} calificaciones</div>
    {#if docente.resumenComentario}
      <div>Resumen: {docente.resumenComentario}</div>
    {/if}
    <div>
      Comentarios:
      {#each docente.comentarios as comentario (comentario.codigo)}
        <div>
          ({comentario.cuatrimestre.anio} - {comentario.cuatrimestre.numero})
          {comentario.contenido}
          {#if comentario.esDeDolly}(Dolly){/if}
        </div>
      {/each}
    </div>
  </div>
{/each}
