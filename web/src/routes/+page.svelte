<script lang="ts">
  let { data } = $props();
</script>

{#each data.materias as materia (materia.codigo)}
  <div>
    <a href={`${materia.codigo}`}>{materia.codigo} - {materia.nombre}</a>
  </div>
{/each}
