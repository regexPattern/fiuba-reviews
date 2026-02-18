<script lang="ts">
  let { data } = $props();
</script>

<svelte:head>
  <title>FIUBA Reviews • {data.materia.codigo} • {data.materia.nombre}</title>
</svelte:head>

Materia no tiene catedras
